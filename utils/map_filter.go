package utils

// FilterMapByKeys 过滤掉键存在于 filteredKeys 中的条目，支持任意 map 类型
func FilterMapByKeys[K comparable, V any](sourceMap map[K]V, filteredKeys []K) map[K]V {
	if len(filteredKeys) == 0 {
		return sourceMap
	}

	// 快速排除集合
	excludeSet := make(map[K]struct{}, len(filteredKeys))
	for _, key := range filteredKeys {
		excludeSet[key] = struct{}{}
	}

	// 构建结果
	out := make(map[K]V, len(sourceMap))
	for k, v := range sourceMap {
		// 只保留不在排除列表中的键
		if _, excluded := excludeSet[k]; !excluded {
			out[k] = v
		}
	}

	return out
}
