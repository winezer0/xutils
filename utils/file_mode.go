package utils

import (
	"os"
	"strings"
)

// ParseFlagFromMode 将 POSIX 风格的文件打开模式字符串转换为 Go os 包对应的 flag 常量
// 支持的模式（兼容常见文件操作场景）：
//
//	r   - 只读模式（文件必须存在）
//	r+  - 读写模式（文件必须存在）
//	w   - 只写模式（文件不存在则创建，存在则清空）
//	w+  - 读写模式（文件不存在则创建，存在则清空）
//	a   - 追加模式（只写，文件不存在则创建，写入从末尾开始）
//	a+  - 追加模式（读写，文件不存在则创建，写入从末尾开始，读取从开头）
//
// 返回值：
//
//	flag - 对应的 os 包 flag 组合（如 os.O_RDONLY | os.O_CREATE）
func ParseFlagFromMode(mode string) int {
	// 统一转为小写，兼容大小写混合输入（如 "A+"、"W"）
	mode = strings.TrimSpace(strings.ToLower(mode))

	// 定义模式与 flag 的映射关系
	var flag int

	switch mode {
	case "r":
		// 只读：文件必须存在，否则 os.Open 会报错
		flag = os.O_RDONLY
	case "r+":
		// 读写：文件必须存在，否则 os.OpenFile 会报错
		flag = os.O_RDWR
	case "w":
		// 只写：创建（不存在）+ 截断（存在）+ 只写
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	case "w+":
		// 读写：创建（不存在）+ 截断（存在）+ 读写
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	case "a":
		// 追加（只写）：创建（不存在）+ 追加 + 只写
		flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	case "a+":
		// 追加（读写）：创建（不存在）+ 追加 + 读写
		flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	default:
		flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	}

	return flag
}

// ParseFlagFromOver 将 Overwrite 风格的文件打开模式字符串转换为 Go os 包对应的 flag 常量
func ParseFlagFromOver(overwrite bool) int {
	if overwrite {
		return ParseFlagFromMode("w+")
	} else {
		return ParseFlagFromMode("a+")
	}
}
