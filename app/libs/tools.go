/**
  create by yy on 2019-08-23
*/

package libs

import (
	"errors"
	"log"
	"reflect"
	"strings"
	"time"
	"unicode"
)

const (
	defaultFromTag = "json"
	defaultToTag   = "json"
)

// 判断一个 元素 是否存在数组(切片中)
func InSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

// 得到当前时间戳
func GetNowTimeStamp() int64 {
	return time.Now().Unix()
}

func GetNowTime(nowTimeStamp int64) string {
	if nowTimeStamp == 0 {
		nowTimeStamp = time.Now().Unix()
	}
	return time.Unix(nowTimeStamp, 0).UTC().Format("2006-01-02 15:04:05")
}

func GetNowTimeMon(nowTimeStamp int64) string {
	if nowTimeStamp == 0 {
		nowTimeStamp = time.Now().Unix()
	}
	return time.Unix(nowTimeStamp, 0).UTC().Format("2006-01-02")
}

func RunSafe(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		fn()
	}()
}

// 结构体数据复制, 主要通过tag转换,
// tag必须一致, 返回转换了的字段(首字母大写)
// 可以指定tag, 不指定则默认为 json - json 指定的时候必须传两个

// 结构体值转换, 一定场景替代Copy函数,
// 添加 s2s:"-" 则此字段会被忽略赋值, 主要用于过滤更新操作时的 主键ID
// 试例:
// type DemoTo struct {
//     Id   int64  `json:"id"`
//     Data string `json:"data"`
// }
//
// type DemoFrom struct {
//     Id   int64  `json:"id" s2s:"-"`
//     Data string `json:"data"`
// }
// 进行赋值操作时, 只会赋值 Data成员
//
// PS: 传入参数from 和 to必须是可寻址的, 也就是必须传入指针
type s2sValue struct {
	Type       string        `json:"type"`
	OriginType reflect.Type  `json:"origin_type"`
	Value      reflect.Value `json:"value"`
}

func Struct2Struct(from interface{}, to interface{}, tags ...string) (err error) {
	var (
		fromTag, toTag string
	)

	if len(tags) > 0 {
		if len(tags) < 2 {
			return errors.New("参数错误")
		}
		fromTag = tags[0]
		toTag = tags[1]
	} else {
		fromTag = defaultFromTag
		toTag = defaultToTag
	}

	fromVt := reflect.TypeOf(from)
	fromVal := reflect.ValueOf(from)
	toVt := reflect.TypeOf(to)
	toVal := reflect.ValueOf(to)

	if kd := fromVt.Kind(); kd != reflect.Struct {
		if kd := fromVt.Elem().Kind(); kd != reflect.Struct {
			return errors.New("the param is not a struct, please make sure")
		} else {
			fromVt = fromVt.Elem()
			fromVal = fromVal.Elem()
		}
	}

	if kd := toVt.Kind(); kd != reflect.Struct {
		if kd := toVt.Elem().Kind(); kd != reflect.Struct {
			return errors.New("the param is not a struct, please make sure")
		} else {
			toVt = toVt.Elem()
			toVal = toVal.Elem()
		}
	}

	// 首先构造 tag map
	copyMap := make(map[string]s2sValue)
	toFieldNum := toVal.NumField()
	for i := 0; i < toFieldNum; i++ {
		tagName := toVt.Field(i).Tag.Get(toTag)
		if tagName == "" {
			continue
		}
		copyMap[tagName] = s2sValue{
			Type:       toVal.Field(i).Type().String(),
			OriginType: toVal.Field(i).Type(),
			Value:      toVal.Field(i),
		}
	}

	fieldNum := fromVal.NumField()
	for i := 0; i < fieldNum; i++ {
		tagName := fromVt.Field(i).Tag.Get(fromTag)
		if tagName == "" {
			continue
		}
		ignoreTag := fromVt.Field(i).Tag.Get("s2s")
		if ignoreTag == "-" {
			continue
		}

		fromType := fromVal.Field(i).Type().String()

		switch fromType {
		case "*string",
			"*int64", "*int32",
			"*int16", "*int8", "*int",
			"*uint64", "*uint32",
			"*uint16", "*uint8", "*uint",
			"*float64", "*float32":
			if fromVal.Field(i).IsNil() {
				continue
			}
		}

		if fromType == copyMap[tagName].Type {
			copyMap[tagName].Value.Set(fromVal.Field(i))
		} else if strings.Replace(fromType, "*", "", -1) == copyMap[tagName].Type {
			copyMap[tagName].Value.Set(fromVal.Field(i).Elem())
		} else if strings.Replace(copyMap[tagName].Type, "*", "", -1) == fromType {
			if copyMap[tagName].Value.IsNil() {
				copyMap[tagName].Value.Set(reflect.New(copyMap[tagName].OriginType.Elem()))
			}
			copyMap[tagName].Value.Elem().Set(fromVal.Field(i))
		}
	}

	return
}

// 下划线写法转为驼峰写法
func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// 首字母小写
func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
