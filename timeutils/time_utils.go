package timeutils

import (
	"time"
)

// 常量定义：Go 的布局模板 (Layout)
const (
	LayoutDate     = "2006-01-02"          // 2026-03-12
	LayoutDateTime = "2006-01-02 15:04:05" // 2026-03-12 13:03:00
	LayoutFileSafe = "20060102_150405"     // 20260312_130300 (适合文件名)
	LayoutCN       = "2006年01月02日 15时04分"
	LayoutISO8601  = "2006-01-02T15:04:05Z07:00"
)

// GetCurrentTimeStr 获取当前时间
func GetCurrentTimeStr(format string) string {
	return time.Now().Format(format)
}

// GetCurrentTimeStrInLocation 获取指定时区的当前时间
func GetCurrentTimeStrInLocation(locationName, format string) (string, error) {
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		return "", err
	}

	// 将当前时间转换到指定时区
	nowInLoc := time.Now().In(loc)
	return nowInLoc.Format(format), nil
}

// GetNow 获取当前时间 (默认本地时区)
func GetNow() time.Time {
	return time.Now()
}

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

// ParseDate 解析日期字符串 (YYYY-MM-DD)
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(LayoutDate, dateStr)
}

// ParseDateTime 解析日期时间字符串 (YYYY-MM-DD HH:mm:ss)
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return time.Parse(LayoutDateTime, dateTimeStr)
}

// GetStartOfDay 获取当天的开始时间 (00:00:00)
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay 获取当天的结束时间 (23:59:59.999999999)
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// AddDays 增加天数 (n 可以是负数表示减少)
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddMonths 增加月份
func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// IsToday 判断是否是今天
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// GetDateRange 获取最近 N 天的起始和结束时间 (包含今天)
// 例如 n=7, 返回 7天前0点 到 今天24点
func GetDateRange(days int) (start, end time.Time) {
	now := time.Now()
	end = GetEndOfDay(now)
	start = GetStartOfDay(AddDays(now, -days+1))
	return start, end
}
