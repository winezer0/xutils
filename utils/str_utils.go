package utils

import (
	"strings"
)

// TruncateString 截断字符串
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
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

// StringHasKey 检查字符串是否包含指定关键字，支持&&和||逻辑
func StringHasKey(content string, key string) bool {
	// 当内容或关键字为空时返回false
	if content == "" || key == "" {
		return false
	}

	contentLower := strings.ToLower(content)
	keyStrLower := strings.ToLower(key)

	if strings.Contains(keyStrLower, "||") {
		// 只要有一个关键词存在于内容中，就返回true
		keys := strings.Split(keyStrLower, "||")
		for _, key := range keys {
			key = strings.TrimSpace(key)
			if key != "" && strings.Contains(contentLower, key) {
				return true
			}
		}
		return false
	} else if strings.Contains(keyStrLower, "&&") {
		// 所有关键词都必须存在于内容中，才返回true
		keys := strings.Split(keyStrLower, "&&")
		for _, key := range keys {
			key = strings.TrimSpace(key)
			if key != "" && !strings.Contains(contentLower, key) {
				return false
			}
		}
		return true
	} else {
		// 普通包含判断
		return strings.Contains(contentLower, keyStrLower)
	}
}

// StringHasOneKey 检查字符串 content 是否包含任意一个 keys 中的关键字
func StringHasOneKey(str string, keys []string, ignoreCase bool) bool {
	if ignoreCase {
		str = strings.ToLower(str)
		keys = ToLowerKeys(keys)
	}
	for _, key := range keys {
		if strings.Contains(str, key) {
			return true
		}
	}
	return false
}

// StringsHasOneKey 检查多个字符串中是否存在任意一个匹配任意一个关键字
func StringsHasOneKey(strList []string, keys []string, ignoreCase bool) bool {
	for _, str := range strList {
		if StringHasOneKey(str, keys, ignoreCase) {
			return true
		}
	}
	return false
}

// StringInStrings 检查字符串是否在切片中
func StringInStrings(str string, strList []string, ignoreCase bool) bool {
	if ignoreCase {
		str = strings.ToLower(str)
		strList = ToLowerKeys(strList)
	}

	for _, v := range strList {
		if str == v {
			return true
		}
	}
	return false
}

// StringsInStrings 判断 filter 中的字符串是否存在于 allowed 中，返回存在和不存在的列表
func StringsInStrings(stringsA, stringsB []string, ignoreCase bool) (existList, notExistList []string) {
	allowedSet := make(map[string]bool)

	if ignoreCase {
		stringsA = ToLowerKeys(stringsA)
		stringsB = ToLowerKeys(stringsB)
	}

	// 构建 stringsB 的 map 集合，用于快速查找
	for _, item := range stringsB {
		allowedSet[item] = true
	}

	// 遍历 stringsA 列表，判断是否存在
	for _, item := range stringsA {
		if allowedSet[item] {
			existList = append(existList, item)
		} else {
			notExistList = append(notExistList, item)
		}
	}

	return existList, notExistList
}
