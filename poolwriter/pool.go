package poolwriter

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/winezer0/xutils/logging"
	"github.com/winezer0/xutils/utils"
)

// StoreTask 表示通用的写入任务，支持存储响应内容和缓存键写入。
// 注意：WriteRaw=true 时始终使用覆盖模式写入，不支持追加。
type StoreTask struct {
	StorePath    string
	StoreContent []byte
	WriteRaw     bool
	Overwrite    bool
}

// NewStoreTask 创建原始数据写入文件任务
func NewStoreTask(storePath string, content []byte, writeRaw, overwrite bool) StoreTask {
	return StoreTask{
		StorePath:    storePath,
		StoreContent: content,
		WriteRaw:     writeRaw,
		Overwrite:    overwrite,
	}
}

// NewStoreLine 创建文本写入任务，使用追加模式。
func NewStoreLine(cachePath string, cacheKey string) StoreTask {
	return StoreTask{
		StorePath:    cachePath,
		StoreContent: []byte(cacheKey),
		WriteRaw:     false,
		Overwrite:    false,
	}
}

// Pool 表示通用写盘 worker 池。
type Pool struct {
	taskCh    chan []StoreTask
	wg        sync.WaitGroup
	failCount atomic.Int64
}

// NewPool 创建写盘池并启动指定数量 worker。
func NewPool(workerNum int, queueSize int) *Pool {
	if workerNum < 1 {
		workerNum = 1
	}
	if queueSize < 1 {
		queueSize = workerNum * 2
	}

	p := &Pool{
		taskCh: make(chan []StoreTask, queueSize),
	}
	for i := 0; i < workerNum; i++ {
		p.wg.Add(1)
		go p.runWorker()
	}
	return p
}

// Submit 投递单个写盘任务。
func (p *Pool) Submit(task StoreTask) {
	p.taskCh <- []StoreTask{task}
}

// StopAndWait 停止接收新任务并等待所有任务处理完成。
func (p *Pool) StopAndWait() {
	close(p.taskCh)
	p.wg.Wait()
}

// GetFailCount 获取任务失败次数。
func (p *Pool) GetFailCount() int64 {
	return p.failCount.Load()
}

// runWorker 执行写盘任务。
func (p *Pool) runWorker() {
	defer p.wg.Done()
	for tasks := range p.taskCh {
		for _, task := range tasks {
			if err := writeTask(task); err != nil {
				p.failCount.Add(1)
				logging.Errorf("writer task failed: %v, store_path=%s", err, task.StorePath)
			}
		}
	}
}

// writeTask 执行单个任务的写盘动作。
func writeTask(task StoreTask) error {
	if task.StorePath == "" {
		return fmt.Errorf("store path is empty")
	}
	if task.WriteRaw {
		return utils.WriteBytes(task.StorePath, task.StoreContent, task.Overwrite)
	} else {
		return utils.WriteLine(task.StorePath, string(task.StoreContent), task.Overwrite)
	}
}
