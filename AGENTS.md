# AGENTS.md

本文件为 AI 编程助手提供项目指导。

## 核心原则（最高优先级）

明确导入项目宪法，确保 AI 在思考任何问题前，都已经加载核心原则。

@./notes/00.constitution.md

## 角色与使命

你是一个资深的程序员，你的职责是协助我完成从需求 → 架构设计 → 任务清单 → 编码实现 → 验收 → 项目上线的全流程开发。

你的所有行动都必须严格遵守上面导入的项目宪法。

## 项目概述

**logo** 是一个轻量级 Go 日志框架（log + go），提供基础日志功能，除 `github.com/lixianmin/got` 外无其他外部依赖。框架通过钩子系统实现扩展性。

## 架构设计

### 核心组件

- **Logger** (`logger.go`)：主日志器实现，包含异步消息处理、可配置过滤级别和钩子管理
- **Message** (`message.go`)：内部消息结构，包含文本、级别和调用栈帧
- **Hook System** (`ihook.go`)：用于扩展日志输出到不同目标的接口
- **Default Logger** (`init.go`)：全局默认日志器及便捷函数

### 钩子实现

- **ConsoleHook** (`console_hook.go`)：控制台输出，支持颜色
- **RollingFileHook** (`rolling_file_hook.go`)：文件输出，支持按大小轮转
- **DingTalk Hook** (`ding/`)：发送日志到钉钉消息平台
- **Lark Hook** (`lark/`)：发送日志到飞书消息平台

### 核心设计模式

- 通过 channel 实现异步消息处理以提升性能
- 基于标志位的日志格式配置
- 通过钩子实现插件化架构
- 全局默认日志器提供便捷使用方式

## 构建、检查和测试命令

```bash
# 构建模块
go build ./...

# 运行所有测试
go test ./...

# 运行测试（详细输出）
go test -v ./...

# 按名称运行单个测试
go test -run TestConsoleHook ./...

# 运行单个测试文件
go test -run TestName ./path/to/package

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# 运行示例
go run examples/main.go
```

## 关键配置

### 日志级别
- `LevelDebug` (1)
- `LevelInfo` (2)
- `LevelWarn` (3)
- `LevelError` (4)

### 格式标志位
- `FlagDate`：显示日期 (1998-10-29)
- `FlagTime`：显示时间 (12:24:00)
- `FlagLongFile`：完整文件路径 (i/am/the/path/file.go:12)
- `FlagShortFile`：简短文件名 (file.go:34)
- `FlagLevel`：显示级别指示器 [D], [I], [W], [E]

### 日志器标志位
- `LogAsyncWrite`：启用异步写入以提升性能（需要显式调用 Flush/Close）

## 代码风格指南

### 文件头

每个源文件以标准化头部开始：

```go
package packagename

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
```

### 导入组织

