package utils

import (
	"fmt"
	"strings"
)

// ErrorsToError 将错误切片转换为单个 error
// 通过字符串拼接将所有错误信息合并
func ErrorsToError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	// 如果只有一个错误，直接返回，保持原始错误类型
	if len(errs) == 1 {
		return errs[0]
	}

	// 多个错误则拼接成一个大字符串
	var sb strings.Builder
	for i, err := range errs {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, err.Error()))
	}

	// 返回一个标准的 error 接口
	return fmt.Errorf("%s", sb.String())
}
