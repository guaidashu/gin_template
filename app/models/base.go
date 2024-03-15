/**
  create by yy on 2019-08-31
*/

package models

import (
	"errors"
	"fmt"
	"gin_template/app/enum"
	"gorm.io/gorm"
	"reflect"
)

type BaseModel interface {
	TableName() string            // 获取数据表名
	getDB() *gorm.DB              // 获取db实例
	getDBWithNoDeleted() *gorm.DB // 获取未被软删除的db实例
	CreateTable() error           // 创建表
	HasTable() bool               // 表是否存在
}

type Model struct {
	CreatedAt int64 `json:"created"`
	UpdatedAt int64 `json:"updated"`
	DeletedAt int64 `json:"deleted"`
}

// 统一是软删除
// 获取未删除db 实例
func getTableWithNoDeleted(gdb *gorm.DB, tableName string) *gorm.DB {
	return gdb.Table(tableName).Unscoped()
}

// 获取所有数据(不论是否被删除) db实例
func getTable(gdb *gorm.DB, tableName string) *gorm.DB {
	return gdb.Table(tableName)
}

func getOrderByStr(orderField, defaultOrder string, orderType int64) (orderBy string) {
	if orderField != "" {
		orderBy = orderField
	} else {
		orderBy = defaultOrder
	}

	if orderType > 0 {
		orderBy = orderBy + " desc"
	} else {
		orderBy = orderBy + " asc"
	}
	return
}

func getOffset(page, pageSize int) int {
	offset := (page - 1) * pageSize
	if offset <= 0 {
		offset = 0
	}

	return offset
}

type separateHandle func(*gorm.DB, interface{}) *gorm.DB

// auto struct gorm.DB
// 1 -> and
// 2 -> or
// 3 -> separate handle 单独处理
// 4 -> like 操作
// separateHandle 单独处理 的字符串数组
func queryAutoWhere(db *gorm.DB, search interface{}, separate ...map[string]separateHandle) (*gorm.DB, error) {
	var (
		err  error
		data interface{}
	)

	vt := reflect.TypeOf(search)
	val := reflect.ValueOf(search)

	if kd := vt.Kind(); kd != reflect.Struct {
		if kd := vt.Elem().Kind(); kd != reflect.Struct {
			return db, errors.New("the param is not a struct, please make sure")
		} else {
			vt = vt.Elem()
			val = val.Elem()
		}
	}

	fieldNum := val.NumField()

	for i := 0; i < fieldNum; i++ {
		handleType := vt.Field(i).Tag.Get(enum.FieldIsHandle)
		if handleType == "" {
			continue
		}
		switch val.Field(i).Type().String() {
		case "string":
			tmp := val.Field(i).String()
			if tmp == "" {
				continue
			}
			data = tmp
		case "*int64", "*int32", "*int16", "*int8", "*int":
			if val.Field(i).IsNil() {
				continue
			}
			data = val.Field(i).Elem().Int()
		case "*uint64", "*uint32", "*uint16", "*uint8", "*uint":
			if val.Field(i).IsNil() {
				continue
			}
			data = val.Field(i).Elem().Uint()
		case "int64", "int32", "int16", "int8", "int":
			data = val.Field(i).Int()
		case "uint64", "uint32", "uint16", "uint8", "uint":
			data = val.Field(i).Uint()
		}

		switch handleType {
		case enum.AutoWhere:
			db = db.Where(vt.Field(i).Tag.Get(enum.FieldHandle), data)
		case enum.AutoOr:
			db = db.Or(vt.Field(i).Tag.Get(enum.FieldHandle), data)
		case enum.AutoCustomHandle:
			// 根据条件 判断大于和小于
			if len(separate) > 0 {
				db = separate[0][vt.Field(i).Tag.Get("json")](db, data)
			}
		case enum.AutoLike:
			db = db.Where(vt.Field(i).Tag.Get(enum.FieldHandle), "%"+fmt.Sprintf("%v", data)+"%")
		}

	}

	return db, err
}

// 获取可以赋空值的字段 map，用于完善 struct2Map
func NewExcept(list ...string) map[string]int64 {
	m := make(map[string]int64)
	for _, v := range list {
		m[v] = 1
	}

	return m
}

