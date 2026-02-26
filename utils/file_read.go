package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadFileByRange 基于代码行号读取文件内容
func ReadFileByRange(file string, startLine, endLine int) string {
	if file == "" || startLine <= 0 || endLine < startLine {
		return ""
	}
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	idx := 0
	for sc.Scan() {
		idx++
		if idx < startLine {
			continue
		}
		if idx > endLine {
			break
		}
		lines = append(lines, sc.Text())
	}
	// simple trim leading/trailing blanks
	i := 0
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}
	j := len(lines) - 1
	for j >= i && strings.TrimSpace(lines[j]) == "" {
		j--
	}
	if i <= j {
		return strings.Join(lines[i:j+1], "\n")
	}
	return ""
}

// ReadFileToBytes 读取指定路径的文件内容，返回字节切片和错误信息
func ReadFileToBytes(path string) ([]byte, error) {
	// 使用 os.ReadFile 一次性读取文件全部内容到内存，返回字节切片和错误
	// os.ReadFile 底层已封装了文件打开、读取、关闭的完整流程，简洁且安全
	content, err := os.ReadFile(path)
	if err != nil {
		// 读取失败时，返回 nil 和原始错误（调用方可通过 errors.Is/As 进一步判断错误类型）
		return nil, err
	}

	if len(content) == 0 {
		// 读取失败时，返回 nil 和原始错误（调用方可通过 errors.Is/As 进一步判断错误类型）
		return nil, fmt.Errorf("file is empty")
	}

	// 读取成功时，返回字节切片和 nil
	return content, nil
}

// IsEmptyFile 检查文件是否为空或不存在
func IsEmptyFile(filename string) bool {
	// Get file info
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		return true
	}
	return false
}
