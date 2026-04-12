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
- 说明开箱即用（init() 自动创建 ConsoleHook）

### 第2段：全局便捷函数
- 展示 Printf 风格和空格拼接风格
- 展示 JSON 结构化日志（JsonI/JsonD/JsonW/JsonE）
- 函数对照表（函数 / 级别 / 用途）
- 说明 formatLog 智能行为

### 第3段：添加文件日志（RollingFileHook）
- 完整可运行示例
- 生成的文件目录结构
- RollingFileHook 全部选项表格（选项/默认值/说明）
- 说明 fallthrough 行为（ERROR 同时写入三个文件）

### 第4段：自定义控制台日志（ConsoleHook）
- 说明默认已自带 ConsoleHook
- 两种方式：设置过滤级别 / 替换全局 Logger
- 解释 SetFuncCallDepth 的作用

### 第5段：异步写入与资源清理
- LogAsyncWrite 的使用和注意事项
- 两种清理模式：简单 defer / 信号退出
- 说明同步模式无需手动处理

### 第6段：消息平台 Hook（钉钉 + 飞书）
- ding.NewTalk + ding.NewHook 完整示例
- lark.NewLark + ding.NewHook 完整示例
- 频率限制说明
- 揭示 lark 复用 ding.Talker 接口

### 第7段：日志级别与格式标志位参考
- 级别对照表（常量/值/说明）
- 标志位对照表（常量/输出示例/说明）
- 推荐组合及输出示例

---

## AGENTS.md 改动

在"项目概述"之后、"架构设计"之前，插入 **"使用指南"** 章节：

- 快速使用推荐模式（完整代码）
- 6 条关键规则：
  1. 优先使用全局便捷函数
  2. Format 支持两种风格
  3. JSON 日志用 key-value 交替传入
  4. 日志标签用方括号
  5. 异步写入必须 Close
  6. RollingFileHook 默认 FilterLevel 为 LevelInfo
- 全局 Logger 获取方式（ILogger 接口 vs *Logger 实例）
- 消息平台 Hook 速查

---

## 不变的部分

- AGENTS.md 中的"架构设计"、"代码风格指南"、"构建命令"等章节保持不变
- notes/ 目录下的文档不做修改
- 项目宪法（00.constitution.md）不做修改
