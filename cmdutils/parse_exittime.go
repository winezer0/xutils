package cmdutils

import (
	"strings"
	"time"
)

// ParseExitTime 解析退出时间或时长(中文注释)
func ParseExitTime(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	// 绝对时间，支持YYYY-MM-DD:HH:MM:SS或YYYYMMDD:HH:MM:SS
	layouts := []string{
		"2006-01-02:15:04:05",
		"20060102:15:04:05",
		"2006-01-02 15:04:05",
	}
	for _, l := range layouts {
		if t, err := time.ParseInLocation(l, s, time.Local); err == nil {
			return t, true
		}
	}
	// 相对时长，支持组合(如10h或20m或30s)
	if dur, err := time.ParseDuration(s); err == nil {
		return time.Now().Add(dur), true
	}
	return time.Time{}, false
}
