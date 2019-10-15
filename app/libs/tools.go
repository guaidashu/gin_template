/**
  create by yy on 2019-08-23
*/

package libs

import "time"

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
