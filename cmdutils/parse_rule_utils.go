package cmdutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// splitAndTrim 按逗号分割，过滤空字符串与首尾空白
func splitAndTrim(s string) []string {
	var res []string
	for _, v := range strings.Split(s, ",") {
		if val := strings.TrimSpace(v); val != "" {
			res = append(res, val)
		}
	}
	return res
}

// readNonEmptyLines 读取文件非空行，自动去除换行符与首尾空白
func readNonEmptyLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file failed: %w", err)
	}
	defer file.Close()

	var res []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			res = append(res, line)
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("read file failed: %w", err)
	}
	return res, nil
}

// deduplicate 去重工具函数：保持顺序，去除重复字符串，空值自动忽略
func deduplicate(slice []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(slice))

	for _, s := range slice {
		// 空值直接跳过
		if s == "" {
			continue
		}
		// 未出现过则保留
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
