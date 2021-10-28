/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/04/2021
 * @Desc: desc
 */

package libs

import (
	"errors"
	"fmt"
	"gin_template/app/enum"
	"github.com/jinzhu/gorm"
	"reflect"
)

func GetPageAndSize(page, pageSize int) (int, int) {
	if page == 0 || pageSize == 0 {
		page = 1
		pageSize = 20
	}

	return page, pageSize
}

func GetOffset(page, pageSize int) int {
	offset := (page - 1) * pageSize
	if offset <= 0 {
		offset = 0
	}

	return offset
}

func getOrderByStr(orderField, defaultOrder string, orderType uint8) (orderBy string) {
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

type separateHandle func(*gorm.DB, interface{}) *gorm.DB

// auto struct gorm.DB
// 1 -> and
// 2 -> or
// 3 -> separate handle 单独处理
// 4 -> like 操作
// separateHandle 单独处理 的字符串数组
func queryAutoWhere(db *gorm.DB, search interface{}, fieldIsHandle, fieldHandle string, separate ...map[string]separateHandle) (*gorm.DB, error) {
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
		handleType := vt.Field(i).Tag.Get(fieldIsHandle)
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
			db = db.Where(vt.Field(i).Tag.Get(fieldHandle), data)
		case enum.AutoOr:
			db = db.Or(vt.Field(i).Tag.Get(fieldHandle), data)
		case enum.AutoCustomHandle:
			// 根据条件 判断大于和小于
			if len(separate) > 0 {
				db = separate[0][vt.Field(i).Tag.Get("json")](db, data)
			}
		case enum.AutoLike:
			db = db.Where(vt.Field(i).Tag.Get(fieldHandle), "%"+fmt.Sprintf("%v", data)+"%")
		}

	}

	return db, err
}
