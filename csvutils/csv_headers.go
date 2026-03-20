package csvutils

import (
	"encoding/csv"
	"fmt"
	"os"
)

// GetCSVHeaders 读取并解析 CSV 文件的首行作为表头信息。
// 该函数会自动处理分隔符检测（如果未指定）以及表头清洗（去重、去空格、处理 BOM）。
//
// 参数说明：
//   - filePath: CSV 文件的绝对或相对路径。
//   - delimiter: CSV 文件的分隔符。如果传入 0，函数将自动尝试检测分隔符。
//
// 返回值：
//   - []string: 清洗后的表头字符串切片。
//   - error: 执行过程中遇到的错误（如文件打开失败、读取失败）。
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

	headers = FixedHeaders(headers)

	return headers, nil
}

// GetCSVSHeaders 从所有 CSV 文件中收集所有唯一的头部字段
// 参数说明：
//   - filePath: CSV 文件的绝对或相对路径。
//   - delimiter: CSV 文件的分隔符。如果传入 0，函数将自动尝试检测分隔符。
//
// 返回值：
//   - []string: 清洗后的表头字符串切片。
//   - []error: 执行过程中遇到的错误（如文件打开失败、读取失败）。
func GetCSVSHeaders(csvFiles []string, delimiter rune) ([]string, []error) {
	seenHeaders := make(map[string]bool)
	allHeaders := make([]string, 0)
	var errors []error

	for _, filePath := range csvFiles {
		// 1. 独立处理每个文件的分隔符
		// 如果外部传入 0，每个文件都会独立调用 DetectCSVDelimiter
		// 如果外部传入具体字符（如 ','），则所有文件强制使用该分隔符
		currentDelimiter := delimiter
		if currentDelimiter == 0 {
			d, err := DetectCSVDelimiter(filePath)
			if err != nil {
				errors = append(errors, fmt.Errorf("detect delimiter failed for %s: %v", filePath, err))
				continue // 跳过此文件
			}
			currentDelimiter = d
		}

		// 2. 获取表头
		headers, err := GetCSVHeaders(filePath, currentDelimiter)
		if err != nil {
			errors = append(errors, fmt.Errorf("read headers failed for %s: %v", filePath, err))
			continue // 跳过此文件，继续处理下一个
		}

		// 3. 去重与合并
		for _, header := range headers {
			key := header
			if !seenHeaders[key] {
				seenHeaders[key] = true
				allHeaders = append(allHeaders, header) // 保留原始格式的输出
			}
		}
	}

	return allHeaders, errors
}
