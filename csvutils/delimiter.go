package csvutils

import (
	"bufio"
	"os"
	"strings"
)

// DetectCSVDelimiter 从 io.Reader 检测分隔符（便于测试）
func DetectCSVDelimiter(filePath string) (rune, error) {
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
