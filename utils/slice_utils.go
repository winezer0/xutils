package utils

import (
	"sort"
	"strings"
	"unicode"
)

// SliceUnique 去除字符串切片中的重复项，支持忽略大小写和跳过空白字符串
// - ignoreCase: 是否忽略大小写（比较时转为小写）
// - skipEmpty: 是否跳过空或仅含空白字符的字符串
func SliceUnique(slice []string, ignoreCase, skipEmpty bool) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(slice))

	for _, item := range slice {
		// 处理 skipEmpty：跳过空或纯空白字符串
		if skipEmpty && strings.TrimSpace(item) == "" {
			continue
		}

		// 确定用于去重比较的键（key）
		key := item
		if ignoreCase {
			key = strings.ToLower(item)
		}

		// 如果未见过该 key，则保留原始 item（不是 key！）
		if !seen[key] {
			seen[key] = true
			out = append(out, item) // 保留原始大小写形式
		}
	}

	return out
}

// SliceSubtract 返回 a - b，即 a 中有但 b 中没有的元素（保持 a 的顺序）
func SliceSubtract(a, b []string) []string {
	// 将 b 转为 map 以实现 O(1) 查找
	bSet := make(map[string]struct{}, len(b))
	for _, item := range b {
		bSet[item] = struct{}{}
	}

	// 遍历 a，只保留不在 bSet 中的元素
	var result []string
	for _, item := range a {
		if _, exists := bSet[item]; !exists {
			result = append(result, item)
		}
	}

	return result
}

// SliceMerge 合并多个字符串切片，并自动去重（复用UniqueSlice的逻辑）
// - slices: 任意数量的待合并字符串切片（支持传入1个或多个切片）
// 返回值: 合并、去重、过滤后的唯一字符串切片
func SliceMerge(slices ...[]string) []string {
	// 1. 合并所有输入切片为一个临时切片
	merged := make([]string, 0)
	for _, s := range slices {
		if s == nil { // 容错：跳过nil切片，避免panic
			continue
		}
		merged = append(merged, s...)
	}

	return merged
}

// SliceSort 对字符串切片进行默认排序 (A-Z, a-z)
func SliceSort(input []string) []string {
	// 复制 防止修改原切片
	result := make([]string, len(input))
	copy(result, input)

	sort.Strings(result)
	return result
}

// SliceSortAsc 对字符串切片进行标准字典序排序，支持升序/降序
// 不修改原切片，返回新切片
// asc: true 升序，false 降序
func SliceSortAsc[T ~string](input []T, asc bool) []T {
	// 复制切片，防止修改原数据
	result := make([]T, len(input))
	copy(result, input)

	// 标准排序
	sort.Slice(result, func(i, j int) bool {
		if asc {
			return result[i] < result[j]
		} else {
			return result[i] > result[j]
		}
	})

	return result
}

// SliceEqualStrict 辅助函数：严格比较字符串切片（顺序 + 内容完全一致）
// 为了兼容性保留，内部可调用 SliceEqual(a, b, false)
func SliceEqualStrict(a, b []string) bool {
	return SliceEqual(a, b, false)
}

// SliceEqual 辅助函数：比较字符串切片
// 参数 sorted:
//   - true: 先排序再比较，适用于文件列表、标签集合等。
//   - false: 严格匹配顺序，但使用高效的拼接比较法。
func SliceEqual(a, b []string, sorted bool) bool {
	// 1. 长度不同，直接不相等
	if len(a) != len(b) {
		return false
	}

	// 2. 如果都为空，相等
	if len(a) == 0 {
		return true
	}

	// 3. 准备切片（如果需要排序，则创建副本；如果不需要排序且原切片可直接用，则优化处理）
	var sliceA, sliceB []string

	if sorted {
		// 对副本进行排序
		sliceA = SliceSort(a)
		sliceB = SliceSort(b)
	} else {
		// 不需要排序，直接使用原切片引用（因为不会修改它们）
		sliceA = a
		sliceB = b
	}

	// 4. 使用特殊分隔符拼接并直接比较
	// 使用 \x00 (Null Byte) 作为分隔符，防止 ["a", "bc"] 和 ["ab", "c"] 这种边界情况误判
	separator := "\x00"
	strA := strings.Join(sliceA, separator)
	strB := strings.Join(sliceB, separator)
	return strA == strB
}

// SliceAppend 尾部追加元素（保持原有顺序）
// 不修改原切片，返回新切片
func SliceAppend[T any](slice []T, elem T) []T {
	newSlice := make([]T, len(slice), len(slice)+1)
	copy(newSlice, slice)
	return append(newSlice, elem)
}

// SliceRemoveFirst 删除【第一个】匹配的元素（保持顺序）
// 不修改原切片，返回新切片
func SliceRemoveFirst[T comparable](slice []T, target T) []T {
	for i, v := range slice {
		if v == target {
			// 构造新切片，跳过该元素
			newSlice := make([]T, len(slice)-1)
			copy(newSlice[:i], slice[:i])
			copy(newSlice[i:], slice[i+1:])
			return newSlice
		}
	}
	// 没找到，返回原切片副本
	return slice
}

// SliceRemoveAll 删除【所有】匹配的元素（保持顺序）
// 不修改原切片，返回新切片
func SliceRemoveAll[T comparable](slice []T, target T) []T {
	newSlice := make([]T, 0, len(slice))
	for _, v := range slice {
		if v != target {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

// SliceSortNatural 对字符串切片执行【自然排序】，仅排序、不去重
// slice : 待排序的字符串切片（如页码列表、文件名列表）
// asc   : 排序方向；true=升序，false=降序
// return: 排序后的新切片，不会修改原切片
func SliceSortNatural(slice []string, asc bool) []string {
	// 复制切片，不修改原数据
	result := make([]string, len(slice))
	copy(result, slice)

	sort.Slice(result, func(i, j int) bool {
		if asc {
			return NaturalLess(result[i], result[j])
		} else {
			return NaturalLess(result[j], result[i])
		}
	})

	return result
}

// NaturalLess 比较两个字符串，按自然顺序（数字按数值比）
func NaturalLess(a, b string) bool {
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
