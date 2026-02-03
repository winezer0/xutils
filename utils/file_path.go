package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// PathExists 检查路径是否存在，并返回是否为目录
func PathExists(path string) (exists bool, isDir bool, err error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		}
		return false, false, err
	}
	return true, info.IsDir(), nil
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	exists, isDir, _ := PathExists(path)
	return exists && !isDir
}

// DirExists 检查目录是否存在
func DirExists(path string) bool {
	exists, isDir, _ := PathExists(path)
	return exists && isDir
}

// GetAllFilePaths 获取指定路径下的所有文件路径
// 如果输入是文件路径，则直接返回包含该文件路径的切片
// 如果输入是目录路径，则返回该目录下所有文件（包括子目录中的文件）的路径
func GetAllFilePaths(path string) ([]string, error) {
	// 获取路径信息
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// 如果是文件，直接返回包含该文件路径的切片
	if !info.IsDir() {
		return []string{path}, nil
	}

	// 如果是目录，遍历所有文件
	var filePaths []string
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 传递错误
		}

		// 只添加文件，跳过目录
		if !info.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return filePaths, nil
}

// PathIsDir 判断路径是否为目录
func PathIsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

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

// EnsureDir 确保目录存在，如果不存在则创建。
// 如果 isFile 为 true，则 dirPath 被视为文件路径，函数会确保其所在目录存在；
// 如果 isFile 为 false，则 dirPath 被视为目录路径，函数会确保该目录存在。
func EnsureDir(dirPath string, isFile bool) error {
	targetDir := dirPath
	if isFile {
		targetDir = filepath.Dir(dirPath)
	}
	return os.MkdirAll(targetDir, 0755)
}

func GetFilesByGlob(path, glob string) ([]string, error) {
	//globMode := "*.yml"
	// 检查路径是否为目录
	isDir, err := PathIsDir(path)
	if err != nil {
		return nil, err
	}

	var files []string
	if !isDir {
		// 如果是文件，直接使用
		if FileExists(path) {
			files = []string{path}
		} else {
			return nil, err
		}
	} else {
		matches, _ := filepath.Glob(filepath.Join(path, glob))
		files = append(files, matches...)
	}
	return files, nil
}

// NormalizeFilepath 标准化文件路径，适配跨平台、清理冗余、处理相对路径
// file: 原始文件路径（支持Windows/Linux/macOS格式，支持相对/绝对路径）
// 返回值：标准化后的路径（统一使用系统原生分隔符，清理冗余，解析相对路径）
func NormalizeFilepath(file string) string {
	// 1. 空输入兜底
	if file == "" {
		return ""
	}

	// 2. 先统一将反斜杠替换为正斜杠（消除Windows分隔符差异）
	normalized := strings.ReplaceAll(file, "\\", "/")

	// 3. 使用标准库 filepath.Clean 做核心标准化：
	//    - 清理连续分隔符（a//b → a/b）
	//    - 解析.和..（a/../b → b）
	//    - 清理首尾分隔符（/a/b/ → a/b，绝对路径除外如/root/ → /root）
	//    - 自动适配系统分隔符（Windows返回\，Linux/macOS返回/）
	cleaned := filepath.Clean(normalized)

	// 4. 可选：若需强制统一返回/（无论系统），取消下方注释
	cleaned = strings.ReplaceAll(cleaned, "\\", "/")
	return cleaned
}
