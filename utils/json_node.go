package utils

import "strings"

// JSONValueByKeyPath 根据指定的键路径（keyPath）从嵌套的 JSON 结构（map[string]interface{}）中获取对应的值
//
// 功能说明：
// 该函数用于解析嵌套的 JSON 结构（以 map[string]interface{} 为基础的 Go 数据结构），
// 通过点分隔的键路径（如 "user.info.name"）逐层查找对应的值，支持根路径（"" 或 "."）直接返回根节点。
// 如果路径中任意一级不存在、或对应层级不是 map[string]interface{} 类型，均返回 nil。
//
// 参数：
//
//	root     - 根节点，通常是解析 JSON 后得到的 map[string]interface{} 类型数据
//	keyPath  - 点分隔的键路径，例如 "a.b.c"；空字符串 "" 或 "." 表示根节点
//	linkSymbol  - 分隔符 默认.
//
// 返回值：
//
//	找到对应路径的值则返回该值（interface{} 类型，需自行类型断言）；
//	// 无效路径返回 nil
//	// 根路径返回整个结构
func JSONValueByKeyPath(root interface{}, keyPath, linkSymbol string) interface{} {
	if len(linkSymbol) == 0 {
		linkSymbol = "."
	}

	if keyPath == "" || keyPath == linkSymbol {
		return root
	}

	parts := strings.Split(keyPath, linkSymbol)
	cur := root
	for _, p := range parts {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil
		}
		v, ok := m[p]
		if !ok {
			return nil
		}
		cur = v
	}
	return cur
}
