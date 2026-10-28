# xutils 日志器可替换方案

范围：整个 xutils。方案已实施：统一 logging 接口、初始化注入、默认静默、日志调用替换和依赖清理均已完成；以下保留设计说明。

## 实施结果

- 实现见 `logging/logging.go`，使用示例见 `README.md`。
- 已移除 zaplogs、zap、multierr、lumberjack 依赖，保留原有构造函数。
- `go test ./...`：135 项测试通过，覆盖 15 个包。
- race 检查：logging、csvwriter、filewriter、poolwriter 通过；cacher 的 `TestCacheManager_AutoSave` 失败。
- 失败原因：`cacher/autosave_test.go:31` 无锁读取 autoSaveTimer，与 `cacher/autosave.go:54` 的后台写入竞争。这两处逻辑本次未修改，属于现有测试问题，未扩大范围修复。
- `git diff --check` 通过。

## 现状

- `go.mod` 直接依赖 `github.com/winezer0/zaplogs v0.0.2`，并间接引入 zap、multierr、lumberjack。
- 5 个包、7 个文件共 8 处直接调用 zaplogs：cacher（2）、csvwriter（1）、filewriter（1）、poolwriter（1）、fileutils（3）。
- `cmdutils/parse_headers.go:17` 使用 fmt.Printf 输出诊断信息，无法随日志器切换。
- poolwriter 在构造时启动 worker，CSVWriter、FileWriter 同样立即启动后台任务，因此日志配置应在创建这些对象之前完成。

## 推荐方案

新增无第三方依赖的公共包 `logging`，由整个 xutils 共用。应用启动时调用 `logging.SetLogger` 指定后端，之后保持配置不变。

最小接口保留现有格式化调用：

```go
type Logger interface {
    Debugf(format string, args ...any)
    Warnf(format string, args ...any)
    Errorf(format string, args ...any)
}
```

提供 `Debugf`、`Warnf`、`Errorf` 包级转发函数，以及 `NopLogger` 和基于函数的适配器 `FuncLogger`。不预先加入当前没有使用的 Info、Fatal、Panic 等方法。

`FuncLogger` 的字段分别接收 Debugf、Warnf、Errorf 函数，空字段视为不输出。这样应用可以传入现有 zaplogs 的包级函数，也可直接传入 zap.SugaredLogger 或 logrus 等满足接口的对象；slog 和标准库 log 通过适配函数接入。适配代码由应用持有，xutils 不导入这些第三方后端。

将现有 8 处调用改为 logging 转发，并将 ParseHeaders 的诊断输出接入 Warnf。保留日志内容、原有错误返回和失败计数，不在这次改动中调整错误处理策略。

## 默认行为与兼容性

建议默认 NopLogger，未配置时不输出；`SetLogger(nil)` 恢复 NopLogger。默认静默会改变当前日志行为，应在发布说明中明确，并提供旧 zaplogs 行为的初始化示例。

保留现有构造函数签名。全局配置适合整个应用使用同一个后端；如果未来需要多个实例分别使用不同日志器，再引入实例配置，本次不扩展所有构造函数。

SetLogger 限定为启动阶段调用；使用同步保护避免并发读写数据竞争，但本次不提供运行中切换的功能承诺。后台调用要求注入的 Logger 支持并发使用。xutils 不负责关闭或 Sync 外部日志器，由创建它的应用负责。

## 实施与验证

1. 新增 logging 接口、空实现、函数适配器和包级转发。
2. 替换 8 处 zaplogs 调用，并接入 ParseHeaders 的诊断输出。
3. 删除 xutils 对 zaplogs 的直接依赖，执行 go mod tidy，核实仅由它引入的间接依赖是否一并移除。
4. 补充 README 初始化示例，说明默认静默、并发要求和生命周期归属。
5. 验证默认/空日志器不崩溃、级别与格式参数正确转发、函数适配、后台任务日志路由；运行 go test ./...，并对新增日志包和受影响后台组件执行适当的 race 检查。

## 另行处理的问题

CSVWriter、FileWriter 的后台写入主要通过日志报告错误，Flush 错误未完整传播。默认静默可能降低这些错误的可见性。建议后续单独设计错误回调或由 Close 返回累计错误，避免与日志器替换混在一次变更中。

## 规范检查

已读取用户提供的 AGENTS 指令及 C:\Users\windows\.codex\AGENTS.md、RTK.md。在仓库根目录、其上级目录及检查的 Codex 根目录中未找到 AGENTS.common.md；仓库不存在 .codegraph 目录。
