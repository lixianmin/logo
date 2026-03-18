# Hook Options 统一重构实施计划

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 统一 Hook 的 Options 模式，创建通用的 HookConfig 结构体和 WithFlag/WithFilterLevel 函数，让所有 Hook 类型共享相同的配置方式

**Architecture:** 创建 HookConfig 公共结构体，包含 Flag 和 FilterLevel 字段。各 Hook 的 options 结构体嵌入 HookConfig。通用的 WithFlag 和 WithFilterLevel 函数操作 HookConfig 指针，实现跨 Hook 类型复用

**Tech Stack:** Go 1.21+, table-driven tests

---

## 文件结构

**新建文件:**
- `hook_config.go` - 公共 HookConfig 结构体和通用 Option 函数
- `hook_config_test.go` - HookConfig 的单元测试
- `rolling_file_hook_option.go` - RollingFileHook 特有的 Option 函数

**修改文件:**
- `console_hook.go:14-56` - 修改 consoleHookOptions 定义和使用
- `rolling_file_hook.go:25-244` - 修改 rollingFileHookOptions 定义和使用
- `ding/hook_option.go:12-24` - 修改 hookOptions 定义和使用
- `ding/hook.go:18-88` - 修改 TalkHook 以使用嵌入的 HookConfig
- `init.go:21` - 更新使用通用 WithFlag
- `logger_test.go:28,46-50,71-74` - 更新测试使用通用 Option 函数

**删除文件:**
- 无（旧函数 WithHookFlag/WithHookFilterLevel 从未实现，无需删除）

---

## Chunk 1: 创建 HookConfig 基础设施

### Task 1: 创建 HookConfig 测试

**Files:**
- Create: `hook_config_test.go`

- [ ] **Step 1: 编写失败的测试 - 测试 HookConfig 结构体**

```go
package logo

import (
	"testing"
)

/********************************************************************
created:    2026-03-18
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestHookConfigDefaults(t *testing.T) {
	var config HookConfig
	
	if config.Flag != 0 {
		t.Errorf("expected default Flag=0, got %d", config.Flag)
	}
	
	if config.FilterLevel != 0 {
		t.Errorf("expected default FilterLevel=0, got %d", config.FilterLevel)
	}
}
```

- [ ] **Step 2: 运行测试验证失败**

Run: `go test -run TestHookConfigDefaults -v`
Expected: FAIL with "undefined: HookConfig"

- [ ] **Step 3: 创建 hook_config.go，定义 HookConfig 结构体**

```go
package logo

/********************************************************************
created:    2026-03-18
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type HookConfig struct {
	Flag        int
	FilterLevel int
}
```

- [ ] **Step 4: 运行测试验证通过**

Run: `go test -run TestHookConfigDefaults -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add hook_config.go hook_config_test.go
git commit -m "feat: add HookConfig struct with default values test"
```

---

### Task 2: 创建通用 WithFlag Option 函数

**Files:**
- Modify: `hook_config.go:9-13`
- Modify: `hook_config_test.go:21-30`

- [ ] **Step 1: 编写失败的测试 - 测试 WithFlag 函数**

在 `hook_config_test.go` 中添加：

```go
func TestWithFlag(t *testing.T) {
	var tests = []struct {
		name      string
		flag      int
		wantFlag  int
	}{
		{"zero flag", 0, 0},
		{"single flag", FlagDate, FlagDate},
		{"combined flags", FlagDate | FlagTime | FlagShortFile, FlagDate | FlagTime | FlagShortFile},
		{"all flags", FlagDate | FlagTime | FlagLongFile | FlagShortFile | FlagLevel, FlagDate | FlagTime | FlagLongFile | FlagShortFile | FlagLevel},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config HookConfig
			var opt = WithFlag(tt.flag)
			opt(&config)
			
			if config.Flag != tt.wantFlag {
				t.Errorf("WithFlag(%d): got Flag=%d, want %d", tt.flag, config.Flag, tt.wantFlag)
			}
		})
	}
}
```

- [ ] **Step 2: 运行测试验证失败**

