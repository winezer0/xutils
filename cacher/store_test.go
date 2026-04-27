package cacher

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestCacheManager_Normal 验证基础的读写与持久化流程。
func TestCacheManager_Normal(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer func() { _ = cm.Close() }()

	if err := cm.Set("key1", "value1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	val, ok := cm.Get("key1")
	if !ok || val != "value1" {
		t.Fatalf("Expected value1, got %v", val)
	}
	if err := cm.SaveCache(); err != nil {
		t.Fatalf("SaveCache failed: %v", err)
	}
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		t.Fatal("Expected cache file to exist")
	}
	if s, ok := cm.GetString("key1"); !ok || s != "value1" {
		t.Fatal("GetString failed")
	}
}

// TestCacheManager_GetAs 验证结构体反序列化读取。
func TestCacheManager_GetAs(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_getas.json")
	defer os.Remove(cacheFile)

	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	cm := NewCacheManager(cacheFile)
	defer func() { _ = cm.Close() }()
	if err := cm.Set("user", User{Name: "Alice", Age: 30}); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	var user User
	ok, err := cm.GetAs("user", &user)
	if err != nil || !ok {
		t.Fatalf("GetAs failed: %v", err)
	}
	if user.Name != "Alice" || user.Age != 30 {
		t.Fatalf("Unexpected user: %+v", user)
	}
}

// TestCacheManager_GetAs_Direct 验证同类型值可直接赋值返回。
func TestCacheManager_GetAs_Direct(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_getas_direct.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer func() { _ = cm.Close() }()
	if err := cm.Set("str_key", "hello"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	var value string
	ok, err := cm.GetAs("str_key", &value)
	if err != nil || !ok || value != "hello" {
		t.Fatalf("GetAs string failed: %v, %v", err, value)
	}
}

// TestCacheManager_MaxSize 验证条目数上限。
func TestCacheManager_MaxSize(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_maxsize.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManagerWithSeconds(cacheFile, 10, 1, 10*1024*1024)
	defer func() { _ = cm.Close() }()
	if err := cm.Set("k1", "v1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if err := cm.Set("k2", "v2"); !errors.Is(err, ErrCacheFull) {
		t.Fatalf("Expected ErrCacheFull, got %v", err)
	}
}

// TestCacheManager_Del 验证删除键值行为。
func TestCacheManager_Del(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_del.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer func() { _ = cm.Close() }()
	if err := cm.Del("missing"); !errors.Is(err, ErrCacheKeyNotFound) {
		t.Fatalf("Expected ErrCacheKeyNotFound, got %v", err)
	}
	if err := cm.Set("k1", "v1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if err := cm.Del("k1"); err != nil {
		t.Fatalf("Del failed: %v", err)
	}
	if _, ok := cm.Get("k1"); ok {
		t.Fatal("Expected key to be deleted")
	}
}

// TestCacheManager_GetAs_ErrorCases 验证 GetAs 错误路径。
func TestCacheManager_GetAs_ErrorCases(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_getas_err.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer func() { _ = cm.Close() }()

	var value string
	if ok, err := cm.GetAs("missing", &value); !errors.Is(err, ErrCacheKeyNotFound) || ok {
		t.Fatalf("Expected ErrCacheKeyNotFound, got %v", err)
	}
	if err := cm.Set("k1", "v1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if ok, err := cm.GetAs("k1", "notptr"); !errors.Is(err, ErrCacheInvalidValue) || ok {
		t.Fatalf("Expected ErrCacheInvalidValue, got %v", err)
	}
}

// TestCacheManager_MaxDataSize 验证字节数上限。
func TestCacheManager_MaxDataSize(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_maxdatasize.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManagerWithSeconds(cacheFile, 10, 100, 100)
	defer func() { _ = cm.Close() }()
	if err := cm.Set("small", "small"); err != nil {
		t.Fatalf("Set small data failed: %v", err)
	}
	if err := cm.Set("large", string(make([]byte, 200))); !errors.Is(err, ErrCacheFull) {
		t.Fatalf("Expected ErrCacheFull, got %v", err)
	}
}

// TestCacheManager_DataSizeCalculation 验证更新与删除会同步刷新大小统计。
func TestCacheManager_DataSizeCalculation(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_datacalculation.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer func() { _ = cm.Close() }()
	if err := cm.Set("key1", "test1"); err != nil {
		t.Fatalf("Set key1 failed: %v", err)
	}
	if err := cm.Set("key1", "test2"); err != nil {
		t.Fatalf("Update key1 failed: %v", err)
	}
	if err := cm.Del("key1"); err != nil {
		t.Fatalf("Del key1 failed: %v", err)
	}
	if _, ok := cm.Get("key1"); ok {
		t.Fatal("Expected key1 to be deleted")
	}
	if cm.state.currentSize != 0 {
		t.Fatalf("Expected currentSize to be 0, got %d", cm.state.currentSize)
	}
}
