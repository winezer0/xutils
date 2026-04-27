package cacher

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestCacheManager_LoadCacheSizeLimit 验证加载时会按大小限制裁剪缓存。
func TestCacheManager_LoadCacheSizeLimit(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_loadsize.json")
	defer os.Remove(cacheFile)

	largeData := map[string]interface{}{
		"large": string(make([]byte, 500)),
		"small": "test",
	}
	data, _ := json.Marshal(largeData)
	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		t.Fatalf("Write cache file failed: %v", err)
	}

	cm := NewCacheManagerWithSeconds(cacheFile, 10, 100, 100)
	defer func() { _ = cm.Close() }()
	if _, ok := cm.Get("large"); ok {
		t.Fatal("Expected large data to be removed due to size limit")
	}
	if !cm.state.modified {
		t.Fatal("Expected trimmed cache to remain marked as modified")
	}
}

// TestCacheManager_ConcurrentManagersSameFile_Save 验证多实例并发写同一文件不会损坏数据。
func TestCacheManager_ConcurrentManagersSameFile_Save(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_concurrent_samefile.json")
	defer os.Remove(cacheFile)

	cm1 := NewCacheManagerWithSeconds(cacheFile, 3600, 10000, 10*1024*1024)
	cm2 := NewCacheManagerWithSeconds(cacheFile, 3600, 10000, 10*1024*1024)
	defer func() { _ = cm1.Close() }()
	defer func() { _ = cm2.Close() }()

	if err := cm1.Set("k1", "v1"); err != nil {
		t.Fatalf("Set cm1 failed: %v", err)
	}
	if err := cm2.Set("k2", "v2"); err != nil {
		t.Fatalf("Set cm2 failed: %v", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		errCh <- cm1.SaveCache()
	}()
	go func() {
		defer wg.Done()
		errCh <- cm2.SaveCache()
	}()
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("SaveCache failed: %v", err)
		}
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if len(payload) == 0 {
		t.Fatal("Expected cache file to contain data")
	}
}

// TestPersistCacheFile_Overwrite 验证持久化函数可正确覆盖旧文件。
func TestPersistCacheFile_Overwrite(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_persist_overwrite.json")
	defer os.Remove(cacheFile)

	if err := os.WriteFile(cacheFile, []byte(`{"old":1}`), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if err := persistCacheFile(cacheFile, []byte(`{"new":2}`)); err != nil {
		t.Fatalf("persistCacheFile failed: %v", err)
	}
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(data) != `{"new":2}` {
		t.Fatalf("unexpected content: %s", string(data))
	}
}

// TestWriteCacheFileByRename 验证非 Windows 分支的临时文件替换逻辑。
func TestWriteCacheFileByRename(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip rename-path test on windows")
	}

	cacheFile := filepath.Join(os.TempDir(), "test_cache_rename_path.json")
	defer os.Remove(cacheFile)
	if err := writeCacheFileByRename(cacheFile, []byte(`{"rename":true}`)); err != nil {
		t.Fatalf("writeCacheFileByRename failed: %v", err)
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(data) != `{"rename":true}` {
		t.Fatalf("unexpected content: %s", string(data))
	}
}

// TestCacheManager_LoadCache_InvalidJSON 验证损坏 JSON 会返回解析错误。
func TestCacheManager_LoadCache_InvalidJSON(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_invalid_json.json")
	defer os.Remove(cacheFile)

	if err := os.WriteFile(cacheFile, []byte(`{"broken":`), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	cm := NewCacheManagerWithConfig(Config{
		CacheFile:    cacheFile,
		SaveInterval: time.Hour,
		MaxEntries:   100,
		MaxDataBytes: 10 * 1024 * 1024,
	})
	if err := cm.LoadCache(); err == nil {
		t.Fatal("Expected LoadCache to fail for invalid JSON")
	}
}

// TestCacheManager_SaveCache_CreatesDir 验证保存时会自动创建缺失目录。
func TestCacheManager_SaveCache_CreatesDir(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "xutils-cacher", "nested", "cache.json")
	_ = os.Remove(cacheFile)
	_ = os.RemoveAll(filepath.Dir(filepath.Dir(cacheFile)))
	defer os.RemoveAll(filepath.Join(os.TempDir(), "xutils-cacher"))

	cm := NewCacheManager(cacheFile)
	defer func() { _ = cm.Close() }()
	if err := cm.Set("key", "value"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if err := cm.SaveCache(); err != nil {
		t.Fatalf("SaveCache failed: %v", err)
	}
	if _, err := os.Stat(cacheFile); err != nil {
		t.Fatalf("Expected cache file to exist: %v", err)
	}
}

// TestCacheManager_SaveCacheSnapshotDoesNotBlockSet 验证锁外写盘不会阻塞后续写入。
func TestCacheManager_SaveCacheSnapshotDoesNotBlockSet(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_snapshot_save.json")
	defer os.Remove(cacheFile)

	originalPersist := persistCacheFileFunc
	defer func() {
		persistCacheFileFunc = originalPersist
	}()

	started := make(chan struct{}, 1)
	release := make(chan struct{})
	persistCacheFileFunc = func(file string, data []byte) error {
		select {
		case started <- struct{}{}:
		default:
		}
		<-release
		return originalPersist(file, data)
	}

	cm := NewCacheManagerWithSeconds(cacheFile, 3600, 10000, 10*1024*1024)
	defer func() { _ = cm.Close() }()
	if err := cm.Set("key", "v1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	saveDone := make(chan error, 1)
	go func() {
		saveDone <- cm.SaveCache()
	}()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("Expected SaveCache to enter persist phase")
	}

	setDone := make(chan error, 1)
	go func() {
		setDone <- cm.Set("key", "v2")
	}()

	select {
	case err := <-setDone:
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Expected Set not to block on disk I/O")
	}

	close(release)
	if err := <-saveDone; err != nil {
		t.Fatalf("SaveCache failed: %v", err)
	}
	if value, ok := cm.GetString("key"); !ok || value != "v2" {
		t.Fatalf("Expected in-memory value v2, got %v", value)
	}
	if !cm.state.modified {
		t.Fatal("Expected state to remain modified after concurrent update during save")
	}
	if err := cm.SaveCache(); err != nil {
		t.Fatalf("Second SaveCache failed: %v", err)
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	var payload map[string]string
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if payload["key"] != "v2" {
		t.Fatalf("Expected persisted value v2, got %s", payload["key"])
	}
}

// TestCacheManager_LoadCache_HalfWrittenJSON 验证半写入 JSON 也会返回解析错误。
func TestCacheManager_LoadCache_HalfWrittenJSON(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_half_written.json")
	defer os.Remove(cacheFile)

	if err := os.WriteFile(cacheFile, []byte(`{"key":"value"`), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	cm := NewCacheManager(cacheFile)
	if err := cm.LoadCache(); err == nil {
		t.Fatal("Expected LoadCache to fail for half-written JSON")
	}
}
