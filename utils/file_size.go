package utils

import (
	"fmt"
	"os"
)

// FormatFileSize 将文件大小转换为人类可读的字符串
// 例如: 1024 -> 1.00 KB
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// GetFileSize 检查文件大小 文件占用的字节数。
func GetFileSize(filename string) (int64, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// BytesToKB 将字节数转换为 KB (浮点数)
// 1 KB = 1024 Bytes
func BytesToKB(bytes int64) float64 {
	return float64(bytes) / 1024.0
}

// BytesToMB 将字节数转换为 MB (浮点数)
// 1 MB = 1024 * 1024 Bytes = 1,048,576 Bytes
func BytesToMB(bytes int64) float64 {
	return float64(bytes) / (1024.0 * 1024.0)
}
