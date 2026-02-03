package cacher

import (
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
	defer cm.Close()

	// Test Set and Get
	cm.Set("key1", "value1")
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

func TestCacheManager_EmptyFile(t *testing.T) {
	cm := NewCacheManager("")
	defer cm.Close()

	// Test Set (should do nothing)
	cm.Set("key1", "value1")

	// Test Get (should return false)
	val, ok := cm.Get("key1")
	if ok {
		t.Error("Expected key1 to NOT exist when cacheFile is empty")
	}
	if val != nil {
		t.Errorf("Expected nil value, got %v", val)
	}

	// Test Save (should do nothing, no error)
	err := cm.SaveCache()
	if err != nil {
		t.Errorf("SaveCache with empty file returned error: %v", err)
	}

	// Test Clear (should do nothing, no error)
	err = cm.Clear()
	if err != nil {
		t.Errorf("Clear with empty file returned error: %v", err)
	}
}

func TestCacheManager_GetAs(t *testing.T) {
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_getas.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer cm.Close()

	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	u := User{Name: "Alice", Age: 30}
	cm.Set("user", u)

	var u2 User
	ok := cm.GetAs("user", &u2)
	if !ok {
		t.Error("GetAs failed")
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
	cm.Close()
	cm.Close() // Should not panic
	cm.Close() // Should not panic

	// Test empty file case
	cmEmpty := NewCacheManager("")
	cmEmpty.Close()
	cmEmpty.Close() // Should not panic
}

func TestCacheManager_GetAs_Direct(t *testing.T) {
	tmpDir := os.TempDir()
	cacheFile := filepath.Join(tmpDir, "test_cache_getas_direct.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManager(cacheFile)
	defer cm.Close()

	// Case 1: Same type (Direct assignment optimization)
	cm.Set("str_key", "hello")
	var s string
	if !cm.GetAs("str_key", &s) {
		t.Error("GetAs string failed")
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
	cm.Set("user", u)

	var u2 User
	if !cm.GetAs("user", &u2) {
		t.Error("GetAs struct failed")
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
	cm.Close() // Should return immediately
}
