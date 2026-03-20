package csvutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestReadBigCSVToDictsWithCallback 测试：回调版读取（逐行处理）
func TestReadBigCSVToDictsWithCallback(t *testing.T) {
	csv := "" +
		"id,name,score\n" +
		"10,David,88\n" +
		"11,Eva,92\n"
	filePath := makeTempCSV(t, "callback.csv", csv)

	var got []map[string]interface{}
	err := ReadBigCSVToDictsWithCallback(filePath, ',', true, false, func(m map[string]interface{}) error {
		got = append(got, m)
		return nil
	})
	if err != nil {
		t.Fatalf("read with callback failed: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("rows count = %d, want 2", len(got))
	}
	if fmt.Sprintf("%v", got[0]["name"]) != "David" {
		t.Fatalf("first name = %v, want David", got[0]["name"])
	}
}

// ExampleReadBigCSVToDictsWithCallback 使用示例：回调处理并打印计数
func ExampleReadBigCSVToDictsWithCallback() {
	csv := "" +
		"id,name,score\n" +
		"10,David,88\n" +
		"11,Eva,92\n"
	path := filepath.Join(os.TempDir(), "example_cb.csv")
	_ = os.WriteFile(path, []byte(csv), 0644)
	count := 0
	_ = ReadBigCSVToDictsWithCallback(path, ',', true, false, func(m map[string]interface{}) error {
		count++
		return nil
	})
	fmt.Println("rows:", count)
	// Output: rows: 2
}
