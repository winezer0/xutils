package csvutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestReadCSVToDicts 测试：一次性读取至字典列表
func TestReadCSVToDicts(t *testing.T) {
	csv := "" +
		"id,name,score\n" +
		"31,Iris,77\n" +
		"32,Jack,81\n"
	filePath := makeTempCSV(t, "dicts.csv", csv)

	headers, dicts, err := ReadCSVToDicts(filePath)
	if err != nil {
		t.Fatalf("ReadCSVToDicts failed: %v", err)
	}
	if len(headers) != 3 || headers[0] != "id" || headers[1] != "name" || headers[2] != "score" {
		t.Fatalf("headers = %v, want [id name score]", headers)
	}
	if len(dicts) != 2 {
		t.Fatalf("dicts count = %d, want 2", len(dicts))
	}
	if dicts[1]["name"] != "Jack" {
		t.Fatalf("second name = %s, want Jack", dicts[1]["name"])
	}
}

// TestReadCSVBytesToDicts 测试：从字节切片读取 CSV
func TestReadCSVBytesToDicts(t *testing.T) {
	csv := "" +
		"id,name,score\n" +
		"31,Iris,77\n" +
		"32,Jack,81\n"

	headers, dicts, err := ReadCSVBytesToDicts([]byte(csv))
	if err != nil {
		t.Fatalf("ReadCSVBytesToDicts failed: %v", err)
	}
	if len(headers) != 3 || headers[0] != "id" || headers[1] != "name" || headers[2] != "score" {
		t.Fatalf("headers = %v, want [id name score]", headers)
	}
	if len(dicts) != 2 {
		t.Fatalf("dicts count = %d, want 2", len(dicts))
	}
	if dicts[1]["name"] != "Jack" {
		t.Fatalf("second name = %s, want Jack", dicts[1]["name"])
	}
}

// TestReadCSVBytesToDictsEmpty 测试：空 CSV 字节
func TestReadCSVBytesToDictsEmpty(t *testing.T) {
	csv := ""

	_, _, err := ReadCSVBytesToDicts([]byte(csv))
	if err == nil {
		t.Error("期望错误，但没有返回错误")
	}
}

// TestReadCSVBytesToDictsOnlyHeader 测试：只有表头的 CSV
func TestReadCSVBytesToDictsOnlyHeader(t *testing.T) {
	csv := "" +
		"id,name,score\n"

	headers, dicts, err := ReadCSVBytesToDicts([]byte(csv))
	if err != nil {
		t.Fatalf("ReadCSVBytesToDicts failed: %v", err)
	}
	if len(headers) != 3 {
		t.Fatalf("headers count = %d, want 3", len(headers))
	}
	if len(dicts) != 0 {
		t.Fatalf("dicts count = %d, want 0", len(dicts))
	}
}

// TestReadCSVBytesToDictsDifferentDelimiter 测试：使用分号分隔符
func TestReadCSVBytesToDictsDifferentDelimiter(t *testing.T) {
	csv := "" +
		"id;name;score\n" +
		"31;Iris;77\n" +
		"32;Jack;81\n"

	headers, dicts, err := ReadCSVBytesToDicts([]byte(csv))
	if err != nil {
		t.Fatalf("ReadCSVBytesToDicts failed: %v", err)
	}
	if len(headers) != 3 || headers[0] != "id" || headers[1] != "name" || headers[2] != "score" {
		t.Fatalf("headers = %v, want [id name score]", headers)
	}
	if len(dicts) != 2 {
		t.Fatalf("dicts count = %d, want 2", len(dicts))
	}
}

// ExampleReadCSVBytesToDicts 使用示例：从字节切片读取 CSV
func ExampleReadCSVBytesToDicts() {
	csv := "" +
		"id,name,score\n" +
		"31,Iris,77\n"
	headers, _, _ := ReadCSVBytesToDicts([]byte(csv))
	fmt.Println("header:", headers[0], headers[1], headers[2])
	// Output: header: id name score
}

// ExampleReadCSVToDicts 使用示例：一次性读取并打印表头
func ExampleReadCSVToDicts() {
	csv := "" +
		"id,name,score\n" +
		"31,Iris,77\n"
	path := filepath.Join(os.TempDir(), "example_dicts.csv")
	_ = os.WriteFile(path, []byte(csv), 0644)
	headers, _, _ := ReadCSVToDicts(path)
	fmt.Println("header:", headers[0], headers[1], headers[2])
	// Output: header: id name score
}
