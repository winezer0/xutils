package cacher

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// CacheManager 表示单个缓存文件对应的管理器。
type CacheManager struct {
	state *cacheManagerState
}

type cacheManagerState struct {
	cacheFile string
	cacheData map[string]interface{}
	cacheMux  sync.RWMutex

	modified       bool
	version        uint64
	saveInProgress atomic.Bool
	saveMux        sync.Mutex

	maxEntries    int
	maxDataBytes  int64
	currentSize   int64
	saveInterval  time.Duration
	timerMux      sync.Mutex
	autoSaveTimer *time.Timer

	disableAutoSave bool
	closed          bool
	closeOnce       sync.Once
}

var cacheFileLocks sync.Map

var (
	ErrCacheDisabled     = errors.New("cache disabled")
	ErrCacheFull         = errors.New("cache size limit reached")
	ErrCacheKeyNotFound  = errors.New("cache key not found")
	ErrCacheInvalidValue = errors.New("cache target must be non-nil pointer")
)

// getState 获取内部状态指针，兼容 CacheManager 的零值与 nil 指针场景。
func (m *CacheManager) getState() *cacheManagerState {
	if m == nil {
		return nil
	}
	return m.state
}

// getCacheFileLock 获取指定缓存文件路径对应的互斥锁。
func getCacheFileLock(cacheFile string) *sync.Mutex {
	if cacheFile == "" {
		return nil
	}
	actual, _ := cacheFileLocks.LoadOrStore(cacheFile, &sync.Mutex{})
	return actual.(*sync.Mutex)
}

// Close 关闭缓存管理器，并在退出前执行一次最终保存。
func (m *CacheManager) Close() error {
	state := m.getState()
	if state == nil {
		return nil
	}

	var closeErr error
	state.closeOnce.Do(func() {
		state.timerMux.Lock()
		state.closed = true
		if state.autoSaveTimer != nil {
			state.autoSaveTimer.Stop()
			state.autoSaveTimer = nil
		}
		state.timerMux.Unlock()
		closeErr = m.SaveCache()
	})
	return closeErr
}
