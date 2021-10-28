/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 2021/5/11 18:02
 * @Desc: 帮助文件
 */

package helper

import (
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"strconv"

	"github.com/shopspring/decimal"
)

// 打乱切片
func RandomSlice(old []int64) []int64 {
	for i := 0; i < len(old); i++ {
		num := rand.Intn(i + 1)
		old[i], old[num] = old[num], old[i]
	}
	return old
}

// map转换
func MapInterface2string(m map[string]interface{}) map[string]string {
	ret := make(map[string]string, len(m))
	for k, v := range m {
		ret[k] = fmt.Sprint(v)
	}
	return ret
}

// InetNtoA IP长整形转字符串
func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// InetAtoN IP字符串转长整形
func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

// Float32ToString 浮点数转字符串
func Float32ToString(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

// Float64ToString 浮点数转字符串
func Float64ToString(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// CompareFloat 比较浮点数大小
func CompareFloat(f1, f2 interface{}) (int, error) {
	d1, err := toDecimal(f1)
	if err != nil {
		return 0, err
	}

	d2, err := toDecimal(f2)
	if err != nil {
		return 0, err
	}

	return d1.Cmp(d2), nil
}

func toDecimal(f interface{}) (decimal.Decimal, error) {
	var (
		err error
		d   decimal.Decimal
	)

	switch v := f.(type) {
	case float32:
		d = decimal.NewFromFloat32(v)
	case float64:
		d = decimal.NewFromFloat(v)
	case int8:
		d = decimal.NewFromInt(int64(v))
	case uint8:
		d = decimal.NewFromInt(int64(v))
	case int16:
		d = decimal.NewFromInt(int64(v))
	case uint16:
		d = decimal.NewFromInt(int64(v))
	case int32:
		d = decimal.NewFromInt32(v)
	case uint32:
		d = decimal.NewFromInt32(int32(v))
	case int64:
		d = decimal.NewFromInt(v)
	case uint64:
		d = decimal.NewFromInt(int64(v))
	case int:
		d = decimal.NewFromInt(int64(v))
	case uint:
		d = decimal.NewFromInt(int64(v))
	case string:
		d, err = decimal.NewFromString(v)
	case []byte:
		d, err = decimal.NewFromString(string(v))
	}

	return d, err
}