var structMap = map[reflect.Kind]reflect.Kind{
	reflect.Bool:       reflect.Bool,
	reflect.Int:        reflect.Int,
	reflect.Int8:       reflect.Int8,
	reflect.Int16:      reflect.Int16,
	reflect.Int32:      reflect.Int32,
	reflect.Int64:      reflect.Int64,
	reflect.Uint:       reflect.Uint,
	reflect.Uint8:      reflect.Uint8,
	reflect.Uint16:     reflect.Uint16,
	reflect.Uint32:     reflect.Uint32,
	reflect.Uint64:     reflect.Uint64,
	reflect.Uintptr:    reflect.Uintptr,
	reflect.Float32:    reflect.Float32,
	reflect.Float64:    reflect.Float64,
	reflect.Complex64:  reflect.Complex64,
	reflect.Complex128: reflect.Complex128,
	reflect.String:     reflect.String,
}

type Struct2MapValue struct {
	Value reflect.Value
	Type  reflect.Type
}

// 字段tag 加上 s2s:"-" 则会被忽略
// 成员属性必须带有 json tag标签才能获取
// 例:
//
//	type structToMap struct {
//		   Data int64 `json:"data" struct2map:"-"`
//	}
//
// except: 可以赋空值的字段， 通过 NewExcept 获取(参数是定义的成员名，不是tag名，需要注意)
// 样例：result, err := struct2Map(data, NewExcept("DataString"), nil)
// 从结构体转为 map (适用于 数据库更新数据)
func struct2Map(v interface{}, except map[string]int64, maps map[string]interface{}, originValue ...Struct2MapValue) (data map[string]interface{}, err error) {
	var (
		vt  reflect.Type
		val reflect.Value
	)

	if except == nil {
		except = make(map[string]int64)
	}

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

	kd := vt.Kind()
	if kd != reflect.Struct {
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
		fieldIgnore := vt.Field(i).Tag.Get("s2s")
		if fieldIgnore == "-" {
			continue
		}

		fieldName := vt.Field(i).Tag.Get("json")
		if fieldName == "" {
			data, _ = struct2Map(val.Field(i), except, data, Struct2MapValue{
				Value: val.Field(i),
				Type:  vt.Field(i).Type,
			})
			continue
		}

		_, ok := except[vt.Field(i).Name]
		d, cycle := getStructData(val, i, ok, val.Field(i).Type().String())
		if cycle {
			data, _ = struct2Map(val.Field(i), except, data, Struct2MapValue{
				Value: val.Field(i),
				Type:  vt.Field(i).Type,
			})
		}
		if d != nil {
			data[fieldName] = d
		}
	}

	return
}

func getStructData(val reflect.Value, i int, ok bool, typeName string) (interface{}, bool) {
	switch typeName {
	case "string":
		if val.Field(i).String() != "" || ok {
			return val.Field(i).String(), false
		}
	case "*string":
		if val.Field(i).IsNil() {
			return nil, false
		}
		return val.Field(i).Elem().String(), false

	case "*int64", "*int32", "*int16", "*int8", "*int":
		if val.Field(i).IsNil() {
			return nil, false
		}
		return val.Field(i).Elem().Int(), false

	case "*uint64", "*uint32", "*uint16", "*uint8", "*uint":
		if val.Field(i).IsNil() {
			return nil, false
		}
		return val.Field(i).Elem().Uint(), false

	case "uint64", "uint32", "uint16", "uint8", "uint":
		if val.Field(i).Uint() != 0 || ok {
			return val.Field(i).Uint(), false
		}

	case "int64", "int32", "int16", "int8", "int":
		if val.Field(i).Int() != 0 || ok {
			return val.Field(i).Int(), false
		}
	case "float64", "float32":
		if val.Field(i).Float() != 0 || ok {
			return val.Field(i).Float(), false
		}
	case "*float64", "*float32":
		if val.Field(i).IsNil() {
			return nil, false
		}
		return val.Field(i).Elem().Float(), false
	case "time.Time", "*time.Time":
		return nil, false
	default:
		k := val.Field(i).Kind()
		if _, existKd := structMap[k]; existKd {
			return getStructData(val, i, ok, k.String())
		}

		return nil, true
	}

	return nil, false
}
