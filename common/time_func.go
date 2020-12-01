//
//	时间相关的函数
//
package common

import (
	"fmt"
	"strings"
	"time"
)

// MonthDayNum t 所在时间的月份总天数
func MonthDayNum(t time.Time) int {
	isLeapYear := isLeap(t.Year())

	month := t.Month()
	switch month {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.February:
		if isLeapYear {
			return 29
		}

		return 28
	default:
		return 30
	}
}

func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	hours := diff.Hours()
	if hours < 1.0 {
		return fmt.Sprintf("约 %.0f 分钟前", diff.Minutes())
	}

	if hours < 24.0 {
		return fmt.Sprintf("约 %.0f 小时前", hours)
	}

	if hours < 72.0 {
		return fmt.Sprintf("约 %.0f 天前", hours/24.0)
	}

	// 同一年，不用年份
	if now.Year() == t.Year() {
		return t.Format("01-02 15:04")
	}

	return t.Format("2006-01-02")
}

// 是否闰年
func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

//日期字符串转时间戳
func Date2Unix(datestr string) int64 {
	timeLayout := "2006-01-02 15:04:05"                          //转化所需模板
	loc, _ := time.LoadLocation("Local")                         //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, datestr, loc) //使用模板在对应时区转化为time.time类型
	return theTime.Unix()
}

//反转日月年    3/4/2020 xx:xx:xx  -> 2020-4-3 xx:xx:xx
func ReverseDate(datestr string) string {
	temp1 := strings.Split(datestr, " ")
	if len(temp1) > 1 {
		temp2 := strings.Split(temp1[0], "/")
		temp3 := []string{}
		for i := len(temp2) - 1; i >= 0; i-- {
			temp3 = append(temp3, temp2[i])
		}
		temp4 := strings.Join(temp3, "-")

		return temp4 + " " + temp1[1]
	}
	return ""
}
