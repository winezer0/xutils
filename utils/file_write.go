package utils

import (
	"bufio"
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

// WriteLine 写入单行内容到文件，自动在内容后添加换行符
func WriteLine(filename, line string, overwrite bool) error {
	EnsureDir(filename, true)
	// 自动在内容末尾添加换行符(确保每行独立)
	content := line + "\n"
	// 确定文件打开模式
	flag := os.O_WRONLY | os.O_CREATE // 基础模式：可写 + 不存在则创建
	if overwrite {
		flag |= os.O_TRUNC // 覆盖模式：清空原有内容
	} else {
		flag |= os.O_APPEND // 追加模式：在末尾添加
	}
	// 打开文件(权限：644 表示 owner 可读写，其他可读取)
	file, err := os.OpenFile(filename, flag, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入带换行符的内容
	_, err = file.WriteString(content)
	return err
}

// WriteLines 将字符串切片逐行写入文件
// filename: 文件路径
// lines: 要写入的字符串切片
// appendMode:
//   - true: 覆盖模式 (如果文件存在则清空重写，不存在则创建)
//   - false:  追加模式 (在文件末尾添加内容，不存在则创建)
func WriteLines(filename string, lines []string, overwrite bool) error {
	EnsureDir(filename, true)
	// 确定文件打开模式
	flag := os.O_WRONLY | os.O_CREATE // 基础模式：可写 + 不存在则创建
	if overwrite {
		flag |= os.O_TRUNC // 覆盖模式：清空原有内容
	} else {
		flag |= os.O_APPEND // 追加模式：在末尾添加
	}
	// 打开文件
	file, err := os.OpenFile(filename, flag, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close() // 确保文件关闭
	// 创建带缓冲的写入器
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		// 写入内容并添加换行符
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}
	}
	// 重要：将缓冲区的内容刷新到磁盘
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}
	return nil
}
