package cacher

import (
	"time"

	"github.com/winezer0/xutils/logging"
)

const (
	defaultSaveInterval = 10 * time.Second
	defaultMaxEntries   = 10000
	defaultMaxDataBytes = 100 * 1024 * 1024
)

// Config 定义缓存管理器的可配置项。
type Config struct {
	CacheFile       string
	SaveInterval    time.Duration
	MaxEntries      int
	MaxDataBytes    int64
	DisableAutoSave bool
}

// NewCacheManager 使用默认配置创建缓存管理器。
func NewCacheManager(file string) *CacheManager {
	return NewCacheManagerWithConfig(defaultConfig(file))
}

// NewCacheManagerWithSeconds 使用秒级保存间隔创建缓存管理器。
func NewCacheManagerWithSeconds(file string, saveIntervalSeconds int, maxEntries int, maxDataBytes int64) *CacheManager {
	return NewCacheManagerWithConfig(Config{
		CacheFile:    file,
		SaveInterval: time.Duration(saveIntervalSeconds) * time.Second,
		MaxEntries:   maxEntries,
		MaxDataBytes: maxDataBytes,
	})
}

// NewCacheManagerWithOptions 兼容旧版位置参数构造方式。
func NewCacheManagerWithOptions(file string, saveInterval time.Duration, maxSize int, maxDataSize int64) *CacheManager {
	return NewCacheManagerWithConfig(Config{
		CacheFile:    file,
		SaveInterval: saveInterval,
		MaxEntries:   maxSize,
		MaxDataBytes: maxDataSize,
	})
}

// NewCacheManagerWithConfig 使用显式配置创建缓存管理器。
func NewCacheManagerWithConfig(cfg Config) *CacheManager {
	cfg = normalizeConfig(cfg)
	state := &cacheManagerState{
		cacheFile:       cfg.CacheFile,
		cacheData:       make(map[string]interface{}),
		saveInterval:    cfg.SaveInterval,
		maxEntries:      cfg.MaxEntries,
		maxDataBytes:    cfg.MaxDataBytes,
		disableAutoSave: cfg.DisableAutoSave,
	}
	manager := &CacheManager{state: state}
	if cfg.CacheFile == "" {
		return manager
	}
	if err := manager.LoadCache(); err != nil {
		logging.Warnf("load cache file error: %v", err)
	}
	return manager
}

// defaultConfig 返回缓存管理器的默认配置。
func defaultConfig(file string) Config {
	return Config{
		CacheFile:    file,
		SaveInterval: defaultSaveInterval,
		MaxEntries:   defaultMaxEntries,
		MaxDataBytes: defaultMaxDataBytes,
	}
}

// normalizeConfig 统一填充默认值，并兼容历史调用中裸整数的时间间隔。
func normalizeConfig(cfg Config) Config {
	defaults := defaultConfig(cfg.CacheFile)
	cfg.SaveInterval = normalizeSaveInterval(cfg.SaveInterval)
	if cfg.MaxEntries <= 0 {
		cfg.MaxEntries = defaults.MaxEntries
	}
	if cfg.MaxDataBytes <= 0 {
		cfg.MaxDataBytes = defaults.MaxDataBytes
	}
	return cfg
}

// normalizeSaveInterval 统一处理保存间隔，兼容误传裸整数的场景。
func normalizeSaveInterval(saveInterval time.Duration) time.Duration {
	if saveInterval <= 0 {
		return defaultSaveInterval
	}
	if saveInterval < time.Millisecond {
		return saveInterval * time.Second
	}
	return saveInterval
}
