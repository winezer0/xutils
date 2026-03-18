package utils

import (
	"sort"
	"unicode"
)

// GetMapKeys 从 map[string]bool 提取所有 key
// 参数:
// - m: 输入的映射
// - sorted: 是否进行排序
// 返回值:
// - []string: 映射中所有的键组成的切片
func GetMapKeys(m map[string][]string, sorted bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	if sorted {
		sort.Strings(keys)
	}

	return keys
}

// GetMapSortedKeys 返回给定 map 的所有字符串键，并按指定顺序排序。
// 参数 m 是任意值类型的 map，其键必须为 string 类型。
// 参数 ascending 控制排序方向：true 表示升序（字典序），false 表示降序。
// 返回值是一个包含排序后键的字符串切片。
func GetMapSortedKeys[V any](m map[string]V, ascending bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	if ascending {
		sort.Strings(keys)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}

	return keys
}

// GetMapKeysWithNaturalSorted 返回给定 map 的所有字符串键，并按指定顺序排序。
// 参数 m 是任意值类型的 map，其键必须为 string 类型。
// 参数 asc 控制排序方向：true 表示升序（字典序），false 表示降序。
// 返回值是一个包含排序后键的字符串切片。
func GetMapKeysWithNaturalSorted[V any](m map[string]V, asc bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	if asc {
		sort.Slice(keys, func(i, j int) bool {
			return naturalLess(keys[i], keys[j])
		})
	} else {
		sort.Slice(keys, func(i, j int) bool {
			return naturalLess(keys[j], keys[i]) // 注意：反转比较
		})
	}

	return keys
}

// naturalLess 比较两个字符串，按自然顺序（数字按数值比）
func naturalLess(a, b string) bool {
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		// 如果当前位置都是数字，提取完整数字进行数值比较
		if unicode.IsDigit(rune(a[i])) && unicode.IsDigit(rune(b[j])) {
			// 提取 a 中的数字
			startA := i
			for i < len(a) && unicode.IsDigit(rune(a[i])) {
				i++
			}
			numA := a[startA:i]

			// 提取 b 中的数字
			startB := j
			for j < len(b) && unicode.IsDigit(rune(b[j])) {
				j++
			}
			numB := b[startB:j]

			// 去除前导零（可选，但推荐）
			// 比较长度：更长的数字更大（避免 "001" vs "1"）
			if len(numA) != len(numB) {
				return len(numA) < len(numB)
			}
			// 长度相同，字典序即数值序
			if numA != numB {
				return numA < numB
			}
		} else {
			// 非数字部分，逐字符比较
			if a[i] != b[j] {
				return a[i] < b[j]
			}
			i++
			j++
		}
	}
	// 若前面都相等，短的排前面
	return len(a) < len(b)
}
