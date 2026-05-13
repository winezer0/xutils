package utils

import (
	"fmt"
	"strconv"
	"time"
)

// TruncateString 定义截断函数，保留前n个字符，超出部分省略
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// TruncateBytes 对字节数组内容进行截断，保留前 maxLen 个字符，超出部分显示 ...
// 不依赖 TruncateString，独立实现
func TruncateBytes(b []byte, maxLen int) string {
	// 长度合法直接返回字符串
	if len(b) <= maxLen {
		return string(b)
	}
	// 超出长度：截取前 maxLen-3 个字节 + 省略号
	return string(b[:maxLen-3]) + "..."
}

// ToStr 将任意类型 v 转换为字符串。
// - 基本类型（int, float, bool 等）→ 语义化字符串（如 123 → "123"）
// - string / []byte → 直接返回（不加引号）
// - time.Time → RFC3339 格式
// - 其他类型（struct, slice, map 等）→ 使用 fmt.Sprint（不带格式）
func ToStr(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case int:
		return strconv.Itoa(x)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case float32:
		return strconv.FormatFloat(float64(x), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)
	case bool:
		return strconv.FormatBool(x)
	case time.Time:
		if x.IsZero() {
			return ""
		}
		return x.Format(time.RFC3339)
	case nil:
		return ""
	default:
		// 包括 struct, slice, map, pointer, channel 等
		return fmt.Sprint(x)
	}
}
