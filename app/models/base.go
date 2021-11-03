/**
  create by yy on 2019-08-31
*/

package models

import (
	"errors"
	"fmt"
	"gin_template/app/enum"
	"github.com/jinzhu/gorm"
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
	db := gdb.Table(tableName)
	if db != nil {
		db = db.Where("deleted = 0")
	}
	return db
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

type Struct2MapValue struct {
	Value reflect.Value
	Type  reflect.Type
}

// 字段tag 加上 s2s:"-" 则会被忽略
// 例:
// type structToMap struct {
//	   Data int64 `json:"data" struct2map:"-"`
// }
// 从结构体转为 map (适用于 数据库更新数据)
func struct2Map(v interface{}, maps map[string]interface{}, originValue ...Struct2MapValue) (data map[string]interface{}, err error) {
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
		fieldIgnore := vt.Field(i).Tag.Get("s2s")
		if fieldIgnore == "-" {
			continue
		}

		fieldName := vt.Field(i).Tag.Get("json")
		if fieldName == "" {
			data, _ = struct2Map(val.Field(i), data, Struct2MapValue{
				Value: val.Field(i),
				Type:  vt.Field(i).Type,
			})
			continue
		}

		switch val.Field(i).Type().String() {
		case "string":
			if val.Field(i).String() != "" {
				data[fieldName] = val.Field(i).String()
			}

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
			if val.Field(i).Uint() > 0 {
				data[fieldName] = val.Field(i).Uint()
			}

		case "int64", "int32", "int16", "int8", "int":
			if val.Field(i).Int() > 0 {
				data[fieldName] = val.Field(i).Int()
			}
		default:
			data, _ = struct2Map(val.Field(i), data, Struct2MapValue{
				Value: val.Field(i),
				Type:  vt.Field(i).Type,
			})
		}
	}

	return
}
