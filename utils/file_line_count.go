package utils

import (
	"bufio"
	"io"
	"os"
)

// CountFileLines 统计文件的行数（适用于大文件）
func CountFileLines(filePath string, ignoreBlank bool) (int, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	lines, err := countIOLines(file, ignoreBlank)
	if err != nil {
		return 0, err
	}
	return int(lines), nil
}

// countIOLines 统计文件行数
// 当 ignoreBlank 为 true 时，仅统计包含非空白字符的行
func countIOLines(r io.Reader, ignoreBlank bool) (int64, error) {
	br := bufio.NewReaderSize(r, 4<<20)
	buf := make([]byte, 4<<20)
	var lines int64
	var hasContent bool
	for {
		n, rerr := br.Read(buf)
		if n > 0 {
			b := buf[:n]
			for i := 0; i < len(b); i++ {
				c := b[i]
				if c == '\n' {
					if !ignoreBlank || hasContent {
						lines++
					}
					hasContent = false
				} else if c == '\r' || c == ' ' || c == '\t' {
					// 空白字符，忽略
				} else {
					hasContent = true
				}
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			return 0, rerr
		}
	}
	if (!ignoreBlank && hasContent) || (ignoreBlank && hasContent) {
		lines++
	}
	return lines, nil
}
