package utils

import (
	"io"
	"os"
	"path/filepath"
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

// PathIsDir 判断路径是否为目录
func PathIsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
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

// IsEmptyFile 检查文件是否为空或不存在
func IsEmptyFile(filename string) bool {
	// Get file info
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		return true
	}
	return false
}

// IsDirEmpty 判断目录是否为空（无任何子项）
func IsDirEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	return err == io.EOF, nil
}
