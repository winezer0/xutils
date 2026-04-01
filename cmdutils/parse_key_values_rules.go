package cmdutils

import (
	"fmt"
	"strings"
)

type KeyValuesRule struct {
	Key    string
	Values map[string]bool
}

// ParseKeyValuesRules 批量解析 key=value 规则
// 优化点：相同 key 自动合并所有 values（去重），不再忽略
func ParseKeyValuesRules(rawRules []string) ([]KeyValuesRule, error) {
	// 使用 map 暂存：key => 合并后的 valueMap
	ruleMap := make(map[string]map[string]bool)

	for _, raw := range rawRules {
		// 解析单条规则
		rule, err := ParseKeyValuesRule(raw)
		if err != nil {
			return nil, err
		}

		// 核心优化：相同 key 自动合并
		if existingMap, exists := ruleMap[rule.Key]; exists {
			// 把新值合并到已有的 map 中（自动去重）
			for v := range rule.Values {
				existingMap[v] = true
			}
		} else {
			// 第一次出现，直接存入
			ruleMap[rule.Key] = rule.Values
		}
	}

	// 转换为最终的 []KeyValuesRule 返回
	var rules []KeyValuesRule
	for key, values := range ruleMap {
		rules = append(rules, KeyValuesRule{
			Key:    key,
			Values: values,
		})
	}

	return rules, nil
}

// ParseKeyValuesRule 解析单条 key=value 规则（不变）
func ParseKeyValuesRule(raw string) (KeyValuesRule, error) {
	parts := strings.SplitN(raw, "=", 2)
	if len(parts) != 2 {
		return KeyValuesRule{}, fmt.Errorf("invalid rule format: %s", raw)
	}

	key := parts[0]
	valStr := parts[1]

	values, err := ParseValueRule(valStr)
	if err != nil {
		return KeyValuesRule{}, fmt.Errorf("error loading values for key %s: %v", key, err)
	}

	// 转成 map 用于快速查找 & 去重
	valueMap := make(map[string]bool)
	for _, v := range values {
		valueMap[v] = true
	}

	return KeyValuesRule{
		Key:    key,
		Values: valueMap,
	}, nil
}
