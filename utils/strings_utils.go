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
