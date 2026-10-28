# xutils
my golang utils package

## 日志器

默认不输出日志。应用应在创建缓存、写入器或任务池之前调用
`logging.SetLogger`，初始化后保持配置不变。原有构造函数无需调整。

例如，使用 Go 标准库日志：

```go
import (
    "log"
    "github.com/winezer0/xutils/logging"
)

func init() {
    logging.SetLogger(logging.FuncLogger{
        WarnFunc: log.Printf,
        ErrorFunc: log.Printf,
    })
}
```

使用 zaplogs 时，由应用自行安装并导入它，再配置：

```go
logging.SetLogger(logging.FuncLogger{
    DebugFunc: zaplogs.Debugf,
    WarnFunc: zaplogs.Warnf,
    ErrorFunc: zaplogs.Errorf,
})
```

实现 `Debugf`、`Warnf`、`Errorf` 方法的对象（例如 zap 的 SugaredLogger）
可以直接传入。slog 可通过函数适配器配合 `fmt.Sprintf` 转换格式化参数。
空的适配函数不输出对应级别；`logging.SetLogger(nil)` 恢复静默。
请传入 nil 本身，不要传入包含 nil 指针的 Logger 接口。

日志器必须支持并发调用，刷新和关闭由应用负责。后台写入失败等诊断
也遵循默认静默设置；需要观察这些错误时，请显式配置日志器。
