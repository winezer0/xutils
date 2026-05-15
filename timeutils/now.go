package timeutils

import "time"

// GetNow 获取当前时间 (默认本地时区)
func GetNow() time.Time {
	return time.Now()
}

// GetNowWithFormat 获取当前时间
func GetNowWithFormat(format string) string {
	return time.Now().Format(format)
}

// GetNowWithLocation 获取指定时区的当前时间
func GetNowWithLocation(locationName, format string) (string, error) {
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		return "", err
	}

	// 将当前时间转换到指定时区
	nowInLoc := time.Now().In(loc)
	return nowInLoc.Format(format), nil
}

// GetNowTimestamp 返回当前 Unix 时间戳（秒）
func GetNowTimestamp() int64 {
	return time.Now().Unix()
}

// GetNowTimestampMs 返回当前 Unix 时间戳（毫秒）
func GetNowTimestampMs() int64 {
	return time.Now().UnixNano() / 1e6
}

// GetNowTimestampNs 返回当前 Unix 时间戳（纳秒）
func GetNowTimestampNs() int64 {
	return time.Now().UnixNano()
}

// IsToday 判断是否是今天
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// GetStartOfDay 获取当天的开始时间 (00:00:00)
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay 获取当天的结束时间 (23:59:59.999999999)
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// GetRecentlyDateRange 获取最近 N 天的起始和结束时间 (包含今天)
// 例如 n=7, 返回 7天前0点 到 今天24点
func GetRecentlyDateRange(days int) (start, end time.Time) {
	now := time.Now()
	end = GetEndOfDay(now)
	start = GetStartOfDay(TimeObjAddDays(now, -days+1))
	return start, end
}
