package logging

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewLogConfig 测试日志配置创建
func TestNewLogConfig(t *testing.T) {
	config := NewLogConfig("debug", "test.log", "TLCM")
	if config.Level != "debug" {
		t.Errorf("Expected level 'debug', got '%s'", config.Level)
	}
	if config.LogFile != "test.log" {
		t.Errorf("Expected log file 'test.log', got '%s'", config.LogFile)
	}
	if config.ConsoleFormat != "TLCM" {
		t.Errorf("Expected console format 'TLCM', got '%s'", config.ConsoleFormat)
	}
	if config.MaxSize != 100 {
		t.Errorf("Expected max size 100, got %d", config.MaxSize)
	}
	if config.MaxBackups != 10 {
		t.Errorf("Expected max backups 10, got %d", config.MaxBackups)
	}
	if config.MaxAge != 30 {
		t.Errorf("Expected max age 30, got %d", config.MaxAge)
	}
	if !config.Compress {
		t.Error("Expected compress to be true")
	}
}

// TestNewLogConfigEmpty 测试空配置创建
func TestNewLogConfigEmpty(t *testing.T) {
	config := NewLogConfigEmpty()
	if config.Level != "info" {
		t.Errorf("Expected level 'info', got '%s'", config.Level)
	}
	if config.LogFile != "" {
		t.Errorf("Expected empty log file, got '%s'", config.LogFile)
	}
	if config.ConsoleFormat != "LCM" {
		t.Errorf("Expected console format 'LCM', got '%s'", config.ConsoleFormat)
	}
	if config.MaxSize != 100 {
		t.Errorf("Expected max size 100, got %d", config.MaxSize)
	}
	if config.MaxBackups != 3 {
		t.Errorf("Expected max backups 3, got %d", config.MaxBackups)
	}
	if config.MaxAge != 30 {
		t.Errorf("Expected max age 30, got %d", config.MaxAge)
	}
	if !config.Compress {
		t.Error("Expected compress to be true")
	}
}

// TestCreateLogger 测试创建日志器
func TestCreateLogger(t *testing.T) {
	tmpDir := os.TempDir()
	logFile := filepath.Join(tmpDir, "test_logger.log")
	defer os.Remove(logFile)

	config := NewLogConfig("debug", logFile, "TLCM")
	logger, err := CreateLogger("test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		_ = CloseAll()
	}()

	// 测试日志输出
	logger.Debugf("Debug message")
	logger.Infof("Info message")
	logger.Warnf("Warn message")
	logger.Errorf("Error message")

	// 验证日志器已创建
	if logger == nil {
		t.Error("Expected logger to be created")
	}
}

// TestCreateLoggerDuplicate 测试重复创建日志器
func TestCreateLoggerDuplicate(t *testing.T) {
	config := NewLogConfig("debug", "", "TLCM")
	_, err := CreateLogger("duplicate", config)
	if err != nil {
		t.Fatalf("Failed to create first logger: %v", err)
	}
	defer func() {
		_ = CloseAll()
	}()

	// 尝试创建同名日志器
	_, err = CreateLogger("duplicate", config)
	if err == nil {
		t.Error("Expected error when creating duplicate logger")
	}
}

// TestGetLogger 测试获取日志器
func TestGetLogger(t *testing.T) {
	config := NewLogConfig("debug", "", "TLCM")
	_, err := CreateLogger("gettest", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		_ = CloseAll()
	}()

	logger, exists := GetLogger("gettest")
	if !exists {
		t.Error("Expected logger to exist")
	}
	if logger == nil {
		t.Error("Expected logger to be not nil")
	}

	// 测试获取不存在的日志器
	_, exists = GetLogger("nonexistent")
	if exists {
		t.Error("Expected logger to not exist")
	}
}

// TestCloseAll 测试关闭所有日志器
func TestCloseAll(t *testing.T) {
	config := NewLogConfig("debug", "", "TLCM")
	_, err := CreateLogger("close1", config)
	if err != nil {
		t.Fatalf("Failed to create logger 1: %v", err)
	}

	_, err = CreateLogger("close2", config)
	if err != nil {
		t.Fatalf("Failed to create logger 2: %v", err)
	}

	err = CloseAll()
	if err != nil {
		t.Errorf("Failed to close all loggers: %v", err)
	}

	// 验证日志器已被清空
	_, exists := GetLogger("close1")
	if exists {
		t.Error("Expected logger to be closed")
	}
}

// TestInitLogger 测试初始化默认日志器
func TestInitLogger(t *testing.T) {
	config := NewLogConfig("debug", "", "TLCM")
	err := InitLogger(config)
	if err != nil {
		t.Fatalf("Failed to init logger: %v", err)
	}
	defer func() {
		_ = CloseAll()
	}()

	// 验证默认日志器已创建
	if defaultLogger == nil {
		t.Error("Expected defaultLogger to be initialized")
	}
}

// TestGlobalLogFunctions 测试全局日志函数
func TestGlobalLogFunctions(t *testing.T) {
	defer func() {
		_ = CloseAll()
	}()

	// 测试在未初始化情况下自动初始化
	Infof("Auto-init info message")
	Warnf("Auto-init warn message")
	Errorf("Auto-init error message")
	Debugf("Auto-init debug message")

	Info("Auto-init info")
	Warn("Auto-init warn")
	Error("Auto-init error")
	Debug("Auto-init debug")

	// 验证默认日志器已创建
	if defaultLogger == nil {
		t.Error("Expected defaultLogger to be auto-initialized")
	}
}

// TestSync 测试同步日志
func TestSync(t *testing.T) {
	config := NewLogConfig("debug", "", "TLCM")
	err := InitLogger(config)
	if err != nil {
		t.Fatalf("Failed to init logger: %v", err)
	}
	defer func() {
		_ = CloseAll()
	}()

	err = Sync()
	if err != nil {
		t.Errorf("Failed to sync logger: %v", err)
	}
}

// TestLoggerWithFileOutput 测试带文件输出的日志器
func TestLoggerWithFileOutput(t *testing.T) {
	tmpDir := os.TempDir()
	logFile := filepath.Join(tmpDir, "test_file_output.log")
	defer os.Remove(logFile)

	config := NewLogConfig("debug", logFile, "off")
	logger, err := CreateLogger("filetest", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		_ = CloseAll()
	}()

	logger.Infof("Test file output message")
	_ = logger.Sync()

	// 验证日志文件已创建
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Expected log file to be created")
	}
}

// TestEnsureDefaultLogger 测试自动初始化默认日志器
func TestEnsureDefaultLogger(t *testing.T) {
	// 注意：由于使用了 sync.Once，无法在测试中重置 defaultLogger
	// 这里只验证全局日志函数能正常工作

	// 调用全局日志函数
	Info("Test auto-init message")

	// 验证 defaultLogger 已被初始化
	if defaultLogger == nil {
		t.Error("Expected defaultLogger to be auto-initialized")
	}

	defer func() {
		_ = CloseAll()
	}()
}

// TestLoggerConcurrent 测试日志器并发安全性
func TestLoggerConcurrent(t *testing.T) {
	config := NewLogConfig("debug", "", "TLCM")
	logger, err := CreateLogger("concurrent", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		_ = CloseAll()
	}()

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Infof("Concurrent message %d", id)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
