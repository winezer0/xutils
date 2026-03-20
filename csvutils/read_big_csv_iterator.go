package csvutils

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

// CSVIterator CSV 迭代器（用于逐行读取大文件）
// 修复点：
// 1. 新增 file 字段保存文件句柄（替代访问 csv.Reader 私有 Input 字段）
// 2. 添加缓冲读取（适配大文件）
// 3. 补充所有缺失的辅助函数
// 4. 增强容错配置（LazyQuotes、TrimLeadingSpace）
type CSVIterator struct {
	reader      *csv.Reader
	file        *os.File // 新增：保存文件句柄，用于关闭
	header      []string
	convertType bool
	lineNum     int
	err         error
}

// NewCSVIterator 创建大文件 CSV 迭代器（修复版）
func NewCSVIterator(filePath string, delimiter rune, convertType bool) (*CSVIterator, error) {
	// 1. 打开文件（保存句柄，用于后续关闭）
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open file failed: %w", err)
	}

	// 2. 手动创建带缓冲的 reader（适配大文件，替代 csv.Reader.BufferSize）
	bufReader := bufio.NewReaderSize(file, 1024*1024) // 1MB 缓冲区

	// 3. 初始化 CSV Reader（增强容错）
	reader := csv.NewReader(bufReader)
	if delimiter != 0 {
		reader.Comma = delimiter
	}
	// 核心容错配置
	reader.FieldsPerRecord = -1    // 允许字段数不一致
	reader.LazyQuotes = true       // 宽松引号解析
	reader.TrimLeadingSpace = true // 自动修剪字段前后空格

	iter := &CSVIterator{
		reader:      reader,
		file:        file, // 保存文件句柄
		convertType: convertType,
		lineNum:     1,
	}

	// 4. 读取表头
	row, err := reader.Read()
	if err != nil {
		file.Close() // 读取失败时关闭文件
		return nil, fmt.Errorf("read header failed: %w", err)
	}
	iter.header = FixedHeaders(row)
	return iter, nil
}

// Next 读取下一行（返回nil表示读取完毕）
func (iter *CSVIterator) Next() map[string]interface{} {
	if iter.err != nil {
		return nil
	}

	row, err := iter.reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			iter.err = err
			return nil
		}
		// 单行读取失败：记录错误，返回nil（不中断迭代）
		iter.err = fmt.Errorf("read row %d failed: %w", iter.lineNum, err)
		return nil
	}

	// 首次读取确定默认表头（无表头场景）
	if len(iter.header) == 0 {
		iter.header = GenDefaultHeaders(len(row))
	}

	iter.lineNum++
	// 行数据转字典
	return RowDataToDict(row, iter.header, iter.convertType)
}

// Error 返回迭代过程中的错误
func (iter *CSVIterator) Error() error {
	return iter.err
}

// Close 关闭迭代器（修复版：直接关闭保存的 file 句柄）
func (iter *CSVIterator) Close() error {
	if iter.file != nil {
		return iter.file.Close()
	}
	return nil
}
