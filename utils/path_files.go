package utils

import (
	"os"
	"path/filepath"
)

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
