package csvutils

import (
	"bufio"
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
	separator, _ := detectCSVSeparator(csvFile)
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

// ReadCSVToDicts 读取 CSV 文件并返回[]string, []map[string]string
// 第一行为 header，后续每行为一条记录
func ReadCSVToDicts(csvFile string) ([]string, []map[string]string, error) {
	separator, _ := detectCSVSeparator(csvFile)
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
	headers = cleanHeaders(headers)
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

// detectCSVSeparator 从 io.Reader 检测分隔符（便于测试）
func detectCSVSeparator(filePath string) (rune, error) {
	separator := ','

	file, err := os.Open(filePath)
	if err != nil {
		return separator, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return separator, nil // 空文件，返回默认
	}
	firstLine := scanner.Text()

	// 候选分隔符（按优先级或常见度排序）
	candidates := []rune{',', '\t', '|', ';'}

	var bestSep rune
	maxCount := -1

	for _, sep := range candidates {
		count := strings.Count(firstLine, string(sep))
		// 只有当分隔符出现 >=1 次才考虑（避免空行干扰）
		if count > maxCount && count > 0 {
			maxCount = count
			bestSep = sep
		}
	}

	// 如果没找到有效分隔符，返回默认
	if bestSep == 0 {
		return separator, nil
	}

	return bestSep, nil
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
