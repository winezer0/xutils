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

// MergeMapsWithOverride 合并多个map，后传入的map会覆盖先传入的map的同名key
// 入参：多个map[string]string（支持nil，会自动跳过）
// 返回：合并后的新map（所有输入map的key都会被包含，同名key取最后一个出现的值）
func MergeMapsWithOverride(maps ...map[string]string) map[string]string {
	// 初始化结果map
	result := make(map[string]string)

	// 遍历所有传入的map
	for _, m := range maps {
		// 跳过nil的map，避免panic
		if m == nil {
			continue
		}
		// 遍历当前map的键值对：直接覆盖（无论是否已存在）
		for k, v := range m {
			result[k] = v // 核心逻辑：存在则覆盖，不存在则新增
		}
	}

	return result
}

// MergeIntfaceMapsNoOverride 合并多个map[string]interface{}，不覆盖已有key（保留第一个出现的值）
// 入参：多个map[string]interface{}（支持nil，会自动跳过）
// 返回：合并后的新map，同名key保留第一个出现的值
func MergeIntfaceMapsNoOverride(maps ...map[string]interface{}) map[string]interface{} {
	// 初始化结果map
	result := make(map[string]interface{})

	// 遍历所有传入的map
	for _, m := range maps {
		// 跳过nil的map，避免遍历nil map触发panic
		if m == nil {
			continue
		}
		// 遍历当前map的键值对，仅当key不存在时才设置
		for k, v := range m {
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	}

	return result
}

// MergeIntfaceMapsWithOverride 合并多个map[string]interface{}，后传入的map会覆盖先传入的map的同名key
// 入参：多个map[string]interface{}（支持nil，会自动跳过）
// 返回：合并后的新map，同名key取最后一个出现的值
func MergeIntfaceMapsWithOverride(maps ...map[string]interface{}) map[string]interface{} {
	// 初始化结果map
	result := make(map[string]interface{})

	// 遍历所有传入的map
	for _, m := range maps {
		// 跳过nil的map，避免遍历nil map触发panic
		if m == nil {
			continue
		}
		// 遍历当前map的键值对，直接覆盖（无论是否已存在）
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}
