package utils

import (
	"bufio"
	"os"
)

// AppendLines 向文件中追加行 (不进行去重)
func AppendLines(filePath string, newLines []string) error {
	if len(newLines) == 0 {
		return nil
	}

	// 追加写入 (a+ 模式)
	outFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	for _, line := range newLines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}
