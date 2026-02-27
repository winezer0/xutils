package utils

import (
	"os"
	"strings"
)

// ReplaceUserHome 替换路径中的用户目录
func ReplaceUserHome(path string) string {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			path = strings.Replace(path, "~", homeDir, 1)
		}
	}
	return path
}
