package utils

import (
	"bufio"
	"fmt"
	"io"
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

// ReadFileToStr 读取指定路径的文件内容，返回字符串和错误信息
func ReadFileToStr(path string) (string, error) {
	b, err := ReadFileToBytes(path)
	return string(b), err
}

// ReadFileToList 读文件到列表 自动忽略空行
func ReadFileToList(input string, ignoreBlanks, cleanUnprint bool) ([]string, error) {
	file, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	// 增加Buffer大小以防行过长，虽然一般hash/plain不会太长
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 移除不可打印字符
		if cleanUnprint {
			line = CleanupUnprintableChars(line)
		}

		// 跳过空白行
		if ignoreBlanks && line == "" {
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// ReadFileFirstLine 读取文件第一行（去除首尾空白），专为大文件优化
// - 仅读取必要内容，内存占用恒定（与文件总大小无关）
// - 正确处理：空文件、无换行符结尾、Windows(\r\n)/Unix(\n)换行符、超长行
// - 错误时静默返回空字符串（与原函数行为一致）
func ReadFileFirstLine(path string, trimSpace bool) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')

	// 仅当发生非EOF错误时返回空（如IO错误）
	if err != nil && err != io.EOF {
		return "", err
	}

	// TrimSpace 自动处理 \r, \n, 空格等空白字符（兼容Windows/Unix）
	if trimSpace {
		line = strings.TrimSpace(line)
	}
	return line, nil
}
