package utils

import (
	"sort"
)

// GetMapKeys 泛型辅助函数：从任意 value 类型的 map[string]V 提取所有 key
// 类型参数:
// - V: map 中 value 的类型
// 参数:
// - m: 输入的映射
// 返回值:
// - []string: 映射中所有的键组成的切片
func GetMapKeys[V any](m map[string]V) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// GetMapSortedKeys 返回给定 map 的所有字符串键，并按指定顺序排序。
// 参数 m 是任意值类型的 map，其键必须为 string 类型。
// 参数 ascending 控制排序方向：true 表示升序（字典序），false 表示降序。
// 返回值是一个包含排序后键的字符串切片。
func GetMapSortedKeys[V any](m map[string]V, ascending bool) []string {
	keys := GetMapKeys(m)
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
	mapKeys := GetMapKeys(m)
	return SliceSortNatural(mapKeys, asc)
}
