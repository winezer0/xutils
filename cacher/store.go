package cacher

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/winezer0/xutils/utils"
)

// Clear 清空所有缓存并删除缓存文件。
func (m *CacheManager) Clear() error {
	state := m.getState()
	if state == nil || state.cacheFile == "" {
		return nil
	}

	state.saveMux.Lock()
	defer state.saveMux.Unlock()

	if fileLock := getCacheFileLock(state.cacheFile); fileLock != nil {
		fileLock.Lock()
		defer fileLock.Unlock()
	}

	state.cacheMux.Lock()
	state.cacheData = make(map[string]interface{})
	state.currentSize = 0
	state.modified = false
	state.version++
	state.cacheMux.Unlock()

	if utils.FileExists(state.cacheFile) {
		if err := os.Remove(state.cacheFile); err != nil {
			return fmt.Errorf("failed to remove cache file: %w", err)
		}
	}

	m.stopAutoSaveTimer()
	return nil
}

// Set 设置指定键的缓存值。
func (m *CacheManager) Set(key string, value interface{}) error {
	state := m.getState()
	if state == nil || state.cacheFile == "" {
		return nil
	}

	valueSize := entrySize(key, value)
	if valueSize == 0 {
		if _, err := json.Marshal(value); err != nil {
			return fmt.Errorf("marshal value error: %w", err)
		}
	}

	shouldSchedule := false
	state.cacheMux.Lock()
	if cur, exists := state.cacheData[key]; exists {
		if reflect.DeepEqual(cur, value) {
			state.cacheMux.Unlock()
			return nil
		}
		state.currentSize -= entrySize(key, cur)
	} else if state.maxEntries > 0 && len(state.cacheData) >= state.maxEntries {
		state.cacheMux.Unlock()
		return ErrCacheFull
	}

	if state.maxDataBytes > 0 && state.currentSize+valueSize > state.maxDataBytes {
		state.cacheMux.Unlock()
		return ErrCacheFull
	}

	state.cacheData[key] = value
	state.currentSize += valueSize
	state.modified = true
	state.version++
	shouldSchedule = !state.disableAutoSave
	state.cacheMux.Unlock()

	if shouldSchedule {
		m.scheduleAutoSave()
	}
	return nil
}

// Del 删除指定键的缓存值。
func (m *CacheManager) Del(key string) error {
	state := m.getState()
	if state == nil || state.cacheFile == "" {
		return nil
	}

	shouldSchedule := false
	state.cacheMux.Lock()
	cur, exists := state.cacheData[key]
	if !exists {
		state.cacheMux.Unlock()
		return ErrCacheKeyNotFound
	}

	state.currentSize -= entrySize(key, cur)
	if state.currentSize < 0 {
		state.currentSize = 0
	}
	delete(state.cacheData, key)
	state.modified = true
	state.version++
	shouldSchedule = !state.disableAutoSave
	state.cacheMux.Unlock()

	if shouldSchedule {
		m.scheduleAutoSave()
	}
	return nil
}

// Get 获取指定键的原始缓存值。
func (m *CacheManager) Get(key string) (data interface{}, exists bool) {
	state := m.getState()
	if state == nil || state.cacheFile == "" {
		return nil, false
	}

	state.cacheMux.RLock()
	defer state.cacheMux.RUnlock()
	value, ok := state.cacheData[key]
	return value, ok
}

// GetAs 将缓存值安全反序列化到目标指针中。
func (m *CacheManager) GetAs(key string, target interface{}) (success bool, hasErr error) {
	state := m.getState()
	if state == nil || state.cacheFile == "" {
		return false, nil
	}

	state.cacheMux.RLock()
	raw, exists := state.cacheData[key]
	state.cacheMux.RUnlock()
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

	jsonBytes, err := toJSONBytes(raw)
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(jsonBytes, target); err != nil {
		return false, fmt.Errorf("unmarshal cache value error: %w", err)
	}
	return true, nil
}

// GetString 获取字符串类型缓存值。
func (m *CacheManager) GetString(key string) (data string, success bool) {
	state := m.getState()
	if state == nil || state.cacheFile == "" {
		return "", false
	}

	state.cacheMux.RLock()
	defer state.cacheMux.RUnlock()
	value, ok := state.cacheData[key]
	if !ok {
		return "", false
	}
	text, ok := value.(string)
	return text, ok
}

// GetBool 获取布尔类型缓存值。
func (m *CacheManager) GetBool(key string) (data bool, success bool) {
	state := m.getState()
	if state == nil || state.cacheFile == "" {
		return false, false
	}

	state.cacheMux.RLock()
	defer state.cacheMux.RUnlock()
	value, ok := state.cacheData[key]
	if !ok {
		return false, false
	}
	boolean, ok := value.(bool)
	return boolean, ok
}

// GetInt 获取整型缓存值。
func (m *CacheManager) GetInt(key string) (data int, success bool) {
	state := m.getState()
	if state == nil || state.cacheFile == "" {
		return 0, false
	}

	state.cacheMux.RLock()
	defer state.cacheMux.RUnlock()
	value, ok := state.cacheData[key]
	if !ok {
		return 0, false
	}
	number, ok := value.(int)
	return number, ok
}

// toJSONBytes 将任意缓存值转换成 JSON 字节切片。
func toJSONBytes(value interface{}) ([]byte, error) {
	switch val := value.(type) {
	case []byte:
		return val, nil
	case json.RawMessage:
		return []byte(val), nil
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("marshal cache value error: %w", err)
		}
		return data, nil
	}
}
