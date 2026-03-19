package utils

import (
	"os"
	"path/filepath"
)

// normalizeExts 规范化扩展名集合 若包含 * 则返回空集合并视为匹配所有文件
func normalizeExts(exts []string) (map[string]struct{}, bool) {
	m := make(map[string]struct{})
	for _, e := range exts {
		if e == "*" {
			return map[string]struct{}{}, true
		}
		if e == "" {
			continue
		}
		if e[0] != '.' {
			e = "." + e
		}
		m[e] = struct{}{}
	}
	if len(m) == 0 {
		m[".txt"] = struct{}{}
	}
	return m, false
}

// matchExts 判断文件是否匹配扩展名
func matchExts(path string, allowAll bool, exts map[string]struct{}) bool {
	if allowAll {
		return true
	}
	ext := filepath.Ext(path)
	_, ok := exts[ext]
	return ok
}

// CollectFiles 收集输入路径中的目标文件（支持目录递归与单文件）
func collectFiles(paths []string, allowAll bool, exts map[string]struct{}) ([]string, error) {
	var files []string
	for _, p := range paths {
		abs := p
		if !filepath.IsAbs(p) {
			a, err := filepath.Abs(p)
			if err == nil {
				abs = a
			}
		}
		fi, err := os.Stat(abs)
		if err != nil {
			return nil, &os.PathError{Op: "stat", Path: abs, Err: err}
		}

		if fi.IsDir() {
			err = filepath.WalkDir(abs, func(wp string, d os.DirEntry, werr error) error {
				if werr != nil {
					return werr
				}
				if d.IsDir() {
					return nil
				}
				if matchExts(wp, allowAll, exts) {
					files = append(files, wp)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		} else {
			// 如果是具体文件路径，直接添加，忽略扩展名检查
			files = append(files, abs)
		}
	}
	return files, nil
}

// CollectFiles 收集输入路径中的目标文件（支持目录递归与单文件）
func CollectFiles(paths, exts []string) ([]string, error) {
	formatedExts, allowAll := normalizeExts(exts)
	files, err := collectFiles(paths, allowAll, formatedExts)
	files = UniqueSlice(files, false, true)
	return files, err
}
