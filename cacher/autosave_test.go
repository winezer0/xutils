package cacher

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

// TestCacheManager_AutoSave 验证数据变更后会按计时器自动落盘。
func TestCacheManager_AutoSave(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_autosave.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManagerWithConfig(Config{
		CacheFile:    cacheFile,
		SaveInterval: 30 * time.Millisecond,
		MaxEntries:   100,
		MaxDataBytes: 10 * 1024 * 1024,
	})
	defer func() { _ = cm.Close() }()

	if cm.state.autoSaveTimer != nil {
		t.Fatal("Expected autoSaveTimer to be nil before data changes")
	}
	if err := cm.Set("key1", "value1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if cm.state.autoSaveTimer == nil {
		t.Fatal("Expected autoSaveTimer to be scheduled after data changes")
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(cacheFile)
		if err == nil && len(data) > 0 {
			var payload map[string]any
			if err := json.Unmarshal(data, &payload); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if payload["key1"] == "value1" {
				return
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("Expected auto save to persist cache data")
}

// TestCacheManager_CloseDuringAutoSave_OnlyPersistOnce 验证 Close 与自动保存并发时不会重复落盘。
func TestCacheManager_CloseDuringAutoSave_OnlyPersistOnce(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_close_autosave.json")
	defer os.Remove(cacheFile)

	originalPersist := persistCacheFileFunc
	defer func() {
		persistCacheFileFunc = originalPersist
	}()

	var count atomic.Int32
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	persistCacheFileFunc = func(file string, data []byte) error {
		count.Add(1)
		select {
		case started <- struct{}{}:
		default:
		}
		<-release
		return originalPersist(file, data)
	}

	cm := NewCacheManagerWithConfig(Config{
		CacheFile:    cacheFile,
		SaveInterval: 20 * time.Millisecond,
		MaxEntries:   100,
		MaxDataBytes: 10 * 1024 * 1024,
	})
	if err := cm.Set("key", "value"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("Expected auto save to start")
	}

	done := make(chan error, 1)
	go func() {
		done <- cm.Close()
	}()

	time.Sleep(50 * time.Millisecond)
	close(release)
	if err := <-done; err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	if count.Load() != 1 {
		t.Fatalf("Expected one persist call, got %d", count.Load())
	}
}

// TestCacheManager_CopyClose_NoHang 验证值拷贝后的 Close 不会死锁或 panic。
func TestCacheManager_CopyClose_NoHang(t *testing.T) {
	cacheFile := filepath.Join(os.TempDir(), "test_cache_copyclose.json")
	defer os.Remove(cacheFile)

	cm := NewCacheManagerWithSeconds(cacheFile, 3600, 10000, 10*1024*1024)
	if err := cm.Set("k1", "v1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	cmCopy := *cm
	done := make(chan struct{})
	panicCh := make(chan any, 1)
	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				panicCh <- recovered
				return
			}
			close(done)
		}()
		_ = cmCopy.Close()
		_ = cm.Close()
	}()

	select {
	case recovered := <-panicCh:
		t.Fatalf("Close panicked: %v", recovered)
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Close timeout")
	}
}
