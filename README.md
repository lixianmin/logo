# logo

轻量级 Go 日志框架（log + go），除 `github.com/lixianmin/got` 外无其他外部依赖。通过钩子系统实现扩展性。

```bash
go get github.com/lixianmin/logo
```

## 快速开始

```go
package main

import "github.com/lixianmin/logo"

func main() {
    defer logo.GetLogger().(*logo.Logger).Close()

    logo.Info("服务启动")
    logo.Warn("磁盘空间不足, 剩余: %dGB", 2)
    logo.Error("连接失败: %s", err)
}
```

无需任何配置即可使用。logo 在 `init()` 时自动创建全局默认 Logger 并添加 ConsoleHook（默认格式：`FlagDate | FlagTime | FlagShortFile | FlagLevel`，FilterLevel=LevelInfo），开箱即用。

> **注意**：`GetLogger()` 返回 `ILogger` 接口（只有 Debug/Info/Warn/Error），需类型断言 `GetLogger().(*Logger)` 才能调用 AddHook/SetFilterLevel/Close 等方法。

---

## 全局便捷函数

logo 提供包级别的便捷函数，直接使用，无需持有 Logger 实例：

```go
// Printf 风格（包含 % 格式化动词时自动格式化）
logo.Info("连接 %s:%d 成功", host, port)
logo.Error("请求失败, status=%d, err=%q", code, err)

// 空格拼接风格（不含 % 时自动用空格拼接参数）
logo.Info("插件已启动:", pluginName, "路径:", entryPath)

// JSON 结构化日志（交替传入 key-value 对）
logo.JsonI("sessionId", id, "elapsed", time.Since(start))
logo.JsonD("command", cmd, "output", truncated)
logo.JsonW("timeout", "id", reqId, "seconds", timeout)
logo.JsonE("err", err, "url", url)
```

| 函数 | 级别 | 用途 |
|------|------|------|
| `Debug` / `JsonD` | DEBUG | 开发调试信息 |
| `Info` / `JsonI` | INFO | 常规运行信息 |
| `Warn` / `JsonW` | WARN | 警告（不影响运行但需关注） |
| `Error` / `JsonE` | ERROR | 错误（需要处理） |

> **提示**：`Info("hello", name)` 等价于 `Info("hello %v", name)`。只要 format 字符串中不含 `%`（且不含 `%%`），logo 会自动补充空格分隔的 `%v` 占位符。

---

## 添加文件日志

默认只输出到控制台。添加 `RollingFileHook` 即可同时写文件，支持按大小自动归档：

```go
package main

import (
    "github.com/lixianmin/logo"
)

func main() {
    var logger = logo.GetLogger().(*logo.Logger)
    defer logger.Close()

    const flag = logo.FlagDate | logo.FlagTime | logo.FlagShortFile | logo.FlagLevel
    var fileHook = logo.NewRollingFileHook(
        logo.WithHookFlag(flag),
        logo.WithDirName("logs"),
        logo.WithMaxFileSize(10*1024*1024), // 10MB
    )
    logger.AddHook(fileHook)

    logo.Info("这条日志会同时输出到控制台和文件")
}
```

生成的文件结构：
```
logs/
├── info.log          # INFO 级别日志
├── warn.log          # WARN 级别日志
├── error.log         # ERROR 级别日志
└── archive/          # 超过大小限制后自动归档到此目录
    ├── info2026-04-12_001.log
    └── warn2026-04-12_001.log
```

### RollingFileHook 选项

| 选项 | 默认值 | 说明 |
|------|--------|------|
| `WithHookFlag(flag)` | `FlagNone` | 日志格式标志位 |
| `WithHookFilterLevel(level)` | `LevelInfo` | 最低输出级别 |
| `WithDirName(dir)` | `"logs"` | 日志目录 |
| `WithFileNamePrefix(prefix)` | `""` | 文件名前缀，如 `"myapp-"` → `myapp-info.log` |
| `WithMaxFileSize(size)` | `10MB` | 单文件最大字节数，超过后自动归档 |
| `WithExpireTime(d)` | `7天` | 归档文件保留时长 |
| `WithCheckRollingInterval(n)` | `1024` | 每写入 n 条检查一次文件大小 |

> **注意**：RollingFileHook 会按级别分流写入，且高级别日志会级联写入所有低级别文件。默认 FilterLevel=LevelInfo 时，一条 ERROR 日志会同时写入 error.log、warn.log、info.log。

---

## 自定义控制台日志

全局默认 Logger 已自带一个 ConsoleHook（FilterLevel=LevelInfo）。如需自定义：

```go
// 方式1：设置全局 Logger 的过滤级别（会影响所有已注册 Hook）
logo.GetLogger().(*logo.Logger).SetFilterLevel(logo.LevelWarn)

// 方式2：创建自定义 Logger 并替换全局默认
var logger = logo.NewLogger()
logger.SetFuncCallDepth(5)

const flag = logo.FlagDate | logo.FlagTime | logo.FlagShortFile | logo.FlagLevel
var console = logo.NewConsoleHook(
    logo.WithFlag(flag),
    logo.WithFilterLevel(logo.LevelDebug),
)
logger.AddHook(console)

logo.SetLogger(logger)
defer logger.Close()
```

