package csvutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestCSVIterator_ReadAll 测试：迭代器逐行读取（含缺失字段容错）
func TestCSVIterator_ReadAll(t *testing.T) {
	// 准备CSV（第三行缺少score列，用于验证容错）
	csv := "" +
		"id,name,score\n" +
		"1,Alice,90\n" +
		"2,Bob,85\n" +
		"3,Carol\n"
	filePath := makeTempCSV(t, "iter.csv", csv)

	iter, err := NewCSVIterator(filePath, ',', false)
	if err != nil {
		t.Fatalf("new iterator failed: %v", err)
	}
	defer iter.Close()

	var rows []map[string]interface{}
	for {
		d := iter.Next()
		if d == nil {
			break
		}
		rows = append(rows, d)
	}
	if err := iter.Error(); err != nil && err.Error() != "EOF" {
		t.Fatalf("iterator error: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows count = %d, want 3", len(rows))
	}
	// 验证键存在（容错场景下score可能为空字符串）
	last := rows[2]
	if _, ok := last["id"]; !ok {
		t.Fatalf("missing key 'id'")
	}
	if _, ok := last["name"]; !ok {
		t.Fatalf("missing key 'name'")
	}
	if _, ok := last["score"]; !ok {
		t.Fatalf("missing key 'score'")
	}
}

// ExampleCSVIterator 使用示例：迭代器读取并打印行数
func ExampleCSVIterator() {
	csv := "" +
		"id,name,score\n" +
		"1,Alice,90\n" +
		"2,Bob,85\n"
	path := filepath.Join(os.TempDir(), "example_iter.csv")
	_ = os.WriteFile(path, []byte(csv), 0644)
	iter, _ := NewCSVIterator(path, ',', false)
	defer iter.Close()
	count := 0
	for {
		m := iter.Next()
		if m == nil {
			break
		}
		count++
	}
	fmt.Println("rows:", count)
	// Output: rows: 2
}
