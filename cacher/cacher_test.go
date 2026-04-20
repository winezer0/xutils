package cacher

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCacheManager_Normal(t *testing.T) {
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache.json")
	defer os.Remove(cacheFile)

	// Clean up if exists
	os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer func() {
		if err := cm.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	// Test Set and Get
	if err := cm.Set("key1", "value1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	val, ok := cm.Get("key1")
	if !ok {
		t.Error("Expected key1 to exist")
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	// Test Persistence (Save)
	err := cm.SaveCache()
	if err != nil {
		t.Errorf("SaveCache failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		t.Error("Cache file was not created")
	}

	// Test GetString
	s, ok := cm.GetString("key1")
	if !ok || s != "value1" {
		t.Error("GetString failed")
	}
}

func TestCacheManager_GetAs(t *testing.T) {
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_getas.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer func() {
		if err := cm.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	u := User{Name: "Alice", Age: 30}
	if err := cm.Set("user", u); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	var u2 User
	ok, err := cm.GetAs("user", &u2)
	if err != nil || !ok {
		t.Errorf("GetAs failed: %v", err)
	}
	if u2.Name != "Alice" || u2.Age != 30 {
		t.Errorf("GetAs returned incorrect struct: %+v", u2)
	}
}

func TestCacheManager_CloseIdempotent(t *testing.T) {
	// Test normal case
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_close.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	// Call Close multiple times
	_ = cm.Close()
	_ = cm.Close() // Should not panic
	_ = cm.Close() // Should not panic

	// Test empty file case
	cmEmpty := NewCacheManager("")
	_ = cmEmpty.Close()
	_ = cmEmpty.Close() // Should not panic
}

func TestCacheManager_GetAs_Direct(t *testing.T) {
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_getas_direct.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer func() {
		if err := cm.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	// Case 1: Same type (Direct assignment optimization)
	if err := cm.Set("str_key", "hello"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	var s string
	if ok, err := cm.GetAs("str_key", &s); err != nil || !ok {
		t.Errorf("GetAs string failed: %v", err)
	}
	if s != "hello" {
		t.Errorf("Expected 'hello', got %v", s)
	}

	// Case 2: JSON fallback (Struct conversion)
	// Note: When loading from memory set with struct, raw is struct.
	// When loading from file (JSON), raw is map[string]interface{}.
	// Here we test in-memory struct-to-struct copy if optimization works,
	// or struct-to-struct via JSON if types differ slightly (though usually types must match for direct assign).

	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	u := User{Name: "Bob", Age: 40}
	if err := cm.Set("user", u); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	var u2 User
	if ok, err := cm.GetAs("user", &u2); err != nil || !ok {
		t.Errorf("GetAs struct failed: %v", err)
	}
	if u2 != u {
		t.Errorf("Expected %+v, got %+v", u, u2)
	}
}

func TestCacheManager_AutoSave(t *testing.T) {
	// Skip if we can't wait for 10s or mock time.
	// Since autoSaveWorker runs every 10s, it's slow to test.
	// We'll skip precise timing test but just ensure start/stop works.
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_autosave.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	if err := cm.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestCacheManager_MaxSize(t *testing.T) {
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_maxsize.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManagerWithOptions(cacheFile, 10, 1, 10*1024*1024)
	defer func() {
		if err := cm.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	if err := cm.Set("k1", "v1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if err := cm.Set("k2", "v2"); !errors.Is(err, ErrCacheFull) {
		t.Errorf("Expected ErrCacheFull, got %v", err)
	}
}

func TestCacheManager_Del(t *testing.T) {
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_del.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer func() {
		if err := cm.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	if err := cm.Del("missing"); !errors.Is(err, ErrCacheKeyNotFound) {
		t.Errorf("Expected ErrCacheKeyNotFound, got %v", err)
	}
	if err := cm.Set("k1", "v1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if err := cm.Del("k1"); err != nil {
		t.Fatalf("Del failed: %v", err)
	}
	if _, ok := cm.Get("k1"); ok {
		t.Error("Expected key to be deleted")
	}
}

func TestCacheManager_GetAs_ErrorCases(t *testing.T) {
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_getas_err.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer func() {
		if err := cm.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	var v string
	if ok, err := cm.GetAs("missing", &v); !errors.Is(err, ErrCacheKeyNotFound) || ok {
		t.Errorf("Expected ErrCacheKeyNotFound, got %v", err)
	}
	if err := cm.Set("k1", "v1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if ok, err := cm.GetAs("k1", "notptr"); !errors.Is(err, ErrCacheInvalidValue) || ok {
		t.Errorf("Expected ErrCacheInvalidValue, got %v", err)
	}
}

func TestCacheManager_MaxDataSize(t *testing.T) {
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_maxdatasize.json")
	defer os.Remove(cacheFile)

	// 创建一个小的 maxDataSize 限制（100字节）
	cm := NewCacheManagerWithOptions(cacheFile, 10, 100, 100)
	defer func() {
		if err := cm.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	// 测试添加小数据
	smallData := "small"
	if err := cm.Set("small", smallData); err != nil {
		t.Fatalf("Set small data failed: %v", err)
	}

	// 测试添加大数据（应该失败）
	largeData := string(make([]byte, 200)) // 200字节数据
	if err := cm.Set("large", largeData); !errors.Is(err, ErrCacheFull) {
		t.Errorf("Expected ErrCacheFull for large data, got %v", err)
	}
}

func TestCacheManager_DataSizeCalculation(t *testing.T) {
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_datacalculation.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer func() {
		if err := cm.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	// 添加数据
	data1 := "test1"
	if err := cm.Set("key1", data1); err != nil {
		t.Fatalf("Set key1 failed: %v", err)
	}

	// 更新数据
	data2 := "test2"
	if err := cm.Set("key1", data2); err != nil {
		t.Fatalf("Update key1 failed: %v", err)
	}

	// 删除数据
	if err := cm.Del("key1"); err != nil {
		t.Fatalf("Del key1 failed: %v", err)
	}

	// 验证删除后数据不存在
	if _, ok := cm.Get("key1"); ok {
		t.Error("Expected key1 to be deleted")
	}
}

func TestCacheManager_LoadCacheSizeLimit(t *testing.T) {
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_loadsize.json")
	defer os.Remove(cacheFile)

	// 创建一个包含大数据的缓存文件
	largeData := map[string]interface{}{
		"large": string(make([]byte, 500)),
		"small": "test",
	}

	// 写入文件
	data, _ := json.Marshal(largeData)
	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		t.Fatalf("Write cache file failed: %v", err)
	}

	// 使用小的 maxDataSize 加载
	cm := NewCacheManagerWithOptions(cacheFile, 10, 100, 100)
	defer func() {
		if err := cm.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	// 验证缓存被限制
	if _, ok := cm.Get("large"); ok {
		t.Error("Expected large data to be removed due to size limit")
	}
}
