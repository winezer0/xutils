package csvutils

import (
	"encoding/csv"
	"fmt"
	"github.com/winezer0/xutils/utils"
	"os"
)

// WriteDictsToCSV 修复版：解决空分隔符+字段数不一致问题
// 核心修复：
// 1. 空分隔符默认设为逗号（避免 csv.Writer 报错）
// 2. 兼容字段数不一致场景（CSV 读取时允许字段数不匹配）
func WriteDictsToCSV(filePath string, dicts []map[string]interface{}, header []string, delimiter rune, overwrite bool) error {
	// 边界处理：无数据时直接返回
	if len(dicts) == 0 {
		return nil
	}

	// 修复1：空分隔符默认设为逗号（解决 invalid delimiter 错误）
	if delimiter == 0 {
		delimiter = ','
	}

	// 1. 处理表头：未指定则从第一个字典生成
	usedHeader, err := GetCSVHeaderFromDicts(dicts, header, false)
	if err != nil {
		return fmt.Errorf("find dict list to header failed: %w", err)
	}

	// 2. 字典列表转 CSV 数据行
	rows, err := dictListToRows(dicts, usedHeader)
	if err != nil {
		return fmt.Errorf("convert dict list to rows failed: %w", err)
	}

	// 3. 确定文件打开模式
	flag := utils.ParseFlagFromOver(overwrite)

	// 4. 打开文件
	file, err := os.OpenFile(filePath, flag, 0644)
	if err != nil {
		return fmt.Errorf("open file %s failed: %w", filePath, err)
	}
	defer file.Close()

	// 5. 初始化 CSV Writer
	writer := csv.NewWriter(file)
	writer.Comma = delimiter
	defer writer.Flush()

	// 6. 调用 shouldWriteHeader 判断是否写表头
	needWriteHeader := shouldWriteHeader(filePath, usedHeader, overwrite, delimiter)

	// 7. 写入表头（若需要）
	if needWriteHeader {
		if err := writer.Write(usedHeader); err != nil {
			return fmt.Errorf("write header to %s failed: %w", filePath, err)
		}
	}

	// 8. 写入数据行
	if len(rows) > 0 {
		if err := writer.WriteAll(rows); err != nil {
			return fmt.Errorf("write data rows to %s failed: %w", filePath, err)
		}
	}

	// 9. 检查 Writer 内部错误
	if err := writer.Error(); err != nil {
		return fmt.Errorf("csv writer flush error for %s: %w", filePath, err)
	}

	return nil
}

// ------------------- 依赖的辅助函数（不变） -------------------
func dictListToRows(dicts []map[string]interface{}, header []string) ([][]string, error) {
	rows := make([][]string, 0, len(dicts))
	for _, dict := range dicts {
		row := make([]string, len(header))
		for colIdx, key := range header {
			val, ok := dict[key]
			if !ok {
				row[colIdx] = ""
				continue
			}
			row[colIdx] = valueToString(val)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func valueToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}
