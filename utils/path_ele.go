package utils

import (
	"path/filepath"
	"strings"
)

// GetPathLastDir 从文件路径获取目录名称 // 如果路径是文件，则返回其所在目录名 // 如果路径是目录，则返回该目录名
func GetPathLastDir(path string) string {
	// 规范化路径
	path = filepath.Clean(path)

	// 检查路径是否是目录
	isDir, err := PathIsDir(path)
	if err != nil {
		// 如果出错，直接返回路径的基础名称
		return filepath.Base(path)
	}

	if isDir {
		// 如果是目录，返回目录名
		return filepath.Base(path)
	} else {
		// 如果是文件，返回其父目录名
		return filepath.Base(filepath.Dir(path))
	}
}

// GetFileBaseName 从文件路径获取文件的基础名称
func GetFileBaseName(path, defaultName string, keepDot bool) string {
	base := filepath.Base(path)
	if base == "" {
		base = defaultName
	}

	if !keepDot {
		dot := strings.LastIndex(base, ".")
		if dot > 0 {
			base = base[:dot]
		}
	}

	return base
}
