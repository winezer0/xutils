package timeutils

// 常量定义：Go 的布局模板 (Layout)
const (
	LayoutDate     = "2006-01-02"          // 2026-03-12
	LayoutDateTime = "2006-01-02 15:04:05" // 2026-03-12 13:03:00
	LayoutFileSafe = "20060102_150405"     // 20260312_130300 (适合文件名)
	LayoutCN       = "2006年01月02日 15时04分"
	LayoutISO8601  = "2006-01-02T15:04:05Z07:00"
)