按以下顺序分组导入，组间用空行分隔：
1. 标准库包
2. 第三方包 (github.com/lixianmin/got/*)
3. 本地包 (github.com/lixianmin/logo/*)

```go
import (
    "fmt"
    "os"
    "sync"

    "github.com/lixianmin/got/loom"
    "github.com/lixianmin/got/convert"

    "github.com/lixianmin/logo/tools"
)
```

### 命名规范

- **公开类型**：PascalCase（`Logger`、`ConsoleHook`、`RollingFileHook`）
- **私有类型**：camelCase（`loggerFetus`、`messageChan`）
- **接口**：`I` 前缀（`IHook`、`ILogger`）
- **Args 结构体**：`TypeNameArgs` 模式（`ConsoleHookArgs`、`RollingFileHookArgs`）
- **Options 模式**：`TypeNameOption` 函数类型（`LoggerOption`、`ConsoleHookOption`）
- **常量**：公开的使用 PascalCase（`LevelDebug`、`FlagDate`），私有的使用 camelCase
- **接收者名称**：方法接收者统一使用 `my`（不要使用 `this` 或 `self`）

```go
func (my *Logger) Info(format string, args ...any) {
    // 实现
}
```

### 结构体定义

- 相关字段用空行分组
- 适当使用嵌入接口/类型
- 结构体类型放在构造函数之前

```go
type RollingFileHook struct {
    wc        loom.WaitClose
    args      rollingFileHookOptions
    formatter *MessageFormatter

    files [LevelMax]struct {
        *os.File
        checkRollingCount int64
    }
}
```

### 构造函数

- 使用 `NewTypeName` 命名模式
- 通过 `...Option` 参数接受可选配置
- 在应用选项前初始化默认值
- 返回指针类型

```go
func NewLogger(opts ...LoggerOption) *Logger {
    var options = loggerOptions{
        BufferSize: 4096,
    }

    for _, opt := range opts {
        opt(&options)
    }

    var my = &Logger{
        // 初始化
    }

    return my
}

func NewConsoleHookWithOptions(opts ...ConsoleHookOption) *ConsoleHook {
    var options = consoleHookOptions{
        FilterLevel: LevelInfo,
    }

    for _, opt := range opts {
        opt(&options)
    }

    var my = &ConsoleHook{
        args:      options,
        formatter: NewMessageFormatter(options.Flag, levelHintsConsole),
    }

    return my
}
```

### 错误处理

- 适当时使用 `_` 显式忽略错误
- 对于非关键错误，使用 `checkPrintError(err)` 模式
- 当调用者需要处理错误时，从公开方法返回错误

```go
// 显式忽略错误
_ = os.MkdirAll(args.DirName, os.ModePerm)

// 非关键错误日志
var err = my.openLogFile(level)
checkPrintError(err)

// 返回错误供调用者处理
func (my *Logger) Close() error {
    return my.wc.Close(func() error {
        my.Flush()
        return nil
    })
}
```

### 注释和文档

- 注释主要使用中文
- 不使用 Go 风格的导出函数/类型文档注释（与标准 Go 规范不同）
- 对复杂逻辑使用行内注释解释
- 暂时注释掉未使用的代码而不是删除，以备后用

### 测试规范

- 测试文件：`*_test.go` 与源文件放在一起
- 测试函数：`TestFunctionName` 模式
- 使用 `t.Fatal()` 表示测试失败（不带消息）
- 测试中使用 `defer` 进行清理

```go
func TestConsoleHook(t *testing.T) {
    var l = NewLogger()
    defer l.Close()

    // 测试实现
    if condition {
        t.Fatal()
    }
}
```

### 并发模式

- 使用 channel 进行异步消息处理
- 使用 `sync/atomic` 进行简单原子操作
- 使用 `sync.Pool` 进行缓冲区复用
- 通过 `loom.WaitClose` 实现带有正确清理的 `Close()` 方法

```go
func (my *Logger) goLoop(later loom.Later) {
    var closeChan = my.wc.C()

    for {
        select {
        case message := <-my.messageChan:
            fetus.WriteMessage(message)
        case <-closeChan:
            return
        }
    }
}
```

### 钩子系统

钩子实现 `IHook` 接口：

```go
type IHook interface {
    SetFilterLevel(level int)
    Write(message Message)
}
```

- 处理前检查过滤级别
- 每个钩子管理自己的输出目标
- 钩子应该是线程安全的

### 必须遵循的关键模式

1. **Options 模式**：用于可配置的构造函数（推荐使用 `NewXxxWithOptions(opts ...Option)` 形式）
2. **Args 结构体**：用于复杂构造函数参数（保留以兼容旧代码）
3. **异步处理**：通过 channel 和 `goLoop` 方法实现
4. **资源清理**：使用 `Close()` 方法和 `defer`
5. **空值检查**：使用 `std.IsNil()` 进行接口空值检测

## 重要注意事项

- 框架依赖作者的 `got` 工具库提供并发原语和工具
- 异步日志需要通过 `Close()` 或 `Flush()` 正确清理以避免数据丢失
- 钩子实现应该是线程安全的，因为它们可能被并发调用
- 消息平台钩子（钉钉、飞书）包含速率限制和消息分割逻辑
