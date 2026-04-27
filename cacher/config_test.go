package cacher

import (
	"testing"
	"time"
)

// TestNormalizeSaveInterval_BareInt 验证裸整数保存间隔会按秒解释。
func TestNormalizeSaveInterval_BareInt(t *testing.T) {
	got := normalizeSaveInterval(10)
	if got != 10*time.Second {
		t.Fatalf("Expected 10s, got %v", got)
	}
}

// TestNormalizeSaveInterval_ExplicitDuration 验证显式时间单位不会被重写。
func TestNormalizeSaveInterval_ExplicitDuration(t *testing.T) {
	got := normalizeSaveInterval(30 * time.Millisecond)
	if got != 30*time.Millisecond {
		t.Fatalf("Expected 30ms, got %v", got)
	}
}

// TestNewCacheManagerWithSeconds 验证秒级构造函数的参数语义。
func TestNewCacheManagerWithSeconds(t *testing.T) {
	cm := NewCacheManagerWithSeconds("test.json", 5, 123, 456)
	defer func() { _ = cm.Close() }()

	if cm.state.saveInterval != 5*time.Second {
		t.Fatalf("Expected 5s, got %v", cm.state.saveInterval)
	}
	if cm.state.maxEntries != 123 {
		t.Fatalf("Expected maxEntries 123, got %d", cm.state.maxEntries)
	}
	if cm.state.maxDataBytes != 456 {
		t.Fatalf("Expected maxDataBytes 456, got %d", cm.state.maxDataBytes)
	}
}

// TestNewCacheManagerWithConfig_DisableAutoSave 验证配置项可以关闭自动保存。
func TestNewCacheManagerWithConfig_DisableAutoSave(t *testing.T) {
	cm := NewCacheManagerWithConfig(Config{
		CacheFile:       "test.json",
		SaveInterval:    time.Second,
		MaxEntries:      10,
		MaxDataBytes:    1024,
		DisableAutoSave: true,
	})
	defer func() { _ = cm.Close() }()

	if err := cm.Set("key", "value"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if cm.state.autoSaveTimer != nil {
		t.Fatal("Expected autoSaveTimer to remain nil when auto save is disabled")
	}
}

// TestCacheManager_EmptyFileDisabled 验证空文件路径场景保持禁用缓存语义。
func TestCacheManager_EmptyFileDisabled(t *testing.T) {
	cm := NewCacheManagerWithConfig(Config{})
	if err := cm.Set("key", "value"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if _, ok := cm.Get("key"); ok {
		t.Fatal("Expected disabled cache to return no data")
	}
	if err := cm.SaveCache(); err != nil {
		t.Fatalf("SaveCache failed: %v", err)
	}
	if err := cm.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

// TestCacheManager_CloseIdempotent 验证 Close 可安全重复调用。
func TestCacheManager_CloseIdempotent(t *testing.T) {
	cm := NewCacheManagerWithConfig(Config{})
	_ = cm.Close()
	_ = cm.Close()
	_ = cm.Close()
}
