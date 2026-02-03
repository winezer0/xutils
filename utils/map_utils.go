package utils

import (
	"fmt"
	"sort"
	"strings"
)

// MapKeys 辅助函数：从 map[string]bool 提取所有 key
// 参数:
// - m: 输入的映射
// 返回值:
// - []string: 映射中所有的键组成的切片
func MapKeys(m map[string]bool) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

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

// JoinMapSlice 从 map 中提取切片/数组/基础类型字段，转换为换行拼接的字符串
// 核心用途：兼容AI返回的数组/字符串/数字等格式，统一转为标准字符串
// 入参：
//
//	m - 源map
//	keys - 字段名列表（支持多别名，按顺序匹配）
//
// 返回：拼接后的字符串，无有效值时返回空字符串
func JoinMapSlice(m map[string]interface{}, keys ...string) string {
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
