package csvutils

import (
	"encoding/csv"
	"fmt"
	"os"
)

// GetCSVHeaders 获取所有CSV文件的头部信息
func GetCSVHeaders(filePath string, delimiter rune) ([]string, error) {
	// 处理空分隔符：默认设为逗号
	if delimiter == 0 { // 核心判断逻辑
		delimiter, _ = DetectCSVDelimiter(filePath)
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1 // 允许每行字段数不一致
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	headers = fixedHeaders(headers)

	return headers, nil
}

// GetCSVSHeaders 从所有CSV文件中收集所有唯一的头部字段
func GetCSVSHeaders(csvFiles []string, delimiter rune) ([]string, error) {
	allHeaders := make([]string, 0)
	seenHeaders := make(map[string]bool)

	for _, filePath := range csvFiles {
		headers, err := GetCSVHeaders(filePath, delimiter)
		if err != nil {
			return nil, fmt.Errorf("error reading headers from %s: %v", filePath, err)
		}

		for _, header := range headers {
			if !seenHeaders[header] {
				seenHeaders[header] = true
				allHeaders = append(allHeaders, header)
			}
		}
	}

	return allHeaders, nil
}
