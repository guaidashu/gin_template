package datetime

import (
	"time"

	"github.com/jinzhu/now"
)

const (
	LayoutForDatetime = "2006-01-02 15:04:05"
	LayoutForDate     = "2006-01-02"
	LayoutForTime     = "15:04:05"
)

// TodayFirstSecond 获取今日第一秒
func TodayFirstSecond() time.Time {
	n := time.Now()

	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
}

// GetExpiredByDay 计算今日还剩多少秒
func GetExpiredByDay() (expire time.Duration) {
	t := time.Now()
	zeroTime := time.Date(t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0, t.Location())
	sec := zeroTime.AddDate(0, 0, 1)
	expire = sec.Sub(t)
	return
}

// TodayLastSecond 获取今日最后一秒
func TodayLastSecond() time.Time {
	n := time.Now()

	return time.Date(n.Year(), n.Month(), n.Day(), 23, 59, 59, 0, n.Location())
}

// TodayDate 获取今日日期
func TodayDate() string {
	return time.Now().Format(LayoutForDate)
}

// NowDatetime 获取当前的日期时间
func NowDatetime() string {
	return time.Now().Format(LayoutForDatetime)
}

// Format: 2020060612
func GetHourDigit(t time.Time) int {
	var year, month, day = t.Date()
	var hour = t.Hour()
	return year*1000000 + int(month)*10000 + day*100 + hour
}

// Return: begin / end timestamp
func GetHourBeginAndEndTimestamp(t time.Time) (int64, int64) {
	var n = now.With(t)
	var begin = n.BeginningOfHour()
	var end = n.EndOfHour()
	return begin.Unix(), end.Unix()
}

// Format: 20200606
func GetDayDigit(t time.Time) int {
	var year, month, day = t.Date()
	return year*10000 + int(month)*100 + day
}

// Return: begin / end timestamp
func GetDayBeginAndEndTimestamp(t time.Time) (int64, int64) {
	var n = now.With(t)
	var begin = n.BeginningOfDay()
	var end = n.EndOfDay()
	return begin.Unix(), end.Unix()
}

func GetZeroStampTime() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t.Unix() - 8*60*60
}

// Format: 20201207(first day of week)
// Start: Monday --> AddDate(0, 0, 1)
func GetWeekDigit(t time.Time) int {
	var year, month, day = now.With(t).BeginningOfWeek().AddDate(0, 0, 1).Date()
	return year*10000 + int(month)*100 + day
}

// Return: begin / end timestamp
func GetWeekBeginAndEndTimestamp(t time.Time) (int64, int64) {
	myConfig := &now.Config{
		WeekStartDay: time.Monday,
		TimeFormats:  []string{"2006-01-02 15:04:05"},
	}
	n := myConfig.With(t)
	end := n.EndOfWeek()
	begin := n.BeginningOfWeek()
	return begin.Unix(), end.Unix()
}

// Format: 202012
func GetMonthDigit(t time.Time) int {
	var year, month, _ = t.Date()
	return year*100 + int(month)
}

// Return: begin / end timestamp
func GetMonthBeginAndEndTimestamp(t time.Time) (int64, int64) {
	var n = now.With(t)
	var begin = n.BeginningOfMonth()
	var end = n.EndOfMonth()
	return begin.Unix(), end.Unix()
}

// a b 时间是否为同一天
func IsSameDay(a, b time.Time) bool {
	var aYear, aMonth, aDay = a.Date()
	var bYear, bMonth, bDay = b.Date()
	return aYear == bYear && aMonth == bMonth && aDay == bDay
}

func GetTodayZeroStamp() int64 {
	timeStr := time.Now().Format("2006-01-02")
	ti, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return ti.Unix()
}

// 获取当天 23:59:59 的时间戳
func GetToday235959TimeStamp() (timeStamp int64) {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	timeStamp = t.Unix() + 86399
	return
}

func GetTodayFull24TimeStamp() (timeStamp int64) {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	timeStamp = t.Unix() + 86400
	return
}

// 获取 当前时间是第几周
func WeekByDate(t time.Time) int {
	yearDay := t.YearDay()
	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
	firstDayInWeek := int(yearFirstDay.Weekday())

	// 今年第一周有几天
	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}

	var week int
	if yearDay <= firstWeekDays {
		week = 1
	} else {
		weekOverDay := (yearDay - firstWeekDays) % 7
		if weekOverDay == 0 {
			week = (yearDay-firstWeekDays)/7 + 1
		} else {
			week = (yearDay-firstWeekDays)/7 + 2
		}
	}

	return week
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 获取现在到今天24时剩下的时间
func GetNow2Today24TimeStamp() (timestamp int64) {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	timestamp = t.Unix() + 86399 - time.Now().Unix()
	return
}
