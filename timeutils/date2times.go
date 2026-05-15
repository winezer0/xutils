package timeutils

import "time"

// DateTimeToTimestamp 将指定日期字符串转为 秒级时间戳
// 日期格式示例："2025-01-01 12:00:00"
func DateTimeToTimestamp(dateStr string) (int64, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, time.Local)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// DateTimeToTimestampMs 将指定日期字符串转为 毫秒级时间戳
func DateTimeToTimestampMs(dateStr string) (int64, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, time.Local)
	if err != nil {
		return 0, err
	}
	return t.UnixNano() / 1e6, nil
}

// DateTimeToTimestampNs 将指定日期字符串转为 纳秒级时间戳
func DateTimeToTimestampNs(dateStr string) (int64, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, time.Local)
	if err != nil {
		return 0, err
	}
	return t.UnixNano(), nil
}
