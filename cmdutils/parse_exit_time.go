package cmdutils

import (
	"fmt"
	"strings"
	"time"
)

// ParseExitTime parses a string into a time.Time.
// Supported formats:
//   - "YYYYMMDD:HH:MM:SS"      e.g. "20251231:23:59:59"
//   - "YYYY-MM-DD:HH:MM:SS"    e.g. "2025-12-31:23:59:59"
//   - "YYYY/MM/DD:HH:MM:SS"    e.g. "2025/12/31:23:59:59"
//   - Duration strings         e.g. "10h", "30m", "5s"
//
// Returns:
//   - (zero time, nil) if input is empty (treated as "no exit time")
//   - (parsed time, nil) on success
//   - (zero time, error) if input is non-empty but invalid
func ParseExitTime(exitTimeString string) (*time.Time, error) {
	s := strings.TrimSpace(exitTimeString)
	if s == "" {
		return nil, nil // 表示没有指定退出时间
	}

	layouts := []string{
		"2006-01-02:15:04:05",
		"2006/01/02:15:04:05",
		"20060102:15:04:05",
		"2006-01-02-15-04-05",
		"2006/01/02/15/04/05",
		"20060102150405",
	}

	// 循环解析每个格式 并记录最后的错误信息
	var lastError error
	for _, layout := range layouts {
		if parseTime, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return &parseTime, nil
		} else {
			lastError = err
		}
	}

	// 尝试解析为持续时间
	if dur, err := time.ParseDuration(s); err == nil {
		parseTime := time.Now().Add(dur)
		return &parseTime, nil
	} else {
		lastError = err
	}

	// 如果所有尝试都失败了，则返回最后一个错误
	return nil, fmt.Errorf("parse time %q last error: %v", s, lastError)
}