Run: `go test -run TestWithFlag -v`
Expected: FAIL with "undefined: WithFlag"

- [ ] **Step 3: 实现 WithFlag 函数**

在 `hook_config.go` 中添加：

```go
type HookOption func(*HookConfig)

func WithFlag(flag int) HookOption {
	return func(config *HookConfig) {
		config.Flag = flag
	}
}
```

- [ ] **Step 4: 运行测试验证通过**

Run: `go test -run TestWithFlag -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add hook_config.go hook_config_test.go
git commit -m "feat: add WithFlag option function with table-driven tests"
```

---

### Task 3: 创建通用 WithFilterLevel Option 函数

**Files:**
- Modify: `hook_config.go:15-18`
- Modify: `hook_config_test.go:53-82`

- [ ] **Step 1: 编写失败的测试 - 测试 WithFilterLevel 函数**

在 `hook_config_test.go` 中添加：

```go
func TestWithFilterLevel(t *testing.T) {
	var tests = []struct {
		name            string
		level           int
		wantFilterLevel int
	}{
		{"level debug", LevelDebug, LevelDebug},
		{"level info", LevelInfo, LevelInfo},
		{"level warn", LevelWarn, LevelWarn},
		{"level error", LevelError, LevelError},
		{"level none", LevelNone, LevelNone},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config HookConfig
			var opt = WithFilterLevel(tt.level)
			opt(&config)
			
			if config.FilterLevel != tt.wantFilterLevel {
				t.Errorf("WithFilterLevel(%d): got FilterLevel=%d, want %d", 
					tt.level, config.FilterLevel, tt.wantFilterLevel)
			}
		})
	}
}
```

- [ ] **Step 2: 运行测试验证失败**

Run: `go test -run TestWithFilterLevel -v`
Expected: FAIL with "undefined: WithFilterLevel"

- [ ] **Step 3: 实现 WithFilterLevel 函数**

在 `hook_config.go` 中添加：

```go
func WithFilterLevel(level int) HookOption {
	return func(config *HookConfig) {
		if level > LevelNone {
			config.FilterLevel = level
		}
	}
}
```

- [ ] **Step 4: 运行测试验证通过**

Run: `go test -run TestWithFilterLevel -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add hook_config.go hook_config_test.go
git commit -m "feat: add WithFilterLevel option function with table-driven tests"
```

---

### Task 4: 测试 HookConfig 嵌入到自定义结构体

**Files:**
- Modify: `hook_config_test.go:85-120`

- [ ] **Step 1: 编写测试 - 验证 HookConfig 可以被嵌入并使用通用 Option**

在 `hook_config_test.go` 中添加：

```go
func TestHookConfigEmbedding(t *testing.T) {
	type customHookOptions struct {
		HookConfig
		CustomField string
	}
	
	var tests = []struct {
		name          string
		opts          []HookOption
		wantFlag      int
		wantLevel     int
		wantCustom    string
	}{
		{
			name:       "no options",
			opts:       nil,
			wantFlag:   0,
			wantLevel:  0,
			wantCustom: "",
		},
		{
			name:       "with flag only",
			opts:       []HookOption{WithFlag(FlagDate | FlagTime)},
			wantFlag:   FlagDate | FlagTime,
			wantLevel:  0,
			wantCustom: "",
		},
		{
			name:       "with filter level only",
			opts:       []HookOption{WithFilterLevel(LevelWarn)},
			wantFlag:   0,
			wantLevel:  LevelWarn,
			wantCustom: "",
		},
		{
			name:       "with both options",
			opts:       []HookOption{WithFlag(FlagShortFile), WithFilterLevel(LevelError)},
			wantFlag:   FlagShortFile,
			wantLevel:  LevelError,
			wantCustom: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options = customHookOptions{
				CustomField: tt.wantCustom,
			}
			
			for _, opt := range tt.opts {
				opt(&options.HookConfig)
			}
			
			if options.Flag != tt.wantFlag {
				t.Errorf("got Flag=%d, want %d", options.Flag, tt.wantFlag)
			}
			if options.FilterLevel != tt.wantFilterLevel {
				t.Errorf("got FilterLevel=%d, want %d", options.FilterLevel, tt.wantLevel)
			}
			if options.CustomField != tt.wantCustom {
				t.Errorf("got CustomField=%s, want %s", options.CustomField, tt.wantCustom)
			}
		})
	}
}
```

