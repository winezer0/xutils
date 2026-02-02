package utils

import (
	"os"
)

// SaveToFile 写入文件内容，如果目录不存在会自动创建
func SaveToFile(filePath string, data []byte) error {
	// 确保目录存在
	if err := EnsureDir(filePath, true); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, data, 0644)
}
