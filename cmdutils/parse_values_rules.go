package cmdutils

import (
	"strings"
)

// ParseValueRule 解析单条原始规则字符串，返回去重去空后的字符串列表
// 支持格式：
//  1. list:aaa,bbb,ccc 按逗号分割（无前缀默认走该逻辑）
//  2. file:1.txt 读取文件非空行
//  3. aaa,bbb 无前缀，默认按 list 处理
//
// 返回：去重、去空、去空白后的值列表、错误
func ParseValueRule(rawRule string) ([]string, error) {
	var values []string
	var err error

	// 根据前缀判断处理模式
	switch {
	case strings.HasPrefix(rawRule, "file:"):
		// 文件模式：读取文件所有非空行
		path := strings.TrimPrefix(rawRule, "file:")
		values, err = readNonEmptyLines(path)
	default:
		// 列表模式：处理 list: 前缀 或 无前缀
		content := strings.TrimPrefix(rawRule, "list:")
		values = splitAndTrim(content)
	}

	if err != nil {
		return nil, err
	}

	// 单条结果内部去重
	return deduplicate(values), nil
}

// ParseValueRules 批量解析多条规则，返回全局去重、去空的结果列表
// 入参：多条规则字符串切片
// 返回：所有规则解析后合并、去重、去空的值列表、错误
func ParseValueRules(rawRules []string) ([]string, error) {
	var allValues []string

	// 遍历解析每一条规则
	for _, rule := range rawRules {
		vals, err := ParseValueRule(rule)
		if err != nil {
			return nil, err
		}
		allValues = append(allValues, vals...)
	}

	// 最终全局去重
	return deduplicate(allValues), nil
}
