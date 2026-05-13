package poolwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPoolWriteStore 验证统一池能够正确写入响应文件。
func TestPoolWriteStore(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "resp.txt")

	p := NewPool(1, 8)
	p.Submit(StoreTask{
		StorePath:    storePath,
		StoreContent: []byte("response body"),
		WriteRaw:     false,
	})
	p.StopAndWait()

	storeData, err := os.ReadFile(storePath)
	if err != nil {
		t.Fatalf("read store file failed: %v", err)
	}
	if string(storeData) != "response body\n" && string(storeData) != "response body" {
		t.Fatalf("unexpected store content: %q", string(storeData))
	}
}

// TestPoolWriteCache 验证统一池能够正确写入缓存文件。
func TestPoolWriteCache(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "cache.txt")

	p := NewPool(1, 8)
	p.Submit(StoreTask{
		StorePath:    cachePath,
		StoreContent: []byte("task_hash_1"),
		WriteRaw:     false,
	})
	p.StopAndWait()

	cacheData, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatalf("read cache file failed: %v", err)
	}
	cacheText := strings.TrimSpace(string(cacheData))
	if cacheText != "task_hash_1" {
		t.Fatalf("unexpected cache content: %q", cacheText)
	}
}

// TestPoolWriteMultipleTasks 验证统一池可以处理多个不同类型的写入任务。
func TestPoolWriteMultipleTasks(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "resp.txt")
	cachePath := filepath.Join(tmpDir, "cache.txt")
	errorCachePath := filepath.Join(tmpDir, "error.cache")

	p := NewPool(2, 16)

	// 投递存储任务
	p.Submit(StoreTask{
		StorePath:    storePath,
		StoreContent: []byte("response body"),
		WriteRaw:     false,
	})

	// 投递缓存任务
	p.Submit(StoreTask{
		StorePath:    cachePath,
		StoreContent: []byte("task_hash_1"),
		WriteRaw:     false,
	})

	// 投递错误缓存任务
	p.Submit(StoreTask{
		StorePath:    errorCachePath,
		StoreContent: []byte("task_hash_2"),
		WriteRaw:     false,
	})

	p.StopAndWait()

	// 验证存储文件
	storeData, err := os.ReadFile(storePath)
	if err != nil {
		t.Fatalf("read store file failed: %v", err)
	}
	if string(storeData) != "response body\n" && string(storeData) != "response body" {
		t.Fatalf("unexpected store content: %q", string(storeData))
	}

	// 验证缓存文件
	cacheData, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatalf("read cache file failed: %v", err)
	}
	cacheText := strings.TrimSpace(string(cacheData))
	if cacheText != "task_hash_1" {
		t.Fatalf("unexpected cache content: %q", cacheText)
	}

	// 验证错误缓存文件
	errorCacheData, err := os.ReadFile(errorCachePath)
	if err != nil {
		t.Fatalf("read error cache file failed: %v", err)
	}
	errorCacheText := strings.TrimSpace(string(errorCacheData))
	if errorCacheText != "task_hash_2" {
		t.Fatalf("unexpected error cache content: %q", errorCacheText)
	}
}

// TestPoolEmptyPath 验证空路径不会报错。
func TestPoolEmptyPath(t *testing.T) {
	p := NewPool(1, 8)
	p.Submit(StoreTask{
		StorePath:    "",
		StoreContent: []byte("should not write"),
		WriteRaw:     false,
	})
	p.StopAndWait()
}

// TestPoolWriteRaw 验证二进制写入模式。
func TestPoolWriteRaw(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "resp.bin")

	p := NewPool(1, 8)
	p.Submit(StoreTask{
		StorePath:    storePath,
		StoreContent: []byte{0x00, 0x01, 0x02, 0x03},
		WriteRaw:     true,
	})
	p.StopAndWait()

	storeData, err := os.ReadFile(storePath)
	if err != nil {
		t.Fatalf("read store file failed: %v", err)
	}
	if len(storeData) != 4 {
		t.Fatalf("unexpected store content length: %d", len(storeData))
	}
}

// TestPoolWriteWithOverwrite 验证 Overwrite 参数可以覆盖已存在的文件。
func TestPoolWriteWithOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "resp.txt")

	// 先写入初始内容
	p := NewPool(1, 8)
	p.Submit(StoreTask{
		StorePath:    storePath,
		StoreContent: []byte("initial content"),
		WriteRaw:     false,
		Overwrite:    false,
	})
	p.StopAndWait()

	// 验证初始内容
	initialData, err := os.ReadFile(storePath)
	if err != nil {
		t.Fatalf("read store file failed: %v", err)
	}
	if !strings.Contains(string(initialData), "initial content") {
		t.Fatalf("unexpected initial content: %q", string(initialData))
	}

	// 使用 Overwrite=false 再次写入，应该追加
	p2 := NewPool(1, 8)
	p2.Submit(StoreTask{
		StorePath:    storePath,
		StoreContent: []byte("appended content"),
		WriteRaw:     false,
		Overwrite:    false,
	})
	p2.StopAndWait()

	appendData, err := os.ReadFile(storePath)
	if err != nil {
		t.Fatalf("read store file failed: %v", err)
	}
	if !strings.Contains(string(appendData), "initial content") || !strings.Contains(string(appendData), "appended content") {
		t.Fatalf("expected both initial and appended content, got: %q", string(appendData))
	}

	// 使用 Overwrite=true 再次写入，应该覆盖
	p3 := NewPool(1, 8)
	p3.Submit(StoreTask{
		StorePath:    storePath,
		StoreContent: []byte("overwritten content"),
		WriteRaw:     false,
		Overwrite:    true,
	})
	p3.StopAndWait()

	overwriteData, err := os.ReadFile(storePath)
	if err != nil {
		t.Fatalf("read store file failed: %v", err)
	}
	if !strings.Contains(string(overwriteData), "overwritten content") {
		t.Fatalf("expected overwritten content, got: %q", string(overwriteData))
	}
}

// TestPoolFailCount 验证失败统计功能。
func TestPoolFailCount(t *testing.T) {
	p := NewPool(1, 8)

	// 投递空路径任务（不会失败）
	p.Submit(StoreTask{
		StorePath:    "",
		StoreContent: []byte("test"),
	})

	// 投递无效路径任务（会失败）
	p.Submit(StoreTask{
		StorePath:    string([]byte{0x00}),
		StoreContent: []byte("test"),
	})

	p.StopAndWait()

	// 验证失败计数
	failCount := p.GetFailCount()
	if failCount < 1 {
		t.Fatalf("expected at least 1 failure, got: %d", failCount)
	}
}
