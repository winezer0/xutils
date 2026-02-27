package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppendLines(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "hashcracker_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "test_append.txt")

	// 1. 测试追加到新文件
	lines1 := []string{"line1", "line2"}
	if err := AppendLines(filePath, lines1); err != nil {
		t.Fatalf("AppendLines failed: %v", err)
	}

	content1, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	expected1 := "line1\nline2\n"
	if string(content1) != expected1 {
		t.Errorf("Expected content %q, got %q", expected1, string(content1))
	}

	// 2. 测试追加到现有文件
	lines2 := []string{"line3", "line4"}
	if err := AppendLines(filePath, lines2); err != nil {
		t.Fatalf("AppendLines failed: %v", err)
	}

	content2, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	expected2 := "line1\nline2\nline3\nline4\n"
	if string(content2) != expected2 {
		t.Errorf("Expected content %q, got %q", expected2, string(content2))
	}

	// 3. 测试追加空列表
	if err := AppendLines(filePath, []string{}); err != nil {
		t.Fatalf("AppendLines failed with empty list: %v", err)
	}
	content3, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content3) != expected2 {
		t.Errorf("Content changed after appending empty list")
	}
}
