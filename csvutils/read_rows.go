package csvutils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

// ReadCSV2Rows 修复版：表头固定取第一行，skipRows 仅控制数据行跳过
// 核心修正：
// 1. 表头永远取 CSV 第一行（符合通用 CSV 读写习惯）
// 2. skipRows 仅跳过数据行，不影响表头提取
func ReadCSV2Rows(filePath string, delimiter rune, skipRows int) (header []string, rows [][]string, err error) {
	// 1. 边界校验
	if skipRows < 0 {
		return nil, nil, errors.New("skipRows cannot be negative")
	}

	// 2. 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("file not found: %s", filePath)
		}
		if os.IsPermission(err) {
			return nil, nil, fmt.Errorf("permission denied to open file: %s", filePath)
		}
		return nil, nil, fmt.Errorf("open file failed: %w", err)
	}
	defer file.Close()

	// 3. 初始化 Reader
	reader := csv.NewReader(file)
	if delimiter != 0 {
		reader.Comma = delimiter
	}
	reader.FieldsPerRecord = -1 // 允许字段数不一致

	// 4. 逐行读取（核心修正：先读表头，再读数据行）
	var (
		lineNum    = 0
		readHeader = false // 是否已读取表头
	)

	for {
		row, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, nil, fmt.Errorf("read csv failed at line %d: %w", lineNum+1, err)
		}

		lineNum++

		// 第一步：读取表头（第一行固定为表头）
		if !readHeader {
			header = append([]string{}, row...) // 深拷贝
			readHeader = true
			// 若 skipRows=0，表头行不加入数据行；否则继续判断
			if skipRows == 0 {
				continue
			}
		}

		// 第二步：处理数据行（跳过指定行数）
		// 数据行行号 = lineNum - 1（因为表头占了第一行）
		dataLineNum := lineNum - 1
		if dataLineNum > skipRows {
			rows = append(rows, append([]string{}, row...)) // 深拷贝
		}
	}

	// 5. 边界处理
	if lineNum == 0 {
		return nil, nil, fmt.Errorf("csv file %s is empty", filePath)
	}
	// 无表头场景（文件只有一行时，表头=该行，数据行空）
	if !readHeader {
		header = nil
	}

	return header, rows, nil
}
