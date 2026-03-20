package csvutils

import (
	"fmt"
	"strconv"
	"strings"
)

// FindColumnIndex 根据提供的键（列名或数字索引）解析并返回目标列在表头中的位置（0-based 索引）。
// 该函数支持两种查找模式：
// 1. 数字模式：如果 key 是纯数字字符串（如 "1", "2"），则按位置查找。
//   - 注意：输入索引基于 1（即 "1" 代表第一列），返回索引基于 0。
//   - 会进行边界检查，确保索引在有效范围内。
//
// 2. 名称模式：如果 key 是字符串，则遍历表头进行精确匹配。
//
// 参数说明：
//   - header: CSV 文件的表头切片。
//   - key: 查找键，可以是列名（如 "user_id"）或从 1 开始的列序号字符串（如 "3"）。
//
// 返回值：
//   - int: 目标列在切片中的索引（0-based）。
//   - error: 如果索引越界或列名不存在时返回错误。
func FindColumnIndex(header []string, key string) (int, error) {
	if idx, err := strconv.Atoi(key); err == nil {
		if idx <= 0 || idx > len(header) {
			return -1, fmt.Errorf("invalid key column index: %d, should be between 1 and %d", idx, len(header))
		}
		return idx - 1, nil
	}
	for i, col := range header {
		if col == key {
			return i, nil
		}
	}
	return -1, fmt.Errorf("column name '%s' not found in header", key)
}

// GetKeyIndices 批量获取指定关键字段在表头中对应的列索引切片。
// 该函数主要用于确定需要提取或处理的列的位置。
//
// 逻辑说明：
// 1. 空值处理：如果 keyFields 为空，默认返回所有列的索引（0 到 len(header)-1）。
// 2. 构建映射：为了提高查找效率，先将表头字段映射到其索引位置（Map[字段名]索引）。
// 3. 顺序匹配：按照 keyFields 指定的顺序，依次查找对应的索引，确保返回结果的顺序与请求一致。
//
// 参数说明：
//   - header: CSV 文件的表头切片。
//   - keyFields: 需要获取索引的字段名列表。
//
// 返回值：
//   - []int: 字段对应的索引切片。
//   - error: 如果请求的字段在表头中不存在，返回错误。
func GetKeyIndices(header []string, keyFields []string) ([]int, error) {
	if len(keyFields) == 0 {
		// 全字段去重: 返回所有列索引
		indices := make([]int, len(header))
		for i := range header {
			indices[i] = i
		}
		return indices, nil
	}

	// 按指定字段名查找索引
	headerMap := make(map[string]int, len(header))
	for idx, field := range header {
		headerMap[strings.TrimSpace(field)] = idx
	}

	indices := make([]int, 0, len(keyFields))
	for _, field := range keyFields {
		field = strings.TrimSpace(field)
		idx, ok := headerMap[field]
		if !ok {
			return nil, fmt.Errorf("title [%s] does not in header: %v", field, headerMap)
		}
		indices = append(indices, idx)
	}
	return indices, nil
}
