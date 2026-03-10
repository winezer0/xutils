package cacher

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/winezer0/xutils/logging"
	"github.com/winezer0/xutils/utils"
)

type CacheManager struct {
	cacheFile string
	cacheData map[string]interface{}
	cacheMux  sync.RWMutex
	modified  bool

	maxSize      int
	saveInterval time.Duration
	stopChan     chan struct{}
	waitGroup    sync.WaitGroup
	closeOnce    sync.Once
}

var (
	ErrCacheDisabled     = errors.New("cache disabled")
	ErrCacheFull         = errors.New("cache size limit reached")
	ErrCacheKeyNotFound  = errors.New("cache key not found")
	ErrCacheInvalidValue = errors.New("cache target must be non-nil pointer")
)

func NewCacheManager(file string) *CacheManager {
	return NewCacheManagerWithOptions(file, 10*time.Second, 10000)
}

func NewCacheManagerWithOptions(file string, saveInterval time.Duration, maxSize int) *CacheManager {
	if saveInterval <= 0 {
		saveInterval = 10 * time.Second
	}
	if maxSize <= 0 {
		maxSize = 10000
	}
	m := &CacheManager{
		cacheFile:    file,
		cacheData:    make(map[string]interface{}),
		stopChan:     make(chan struct{}),
		saveInterval: saveInterval,
		maxSize:      maxSize,
	}

	if file == "" {
		return m
	}

	// 加载缓存
	if err := m.LoadCache(); err != nil {
		logging.Warnf("load cache file error: %v", err)
	}

	// 启动定时保存缓存数据
	m.waitGroup.Add(1)
	go m.autoSaveWorker()
	return m
}

func (m *CacheManager) LoadCache() error {
	if m.cacheFile == "" {
		return nil
	}
	m.cacheMux.Lock()
	defer m.cacheMux.Unlock()

	data, err := os.ReadFile(m.cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read cache file error: %w", err)
	}
	if len(data) == 0 {
		m.cacheData = make(map[string]interface{})
		return nil
	}
	if err := json.Unmarshal(data, &m.cacheData); err != nil {
		return fmt.Errorf("parse cache json error: %w", err)
	}
	if m.maxSize > 0 && len(m.cacheData) > m.maxSize {
		for k := range m.cacheData {
			if len(m.cacheData) <= m.maxSize {
				break
			}
			delete(m.cacheData, k)
		}
		return ErrCacheFull
	}
	return nil
}

func (m *CacheManager) SaveCache() error {
	if m.cacheFile == "" {
		return nil
	}
	m.cacheMux.Lock()
	defer m.cacheMux.Unlock()

	if !m.modified {
		return nil
	}

	data, err := utils.ToJSONBytes(m.cacheData)
	if err != nil {
		return fmt.Errorf("serialize cache data error: %w", err)
	}

	if err := utils.EnsureDir(m.cacheFile, true); err != nil {
		return fmt.Errorf("ensure cache dir error: %w", err)
	}

	dir := filepath.Dir(m.cacheFile)
	base := filepath.Base(m.cacheFile)
	tmpFile, err := os.CreateTemp(dir, base+".tmp-")
	if err != nil {
		return fmt.Errorf("create temp cache file error: %w", err)
	}
	tmpName := tmpFile.Name()
	if _, err = tmpFile.Write(data); err != nil {
		tmpFile.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write temp cache file error: %w", err)
	}
	if err = tmpFile.Sync(); err != nil {
		tmpFile.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("sync temp cache file error: %w", err)
	}
	if err = tmpFile.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close temp cache file error: %w", err)
	}
	if utils.FileExists(m.cacheFile) {
		if err := os.Remove(m.cacheFile); err != nil {
			_ = os.Remove(tmpName)
			return fmt.Errorf("remove old cache file error: %w", err)
		}
	}
	if err := os.Rename(tmpName, m.cacheFile); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("replace cache file error: %w", err)
	}
	m.modified = false
	return nil
}

func (m *CacheManager) autoSaveWorker() {
	defer m.waitGroup.Done()
	if m.cacheFile == "" {
		return
	}
	ticker := time.NewTicker(m.saveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			if err := m.SaveCache(); err != nil {
				logging.Warnf("save cache error: %v", err)
			}
		}
	}
}

