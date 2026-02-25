package cmdutils

import "strings"

// SplitComma 按逗号分割字符串，去除每项前后的空白字符，并过滤空字符串。
// 如果输入不含逗号，它将返回一个包含修剪后原字符串的单元素切片（除非原字符串修剪后为空）。
func SplitComma(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

// FormatCmdsComma 格式化命令行参数列表。
// 它会将列表中每个字符串按逗号分割，并合并成一个扁平化的字符串切片。
// 同时自动去除各项周围的空白字符并忽略空项。
func FormatCmdsComma(raw []string) []string {
	var formatted []string
	for _, item := range raw {
		// 直接对每一项使用 SplitComma
		// 无论 item 是否包含逗号，SplitComma 都能正确处理：
		// 1. "a" -> ["a"]
		// 2. "a,b" -> ["a", "b"]
		// 3. "  " -> [] (自动忽略)
		parts := SplitComma(item)
		formatted = append(formatted, parts...)
	}
	return formatted
}
