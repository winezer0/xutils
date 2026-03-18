package csvutils

import (
	"fmt"
	"strconv"
	"strings"
)

// FindColumnIndex 根据列名或索引解析目标列位置
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

// GetKeyIndices 获取关键字段对应的列索引（基于表头）
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
