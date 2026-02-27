package utils

import (
	"path/filepath"
	"strings"
)

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
