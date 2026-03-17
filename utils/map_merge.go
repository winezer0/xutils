package utils

// MergeMapsNoOverride 合并两个map 不覆盖已有 key（保留第一个值）
func MergeMapsNoOverride(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		if m == nil {
			continue
		}
		for k, v := range m {
			if _, exists := result[k]; !exists {
				result[k] = v // 只有不存在时才设置
			}
		}
	}
	return result
}
