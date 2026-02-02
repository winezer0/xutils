package progress

import (
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
)

// NewProcessBarByTotalTask 初始化进度显示条
func NewProcessBarByTotalTask(total int64, desc string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(
		total,                                     // 任务总数
		progressbar.OptionSetDescription(desc),    // 进度条描述信息
		progressbar.OptionSetWidth(15),            // 进度条宽度
		progressbar.OptionThrottle(1*time.Second), // 更新频率限制，避免刷新过快
		progressbar.OptionShowCount(),             // 显示已完成/总任务数
		progressbar.OptionShowIts(),               // 显示每秒完成任务数
		progressbar.OptionSetPredictTime(true),    // 启用自动ETA预测
		progressbar.OptionOnCompletion(func() { // 完成时的回调函数
			fmt.Println() // 完成后换行，保持输出整洁
		}),
		progressbar.OptionSpinnerType(14), // 设置进度条动画样式
		progressbar.OptionFullWidth(),     // 启用全屏宽度
		progressbar.OptionSetTheme(progressbar.Theme{ // 自定义进度条样式
			Saucer:        "=", // 已完成部分的填充字符
			SaucerHead:    ">", // 进度条头部字符
			SaucerPadding: " ", // 未完成部分的填充字符
			BarStart:      "[", // 进度条起始字符
			BarEnd:        "]", // 进度条结束字符
		}),
	)
	return bar
}
