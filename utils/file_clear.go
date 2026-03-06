package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// ClearPaths 根据需要清除文件或目录：
// - 对于文件：仅删除（不重新构建）
// - 对于目录：首先删除，如果 rebuild 为 true 则重新构建
// paths：要处理的路径列表（文件/目录）
// rebuild：在删除操作后是否重新构建目录（仅适用于目录）
// perm：重新构建目录时的权限设置（使用 0 可保留原始权限）
// return：处理过程中出现的汇总错误
func ClearPaths(paths []string, rebuild bool, perm os.FileMode) (removed []string, errors []error) {
	errors = make([]error, 0, len(paths))
	removed = make([]string, 0, len(paths))

	for _, path := range paths {
		cleanPath := filepath.Clean(path)
		// 获取路径信息以确定类型（文件/目录）
		info, err := os.Stat(cleanPath)

		// Case 1：路径不存在
		if os.IsNotExist(err) {
			errMsg := fmt.Errorf("path:%s does not exist", cleanPath)
			errors = append(errors, errMsg)
			continue
		}

		// Case 2: 未能获取路径信息（出现其他错误）
		if err != nil {
			errMsg := fmt.Errorf("Failed to get info for path %s: %v; ", cleanPath, err)
			errors = append(errors, errMsg)
			continue
		}

		// Case 3: 路径为 文件 - 仅删除，不重建
		if !info.IsDir() {
			if err := os.Remove(cleanPath); err != nil {
				errMsg := fmt.Errorf("Failed to delete file %s: %v; ", cleanPath, err)
				errors = append(errors, errMsg)
			} else {
				removed = append(removed, cleanPath)
			}
			continue
		}

		// Case 4: 路径是一个目录 - 通过正确的权限逻辑进行删除/重建操作
		// Step 1: 获取原始目录权限（在删除之前）
		originalPerm := info.Mode().Perm()
		// Step 2: 确定最终的重建许可
		var finalPerm os.FileMode
		if perm == 0 {
			// 若 perm 的值为 0（即为设计要求），则使用原始权限。
			finalPerm = originalPerm
		} else {
			// 如果 perm 不等于 0 ，则使用指定的权限。
			finalPerm = perm
		}

		// Step 3: 删除整个目录（包括其中的所有内容）
		if err := os.RemoveAll(cleanPath); err != nil {
			errMsg := fmt.Errorf("Failed to delete dir %s: %v; ", cleanPath, err)
			errors = append(errors, errMsg)
			continue
		} else {
			removed = append(removed, cleanPath)
		}

		// Step 4: 如需则重新构建目录（仅适用于目录）
		if rebuild {
			// 如果原始目录不存在（极端情况），则使用 0755 作为备用设置。
			if finalPerm == 0 {
				finalPerm = 0755
			}
			if err := os.MkdirAll(cleanPath, finalPerm); err != nil {
				errMsg := fmt.Errorf("Failed to rebuild directory %s (permission: %o): %v; ", cleanPath, finalPerm, err)
				errors = append(errors, errMsg)
			} else {
				fmt.Printf("Directory %s has been rebuilt (permission: %o)\n", cleanPath, finalPerm)
			}
		}
	}

	return removed, errors
}

// ClearFiles 批量删除指定的文件。
//
// 参数:
//   - paths: 可变参数，一个或多个需要删除的文件完整路径。
//
// 返回:
//   - error: 如果在删除过程中遇到错误（如权限不足、文件不存在等），返回汇总错误；若全部成功则返回 nil。
func ClearFiles(paths ...string) (removed []string, errors []error) {
	if len(paths) == 0 {
		return nil, nil
	}
	filesExist, _ := CheckFilesExist(paths)
	return ClearPaths(filesExist, false, 0)
}
