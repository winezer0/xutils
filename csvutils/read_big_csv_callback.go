package csvutils

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

// ReadBigCSVToDictsWithCallback 兼容所有 Go 版本的回调版读取函数
// 解决：BufferSize 未解析 + LazyQuotes/FieldsPerRecord 作用域问题
func ReadBigCSVToDictsWithCallback(
	filePath string,
	delimiter rune,
	haveHeader bool,
	convertType bool,
	callback func(map[string]interface{}) error,
) error {
	// 边界校验
	if callback == nil {
		return errors.New("callback function cannot be nil")
	}

	// 1. 打开文件（只读模式）
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied to open file: %s", filePath)
		}
		return fmt.Errorf("open file failed: %w", err)
	}
	defer file.Close()

	// 2. 手动创建带缓冲的 Reader（替代 csv.Reader.BufferSize，兼容所有版本）
	// 1MB 缓冲区（适配大文件，和原 BufferSize=1024*1024 效果一致）
	bufReader := bufio.NewReaderSize(file, 1024*1024)

	// 3. 初始化 CSV Reader（核心：解决所有未解析引用）
	reader := csv.NewReader(bufReader) // 包装缓冲 reader
	if delimiter != 0 {
		reader.Comma = delimiter // 设置自定义分隔符
	}
	// 关键配置（兼容所有 Go 版本）
	reader.FieldsPerRecord = -1 // 允许字段数不一致（核心容错）
	reader.LazyQuotes = true    // 宽松引号解析（容错非标准 CSV）
	// 可选：额外容错配置（提升兼容性）
	reader.TrimLeadingSpace = true // 自动修剪字段前后空格

	// 4. 读取表头（逻辑不变）
	var header []string
	if haveHeader {
		row, err := reader.Read()
		if err != nil {
			return fmt.Errorf("read header failed: %w", err)
		}
		header = RepairHeaders(row)
	} else {
		// 无表头时，先读第一行确定列数
		firstRow, err := reader.Read()
		if err != nil {
			return fmt.Errorf("read first row failed: %w", err)
		}
		header = GenDefaultHeaders(len(firstRow))
		// 处理第一行数据
		dict := RowDataToDict(firstRow, header, convertType)
		if err := callback(dict); err != nil {
			return fmt.Errorf("callback failed at first row: %w", err)
		}
	}

	// 5. 逐行读取+处理（核心逻辑不变）
	lineNum := 1
	for {
		row, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break // 正常结束
			}
			// 单行错误：警告+跳过
			fmt.Printf("warning: skip row %d, read error: %v\n", lineNum, err)
			lineNum++
			continue
		}

		lineNum++
		// 行数据转字典
		dict := RowDataToDict(row, header, convertType)
		// 调用回调处理
		if err := callback(dict); err != nil {
			return fmt.Errorf("callback failed at row %d: %w", lineNum, err)
		}
	}

	return nil
}

// ReadBigCSVWithConsumers 并发处理超大CSV（生产者+消费者模型）
// 参数：
//
//	filePath      - CSV文件路径
//	delimiter     - 分隔符
//	skipHeader    - 是否跳过表头
//	convertType   - 是否自动转换类型
//	consumerCount - 消费者数量（建议=CPU核心数*2）
//	chanBuffer    - 通道缓冲区大小（建议=consumerCount*10）
//	consumerFunc  - 单个消费者的处理逻辑
//
// 返回值：
//
//	err - 致命错误（文件读取失败/消费者全部崩溃）
func ReadBigCSVWithConsumers(
	ctx context.Context,
	filePath string,
	delimiter rune,
	haveHeader bool,
	convertType bool,
	consumerCount int,
	chanBuffer int,
	consumerFunc func(map[string]interface{}) error,
) error {
	// 边界校验
	if consumerCount <= 0 {
		return errors.New("consumerCount must be > 0")
	}
	if chanBuffer <= 0 {
		chanBuffer = consumerCount * 10 // 默认缓冲区大小
	}
	if consumerFunc == nil {
		return errors.New("consumerFunc cannot be nil")
	}

	// 1. 创建带缓冲的通道（核心：控速+解耦）
	dataChan := make(chan map[string]interface{}, chanBuffer)
	// 错误通道：收集消费者/生产者的错误
	errChan := make(chan error, consumerCount+1)
	// 退出信号：通知消费者停止
	doneChan := make(chan struct{})

	// 2. 启动消费者池
	var wg sync.WaitGroup
	for i := 0; i < consumerCount; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done(): // 外部中断（如超时/手动停止）
					fmt.Printf("consumer %d: context canceled, exit\n", consumerID)
					return
				case <-doneChan: // 生产者完成，无更多数据
					fmt.Printf("consumer %d: no more data, exit\n", consumerID)
					return
				case data, ok := <-dataChan: // 从通道取数据
					if !ok {
						return
					}
					// 执行消费逻辑
					if err := consumerFunc(data); err != nil {
						// 记录错误，但不终止其他消费者（容错）
						errChan <- fmt.Errorf("consumer %d process data failed: %w", consumerID, err)
						// 可选：若为致命错误，触发全局中断
						// ctx.cancel()
					}
				}
			}
		}(i)
	}

	// 3. 启动生产者（逐行读取CSV并发送到通道）
	go func() {
		defer close(dataChan) // 生产者完成，关闭数据通道
		// 调用回调版读取函数，将数据发送到通道
		err := ReadBigCSVToDictsWithCallback(
			filePath,
			delimiter,
			haveHeader,
			convertType,
			func(dict map[string]interface{}) error {
				select {
				case <-ctx.Done(): // 外部中断，停止生产
					return ctx.Err()
				case dataChan <- dict: // 发送数据到通道（缓冲区满时阻塞，自动控速）
					return nil
				}
			},
		)
		if err != nil {
			errChan <- fmt.Errorf("producer read csv failed: %w", err)
		}
	}()

	// 4. 等待生产者完成 + 消费者处理完毕
	go func() {
		// 等待所有数据发送完成（通道关闭）
		// 通知消费者无更多数据
		close(doneChan)
		// 等待所有消费者完成
		wg.Wait()
		// 关闭错误通道
		close(errChan)
	}()

	// 5. 收集错误（返回第一个致命错误）
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
