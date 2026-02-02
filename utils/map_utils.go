package utils

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
