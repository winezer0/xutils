package csvutils

import (
	"fmt"
	"strconv"
	"strings"
)

// RepairHeaders 用于清洗和标准化 CSV 表头，确保列名唯一且格式正确。
// 处理 BOM 头、空白字符以及重复列名的问题。
func RepairHeaders(headers []string) []string {
	cleaned := make([]string, len(headers))
	seen := make(map[string]int)
	for i, header := range headers {
		// 安全移除 UTF-8 BOM
		if strings.HasPrefix(header, "\ufeff") {
			header = header[3:] // "\ufeff" 在 UTF-8 中占 3 字节
		}

		// 移除首尾的空白字符（空格、制表符、换行符等）
		header = strings.TrimSpace(header)
		// 如果列名为空，则自动生成一个默认列名
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

// GenDefaultHeaders 根据指定的列数生成一套标准化的默认表头切片。
// 当输入数据缺失表头行，或用户选择忽略原有表头时，此函数用于构建基础的列索引（如 col0, col1...），
// 确保后续的数据处理逻辑始终有可用的字段名进行引用。
func GenDefaultHeaders(colCount int) []string {
	header := make([]string, colCount)
	for i := 0; i < colCount; i++ {
		header[i] = fmt.Sprintf("col%d", i)
	}
	return header
}

// RowDataToDict 将单行字符串切片（CSV 行数据）转换为字典格式（map[string]interface{}）。
// 该函数负责将原始的行数据与表头进行映射，并根据配置决定是否进行数据类型推断。
//
// 参数说明：
//   - row: 原始字符串切片，代表一行数据。
//   - header: 表头切片，用于作为字典的键（Key）。
//   - convertType: 是否开启自动类型转换。若为 true，会尝试将字符串解析为 int/float/bool 等类型。
//
// 返回值：
//   - map[string]interface{}: 键值对形式的数据行，键为列名，值为转换后的数据。
func RowDataToDict(row []string, header []string, convertType bool) map[string]interface{} {
	dict := make(map[string]interface{}, len(header))
	for colIdx, key := range header {
		if colIdx >= len(row) {
			dict[key] = nil
			continue
		}
		val := strings.TrimSpace(row[colIdx])
		if convertType {
			dict[key] = convertActualType(val)
		} else {
			dict[key] = val
		}
	}
	return dict
}

// convertActualType 尝试根据字符串内容自动推断并转换为最合适的 Go 数据类型。
// 该函数按照 "空值 -> 布尔 -> 整数 -> 浮点数 -> 字符串" 的优先级顺序进行尝试。
//
// 转换逻辑：
// 1. 空字符串 ("") 会被转换为 nil。
// 2. "true" 或 "false" 会被转换为 bool 类型。
// 3. 符合整数格式的字符串（如 "123"）会被转换为 int64 类型。
// 4. 符合浮点数格式的字符串（如 "3.14"）会被转换为 float64 类型。
// 5. 如果以上都不匹配，则保持原始 string 类型返回。
func convertActualType(val string) interface{} {
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
