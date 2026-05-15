package timeutils

import "time"

// DateToTimeObj 解析日期字符串 (YYYY-MM-DD)
// 作用：把 2025-12-20 这种 YYYY-MM-DD 格式的字符串，解析成 Go 语言能使用的 time.Time 时间对象。
func DateToTimeObj(dateStr string) (time.Time, error) {
	return time.Parse(LayoutDate, dateStr)
}

// DateTimeToTimeObj 解析日期时间字符串 (YYYY-MM-DD HH:mm:ss)
// 把 2025-12-20 15:30:59 这种「年月日 时分秒」字符串，解析成 Go 语言能操作的 time.Time 时间对象。
func DateTimeToTimeObj(dateTimeStr string) (time.Time, error) {
	return time.Parse(LayoutDateTime, dateTimeStr)
}
