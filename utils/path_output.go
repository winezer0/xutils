package utils

import (
	"fmt"
	"github.com/winezer0/xutils/hashutils"
	"os"
	"path/filepath"
	"strings"
)

// BuildProjectFilePath 生成项目缓存文件名或完整路径
// 参数:
//   - projectName: 项目名称 (作为文件名前缀)
//   - projectPath: 项目路径 (用于生成 Hash 以区分不同路径的同名项目)
//   - toolName: 工具标识
//   - suffix: 文件后缀 (默认 "cache")
//   - usePwd:
//     true:  返回当前工作目录下的完整绝对路径 (例如: "/cwd/myproj.a1b2.tool.cache")
//     false: 仅返回生成的文件名 (例如: "myproj.a1b2.tool.cache")
//   - hidden: 是否隐藏文件。
//   - - true: 生成的文件名前添加<.>
//
// 返回:
//   - string: 文件名或完整路径
//   - error: 当 usePwd=true 且无法获取当前工作目录时返回错误
func BuildProjectFilePath(projectName, projectPath, toolName, suffix string, usePwd, hidden bool) string {
	// 1. 默认值处理
	if projectName == "" {
		projectName = "default"
	}

	// 清理输入中的多余点号，防止格式错乱
	if suffix == "" {
		suffix = "cache"
	} else {
		suffix = strings.TrimLeft(suffix, ".")
	}

	// 3. 生成 Hash (保持原有逻辑)
	// 假设 GetStrHashShort 返回的是短字符串，如 "a1b2c3"
	pathHash := hashutils.GetStrHashShort(projectPath)
	if pathHash == "" {
		pathHash = "none_hash"
	}

	// 4. 构建基础文件名: projectName.hash.tool.suffix
	filename := fmt.Sprintf("%s.%s.%s.%s", projectName, pathHash, toolName, suffix)
	// 在文件名前添加dot
	if hidden {
		filename = fmt.Sprintf(".%s", filename)
	}

	// 5. 根据 usePwd 决定是否拼接当前路径
	if usePwd {
		cwd, err := os.Getwd()
		if err != nil {
			return filename
		}
		return filepath.Join(cwd, filename)
	}

	// 如果不需要绝对路径，直接返回文件名
	return filename
}

// BuildOutFilePath 生成文件路径
// 参数说明:
//   - filePath: 原始文件路径 (例如: /data/urls.txt 或 ./logs/app.log)
//   - toolName: 工具名称标识 (例如: "scanner")
//   - ext: 扩展名 (例如: "cache")
//   - usePwd: 是否强制使用当前命令行运行路径 (CWD)。
//   - - true: 缓存文件生成在当前执行命令的目录下，文件名基于 filePath 的 basename。
//   - - false: 缓存文件生成在 filePath 所在的同级目录下。
//   - hidden: 是否隐藏文件。
//   - - true: 生成的文件名前添加<.>
//
// 返回: 完整的缓存文件绝对路径或相对路径
func BuildOutFilePath(filePath, toolName, ext string, usePwd, hidden bool) string {
	var dir string
	var err error

	// 清理 ext 开头可能存在的点号，防止出现双点 ".."
	if len(ext) == 0 {
		ext = "cache"
	} else {
		ext = strings.TrimLeft(ext, ".")
	}

	// 构建新文件名：<原文件名>.<tool>.<ext>
	basename := filepath.Base(filePath)
	filename := fmt.Sprintf("%s.%s.%s", basename, toolName, ext)

	// 在文件名前添加dot
	if hidden {
		filename = fmt.Sprintf(".%s", filename)
	}

	if usePwd {
		// 获取当前工作目录 (Current Working Directory)
		dir, err = os.Getwd()
		if err != nil {
			return filename
		}
	} else {
		// 使用原文件所在的目录
		dir = filepath.Dir(filePath)
	}

	return filepath.Join(dir, filename)
}

// BuildAvoidConflictName 生成避免冲突的文件名 若文件存在则追加 -N 后缀避免冲突(中文注释)
func BuildAvoidConflictName(path string) string {
	if _, err := os.Stat(path); err != nil {
		return path
	}
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	name := base
	ext := ""
	if i := strings.LastIndex(base, "."); i >= 0 {
		name = base[:i]
		ext = base[i:]
	}
	for i := 1; i < 1000; i++ {
		cand := filepath.Join(dir, fmt.Sprintf("%s-%d%s", name, i, ext))
		if _, err := os.Stat(cand); os.IsNotExist(err) {
			return cand
		}
	}
	return path
}
