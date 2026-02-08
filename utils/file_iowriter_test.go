package utils

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// TestNewFileReceiver_File 测试创建文件接收器
func TestNewFileReceiver_File(t *testing.T) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// 创建FileReceiver
	receiver, err := NewFileReceiver(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to create FileReceiver: %v", err)
	}

	// 验证receiver不为nil
	if receiver == nil {
		t.Fatal("FileReceiver should not be nil")
	}

	// 验证writer不为nil
	if receiver.writer == nil {
		t.Fatal("Writer should not be nil")
	}

	// 关闭文件
	if receiver.file != os.Stdout {
		receiver.file.Close()
	}
}

// TestNewFileReceiver_Stdout 测试创建stdout接收器
func TestNewFileReceiver_Stdout(t *testing.T) {
	// 创建FileReceiver
	receiver, err := NewFileReceiver("")
	if err != nil {
		t.Fatalf("Failed to create FileReceiver: %v", err)
	}

	// 验证receiver不为nil
	if receiver == nil {
		t.Fatal("FileReceiver should not be nil")
	}

	// 验证writer不为nil
	if receiver.writer == nil {
		t.Fatal("Writer should not be nil")
	}

	// 验证文件是stdout
	if receiver.file != os.Stdout {
		t.Fatal("File should be os.Stdout")
	}
}

// TestStart 测试Start方法
func TestStart(t *testing.T) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// 创建FileReceiver
	receiver, err := NewFileReceiver(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to create FileReceiver: %v", err)
	}
	defer func() {
		if receiver.file != os.Stdout {
			receiver.file.Close()
		}
	}()

	// 创建通道
	ch := make(chan string)
	done := make(chan struct{})

	// 启动goroutine写入数据
	go func() {
		ch <- "test line 1"
		ch <- "test line 2"
		ch <- "test line 3"
		close(ch)
	}()

	// 启动接收器
	go receiver.Start(ch, done)

	// 等待完成
	<-done

	// 读取文件内容验证
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	expected := "test line 1\ntest line 2\ntest line 3\n"
	if string(content) != expected {
		t.Errorf("Expected content: %q, got: %q", expected, string(content))
	}
}

// TestStart_Stdout 测试写入到stdout
func TestStart_Stdout(t *testing.T) {
	// 保存原始stdout
	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout
	}()

	// 创建管道捕获stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = writer

	// 创建FileReceiver
	receiver, err := NewFileReceiver("")
	if err != nil {
		t.Fatalf("Failed to create FileReceiver: %v", err)
	}

	// 创建通道
	ch := make(chan string)
	done := make(chan struct{})

	// 启动goroutine写入数据
	go func() {
		ch <- "stdout test line"
		close(ch)
	}()

	// 启动接收器
	go receiver.Start(ch, done)

	// 等待完成
	<-done
	// 关闭writer
	writer.Close()

	// 读取捕获的stdout
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}
	reader.Close()

	expected := "stdout test line\n"
	if buf.String() != expected {
		t.Errorf("Expected stdout: %q, got: %q", expected, buf.String())
	}
}
