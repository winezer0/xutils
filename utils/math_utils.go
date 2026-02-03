package utils

// MaxNum 返回两个整数中的较大值
func MaxNum(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinNum 返回两个整数中的较小值
func MinNum(a, b int) int {
	if a < b {
		return a
	}
	return b
}
