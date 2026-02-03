package utils

import (
	"encoding/json"
	"testing"
)

type DeduplicationTestTarget struct {
	Items []DeduplicationTestItem `json:"items"`
}

type DeduplicationTestItem struct {
	Known string `json:"known"`
}

func TestDetectUnknownFields_Deduplication(t *testing.T) {
	// 构造包含20个元素的数组，每个元素都有unknown字段
	var items []map[string]interface{}
	for i := 0; i < 20; i++ {
		item := map[string]interface{}{
			"known":   "value",
			"unknown": "value",
		}
		// 在第5个元素添加 unique_at_5 (索引4)
		if i == 4 {
			item["unique_at_5"] = "val"
		}
		// 在第15个元素添加 unique_at_15 (索引14)
		// 由于限制了只检查前10个元素，这个字段应该不会被检测到
		if i == 14 {
			item["unique_at_15"] = "val"
		}
		items = append(items, item)
	}

	data := map[string]interface{}{
		"items": items,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 模拟实际使用场景：先将JSON反序列化到目标结构体
	var target DeduplicationTestTarget
	if err := json.Unmarshal(jsonBytes, &target); err != nil {
		t.Fatalf("Unmarshal to target failed: %v", err)
	}

	// 调用 DetectUnknownFields
	unknowns := DetectUnknownFields(jsonBytes, &target)

	// 验证结果去重
	seen := make(map[string]bool)
	for _, u := range unknowns {
		if seen[u] {
			t.Errorf("发现重复字段: %s", u)
		}
		seen[u] = true
	}

	// 验证包含期望的字段
	if !seen["items[*].unknown"] {
		t.Errorf("期望包含 items[*].unknown，但未找到。实际结果: %v", unknowns)
	}
	if !seen["items[*].unique_at_5"] {
		t.Errorf("期望包含 items[*].unique_at_5，但未找到。实际结果: %v", unknowns)
	}

	// 验证不包含被忽略的字段 (索引 >= 10)
	if seen["items[*].unique_at_15"] {
		t.Errorf("不应包含 items[*].unique_at_15 (应被 limit=10 忽略)，但找到了。实际结果: %v", unknowns)
	}
}
