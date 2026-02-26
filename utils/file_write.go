package utils

import (
	"fmt"
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

// CreateFile 创建文件，如果目录不存在会自动创建
func CreateFile(filePath string) (*os.File, error) {
	// 确保目录存在
	if err := EnsureDir(filePath, true); err != nil {
		return nil, err
	}
	// 创建文件
	return os.Create(filePath)
}

// WriteToFile 将任意数据写入文本文件
func WriteToFile(filePath string, data interface{}) error {
	// 将任意数据转换为字符串形式
	content := fmt.Sprintf("%+v", data)
	return SaveToFile(filePath, []byte(content))
}
