package csvutils

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"
)

// TestReadCSV2Rows_WithHeader 修正版：匹配表头固定取第一行的逻辑
func TestReadCSV2Rows_WithHeader(t *testing.T) {
	testReadFile := "test_read_header.csv"
	// 1. 准备测试文件
	cleanupTestFile(t, testReadFile)
	defer cleanupTestFile(t, testReadFile)

	// 写入测试内容：表头+2行数据
	file, err := os.Create(testReadFile)
	if err != nil {
		t.Fatalf("create test file failed: %v", err)
	}
	writer := csv.NewWriter(file)
	err = writer.WriteAll([][]string{
		{"name", "age", "city"}, // 第一行：固定表头
		{"张三", "25", "北京"},      // 数据行1
		{"李四", "30", "上海"},      // 数据行2
	})
	if err != nil {
		t.Fatalf("write test csv failed: %v", err)
	}
	writer.Flush()
	file.Close()

	// 2. 测试 skipRows=1（跳过1行数据，表头仍为第一行）
	header, rows, err := ReadCSV2RowsWithSkip(testReadFile, ',', 1)
	if err != nil {
		t.Fatalf("read csv failed (skipRows=1): %v", err)
	}
	expectedHeader := []string{"name", "age", "city"}
	if !reflect.DeepEqual(header, expectedHeader) {
		t.Errorf("header mismatch (skipRows=1): expected %v, got %v", expectedHeader, header)
	}
	// skipRows=1 → 跳过数据行1，仅保留数据行2
	if len(rows) != 1 {
		t.Errorf("rows count mismatch (skipRows=1): expected 1, got %d", len(rows))
	}
	if !reflect.DeepEqual(rows[0], []string{"李四", "30", "上海"}) {
		t.Errorf("first row mismatch (skipRows=1): expected %v, got %v", []string{"李四", "30", "上海"}, rows[0])
	}

	// 3. 测试 skipRows=0（不跳过数据行，表头=第一行，数据行=2行）
	header2, rows2, err := ReadCSV2RowsWithSkip(testReadFile, ',', 0)
	if err != nil {
		t.Fatalf("read csv failed (skipRows=0): %v", err)
	}
	if !reflect.DeepEqual(header2, expectedHeader) {
		t.Errorf("header mismatch (skipRows=0): expected %v, got %v", expectedHeader, header2)
	}
	if len(rows2) != 2 {
		t.Errorf("rows count mismatch (skipRows=0): expected 2, got %d", len(rows2))
	}
	if !reflect.DeepEqual(rows2[0], []string{"张三", "25", "北京"}) {
		t.Errorf("first row mismatch (skipRows=0): expected %v, got %v", []string{"张三", "25", "北京"}, rows2[0])
	}
	if !reflect.DeepEqual(rows2[1], []string{"李四", "30", "上海"}) {
		t.Errorf("second row mismatch (skipRows=0): expected %v, got %v", []string{"李四", "30", "上海"}, rows2[1])
	}

	// 4. 测试 skipRows=2（跳过2行数据，数据行空）
	header3, rows3, err := ReadCSV2RowsWithSkip(testReadFile, ',', 2)
	if err != nil {
		t.Fatalf("read csv failed (skipRows=2): %v", err)
	}
	if !reflect.DeepEqual(header3, expectedHeader) {
		t.Errorf("header mismatch (skipRows=2): expected %v, got %v", expectedHeader, header3)
	}
	if len(rows3) != 0 {
		t.Errorf("rows count mismatch (skipRows=2): expected 0, got %d", len(rows3))
	}
}

// TestReadCSV2Rows_EmptyFile 测试空文件（不变）
func TestReadCSV2Rows_EmptyFile(t *testing.T) {
	cleanupTestFile(t, testReadFile)
	defer cleanupTestFile(t, testReadFile)

	// 创建空文件
	file, err := os.Create(testReadFile)
	if err != nil {
		t.Fatalf("create empty test file failed: %v", err)
	}
	file.Close()

	header, rows, err := ReadCSV2RowsWithSkip(testReadFile, ',', 0)
	if err == nil {
		t.Error("expected error for empty file, but got nil")
	}
	if header != nil {
		t.Errorf("expected nil header for empty file, got %v", header)
	}
	if rows != nil {
		t.Errorf("expected nil rows for empty file, got %v", rows)
	}
}
