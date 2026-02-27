package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ResolveProjectFilePath 根据文件路径和项目路径 获取实际的文件路径
func ResolveProjectFilePath(projectPath, filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("file path is empty")
	}

	// 首先尝试以当前工作目录为基准直接进行查找
	if FileExists(filePath) {
		if abs, err := filepath.Abs(filePath); err == nil {
			return abs, nil
		}
		return filePath, nil
	}

	// 根据项目路径进行尝试
	if projectPath != "" {
		fullPath := filepath.Join(projectPath, filePath)

		// 1. 直接拼接后访问
		if FileExists(fullPath) {
			return fullPath, nil
		}

		// 2. 检查文件路径是否以项目路径的最后一个元素开头
		if possiblePath, exist := resolvePossiblePath(projectPath, filePath); exist {
			return possiblePath, nil
		}
	}

	return "", fmt.Errorf("project path (%s) error or file (%s) not found", projectPath, filePath)
}

// resolvePossiblePath 检查文件路径是否以项目路径的最后一个元素开头
func resolvePossiblePath(projectPath, filePath string) (string, bool) {
	// 2. 检查文件路径是否以项目路径的最后一个元素开头
	// 例如：项目路径为“xxxx/WebFTP_3.6.2”，文件路径为“WebFTP_3.6.2/core/...”
	projectBase := filepath.Base(projectPath)
	filePath = filepath.Clean(filePath)
	prefix := projectBase + string(filepath.Separator)
	if strings.HasPrefix(filePath, prefix) {
		relPath := strings.TrimPrefix(filePath, prefix)
		fullPath := filepath.Join(projectPath, relPath)
		if FileExists(fullPath) {
			return fullPath, true
		}
	}
	return "", false
}
