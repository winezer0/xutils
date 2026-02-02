package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ToRelativePath 将给定的文件路径（filePath）转换为相对于项目根路径（projectPath）的相对路径。
// 路径分隔符统一为 '/'，确保跨平台一致性。
//
// 行为说明：
//   - 若 filePath 为空，返回空字符串。
//   - 若 filePath 已是相对路径，则直接标准化并返回（不依赖 projectPath）。
//   - 若 filePath 是绝对路径且位于 projectPath 目录下，则返回其相对路径（如 "src/main.go"）。
//   - 若 filePath 不在 projectPath 内（如跨盘符或上级目录），或 projectPath 无效，
//     则 fallback 返回标准化后的绝对路径（仍使用 '/' 分隔）。
//   - 所有输出路径均不包含前导 "./"。
//
// 注意：本函数不会修改原始输入，所有路径均为值拷贝处理。
func ToRelativePath(filePath, projectPath string) (string, error) {
	if filePath == "" {
		return "", nil
	}

	filePath = filepath.Clean(filePath)
	projectPath = filepath.Clean(projectPath)

	// 如果是相对路径，直接标准化
	if !filepath.IsAbs(filePath) {
		return ToSlashPath(filePath), nil
	}

	if projectPath == "" {
		return ToSlashPath(filePath), nil
	}

	// 确保 projectPath 是绝对路径
	if !filepath.IsAbs(projectPath) {
		var err error
		projectPath, err = filepath.Abs(projectPath)
		if err != nil {
			return ToSlashPath(filePath), nil
		}
	}

	relPath, err := filepath.Rel(projectPath, filePath)
	if err != nil {
		return ToSlashPath(filePath), nil // fallback
	}

	// 防止返回项目外路径（如 ../../xxx）
	if strings.HasPrefix(relPath, ".."+string(filepath.Separator)) ||
		strings.HasPrefix(relPath, "../") {
		return ToSlashPath(filePath), nil
	}

	return ToSlashPath(relPath), nil
}

// ToAbsolutePath 将给定的文件路径（filePath）转换为绝对路径。
// 路径分隔符统一为 '/'，确保跨平台一致性。
//
// 行为说明：
// - 若 filePath 为空，返回空字符串。
// - 若 filePath 已是绝对路径，则标准化后返回。
// - 若 filePath 是相对路径，则基于 projectPath 拼接成绝对路径。
// - 若 projectPath 为空，则使用当前工作目录（os.Getwd()）作为基准。
// - 若 projectPath 是相对路径，会先转换为绝对路径。
//
// 错误情况：
// - 无法获取当前工作目录（仅当 projectPath 为空且 filePath 为相对路径时）
// - 路径拼接后无法解析为绝对路径（极少见）
//
// 返回的路径始终使用 '/' 分隔，无前导 "./"。
func ToAbsolutePath(filePath, projectPath string) (string, error) {
	if filePath == "" {
		return "", nil
	}

	if filepath.IsAbs(filePath) {
		abs, err := filepath.Abs(filePath)
		if err != nil {
			return "", err
		}
		return ToSlashPath(abs), nil
	}

	if projectPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		projectPath = cwd
	}

	if !filepath.IsAbs(projectPath) {
		var err error
		projectPath, err = filepath.Abs(projectPath)
		if err != nil {
			return "", err
		}
	}

	absPath := filepath.Join(projectPath, filePath)
	absPath, err := filepath.Abs(absPath)
	if err != nil {
		return "", err
	}

	return ToSlashPath(absPath), nil
}

// ToSlashPath 将任意路径中的分隔符统一转换为 '/'
// 同时清理冗余符号（如 ./, ../ 保留，但多余分隔符会被清理）
func ToSlashPath(p string) string {
	if p == "" {
		return ""
	}
	// 先用 filepath.Clean 清理路径（不解析符号链接）
	p = filepath.Clean(p)
	// 转换所有 \ 或 / 为 /
	p = filepath.ToSlash(p)
	// 确保不以 ./ 开头（可选，保持简洁）
	p = strings.TrimPrefix(p, "./")
	return p
}