- [ ] **Step 2: 运行测试验证通过（应该立即通过，因为已有 HookConfig）**

Run: `go test -run TestHookConfigEmbedding -v`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add hook_config_test.go
git commit -m "test: add HookConfig embedding test to verify option pattern works"
```

---

## Chunk 2: 重构 ConsoleHook

### Task 5: 重构 consoleHookOptions 使用嵌入的 HookConfig

**Files:**
- Modify: `console_hook.go:14-56`
- Modify: `logger_test.go:28`

- [ ] **Step 1: 编写测试 - 验证 ConsoleHook 可以使用通用 Option**

测试已存在于 `logger_test.go:22-37`，先验证当前状态：

- [ ] **Step 2: 修改 console_hook.go，定义 consoleHookOptions 嵌入 HookConfig**

```go
type consoleHookOptions struct {
	HookConfig
}

type ConsoleHookOption = HookOption
```

- [ ] **Step 3: 运行 ConsoleHook 测试**

Run: `go test -run TestConsoleHook -v`
Expected: PASS

- [ ] **Step 4: 运行所有测试验证没有破坏任何功能**

Run: `go test -v`
Expected: PASS（除了可能失败的 RollingFileHook 测试）

- [ ] **Step 5: Commit**

```bash
git add console_hook.go
git commit -m "refactor: make consoleHookOptions embed HookConfig"
```

---

### Task 6: 验证 init.go 使用通用 WithFlag

**Files:**
- `init.go:21`

- [ ] **Step 1: 验证 init.go 已经使用通用 WithFlag**

检查 `init.go:21` 已经使用 `WithFlag(flag)`，无需修改

- [ ] **Step 2: 运行构建验证**

Run: `go build`
Expected: 成功，无错误

- [ ] **Step 3: 运行测试验证 init.go 功能正常**

Run: `go test -run TestAutoFlush -v`
Expected: PASS

---

## Chunk 3: 重构 RollingFileHook

### Task 7: 重构 rollingFileHookOptions 使用嵌入的 HookConfig

**Files:**
- Create: `rolling_file_hook_option.go`
- Modify: `rolling_file_hook.go:25-244`
- Modify: `logger_test.go:46-50,71-74`

- [ ] **Step 1: 创建 rolling_file_hook_option.go，定义特有 Option 函数**

```go
package logo

import (
	"time"
)

