package cacher

import (
	"encoding/json"
	"fmt"
	"os"
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

	stopChan  chan struct{}
	waitGroup sync.WaitGroup
	closeOnce sync.Once
}

func NewCacheManager(file string) *CacheManager {
	m := &CacheManager{
		cacheFile: file,
		cacheData: make(map[string]interface{}),
		stopChan:  make(chan struct{}),
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

	if err := utils.LoadJSON(m.cacheFile, &m.cacheData); err != nil {
		if os.IsNotExist(err) {
			m.cacheData = make(map[string]interface{})
			return nil
		}
		return fmt.Errorf("failed to parse cache file [%s]: %w", m.cacheFile, err)
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

	if err := utils.SaveJSON(m.cacheFile, m.cacheData); err != nil {
		return fmt.Errorf("save cache data to file occur error: %w", err)
	}
	m.modified = false
	return nil
}

// autoSaveWorker 和 Close 保持不变
func (m *CacheManager) autoSaveWorker() {
	defer m.waitGroup.Done()
	if m.cacheFile == "" {
		return
	}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			m.SaveCache()
			return
		case <-ticker.C:
			m.SaveCache()
		}
	}
}

func (m *CacheManager) Close() {
	m.closeOnce.Do(func() {
		close(m.stopChan)
		m.waitGroup.Wait()
	})
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
func (m *CacheManager) Set(key string, value interface{}) {
	if m.cacheFile == "" {
		return
	}
	m.cacheMux.Lock()
	defer m.cacheMux.Unlock()

	m.cacheData[key] = value
	m.modified = true
}

// Del 移除指定key的缓存
func (m *CacheManager) Del(key string) {
	if m.cacheFile == "" {
		return
	}
	m.cacheMux.Lock()
	defer m.cacheMux.Unlock()

	delete(m.cacheData, key)
	m.modified = true
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
func (m *CacheManager) GetAs(key string, target interface{}) bool {
	if m.cacheFile == "" {
		return false
	}
	// Step 1: 快速读取原始值（最小化锁时间）
	m.cacheMux.RLock()
	raw, exists := m.cacheData[key]
	m.cacheMux.RUnlock() // 尽早释放读锁

	if !exists {
		return false
	}

	// 尝试直接类型断言（更快）
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false
	}

	// 如果 raw 类型和 target 指向的类型一致，直接赋值
	if reflect.TypeOf(raw) == rv.Elem().Type() {
		rv.Elem().Set(reflect.ValueOf(raw))
		return true
	}

	// Step 2: 在无锁环境下进行 JSON 转换（避免阻塞其他读操作）
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		logging.Debugf("GetAs marshal failed for key %s: %v", key, err)
		return false
	}

	if err := json.Unmarshal(jsonBytes, target); err != nil {
		logging.Debugf("GetAs unmarshal failed for key %s: %v", key, err)
		return false
	}

	return true
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
