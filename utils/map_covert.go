package utils

import (
	"strconv"
)

// GetMapString 从 map 中获取字符串字段，支持多别名匹配
// 入参：
//
//	m - 源map（key为字符串，value为任意类型）
//	keys - 字段名列表（按顺序匹配，优先返回第一个存在且为字符串的字段值）
//
// 返回：匹配到的字符串，无有效值时返回空字符串
func GetMapString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			if s, ok2 := v.(string); ok2 {
				return s
			}
		}
	}
	return ""
}

// GetInt 从 map 中获取整数类型值，仅支持 int 类（int/float64截断/字符串整数）
// 入参：
//
//	m - 源map
//	keys - 字段名列表（支持多别名，按顺序匹配）
//
// 返回：匹配到的整数值，无有效值时返回 0
func GetInt(m map[string]interface{}, keys ...string) int {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			switch tv := v.(type) {
			case int:
				return tv // 原生int类型，直接返回
			case float64:
				return int(tv) // float64截断为int（如9.9→9）
			case string:
				// 仅解析整数格式字符串，非整数返回0
				if n, err := strconv.Atoi(tv); err == nil {
					return n
				}
			}
		}
	}
	return -1
}

// GetMapBool 从 map 中获取布尔类型字段，支持多别名匹配
// 入参：
//
//	m - 源map（key为字符串，value为任意类型）
//	keys - 字段名列表（按顺序匹配，优先返回第一个存在且为布尔的字段值）
//
// 返回：匹配到的布尔值，无有效值时返回 false
func GetMapBool(m map[string]interface{}, keys ...string) bool {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			if b, ok2 := v.(bool); ok2 {
				return b
			}
		}
	}
	return false
}

// GetMapStringSlice 从 map 中获取字符串切片类型字段，支持多别名匹配
// 入参：
//
//	m - 源map（key为字符串，value为任意类型）
//	keys - 字段名列表（按顺序匹配，优先返回第一个存在且为字符串切片的字段值）
//
// 返回：匹配到的字符串切片，无有效值时返回 nil
func GetMapStringSlice(m map[string]interface{}, keys ...string) []string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			// 直接匹配 []string 类型
			if ss, ok2 := v.([]string); ok2 {
				return ss
			}
			// 匹配 []interface{} 类型并转换为 []string
			if ivs, ok2 := v.([]interface{}); ok2 {
				ss := make([]string, 0, len(ivs))
				for _, iv := range ivs {
					if s, ok3 := iv.(string); ok3 {
						ss = append(ss, s)
					}
				}
				return ss
			}
		}
	}
	return nil
}
