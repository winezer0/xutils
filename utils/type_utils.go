package utils

import "strconv"

// TypeConverter 将字符串转换为目标类型 T 。
type TypeConverter[T any] func(string) (T, error)

var (
	// StringConverter 直接返回字符串
	StringConverter = func(s string) (string, error) { return s, nil }
	// IntConverter 转换为 int
	IntConverter = func(s string) (int, error) { return strconv.Atoi(s) }
	// BoolConverter 转换为 bool
	BoolConverter = func(s string) (bool, error) { return strconv.ParseBool(s) }
	// Float64Converter 转换为 float64
	Float64Converter = func(s string) (float64, error) { return strconv.ParseFloat(s, 64) }
)

// ConvertStrings 转换字符串切片为目标类型切片，使用提供的转换函数
func ConvertStrings[T any](strList []string, converter TypeConverter[T]) ([]T, error) {
	result := make([]T, 0, len(strList))
	for _, s := range strList {
		val, err := converter(s)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}
