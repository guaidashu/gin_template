/**
  create by yy on 2019-10-10
*/

package libs

import (
	"errors"
	"fmt"
	"gin_template/app/config"
	"gopkg.in/mgo.v2/bson"
	"os"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

func GetErrorString(err error) string {
	return fmt.Sprintf("error: %v", err)
}

func NewReportError(err error) error {
	if !config.Config.App.DEBUG {
		return err
	}
	_, fileName, line, _ := runtime.Caller(1)
	data := fmt.Sprintf("%v, report in: %v: in line %v", err, fileName, line)
	return errors.New(data)
}

func DebugPrint(format string, values ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	_, _ = fmt.Fprintf(os.Stderr, "[Guaidashu-debug] "+format, values...)
}

type Struct2MapValue struct {
	Value reflect.Value
	Type  reflect.Type
}

// 字段tag 加上 struct2map:"-" 则会被忽略
// 例:
// type structToMap struct {
//	   Data int64 `json:"data" struct2map:"-"`
// }
// 从结构体转为 map (适用于 数据库更新数据)
func Struct2Map(v interface{}, maps map[string]interface{}, originValue ...Struct2MapValue) (data map[string]interface{}, err error) {
	var (
		vt  reflect.Type
		val reflect.Value
	)

	if maps != nil {
		data = maps
	} else {
		data = make(map[string]interface{})
	}

	if len(originValue) > 0 {
		vt = originValue[0].Type
		val = originValue[0].Value
	} else {
		vt = reflect.TypeOf(v)
		val = reflect.ValueOf(v)
	}

	if kd := vt.Kind(); kd != reflect.Struct {
		if kd := vt.Elem().Kind(); kd != reflect.Struct {
			return nil, errors.New("the param is not a struct, please make sure")
		} else {
			vt = vt.Elem()
			val = val.Elem()
			if vt.Kind() != reflect.Struct {
				return nil, errors.New("the param is not a struct, please make sure")
			}
		}
	}

	fieldNum := val.NumField()

	for i := 0; i < fieldNum; i++ {
		fieldIgnore := vt.Field(i).Tag.Get("struct2map")
		if fieldIgnore == "-" {
			continue
		}

		fieldName := vt.Field(i).Tag.Get("json")
		if fieldName == "" {
			data, _ = Struct2Map(val.Field(i), data, Struct2MapValue{
				Value: val.Field(i),
				Type:  vt.Field(i).Type,
			})
			continue
		}

		switch val.Field(i).Type().String() {
		case "string":
			data[fieldName] = val.Field(i).String()

		case "*int64", "*int32", "*int16", "*int8", "*int":
			if val.Field(i).IsNil() {
				continue
			}
			data[fieldName] = val.Field(i).Elem().Int()

		case "*uint64", "*uint32", "*uint16", "*uint8", "*uint":
			if val.Field(i).IsNil() {
				continue
			}
			data[fieldName] = val.Field(i).Elem().Uint()

		case "uint64", "uint32", "uint16", "uint8", "uint":
			data[fieldName] = val.Field(i).Uint()

		case "int64", "int32", "int16", "int8", "int":
			data[fieldName] = val.Field(i).Int()
		default:
			data, _ = Struct2Map(val.Field(i), data, Struct2MapValue{
				Value: val.Field(i),
				Type:  vt.Field(i).Type,
			})
		}
	}

	return
}

func GenerateDataId() string {
	id := bson.NewObjectId().Hex()
	return id
}

// 以下两个函数，建议在做高频转换且不做修改的情况下使用

// bytes转string，不做内存拷贝，二者共用同一份内存数据（转换后修改bytes里面的内容，string同时会做修改）
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// string转bytes，不做内存拷贝，二者共用同一份内存数据（转换后修改bytes里面的内容会panic）
func String2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}