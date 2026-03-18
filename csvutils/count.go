package csvutils

import (
	"encoding/csv"
	"os"
)

// CountCSVLines 统计CSV文件的行数
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
