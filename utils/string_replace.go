package utils

import "strings"

// StringsReplacer 移除或替换字符串中的特定字符。
// 该函数接受一个替换规则映射，允许用户自定义需要替换的字符及其目标字符。
// 同时，它会自动处理连续的空格，确保结果中不包含多余的空格。
//
// 参数说明：
//   - s: 原始字符串。
//   - rules: 替换规则 map，key 为旧字符（如 "\n"），value 为新字符（如 " "）。
//
// 返回值：
//   - string: 处理后的字符串。
func StringsReplacer(s string, rules map[string]string) string {
	if s == "" {
		return s
	}
	// 如果没有提供规则，直接返回
	if len(rules) == 0 {
		return s
	}
	// 将 map 转换为 strings.Replacer 所需的参数格式
	// strings.Replacer 接受 ...string，即 old1, new1, old2, new2...
	replArgs := make([]string, 0, len(rules)*2)
	for oldS, newS := range rules {
		replArgs = append(replArgs, oldS, newS)
	}
	// 创建 Replacer
	// 注意：由于 rules 是动态的，这里每次调用都会创建一个新的 Replacer。
	// 对于高频调用场景，如果 rules 是固定的，建议在外部创建好 Replacer 传入。
	replacer := strings.NewReplacer(replArgs...)
	// 执行替换
	return replacer.Replace(s)
}

// RemoveNewlines 移除字符串中的换行符等空白字符
func RemoveNewlines(s string) string {
	rules := map[string]string{
		"\r": " ",
		"\n": " ",
		"\t": " ",
	}

	// 执行替换
	replaced := StringsReplacer(s, rules)
	return replaced
}
