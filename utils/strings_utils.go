package utils

import "strings"

// CleanBlanks 去除空行
func CleanBlanks(lines []string) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// ToLowerKeys 将排除关键字列表全部转为小写
func ToLowerKeys(keys []string) []string {
	// 显式处理空列表，避免不必要的切片创建（虽然make空切片性能影响极小，但更直观）
	if len(keys) == 0 {
		return []string{}
	}

	lowerKeys := make([]string, len(keys))
	for i, key := range keys {
		lowerKeys[i] = strings.ToLower(key)
	}
	return lowerKeys
}

// ToUpperKeys 将排除关键字列表全部转为大写
func ToUpperKeys(keys []string) []string {
	// 显式处理空列表，避免不必要的切片创建（虽然make空切片性能影响极小，但更直观）
	if len(keys) == 0 {
		return []string{}
	}

	upperKeys := make([]string, len(keys))
	for i, key := range keys {
		upperKeys[i] = strings.ToUpper(key)
	}
	return upperKeys
}

// StringsAffix 为字符串切片 ss 中的每个元素批量添加前缀和后缀。
//
// 参数:
//   - ss:       原始字符串切片。若为 nil，返回 nil；若为空切片，返回空切片。
//   - prefix:   要添加的前缀字符串。如果为空，则不添加前缀。
//   - suffix:   要添加的后缀字符串。如果为空，则不添加后缀。
//   - check:    是否检查元素是否已包含前后缀。
//     若为 true，则仅对尚未包含对应前缀/后缀的元素进行追加，避免重复；
//     若为 false，则无条件追加。
//
// 返回:
//   - 处理后的新字符串切片，原切片不被修改。
func StringsAffix(ss []string, prefix, suffix string, check bool) []string {
	if len(ss) == 0 {
		return []string{}
	}

	if prefix == "" && suffix == "" {
		result := make([]string, len(ss))
		copy(result, ss)
		return result
	}

	result := make([]string, 0, len(ss))
	for _, s := range ss {
		// 初始化最终字符串为原始值
		final := s
		if !check {
			final = prefix + final + suffix
		} else {
			// 处理前缀
			if prefix != "" && !strings.HasPrefix(s, prefix) {
				final = prefix + final
			}
			// 处理后缀
			if suffix != "" && !strings.HasSuffix(final, suffix) {
				final = final + suffix
			}
		}

		result = append(result, final)
	}

	return result
}
