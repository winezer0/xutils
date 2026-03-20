package csvutils

import (
	"encoding/csv"
	"github.com/winezer0/xutils/utils"
	"os"
)

// WriteRowsToCSV 将表头与数据行写入 CSV 文件
// 追加模式下会自动检测表头是否已存在，避免重复写入
// 若表头不匹配则返回错误
func WriteRowsToCSV(filePATH string, header []string, rows [][]string, delimiter rune, overwrite bool) error {
	if len(rows) == 0 && len(header) == 0 {
		return nil
	}

	needWriteHeader := ShouldWriteHeader(filePATH, header, overwrite, delimiter)

	flag := utils.ParseFlagFromOver(overwrite)

	file, err := os.OpenFile(filePATH, flag, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Comma = delimiter
	defer w.Flush()

	if needWriteHeader && len(header) > 0 {
		if err := w.Write(header); err != nil {
			return err
		}
	}

	if len(rows) > 0 {
		if err := w.WriteAll(rows); err != nil {
			return err
		}
	}

	return w.Error()
}
