package csvutils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// detectCSVDelimiter 从 io.Reader 检测 CSV 分隔符（精简版）
// 核心特性：
// 1. 只读取第一行用于分析（不读取所有内容）
// 2. 跳过引号包裹的分隔符（符合 CSV 标准）
// 3. 跳过空白首行，读取第一个非空行分析
// 4. 固定候选分隔符：, \t | ;
func detectCSVDelimiter(r io.Reader) (rune, error) {
	// 1. 入参校验
	if r == nil {
		return 0, errors.New("reader is nil")
	}

	// 2. 创建带缓冲的 reader
	br := bufio.NewReader(r)

	// 3. 读取第一个非空行（只读取第一行）
	firstLine, err := readFirstNonEmptyLine(br)
	if err != nil {
		return ',', fmt.Errorf("read first line failed: %w", err)
	}
	if firstLine == "" {
		return ',', nil // 空文件，返回默认分隔符
	}

	// 4. 固定候选分隔符（按常见度排序）
	candidates := []rune{',', '\t', '|', ';'}

	// 5. 一次遍历首行，统计有效分隔符（跳过引号内的分隔符）
	sepCount := make(map[rune]int)
	inQuotes := false
	for _, char := range firstLine {
		// 处理引号包裹的内容（跳过内部分隔符）
		if char == '"' {
			inQuotes = !inQuotes
			continue
		}
		if inQuotes {
			continue
		}

		// 统计候选分隔符出现次数
		for _, sep := range candidates {
			if char == sep {
				sepCount[sep]++
				break
			}
		}
	}

	// 6. 选择出现次数最多的分隔符
	bestSep := ','
	maxCount := -1
	for sep, count := range sepCount {
		if count > maxCount {
			maxCount = count
			bestSep = sep
		}
	}

	return bestSep, nil
}

// readFirstNonEmptyLine 读取第一个非空行
// 特点：只读取第一行，不读取所有内容
func readFirstNonEmptyLine(br *bufio.Reader) (string, error) {
	for {
		// 读取一行（包括换行符）
		lineBytes, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return "", err
		}

		// 去除首尾空白（包括换行符）
		line := strings.TrimSpace(string(lineBytes))
		if line != "" {
			return line, nil // 找到非空行
		}

		// 已到文件末尾且未找到非空行
		if err == io.EOF {
			break
		}
	}

	return "", nil
}

// DetectCSVDelimiter 从文件路径检测 CSV 分隔符
func DetectCSVDelimiter(filePath string) (rune, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("open file failed: %w", err)
	}
	defer file.Close()

	return detectCSVDelimiter(file)
}

// DetectCSVBytesDelimiter 从字节切片检测 CSV 分隔符
func DetectCSVBytesDelimiter(csvBytes []byte) (rune, error) {
	return detectCSVDelimiter(bytes.NewReader(csvBytes))
}

// FormatCSVDelimiter 格式化csv的分隔符
func FormatCSVDelimiter(delimiters string) int32 {
	delimiter := ','
	if delimiters == "\\t" || strings.HasPrefix(delimiters, "t") {
		delimiter = '\t'
	} else if len(delimiters) == 1 {
		delimiter = rune(delimiters[0])
	}
	return delimiter
}
