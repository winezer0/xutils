package csvutils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

// ReadCSV2Rows 读取 CSV 文件，根据 haveHeader 参数决定是否提取表头。
//
// 核心逻辑：
// - haveHeader = true: 第一行作为表头，后续行为数据。
// - haveHeader = false: 所有行（包括第一行）都作为数据处理，表头返回 nil。
//
// 参数说明：
//   - filePath: CSV 文件路径。
//   - delimiter: 分隔符（0 表示默认逗号）。
//   - haveHeader: 是否包含表头。
//
// 返回值：
//   - header: 表头切片（若无表头则为 nil）。
//   - rows: 数据行切片。
//   - err: 错误信息。
func ReadCSV2Rows(filePath string, delimiter rune, haveHeader bool) (header []string, rows [][]string, err error) {
	// 1. 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("open file failed: %w", err)
	}
	defer file.Close()

	// 2. 初始化 Reader
	reader := csv.NewReader(file)
	if delimiter != 0 {
		reader.Comma = delimiter
	}
	reader.FieldsPerRecord = -1 // 允许字段数不一致

	// 3. 读取所有行
	// 使用 ReadAll 一次性读取，代码最简洁且性能足够好 ReadAll内部已经做好了缓冲处理。
	allRecords, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("read csv failed: %w", err)
	}

	// 4. 处理空文件
	if len(allRecords) == 0 {
		return nil, nil, nil
	}

	// 5. 根据 haveHeader 拆分数据
	if haveHeader {
		// 第一行是表头
		header = allRecords[0]
		// 剩余行是数据
		rows = allRecords[1:]
	} else {
		// 没有表头，全是数据
		header = GenDefaultHeaders(len(allRecords[0]))
		rows = allRecords
	}

	return header, rows, nil
}

// ReadCSV2RowsWithSkip 读取 CSV 文件内容，支持跳过指定数量的数据行。
//
// 核心逻辑说明：
// 1. 表头固定：无论 skipRows 设置为多少，文件的物理第一行永远被视为表头。
// 2. 数据跳过：skipRows 仅作用于表头之后的数据行，用于忽略紧跟在表头后的元数据或注释行。
//
// 参数说明：
//   - filePath: CSV 文件的绝对或相对路径。
//   - delimiter: CSV 分隔符（rune类型）。若传入 0，则使用默认分隔符（通常为逗号）。
//   - skipRows: 需要跳过的数据行数（从表头后的第一行开始计算）。
//     例如：skipRows=1 表示读取表头后，忽略紧接着的 1 行数据。
//
// 返回值：
//   - header: 提取的表头字符串切片。
//   - rows: 过滤后的数据行二维切片。
//   - err: 执行过程中遇到的错误（如文件不存在、权限不足、读取失败等）。
func ReadCSV2RowsWithSkip(filePath string, delimiter rune, skipRows int) (header []string, rows [][]string, err error) {
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
