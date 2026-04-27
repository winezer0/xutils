package cacher

import (
	"time"

	"github.com/winezer0/xutils/logging"
)

// scheduleAutoSave 在数据变更后延迟触发一次自动保存。
func (m *CacheManager) scheduleAutoSave() {
	state := m.getState()
	if state == nil || state.cacheFile == "" || state.saveInterval <= 0 {
		return
	}
	if state.disableAutoSave || state.closed || state.saveInProgress.Load() {
		return
	}

	state.timerMux.Lock()
	defer state.timerMux.Unlock()
	if state.closed {
		return
	}
	if state.autoSaveTimer == nil {
		state.autoSaveTimer = time.AfterFunc(state.saveInterval, m.runAutoSave)
		return
	}
	state.autoSaveTimer.Reset(state.saveInterval)
}

// stopAutoSaveTimer 停止当前自动保存计时器。
func (m *CacheManager) stopAutoSaveTimer() {
	state := m.getState()
	if state == nil {
		return
	}

	state.timerMux.Lock()
	defer state.timerMux.Unlock()
	if state.autoSaveTimer != nil {
		state.autoSaveTimer.Stop()
		state.autoSaveTimer = nil
	}
}

// runAutoSave 执行一次自动保存，并避免与进行中的保存重复落盘。
func (m *CacheManager) runAutoSave() {
	state := m.getState()
	if state == nil {
		return
	}

	state.timerMux.Lock()
	state.autoSaveTimer = nil
	closed := state.closed
	state.timerMux.Unlock()
	if closed {
		return
	}
	if state.saveInProgress.Load() {
		m.scheduleAutoSave()
		return
	}
	if err := m.SaveCache(); err != nil {
		logging.Warnf("save cache error: %v", err)
	}
}
