package utils

import "strings"

// SanitizeFilename 清理文件名中的非法字符，确保生成的文件名在大多数文件系统中是安全的。
// 它会将常见的非法字符（如路径分隔符、空格、冒号）替换为下划线，并处理首尾空白及末尾的点号。
// 如果清理后的文件名为空，则返回提供的默认名称 defaultName。
//
// 参数:
//   - basename: 原始文件名或字符串。
//   - defaultName: 当清理后的文件名为空时使用的备用名称。
//
// 返回:
//   - 清理后的安全文件名字符串。
func SanitizeFilename(basename string, defaultName string) string {
	//strings.NewReplacer 是 Go 标准库中用于批量替换字符串中多个子串的高效工具。 参数：成对出现的字符串，old1, new1, old2, new2, ...
	s := "_"
	r := strings.NewReplacer("/", s, "\\", s, " ", s, ":", s)
	basename = r.Replace(basename)
	basename = strings.TrimSpace(basename)
	basename = strings.TrimRight(basename, ". ")
	if basename == "" {
		basename = defaultName
	}
	return basename
}
