# README 与 AGENTS.md 重写设计文档

**日期**: 2026-04-12
**目标**: 重写 logo 库的 README.md 和 AGENTS.md，使 AI Agent 和人类开发者看到后能直接学会使用，无需额外学习成本。

---

## 背景

当前 README.md 内容简陋，仅包含简介和一段不完整的示例代码。缺少：
- 全局便捷函数的说明
- 各 Hook 的完整使用示例
- 格式标志位和日志级别的参考
- 异步写入与资源清理的指导
- 消息平台 Hook 的使用方法

通过分析 pc 项目（`/Users/xmli/me/code/pc`）中 27 个文件、110+ 处日志调用的真实使用方式，提炼最佳实践写入文档。

---

## 设计决策

| 决策项 | 选择 | 理由 |
|--------|------|------|
| 文档策略 | 最佳实践驱动型（方案 A） | 场景→代码→原理，读者直接抄代码就能用 |
| 目标读者 | 人类 + AI Agent 两者兼顾 | Agent 能提取模式使用，人类能快速上手 |
| 文档范围 | README.md + AGENTS.md | 保持一致性 |
| 语言 | 中文为主 | 与现有文档风格一致 |

---

## README.md 新结构

### 第1段：简介 + 快速开始
- 一句话说清是什么、依赖情况
- 最小可运行示例（全局便捷函数 + defer Close）
- 说明开箱即用：init() 自动创建 ConsoleHook（默认 flag=`FlagDate|FlagTime|FlagShortFile|FlagLevel`，FilterLevel=LevelInfo）
- 明确 ILogger 接口 vs *Logger 实例的区别：`GetLogger()` 返回 `ILogger` 接口（只有 Debug/Info/Warn/Error），需类型断言 `GetLogger().(*Logger)` 才能调用 AddHook/SetFilterLevel/Close 等

### 第2段：全局便捷函数
- 展示 Printf 风格和空格拼接风格
- 展示 JSON 结构化日志（JsonI/JsonD/JsonW/JsonE）
- 函数对照表（函数 / 级别 / 用途）
- 说明 formatLog 智能行为

### 第3段：添加文件日志（RollingFileHook）
- 完整可运行示例
- 生成的文件目录结构
- RollingFileHook 全部选项表格（7个选项：WithHookFlag, WithHookFilterLevel, WithDirName, WithFileNamePrefix, WithMaxFileSize, WithExpireTime, WithCheckRollingInterval，含默认值）
- 说明 fallthrough 行为：高级别日志会级联写入所有低级别文件。默认 FilterLevel=LevelInfo 时，ERROR 会同时写入 error.log、warn.log、info.log

### 第4段：自定义控制台日志（ConsoleHook）
- 说明默认已自带 ConsoleHook（FilterLevel=LevelInfo）
- 两种方式：设置过滤级别 / 替换全局 Logger
- 解释 SetFuncCallDepth 的作用
- 说明 SetFilterLevel 的级联效果：Logger.SetFilterLevel 会同步更新所有已注册 Hook 的 FilterLevel

### 第5段：异步写入与资源清理
- LogAsyncWrite 的使用和注意事项
- WithBufferSize 控制异步消息 channel 缓冲区大小（默认 4096）
- 两种清理模式：简单 defer / 信号退出
- 说明同步模式无需手动处理

### 第6段：消息平台 Hook（钉钉 + 飞书）
- ding.NewTalk + ding.NewHook 完整示例（含 ding.WithFilterLevel 选项）
- lark.NewLark + ding.NewHook 完整示例（含 ding.WithFilterLevel 选项）
- 频率限制说明（钉钉 20条/分钟，飞书 200条/分钟）
- 揭示 lark 复用 ding.Talker 接口
- 简要提及 SendMessage 同步发送接口（绕过限流队列，用于紧急场景）

### 第7段：日志级别与格式标志位参考
- 级别对照表（常量/值/说明）
- 标志位对照表（常量/输出示例/说明）
- 推荐组合及输出示例
- 提及 SetStackLevel(level) 控制 full stack trace 的捕获阈值（默认 LevelError）

---

## AGENTS.md 改动

在"项目概述"之后、"架构设计"之前，插入 **"使用指南"** 章节：

- 快速使用推荐模式（完整代码）
- 7 条关键规则：
  1. 优先使用全局便捷函数
  2. Format 支持两种风格
  3. JSON 日志用 key-value 交替传入
  4. 日志标签用方括号
  5. 异步写入必须 Close
  6. 所有 Hook 默认 FilterLevel 为 LevelInfo（ConsoleHook 和 RollingFileHook 均如此）
  7. `GetLogger()` 返回 ILogger 接口，需类型断言 `.(*Logger)` 才能调用 AddHook/Close
- 全局 Logger 获取方式（ILogger 接口 vs *Logger 实例）
- 消息平台 Hook 速查

---

## 不变的部分

- AGENTS.md 中的"架构设计"、"代码风格指南"、"构建命令"等章节保持不变
- notes/ 目录下的文档不做修改
- 项目宪法（00.constitution.md）不做修改

---

## 已知问题与修复记录

基于 spec review 的 12 个问题，已修复：
1. **ILogger vs *Logger 区别**：在第1段和 AGENTS.md 规则7中明确说明
2. **旧 API 名称**：新设计已使用正确 API（`GetLogger` 而非 `GetDefaultLogger`，Option 模式而非 Args 模式）
3. **Fallthrough 行为**：改为精确描述级联行为而非"写入三个文件"
4. **SetStackLevel**：在第7段中提及
5. **SetFilterLevel 级联效果**：在第4段中说明
6. **ding.WithFilterLevel**：在第6段示例中展示
7. **所有 Hook 默认 FilterLevel**：AGENTS.md 规则6改为"所有 Hook"
8. **WithBufferSize**：在第5段中提及
9. **SendMessage 同步发送**：在第6段中简要提及
10. **init() 默认 flag**：在第1段中说明
