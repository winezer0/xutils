package csvutils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// ReadCSVHeaders 读取 CSV 文件的第一行作为列名（header）
// 如果文件为空或没有 header，返回错误
func ReadCSVHeaders(csvFile string) ([]string, error) {
	separator, _ := DetectCSVDelimiter(csvFile)
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = separator
	headers, err := reader.Read() // 读取第一行
	headers = cleanHeaders(headers)
	if err != nil {
		if errors.Is(err, csv.ErrFieldCount) {
			// 即使字段不规范，只要读到了内容，仍可返回
			// 但通常 Read() 在空文件时返回 io.EOF
			return nil, errors.New("failed to read CSV header")
		}
		return nil, err
	}

	// 可选：去除每列的前后空格（根据需求决定）
	for i, h := range headers {
		headers[i] = strings.TrimSpace(h)
	}

	return headers, nil
}

// cleanHeaders 清理 headers，移除 BOM 和空白
func cleanHeaders(headers []string) []string {
	if headers == nil {
		return nil
	}

	cleaned := make([]string, len(headers))
	for i, h := range headers {
		// 安全移除 UTF-8 BOM
		if strings.HasPrefix(h, "\ufeff") {
			h = h[3:] // "\ufeff" 在 UTF-8 中占 3 字节
		}
		cleaned[i] = strings.TrimFunc(h, unicode.IsSpace)
	}
	return cleaned
}

// GetCSVHeaders 获取所有CSV文件的头部信息
func GetCSVHeaders(filePath string, delimiter rune) ([]string, error) {
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
