package csvutils

import (
	"encoding/csv"
	"os"
)

// CountCSVLines 统计指定 CSV 文件的总行数（记录数）。
// 该函数以只读模式打开文件，利用 csv.Reader 进行流式读取，适用于处理大文件。
//
// 参数说明：
//   - filePath: CSV 文件的绝对或相对路径。
//   - delimiter: CSV 文件的分隔符（如 ',' 或 '\t'）。
//
// 返回值：
//   - int: 文件的总行数。
//   - error: 执行过程中遇到的错误（如文件不存在、权限不足或读取失败）。
func CountCSVLines(filePath string, delimiter rune) (int, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1 // 允许每行字段数不一致

	count := 0
	for {
		_, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return 0, err
		}
		count++
	}

	return count, nil
}
