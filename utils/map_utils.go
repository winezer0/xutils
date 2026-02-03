package utils

import (
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