func (m *CacheManager) Close() error {
	var closeErr error
	m.closeOnce.Do(func() {
		close(m.stopChan)
		m.waitGroup.Wait()
		closeErr = m.SaveCache()
	})
	return closeErr
}

// Clear 清空所有缓存并删除缓存文件
func (m *CacheManager) Clear() error {
	if m.cacheFile == "" {
		return nil
	}
	m.cacheMux.Lock()
	defer m.cacheMux.Unlock()

	// 清空内存
	m.cacheData = make(map[string]interface{})

	// 删除磁盘文件
	if utils.FileExists(m.cacheFile) {
		if err := os.Remove(m.cacheFile); err != nil {
			return fmt.Errorf("failed to remove cache file: %w", err)
		}
	}

	// 关键：重置 modified 状态！ 因为内存和磁盘都已“清空”，视为同步状态
	m.modified = false

	return nil
}

// Set 存任意类型
func (m *CacheManager) Set(key string, value interface{}) error {
	if m.cacheFile == "" {
		return ErrCacheDisabled
	}
	m.cacheMux.Lock()
	defer m.cacheMux.Unlock()

	if cur, exists := m.cacheData[key]; exists {
		if reflect.DeepEqual(cur, value) {
			return nil
		}
	} else if m.maxSize > 0 && len(m.cacheData) >= m.maxSize {
		return ErrCacheFull
	}
	m.cacheData[key] = value
	m.modified = true
	return nil
}

// Del 移除指定key的缓存
func (m *CacheManager) Del(key string) error {
	if m.cacheFile == "" {
		return ErrCacheDisabled
	}
	m.cacheMux.Lock()
	defer m.cacheMux.Unlock()

	if _, exists := m.cacheData[key]; !exists {
		return ErrCacheKeyNotFound
	}
	delete(m.cacheData, key)
	m.modified = true
	return nil
}

// Get 获取值，返回 interface{} 和是否存在
func (m *CacheManager) Get(key string) (interface{}, bool) {
	if m.cacheFile == "" {
		return nil, false
	}
	m.cacheMux.RLock()
	defer m.cacheMux.RUnlock()

	val, exists := m.cacheData[key]
	return val, exists
}

// GetAs 安全地将缓存中的值反序列化为目标类型。
// 注意：target 必须是指针（如 &myStruct）。
func (m *CacheManager) GetAs(key string, target interface{}) (bool, error) {
	if m.cacheFile == "" {
		return false, ErrCacheDisabled
	}
	m.cacheMux.RLock()
	raw, exists := m.cacheData[key]
	m.cacheMux.RUnlock() // 尽早释放读锁

	if !exists {
		return false, ErrCacheKeyNotFound
	}

	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false, ErrCacheInvalidValue
	}

	if reflect.TypeOf(raw) == rv.Elem().Type() {
		rv.Elem().Set(reflect.ValueOf(raw))
		return true, nil
	}

	var jsonBytes []byte
	switch val := raw.(type) {
	case []byte:
		jsonBytes = val
	case json.RawMessage:
		jsonBytes = []byte(val)
	default:
		var err error
		jsonBytes, err = json.Marshal(raw)
		if err != nil {
			return false, fmt.Errorf("marshal cache value error: %w", err)
		}
	}

	if err := json.Unmarshal(jsonBytes, target); err != nil {
		return false, fmt.Errorf("unmarshal cache value error: %w", err)
	}

	return true, nil
}

func (m *CacheManager) GetString(key string) (string, bool) {
	if m.cacheFile == "" {
		return "", false
	}
	m.cacheMux.RLock()
	defer m.cacheMux.RUnlock()
	if v, ok := m.cacheData[key]; ok {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}

func (m *CacheManager) GetBool(key string) (bool, bool) {
	if m.cacheFile == "" {
		return false, false
	}
	m.cacheMux.RLock()
	defer m.cacheMux.RUnlock()
	if v, ok := m.cacheData[key]; ok {
		if b, ok := v.(bool); ok {
			return b, true
		}
	}
	return false, false
}

func (m *CacheManager) GetInt(key string) (int, bool) {
	if m.cacheFile == "" {
		return 0, false
	}
	m.cacheMux.RLock()
	defer m.cacheMux.RUnlock()
	if v, ok := m.cacheData[key]; ok {
		if i, ok := v.(int); ok {
			return i, true
		}
	}
	return 0, false
}
