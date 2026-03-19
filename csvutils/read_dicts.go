package csvutils

import (
	"encoding/csv"
	"fmt"
	"os"
)

// ReadCSVToDicts 读取 CSV 文件并返回[]string, []map[string]string
// 第一行为 header，后续每行为一条记录
func ReadCSVToDicts(csvFile string) ([]string, []map[string]string, error) {
	separator, _ := DetectCSVDelimiter(csvFile)
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = separator
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, nil, fmt.Errorf("CSV content is empty")
	}

	headers := records[0]
	headers = fixedHeaders(headers)
	if len(headers) == 0 {
		return nil, nil, fmt.Errorf("CSV header is empty")
	}

	var dicts []map[string]string
	for i, row := range records[1:] { // 跳过 header
		if len(row) == 0 {
			continue // 跳过空行
		}

		// 补齐缺失字段（如果行比 header 短）
		for len(row) < len(headers) {
			row = append(row, "")
		}
		// 忽略多余字段（如果行比 header 长）
		if len(row) > len(headers) {
			row = row[:len(headers)]
		}

		record := make(map[string]string)
		for j, key := range headers {
			record[key] = row[j]
		}
		dicts = append(dicts, record)

		// 可选：验证字段数量（调试用）
		if len(row) != len(headers) {
			return nil, nil, fmt.Errorf("row %d has %d fields, expected %d", i+2, len(row), len(headers))
		}
	}

	return headers, dicts, nil
}
