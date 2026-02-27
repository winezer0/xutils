package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
)

func TestDeduplicateFile(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "hashcracker_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "test_dedup.txt")

	// 准备包含重复行的文件内容
	content := "line1\nline2\nline1\nline3\nline2\nline4\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// 执行去重
	if err := DeduplicateFile(filePath); err != nil {
		t.Fatalf("DeduplicateFile failed: %v", err)
	}

	// 读取结果验证
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open result file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	expectedLines := []string{"line1", "line2", "line3", "line4"}
	if len(lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}

	// 验证内容（顺序可能保留，也可能不保留，但在 map 实现中通常顺序是不确定的，
	// 不过我的实现是按顺序 append 到 slice 的，所以顺序应该保留）
	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expectedLines[i], line)
		}
	}
}
