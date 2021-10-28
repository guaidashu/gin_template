/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 2021/6/10 15:39
 * @Desc: 抽奖算法
 */

package random

import (
	"gin_template/app/libs/helper"
	"math"
	"math/rand"
	"strings"
)

// Lottery 随机抽奖算法
// ps {"1":10.00001,"2":9.5,"3":15.3} 	奖品概率池
// num 									从中抽取奖品数量
// unique 								每个奖品是否唯一
func Lottery(ps map[string]float64, num int, unique bool) []string {
	var (
		scale float64 = 0
		sum   float64 = 0
		res           = make([]string, 0)
		pool          = make(map[string]float64)
		key   string
		value float64
	)

	if count := len(ps); count == 0 || (count < num && unique) {
		return res
	}

	for _, v := range ps {
		n := helper.Float64ToString(v)

		pos := strings.Index(n, ".")
		if pos < 0 {
			pos = len(n) - 1
		}

		scale = math.Max(scale, math.Pow10(len(n)-pos-1))
	}

	for k, v := range ps {
		pool[k] = v * scale

		sum += pool[k]
	}

	for i := 0; i < num; i++ {
		r := float64(rand.Int63n(int64(sum)))

		for k, v := range pool {
			if r < v {
				key, value = k, v
				break
			} else {
				r -= v
			}
		}

		res = append(res, key)

		if num == 1 {
			break
		}

		if unique {
			sum -= value
			delete(pool, key)
		}
	}

	return res
}