/********************************************************************
created:    2026-03-18
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type rollingFileHookOptions struct {
	HookConfig
	DirName              string
	FileNamePrefix       string
	MaxFileSize          int64
	ExpireTime           time.Duration
	CheckRollingInterval int64
}

type RollingFileHookOption func(*rollingFileHookOptions)

func WithDirName(dirName string) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		if dirName != "" {
			options.DirName = dirName
		}
	}
}

func WithFileNamePrefix(prefix string) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		options.FileNamePrefix = prefix
	}
}

func WithMaxFileSize(size int64) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		if size > 0 {
			options.MaxFileSize = size
		}
	}
}

func WithExpireTime(duration time.Duration) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		if duration > 0 {
			options.ExpireTime = duration
		}
	}
}

func WithCheckRollingInterval(interval int64) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		if interval > 0 {
			options.CheckRollingInterval = interval
		}
	}
}
```

- [ ] **Step 2: 修改 rolling_file_hook.go，更新默认值和初始化逻辑**

修改 `NewRollingFileHook` 函数：

```go
func NewRollingFileHook(opts ...RollingFileHookOption) *RollingFileHook {
	var options = rollingFileHookOptions{
		HookConfig: HookConfig{
			FilterLevel: LevelInfo,
		},
		DirName:              "logs",
		MaxFileSize:          10 * 1024 * 1024,
		ExpireTime:           7 * 24 * time.Hour,
		CheckRollingInterval: 1024,
	}

	for _, opt := range opts {
		opt(&options)
	}

	var my = &RollingFileHook{
		args:      options,
		formatter: NewMessageFormatter(options.Flag, levelHints),
	}

	_ = os.MkdirAll(options.DirName, os.ModePerm)

	for level := options.FilterLevel; level < LevelMax; level++ {
		var err = my.openLogFile(level)
		checkPrintError(err)
	}

	loom.Go(my.goLoop)
	return my
}
```

- [ ] **Step 3: 更新 logger_test.go，移除对不存在函数的调用**

修改 `TestRollingFileHook` 和 `TestFileHookFilterLevel`：

```go
func TestRollingFileHook(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(4)

	var fileHook = NewRollingFileHook(
		WithFlag(flag),
		WithFilterLevel(LevelWarn),
		WithMaxFileSize(16),
		WithCheckRollingInterval(10),
	)

	l.AddHook(fileHook)

	// 测试archive中生成的文件名
	for i := 0; i < 200; i++ {
		l.Info("This is info")
		l.Warn("I am a warning")
		l.Error("Error occurred")
	}

	l.Close()
}

func TestFileHookFilterLevel(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(4)

	var fileHook = NewRollingFileHook(
		WithFlag(flag),
		WithFilterLevel(LevelWarn),
		WithMaxFileSize(16),
	)

	l.AddHook(fileHook)

	l.Info("1 warn level: This is info")
	l.Warn("2 warn level: I am a warning")
	l.Error("3 warn level: Error occurred")

	if fileHook.args.FilterLevel != LevelWarn {
		t.Fatal()
	}

	fileHook.SetFilterLevel(LevelInfo)
	if fileHook.args.FilterLevel != LevelInfo {
		t.Fatal()
	}

	l.Info("4 info level: This is info")
	l.Warn("5 info level: I am a warning")
	l.Error("6 info level: Error occurred")

	_ = l.Close()
}
```

- [ ] **Step 4: 运行 RollingFileHook 测试**

Run: `go test -run TestRollingFileHook -v`
Expected: PASS

Run: `go test -run TestFileHookFilterLevel -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add rolling_file_hook.go rolling_file_hook_option.go logger_test.go
git commit -m "refactor: make rollingFileHookOptions embed HookConfig and use universal options"
```

---

## Chunk 4: 重构 ding Hook

### Task 8: 重构 ding hookOptions 使用嵌入的 HookConfig

**Files:**
- Modify: `ding/hook_option.go:12-24`
- Modify: `ding/hook.go:18-88`

- [ ] **Step 1: 修改 ding/hook_option.go，嵌入 HookConfig**

```go
package ding

import "github.com/lixianmin/logo"

/********************************************************************
created:    2020-02-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type hookOptions struct {
	logo.HookConfig
}

type HookOption func(*hookOptions)

func WithFilterLevel(level int) HookOption {
	return func(options *hookOptions) {
		logo.WithFilterLevel(level)(&options.HookConfig)
	}
}
```

- [ ] **Step 2: 修改 ding/hook.go，使用嵌入的 FilterLevel**

修改 `NewHook` 函数：

```go
func NewHook(talker Talker, opts ...HookOption) *TalkHook {
	if talker == nil {
		panic("Talker should not be null")
	}

	var options = hookOptions{
		HookConfig: logo.HookConfig{
			FilterLevel: logo.LevelInfo,
		},
	}

	for _, opt := range opts {
		opt(&options)
	}

	var my = &TalkHook{
		talker:      talker,
		filterLevel: options.FilterLevel,
	}

	return my
}
```

修改 `SetFilterLevel` 方法（如果需要的话，可以移除 filterLevel 字段，直接使用 HookConfig）：

```go
type TalkHook struct {
	talker      Talker
	args        hookOptions
}

func NewHook(talker Talker, opts ...HookOption) *TalkHook {
	if talker == nil {
		panic("Talker should not be null")
	}

	var options = hookOptions{
		HookConfig: logo.HookConfig{
			FilterLevel: logo.LevelInfo,
		},
	}

	for _, opt := range opts {
		opt(&options)
	}

	var my = &TalkHook{
		talker: talker,
		args:   options,
	}

	return my
}

func (my *TalkHook) SetFilterLevel(level int) {
	if level > logo.LevelNone && level < logo.LevelMax {
		my.args.FilterLevel = level
	}
}

func (my *TalkHook) Write(message logo.Message) {
	var level = message.GetLevel()
	if level < my.args.FilterLevel {
		return
	}

	var text = message.GetText()
	var frames = message.GetFrames()
	if len(frames) > 0 {
		var buffer = make([]byte, 0, 128)
		for i := 1; i < len(frames); i++ {
			buffer = append(buffer, "  \n  "...)
			buffer = tools.AppendFrameInfo(buffer, frames[i])
		}

		var first = frames[0]
		text = fmt.Sprintf("%s:%d [%s()] %s %s", path.Base(first.File), first.Line, getFunctionName(first.Function), text, buffer)
	}

	my.talker.PostMessage(level, "", text)
}
```

- [ ] **Step 3: 运行 ding 测试**

Run: `go test ./ding -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add ding/hook_option.go ding/hook.go
git commit -m "refactor: make ding hookOptions embed HookConfig"
```

---

## Chunk 5: 添加 lark Hook 支持（如果需要）

### Task 9: 检查 lark 是否需要 Hook Option

**Files:**
- Check: `lark/` 目录

- [ ] **Step 1: 检查 lark 目录是否有 Hook 结构**

Run: `grep -r "Hook" lark/`

如果发现 lark 也有类似 ding 的 Hook 结构，则按照 ding 的方式重构。如果没有，跳过此任务。

- [ ] **Step 2: 如果需要，重构 lark Hook**

（根据检查结果决定是否执行）

---

## Chunk 6: 最终验证和清理

### Task 10: 运行完整测试套件

**Files:**
- All test files

- [ ] **Step 1: 运行所有测试**

Run: `go test ./... -v`
Expected: 所有测试 PASS

- [ ] **Step 2: 生成覆盖率报告**

Run: `go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out`
Expected: 覆盖率 >= 80%

如果覆盖率不足 80%，添加更多测试用例。

- [ ] **Step 3: 运行代码格式化**

Run: `go fmt ./...`

- [ ] **Step 4: 运行静态检查**

Run: `go vet ./...`

- [ ] **Step 5: 最终构建验证**

Run: `go build ./...`

- [ ] **Step 6: Commit 最终状态**

```bash
git add .
git commit -m "chore: final cleanup and verify all tests pass with 80%+ coverage"
```

---

## Task 11: 更新文档

**Files:**
- Update: `AGENTS.md` (如果需要)

- [ ] **Step 1: 检查 AGENTS.md 是否需要更新**

如果重构改变了公共 API 或使用模式，更新 AGENTS.md 中的相关示例。

- [ ] **Step 2: Commit 文档更新**

```bash
git add AGENTS.md
git commit -m "docs: update AGENTS.md with new Hook Option pattern"
```

---

## 执行注意事项

1. **TDD 原则**: 严格遵循 Red-Green-Refactor 循环，每一步都先写测试，再实现代码
2. **小步提交**: 每个 Task 完成后立即提交，便于回滚和追踪
3. **测试先行**: 在修改任何生产代码前，确保有对应的测试
4. **覆盖率目标**: 每个新增的功能都要有对应的测试，确保整体覆盖率 >= 80%
5. **向后兼容**: 虽然 AGENTS.md 说不需要兼容旧代码，但要确保所有现有测试都能通过

## 预期结果

完成后，代码应该：
1. 所有测试通过
2. 测试覆盖率 >= 80%
3. ConsoleHook、RollingFileHook、ding.Hook 都使用嵌入的 HookConfig
4. 通用的 WithFlag 和 WithFilterLevel 函数可以在所有 Hook 类型间共享
5. 代码结构清晰，易于扩展新的 Hook 类型
