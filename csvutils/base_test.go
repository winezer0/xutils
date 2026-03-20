package csvutils

import (
	"os"
	"path/filepath"
	"testing"
)

// makeTempCSV 创建临时CSV文件（测试辅助）
func makeTempCSV(t *testing.T, name string, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp csv failed: %v", err)
	}
	return path
}

const testFile = "test_output.csv"

func cleanup() {
	_ = os.Remove(testFile)
}
