package utils

import (
	"bufio"
	"os"
)

// DeduplicateFile 对文件内容进行去重并覆盖原文件
func DeduplicateFile(filePath string) error {
	uniqueLines := make(map[string]struct{})
	var lines []string

	// 读取文件
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	scanner := bufio.NewScanner(f)
	// 增加Buffer大小
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if _, exists := uniqueLines[line]; !exists {
			uniqueLines[line] = struct{}{}
			lines = append(lines, line)
		}
	}
	f.Close()

	if err := scanner.Err(); err != nil {
		return err
	}

	// 覆盖写入
	outFile, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}
