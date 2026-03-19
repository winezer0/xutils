package csvutils

import (
	"fmt"
	"strconv"
	"strings"
)

// ------------------- 辅助函数（完整实现，确保无依赖） -------------------
func fixedHeaders(headers []string) []string {
	cleaned := make([]string, len(headers))
	seen := make(map[string]int)
	for i, header := range headers {
		// 安全移除 UTF-8 BOM
		if strings.HasPrefix(header, "\ufeff") {
			header = header[3:] // "\ufeff" 在 UTF-8 中占 3 字节
		}

		// 移除空白
		header = strings.TrimSpace(header)
		if header == "" {
			header = fmt.Sprintf("col%d", i)
		}

		// 处理重复头
		if cnt, ok := seen[header]; ok {
			seen[header]++
			header = fmt.Sprintf("%s_%d", header, cnt+1)
		} else {
			seen[header] = 1
		}
		cleaned[i] = header
	}
	return cleaned
}

// genDefaultHeader 生成默认头
func genDefaultHeader(colCount int) []string {
	header := make([]string, colCount)
	for i := 0; i < colCount; i++ {
		header[i] = fmt.Sprintf("col%d", i)
	}
	return header
}

// rowToDict []string 转 dict 格式
func rowToDict(row []string, header []string, convertType bool) map[string]interface{} {
	dict := make(map[string]interface{}, len(header))
	for colIdx, key := range header {
		if colIdx >= len(row) {
			dict[key] = nil
			continue
		}
		val := strings.TrimSpace(row[colIdx])
		if convertType {
			dict[key] = autoConvertType(val)
		} else {
			dict[key] = val
		}
	}
	return dict
}

// autoConvertType 自动转换csv 的值的类型
func autoConvertType(val string) interface{} {
	if val == "" {
		return nil
	}
	if val == "true" || val == "false" {
		b, _ := strconv.ParseBool(val)
		return b
	}
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	return val
}