> **提示**：`SetFuncCallDepth(5)` 用于修正调用栈深度，使日志中显示的文件名和行号指向实际调用位置（而非经过多层包装后的位置）。默认全局 Logger 已设置为 5。

> **提示**：`Logger.SetFilterLevel(level)` 会同步更新所有已注册 Hook 的 FilterLevel。

---

## 异步写入

默认每次日志调用都是同步落盘。对性能敏感的场景可开启异步写入：

```go
var logger = logo.NewLogger()
logger.AddFlag(logo.LogAsyncWrite)

// ... 添加 Hook ...

// 开启异步写入后，必须在退出时主动 Close 或 Flush，否则可能丢失日志
defer logger.Close()

// 可选：调整异步消息缓冲区大小（默认 4096）
// var logger = logo.NewLogger(logo.WithBufferSize(8192))
```

### 资源清理最佳实践

```go
// 场景1：使用全局默认 Logger（推荐）
func main() {
    defer logo.GetLogger().(*logo.Logger).Close()
    // ... 业务代码 ...
}

// 场景2：信号退出时主动关闭
func main() {
    var logger = logo.GetLogger().(*logo.Logger)

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        logo.Info("收到退出信号，正在关闭...")
        logger.Close()
        os.Exit(0)
    }()

    // ... 业务代码 ...
}
```

> **注意**：未开启 `LogAsyncWrite` 时，每次日志调用会自动 Flush，无需手动处理。开启后务必在退出前调用 `Close()`（内部会 Flush 所有未写入的消息）。

---

## 消息平台 Hook

### 钉钉（ding.TalkHook）

将日志推送到钉钉群机器人：

```go
import (
    "github.com/lixianmin/logo"
    "github.com/lixianmin/logo/ding"
)

func main() {
    var logger = logo.GetLogger().(*logo.Logger)
    defer logger.Close()

    // 创建钉钉 Talker（token 为钉钉机器人的 access_token）
    var talk = ding.NewTalk("我的服务", "your-dingtalk-token")
    defer talk.Close()

    // 创建 Hook 并添加到 Logger
    var hook = ding.NewHook(talk,
        ding.WithFilterLevel(logo.LevelError), // 只推送 ERROR 级别
    )
    logger.AddHook(hook)

    logo.Error("这条日志会同时输出到控制台和钉钉群")
}
```

钉钉机器人发送频率限制：20 条/分钟。Talk 内置令牌桶限流，无需额外处理。

### 飞书（lark.Lark）

将日志推送到飞书群机器人：

```go
import (
    "github.com/lixianmin/logo"
    "github.com/lixianmin/logo/ding"
    "github.com/lixianmin/logo/lark"
)

func main() {
    var logger = logo.GetLogger().(*logo.Logger)
    defer logger.Close()

    // 创建飞书 Lark（token 为飞书机器人的 hook token）
    var lk = lark.NewLark("我的服务", "your-lark-token")
    defer lk.Close()

    // lark.Lark 实现了 ding.Talker 接口，直接用 ding.NewHook 包装
    var hook = ding.NewHook(lk,
        ding.WithFilterLevel(logo.LevelError),
    )
    logger.AddHook(hook)

    logo.Error("这条日志会推送到飞书群")
}
```

飞书机器人发送频率限制：200 条/分钟。Lark 同样内置令牌桶限流。

> **提示**：`lark.Lark` 实现了 `ding.Talker` 接口，因此用 `ding.NewHook()` 统一包装。`ding.Talk` 和 `lark.Lark` 都支持 `PostMessage`（限流队列发送）和 `SendMessage`（立即发送，绕过队列，用于紧急场景）两种方式。

---

## 日志级别

| 常量 | 值 | 说明 |
|------|----|------|
| `LevelDebug` | 1 | 开发调试 |
| `LevelInfo` | 2 | 常规信息（默认） |
| `LevelWarn` | 3 | 警告 |
| `LevelError` | 4 | 错误 |

级别过滤规则：只输出 `>= FilterLevel` 的日志。例如设为 `LevelWarn` 后，`Debug` 和 `Info` 会被丢弃。

## 格式标志位

通过位运算组合，控制每条日志的输出格式：

```go
const flag = logo.FlagDate | logo.FlagTime | logo.FlagShortFile | logo.FlagLevel
```

| 常量 | 输出示例 | 说明 |
|------|----------|------|
| `FlagDate` | `2026-04-12` | 日期 |
| `FlagTime` | `14:30:00` | 时间 |
| `FlagLongFile` | `pkg/server/handler.go:42` | 完整文件路径和行号 |
| `FlagShortFile` | `handler.go:42` | 文件名和行号 |
| `FlagLevel` | `[I]` | 级别标识（控制台带颜色） |

**推荐组合**：
```go
const flag = logo.FlagDate | logo.FlagTime | logo.FlagShortFile | logo.FlagLevel
// 输出示例: 2026-04-12 14:30:00 handler.go:42 [I] 服务启动
```

## 其他配置

```go
// 控制 full stack trace 的捕获阈值（默认 LevelError，即只有 ERROR 级别才捕获完整调用栈）
logger.SetStackLevel(logo.LevelWarn)
```
