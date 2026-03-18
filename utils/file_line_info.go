package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CodeLine 表示带行号的代码行
type CodeLine struct {
	Number  int    // 行号（从1开始）
	Content string // 行内容
}

// CodeChunk 表示一个代码片段，包含多行带行号的代码
type CodeChunk struct {
	StartLine int    // 起始行号
	EndLine   int    // 结束行号
	Content   string // 包含行号标识的代码内容
}

// readFileWithLines 读取文件内容，按行拆分并生成行号与代码内容的映射
func readFileWithLines(filePath string) ([]CodeLine, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	var lines []CodeLine
	scanner := bufio.NewScanner(file)
	lineNum := 1

	// 逐行读取文件内容
	for scanner.Scan() {
		lineContent := scanner.Text()
		lines = append(lines, CodeLine{
			Number:  lineNum,
			Content: lineContent,
		})
		lineNum++
	}

	// 检查扫描过程中是否发生错误
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %v", err)
	}

	return lines, nil
}

// splitFileToChunks 根据预设分片阈值，将行号映射表拆分为多个代码片段
func splitFileToChunks(lines []CodeLine, chunkSize int) []CodeChunk {
	var chunks []CodeChunk

	// 如果未指定分片大小，使用默认值
	if chunkSize <= 0 {
		chunkSize = 200
	}

	// 如果文件为空，返回空切片
	if len(lines) == 0 {
		return chunks
	}

	// 计算分片数量
	numChunks := (len(lines) + chunkSize - 1) / chunkSize

	// 生成代码片段
	for i := 0; i < numChunks; i++ {
		startIdx := i * chunkSize
		endIdx := (i + 1) * chunkSize
		if endIdx > len(lines) {
			endIdx = len(lines)
		}

		// 构建包含行号标识的代码内容
		var chunkContent strings.Builder
		for _, line := range lines[startIdx:endIdx] {
			chunkContent.WriteString(fmt.Sprintf("%d| %s\n", line.Number, line.Content))
		}

		// 创建代码片段
		chunk := CodeChunk{
			StartLine: lines[startIdx].Number,
			EndLine:   lines[endIdx-1].Number,
			Content:   chunkContent.String(),
		}

		chunks = append(chunks, chunk)
	}

	return chunks
}

// ReadFileToCodeChunk 处理单个文件，生成代码片段
func ReadFileToCodeChunk(filePath string, chunkSize int) ([]CodeChunk, error) {
	// 读取文件内容并生成带行号的代码行
	lines, err := readFileWithLines(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败 %s: %v", filePath, err)
	}

	// 将文件内容拆分为代码片段
	chunks := splitFileToChunks(lines, chunkSize)
	return chunks, nil
}

// ExtractCodeFromFile 从指定文件中提取 [startLine, endLine] 范围内的代码内容，
// 返回格式为：每行以 "行号| 内容" 的形式拼接，末尾带换行符。
func ExtractCodeFromFile(filePath string, startLine, endLine int) (string, error) {
	if startLine < 1 || endLine < 1 || startLine > endLine {
		return "", fmt.Errorf("invalid line range: start=%d, end=%d", startLine, endLine)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var result strings.Builder
	currentLine := 1

	for scanner.Scan() {
		if currentLine > endLine {
			break // 已超出目标范围
		}
		if currentLine >= startLine {
			result.WriteString(fmt.Sprintf("%d| %s\n", currentLine, scanner.Text()))
			//result.WriteString(fmt.Sprintf("%s\n", scanner.Text()))
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	output := result.String()
	if output == "" {
		return "", fmt.Errorf("no content found in line range [%d, %d] in file %s", startLine, endLine, filePath)
	}

	return output, nil
}
