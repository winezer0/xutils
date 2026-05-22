package utils

import "strings"

// StringHasKeyPerfect 终极完美版
// 支持 && || 混合表达式，&& 优先级高于 ||
// ignoreCase: 是否忽略大小写
// trimSpace: 是否对内容、关键字自动去空格
// noneRet: 当 content 或 keyRule 为空时返回的值
func StringHasKeyPerfect(content string, keyRule string, ignoreCase bool, trimSpace bool, noneRet bool) bool {
	// 1. 可选：去空格
	if trimSpace {
		content = strings.TrimSpace(content)
		keyRule = strings.TrimSpace(keyRule)
	}

	// 2. 空值判断 + 自定义返回
	if content == "" || keyRule == "" {
		return noneRet
	}

	// 3. 可选：忽略大小写
	checkContent := content
	checkKey := keyRule
	if ignoreCase {
		checkContent = strings.ToLower(checkContent)
		checkKey = strings.ToLower(checkKey)
	}

	// 4. 先按 || 分割（优先级最低）
	orList := strings.Split(checkKey, "||")

	// 5. 遍历每个 OR 分组，内部必须满足全部 AND
	for _, orItem := range orList {
		// 可选：每个关键词也去空格
		if trimSpace {
			orItem = strings.TrimSpace(orItem)
		}
		if orItem == "" {
			continue
		}

		// 按 && 分割
		andList := strings.Split(orItem, "&&")
		match := true

		for _, kw := range andList {
			if trimSpace {
				kw = strings.TrimSpace(kw)
			}
			if kw == "" {
				continue
			}
			if !strings.Contains(checkContent, kw) {
				match = false
				break
			}
		}

		if match {
			return true
		}
	}

	return false
}

// StringHasKey 检查字符串是否包含指定关键字，支持&&和||逻辑
// 【默认行为】：不区分大小写 + 自动 trim 空格 + 空内容/空关键字返回 false
func StringHasKey(content string, keyStr string) bool {
	// 直接调用完美版，传入和原来逻辑完全一致的默认参数
	return StringHasKeyPerfect(content, keyStr, true, true, false)
}

// StringHasAllKeys 检查字符串是否包含所有关键字
func StringHasAllKeys(str string, keys []string, ignoreCase bool) bool {
	if ignoreCase {
		str = strings.ToLower(str)
		keys = ToLowerKeys(keys)
	}

	for _, key := range keys {
		if !strings.Contains(str, key) {
			return false
		}
	}
	return true
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

// StringMatchKeysExpr 关键词表达式匹配(支持简单 || 或 &&)(中文注释)
func StringMatchKeysExpr(content string, expr string, ignoreCase bool) bool {
	if expr == "" || content == "" {
		return false
	}

	if ignoreCase {
		content = strings.ToLower(content)
		expr = strings.ToLower(expr)
	}

	// 只要有一个关键词存在于内容中，就返回true
	if strings.Contains(expr, "||") {
		keys := strings.Split(expr, "||")
		for _, key := range keys {
			key = strings.TrimSpace(key)
			if key != "" && strings.Contains(content, key) {
				return true
			}
		}
		return false
	}

	// 所有关键词都必须存在于内容中，才返回true
	if strings.Contains(expr, "&&") {
		keys := strings.Split(expr, "&&")
		for _, key := range keys {
			key = strings.TrimSpace(key)
			if key != "" && !strings.Contains(content, key) {
				return false
			}
		}
		return true
	}

	return strings.Contains(content, strings.TrimSpace(expr))
}
