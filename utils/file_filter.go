package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// FilterByExtension 检查文件是否应该根据扩展名被排除
// 返回true表示需要排除，false表示保留
func FilterByExtension(path string, ExPathExts []string) bool {
	// 若排除关键字列表为空，直接返回不排除
	if len(ExPathExts) == 0 {
		return false
	}
	ext := strings.ToLower(filepath.Ext(path))
	for _, excludeExt := range ExPathExts {
		if excludeExt == ext {
			return true
		}
	}
	return false
}

// FilterByPathKeys 检查路径是否包含任何排除关键字，若包含则返回true表示需要排除
func FilterByPathKeys(path string, exPathKeys []string) bool {
	// 若排除关键字列表为空，直接返回不排除
	if len(exPathKeys) == 0 {
		return false
	}
	// 遍历所有排除关键字，检查路径是否包含其中任何一个
	path = strings.ToLower(path)
	for _, key := range exPathKeys {
		if strings.Contains(path, key) {
			return true // 包含关键字，需要排除
		}
	}
	return false // 不包含任何关键字，无需排除
}

// FileIsLarger 检查文件是否超过指定大小（MB）
func FileIsLarger(filePath string, limitSize int) bool {
	if limitSize <= 0 {
		return false
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return info.Size() > int64(limitSize*1024*1024)
}

// GetFilesWithFilter 获取符合条件的文件列表
func GetFilesWithFilter(path string, excludeSuffixes, excludePathKeys []string, limitSize int) ([]string, error) {
	var files []string

	// 获取所有文件列表
	allFile, err := GetAllFilePaths(path)
	if err != nil && len(allFile) == 0 {
		return nil, err
	}

	// 过滤文件
	excludePathKeys = ToLowerKeys(excludePathKeys)
	excludeSuffixes = ToLowerKeys(excludeSuffixes)
	for _, file := range allFile {
		if FilterByExtension(file, excludeSuffixes) {
			continue
		}

		if FilterByPathKeys(file, excludePathKeys) {
			continue
		}

		if FileIsLarger(file, limitSize) {
			continue
		}

		files = append(files, file)
	}
	return files, err
}

// IsAllowedSuffixes 检查文件是否需要根据规则进行处理。
// 如果 suffixes 参数是 "*"，则返回 true。
// 否则，检查 filename 是否以 suffixes 中任意一个后缀结尾（忽略大小写）。
func IsAllowedSuffixes(filename string, suffixesStr string) bool {
	// 如果 suffixes 是 "*", 直接返回 true
	if suffixesStr == "*" {
		return true
	}

	// 将 suffixes 字符串按 "|" 分割成 slice
	suffixes := strings.Split(suffixesStr, "|")

	// 转换 filename 和 suffixes 到小写以便忽略大小写比较
	filenameLower := strings.ToLower(filename)

	for _, suffix := range suffixes {
		// 忽略大小写比较
		if strings.HasSuffix(filenameLower, strings.ToLower(suffix)) {
			return true
		}
	}

	return false
}

func CollectFiles(paths []string, suffixesStr string) ([]string, error) {
	finalFiles := make([]string, 0, len(paths))
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		// 判断路径是否存在
		if exists, _, _ := PathExists(path); !exists {
			continue
		}
		// 获取目录下的所有文件
		files, _ := GetAllFilePaths(path)
		for _, file := range files {
			if IsAllowedSuffixes(file, suffixesStr) {
				finalFiles = append(finalFiles, file)
			}
		}
	}

	finalFiles = UniqueSlice(finalFiles, false, true)
	return finalFiles, nil
}
