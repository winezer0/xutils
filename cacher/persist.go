package cacher

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/winezer0/xutils/utils"
)

var persistCacheFileFunc = persistCacheFile

// LoadCache 从磁盘加载缓存文件。
func (m *CacheManager) LoadCache() error {
	state := m.getState()
	if state == nil || state.cacheFile == "" {
		return nil
	}

	if fileLock := getCacheFileLock(state.cacheFile); fileLock != nil {
		fileLock.Lock()
		defer fileLock.Unlock()
	}

	state.cacheMux.Lock()
	defer state.cacheMux.Unlock()

	data, err := os.ReadFile(state.cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read cache file error: %w", err)
	}
	if len(data) == 0 {
		state.cacheData = make(map[string]interface{})
		state.currentSize = 0
		state.modified = false
		state.version = 0
		return nil
	}

	loaded := make(map[string]interface{})
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("parse cache json error: %w", err)
	}

	currentSize := calculateCacheSize(loaded)
	trimmed := trimLoadedData(loaded, &currentSize, state.maxEntries, state.maxDataBytes)

	state.cacheData = loaded
	state.currentSize = currentSize
	state.modified = trimmed
	if trimmed {
		state.version = 1
		return ErrCacheFull
	}
	state.version = 0
	return nil
}

// SaveCache 将当前缓存安全写入磁盘，并尽量缩短 cacheMux 持有时间。
func (m *CacheManager) SaveCache() error {
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

	snapshot, version, ok := m.prepareSaveSnapshot()
	if !ok {
		return nil
	}

	saveSucceeded := false
	defer func() {
		state.saveInProgress.Store(false)
		if m.completeSave(version, saveSucceeded) {
			m.scheduleAutoSave()
		}
	}()

	data, err := utils.ToJSONBytes(snapshot)
	if err != nil {
		return fmt.Errorf("serialize cache data error: %w", err)
	}
	if state.maxDataBytes > 0 && int64(len(data)) > state.maxDataBytes {
		return fmt.Errorf("serialized cache data exceeds max size limit")
	}
	if err := utils.EnsureDir(state.cacheFile, true); err != nil {
		return fmt.Errorf("ensure cache dir error: %w", err)
	}
	if err := persistCacheFileFunc(state.cacheFile, data); err != nil {
		return err
	}

	saveSucceeded = true
	return nil
}

// prepareSaveSnapshot 在锁内复制一份可持久化快照，并标记保存开始。
func (m *CacheManager) prepareSaveSnapshot() (map[string]interface{}, uint64, bool) {
	state := m.getState()
	state.cacheMux.Lock()
	defer state.cacheMux.Unlock()

	if !state.modified {
		return nil, 0, false
	}

	snapshot := cloneCacheData(state.cacheData)
	version := state.version
	state.saveInProgress.Store(true)
	return snapshot, version, true
}

// completeSave 根据保存结果更新 modified 标记，并决定是否需要重新调度自动保存。
func (m *CacheManager) completeSave(version uint64, saveSucceeded bool) bool {
	state := m.getState()
	state.cacheMux.Lock()
	defer state.cacheMux.Unlock()

	if saveSucceeded && state.version == version {
		state.modified = false
	}
	return state.modified && !state.closed && !state.disableAutoSave
}

// persistCacheFile 持久化缓存文件。
func persistCacheFile(cacheFile string, data []byte) error {
	if runtime.GOOS == "windows" {
		return writeCacheFileDirect(cacheFile, data)
	}
	return writeCacheFileByRename(cacheFile, data)
}

// writeCacheFileDirect 直接覆盖写入缓存文件，并在写入完成后执行 Sync。
func writeCacheFileDirect(cacheFile string, data []byte) error {
	file, err := os.OpenFile(cacheFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open cache file error: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err = file.Write(data); err != nil {
		return fmt.Errorf("write cache file error: %w", err)
	}
	if err = file.Sync(); err != nil {
		return fmt.Errorf("sync cache file error: %w", err)
	}
	if err = file.Close(); err != nil {
		return fmt.Errorf("close cache file error: %w", err)
	}
	return nil
}

// writeCacheFileByRename 先写入临时文件，再通过重命名替换目标文件。
func writeCacheFileByRename(cacheFile string, data []byte) error {
	dir := filepath.Dir(cacheFile)
	base := filepath.Base(cacheFile)
	tmpFile, err := os.CreateTemp(dir, base+".tmp-")
	if err != nil {
		return fmt.Errorf("create temp cache file error: %w", err)
	}
	tmpName := tmpFile.Name()

	if _, err = tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write temp cache file error: %w", err)
	}
	if err = tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("sync temp cache file error: %w", err)
	}
	if err = tmpFile.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close temp cache file error: %w", err)
	}
	if err := os.Rename(tmpName, cacheFile); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("replace cache file error: %w", err)
	}
	return nil
}

// cloneCacheData 返回缓存数据的浅拷贝，用于缩短持锁时间。
func cloneCacheData(source map[string]interface{}) map[string]interface{} {
	target := make(map[string]interface{}, len(source))
	for key, value := range source {
		target[key] = value
	}
	return target
}

// calculateCacheSize 估算缓存数据总大小。
func calculateCacheSize(cacheData map[string]interface{}) int64 {
	size := int64(0)
	for key, value := range cacheData {
		size += entrySize(key, value)
	}
	return size
}

// trimLoadedData 在加载阶段按条目数与字节数限制裁剪数据。
func trimLoadedData(cacheData map[string]interface{}, currentSize *int64, maxEntries int, maxDataBytes int64) bool {
	trimmed := false
	if maxEntries > 0 && len(cacheData) > maxEntries {
		for key := range cacheData {
			if len(cacheData) <= maxEntries {
				break
			}
			*currentSize -= entrySize(key, cacheData[key])
			delete(cacheData, key)
			trimmed = true
		}
	}
	if maxDataBytes > 0 && *currentSize > maxDataBytes {
		for key := range cacheData {
			if *currentSize <= maxDataBytes {
				break
			}
			*currentSize -= entrySize(key, cacheData[key])
			delete(cacheData, key)
			trimmed = true
		}
	}
	if *currentSize < 0 {
		*currentSize = 0
	}
	return trimmed
}

// entrySize 估算单个键值对的字节大小。
func entrySize(key string, value interface{}) int64 {
	data, err := json.Marshal(value)
	if err != nil {
		return 0
	}
	return int64(len(key) + len(data))
}
