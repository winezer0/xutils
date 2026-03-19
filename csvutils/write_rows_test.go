package csvutils

import (
	"encoding/csv"
	"os"
	"testing"
)

// 读取整个文件用于校验
func readAllRows(t *testing.T, filename string, delimiter rune) [][]string {
	t.Helper()
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = delimiter
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// 1. 覆盖模式：正常写入表头+数据
func TestWrite_Overwrite_Normal(t *testing.T) {
	cleanup()
	header := []string{"id", "name", "age"}
	rows := [][]string{{"1", "Alice", "20"}}

	err := WriteRowsToCSV(testFile, header, rows, ',', true)
	if err != nil {
		t.Fatal(err)
	}

	all := readAllRows(t, testFile, ',')
	if len(all) != 2 {
		t.Errorf("want 2 rows, got %d", len(all))
	}
}

// 2. 追加模式：不重复写入表头
func TestWrite_Append_NoDuplicateHeader(t *testing.T) {
	cleanup()
	header := []string{"id", "name"}
	rows1 := [][]string{{"1", "A"}}
	rows2 := [][]string{{"2", "B"}}

	_ = WriteRowsToCSV(testFile, header, rows1, ',', true)
	err := WriteRowsToCSV(testFile, header, rows2, ',', false)
	if err != nil {
		t.Fatal(err)
	}

	all := readAllRows(t, testFile, ',')
	if len(all) != 3 {
		t.Errorf("want 3 rows, got %d", len(all))
	}
}

// 3. 追加模式：表头不匹配 → 不报错
func TestWrite_Append_HeaderMismatch_Error(t *testing.T) {
	cleanup()
	h1 := []string{"id", "name"}
	h2 := []string{"id", "age"}

	_ = WriteRowsToCSV(testFile, h1, nil, ',', true)
	err := WriteRowsToCSV(testFile, h2, nil, ',', false)
	if err != nil {
		t.Error("expected mismatch error")
	}
}

// 4. 空表头：只写数据，不写表头
func TestWrite_EmptyHeader(t *testing.T) {
	cleanup()
	rows := [][]string{{"1", "2"}}
	err := WriteRowsToCSV(testFile, nil, rows, ',', true)
	if err != nil {
		t.Fatal(err)
	}
	all := readAllRows(t, testFile, ',')
	if len(all) != 1 {
		t.Errorf("want 1 row, got %d", len(all))
	}
}

// 5. 空数据：只写表头
func TestWrite_OnlyHeader(t *testing.T) {
	cleanup()
	header := []string{"a", "b"}
	err := WriteRowsToCSV(testFile, header, nil, ',', true)
	if err != nil {
		t.Fatal(err)
	}
	all := readAllRows(t, testFile, ',')
	if len(all) != 1 {
		t.Errorf("want 1 header row, got %d", len(all))
	}
}

// 6. 既无表头也无数据 → 不创建文件/不写入
func TestWrite_NothingToWrite(t *testing.T) {
	cleanup()
	err := WriteRowsToCSV(testFile, nil, nil, ',', true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(testFile); err == nil {
		t.Error("file should not exist")
	}
}

// 7. 自定义分隔符：制表符 \t
func TestWrite_CustomDelimiter_Tab(t *testing.T) {
	cleanup()
	header := []string{"x", "y"}
	rows := [][]string{{"1", "2"}}
	err := WriteRowsToCSV(testFile, header, rows, '\t', true)
	if err != nil {
		t.Fatal(err)
	}
	all := readAllRows(t, testFile, '\t')
	if len(all) != 2 || all[1][0] != "1" {
		t.Error("tab delimiter failed")
	}
}

// 8. 自定义分隔符：分号 ;
func TestWrite_CustomDelimiter_Semicolon(t *testing.T) {
	cleanup()
	err := WriteRowsToCSV(testFile, []string{"a"}, [][]string{{"1"}}, ';', true)
	if err != nil {
		t.Fatal(err)
	}
	all := readAllRows(t, testFile, ';')
	if len(all) != 2 {
		t.Error("semicolon delimiter failed")
	}
}

// 9. 文件不存在时的追加模式 = 新建文件并写入表头+数据
func TestWrite_Append_FileNotExist(t *testing.T) {
	cleanup()
	header := []string{"id"}
	rows := [][]string{{"99"}}
	err := WriteRowsToCSV(testFile, header, rows, ',', false)
	if err != nil {
		t.Fatal(err)
	}
	all := readAllRows(t, testFile, ',')
	if len(all) != 2 {
		t.Error("append on new file should write header+data")
	}
}

// 10. 多次追加：数据累加，表头只出现一次
func TestWrite_MultipleAppend(t *testing.T) {
	cleanup()
	h := []string{"a"}
	_ = WriteRowsToCSV(testFile, h, [][]string{{"1"}}, ',', true)
	_ = WriteRowsToCSV(testFile, h, [][]string{{"2"}}, ',', false)
	_ = WriteRowsToCSV(testFile, h, [][]string{{"3"}}, ',', false)

	all := readAllRows(t, testFile, ',')
	if len(all) != 4 {
		t.Errorf("want 4 rows, got %d", len(all))
	}
}
