package utils

import (
	"fmt"
	"sort"
	"strings"
)

// JoinMap 将 name -> version 的映射格式化为 "name@version" 的字符串列表，并用逗号连接
func JoinMap(items map[string]string) string {
	if len(items) == 0 {
		return ""
	}

	var parts []string
	for name, version := range items {
		if version == "" {
			parts = append(parts, name) // 如果版本为空，只保留名称（或可选地跳过）
			// 或者：parts = append(parts, name+"@unknown") // 根据需求决定
		} else {
			parts = append(parts, name+"@"+version)
		}
	}

	// 可选：排序以保证输出稳定（例如用于日志或测试）
	sort.Strings(parts)
	return strings.Join(parts, ", ")
}

// JoinMaps 从 map 中提取切片/数组/基础类型字段，转换为换行拼接的字符串
// 核心用途：兼容AI返回的数组/字符串/数字等格式，统一转为标准字符串
// 入参：
//
//	m - 源map
//	keys - 字段名列表（支持多别名，按顺序匹配）
//
// 返回：拼接后的字符串，无有效值时返回空字符串
func JoinMaps(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			var parts []string
			switch tv := v.(type) {
			case []string:
				parts = tv // 直接复用字符串数组，无需遍历
			case []interface{}:
				// 遍历任意类型数组，转为字符串后拼接
				for _, it := range tv {
					parts = append(parts, fmt.Sprintf("%v", it))
				}
			default:
				// 非切片类型（字符串/数字/布尔等），直接转为字符串
				return fmt.Sprintf("%v", tv)
			}
			// 用换行符拼接数组，性能优于+=
			return strings.Join(parts, "\n")
		}
	}
	return ""
}

// MapToStrings 将 map[string]string 转为 "k1=v1;k2=v2;..." 字符串
// - 若 keys 为空，则使用 map 中所有 key 并按字典序排序
// - 若 keys 非空，则按 keys 顺序，仅包含存在于 map 中的 key
func MapToStrings(m map[string]string, keys []string) string {
	if m == nil {
		m = make(map[string]string)
	}

	var pairs []string

	if len(keys) == 0 {
		// 提取所有 key 并排序
		allKeys := make([]string, 0, len(m))
		for k := range m {
			allKeys = append(allKeys, k)
		}
		sort.Strings(allKeys)
		for _, k := range allKeys {
			pairs = append(pairs, k+"="+m[k])
		}
	} else {
		// 按 keys 顺序，只取存在的 key
		for _, k := range keys {
			if v, exists := m[k]; exists {
				pairs = append(pairs, k+"="+v)
			}
			// 不存在则跳过（不加入）
		}
	}

	return strings.Join(pairs, ";")
}
