package utils

import (
	"encoding/json"
	"fmt"
	"sort"
)

// DetectUnknownFields 检测JSON数据中是否包含结构体不支持的字段，支持递归检测嵌套字段
func DetectUnknownFields(data []byte, target interface{}) []string {
	// 先解析到map
	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return nil
	}

	// 解析到目标结构体
	var targetStruct map[string]interface{}
	targetData, err := json.Marshal(target)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(targetData, &targetStruct); err != nil {
		return nil
	}

	// 递归比较字段
	unknownFields := make(map[string]struct{})
	compareMaps("", rawMap, targetStruct, unknownFields)

	// 转换为切片并排序
	result := make([]string, 0, len(unknownFields))
	for k := range unknownFields {
		result = append(result, k)
	}
	sort.Strings(result)

	return result
}

// compareValues 递归比较两个值，找出在rawValue中存在但在targetValue中不存在的所有字段路径
func compareValues(prefix string, rawValue, targetValue interface{}, unknownFields map[string]struct{}) {
	// 处理map类型
	rawMap, isRawMap := rawValue.(map[string]interface{})
	targetMap, isTargetMap := targetValue.(map[string]interface{})
	if isRawMap && isTargetMap {
		compareMaps(prefix, rawMap, targetMap, unknownFields)
		return
	}

	// 处理slice类型
	rawSlice, isRawSlice := rawValue.([]interface{})
	targetSlice, isTargetSlice := targetValue.([]interface{})
	if isRawSlice && isTargetSlice {
		// 只比较与targetSlice长度相同的部分
		minLen := len(rawSlice)
		if len(targetSlice) < minLen {
			minLen = len(targetSlice)
		}

		// 限制检测的数组元素数量，避免大数组导致性能问题和输出冗余
		// 既然已经使用了[*]进行去重，这里主要为了性能优化
		const maxSliceInspect = 10
		if minLen > maxSliceInspect {
			minLen = maxSliceInspect
		}

		for i := 0; i < minLen; i++ {
			// 使用 [*] 替换具体的索引，避免输出大量冗余的数组元素字段
			fullKey := fmt.Sprintf("%s[*]", prefix)
			compareValues(fullKey, rawSlice[i], targetSlice[i], unknownFields)
		}
		return
	}
}

// compareMaps 递归比较两个map，找出在rawMap中存在但在targetMap中不存在的所有字段路径
func compareMaps(prefix string, rawMap, targetMap map[string]interface{}, unknownFields map[string]struct{}) {
	for key, rawValue := range rawMap {
		// 构建完整字段路径
		var fullKey string
		if prefix == "" {
			fullKey = key
		} else {
			fullKey = prefix + "." + key
		}

		// 检查字段是否存在于targetMap中
		targetValue, exists := targetMap[key]
		if !exists {
			unknownFields[fullKey] = struct{}{}
			continue
		}

		// 递归比较值
		compareValues(fullKey, rawValue, targetValue, unknownFields)
	}
}
