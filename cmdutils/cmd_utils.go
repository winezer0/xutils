package cmdutils

import "strings"

func ContainsComma(s string) bool { return strings.Contains(s, ",") }

func SplitComma(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func FormatCmdsComma(raw []string) []string {
	var formated []string
	for _, selectedModel := range raw {
		// 按逗号分割多个模型定义
		if ContainsComma(selectedModel) {
			parts := SplitComma(selectedModel)
			formated = append(formated, parts...)
		}
	}
	return formated
}
