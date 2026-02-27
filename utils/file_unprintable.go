package utils

// 自定义IsPrint函数，判断字符是否为可打印字符
func isPrint(r rune) bool {
	// 空格(32)是可打印字符
	if r == 32 {
		return true
	}
	// ASCII可打印字符(33-126)
	if r >= 33 && r <= 126 {
		return true
	}
	// Unicode可打印字符(排除控制字符)
	// 控制字符范围：U+0000-U+001F、U+007F-U+009F、U+2000-U+200F等
	if r <= 0x001F || (r >= 0x007F && r <= 0x009F) {
		return false
	}
	// 过滤零宽空格(U+200B)等不可见控制字符
	if r >= 0x2000 && r <= 0x200F {
		return false
	}
	// 其他Unicode字符(字母、数字、符号等)视为可打印
	return true
}

// CleanupUnprintableChars 移除字符串中的不可打印字符
func CleanupUnprintableChars(line string) string {
	// 预分配切片，避免频繁扩容，提升性能
	filtered := make([]rune, 0, len(line))

	for _, r := range line {
		if isPrint(r) {
			filtered = append(filtered, r)
		}
	}
	return string(filtered)
}
