package timeutils

import "time"

// FormatDate 格式化为日期字符串 (YYYY-MM-DD)
func FormatDate(t time.Time) string {
	return t.Format(LayoutDate)
}

// FormatDateTime 格式化为日期时间字符串 (YYYY-MM-DD HH:mm:ss)
func FormatDateTime(t time.Time) string {
	return t.Format(LayoutDateTime)
}

// FormatFileSafe 格式化为适合文件名的字符串 (YYYYMMDD_HHMMSS)
func FormatFileSafe(t time.Time) string {
	return t.Format(LayoutFileSafe)
}
