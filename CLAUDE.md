# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**logo** is a minimal Go logging framework (log + go) that provides basic logging functionality without external dependencies (except for the author's own `got` utility library). The framework is designed to be simple and extensible through a hook system.

## Architecture

### Core Components

- **Logger** (`logger.go`): Main logger implementation with async message processing, configurable filter levels, and hook management
- **Message** (`message.go`): Internal message structure with text, level, and stack frames
- **Hook System** (`ihook.go`): Interface for extending log output to various destinations
- **Default Logger** (`init.go`): Global logger instance with convenience functions

### Hook Implementations

- **ConsoleHook** (`console_hook.go`): Outputs to console with color support
- **RollingFileHook** (`rolling_file_hook.go`): File output with rotation by size
- **DingTalk Hook** (`ding/`): Sends logs to DingTalk (钉钉) messaging platform
- **Lark Hook** (`lark/`): Sends logs to Lark (飞书) messaging platform

### Key Design Patterns

- Async message processing via channels for performance
- Flag-based configuration for log formatting
- Plugin architecture via hooks for extensibility
- Global default logger for convenience usage

## Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
# Or use the provided script:
bash go.tool.cover.sh
```

### Building
```bash
# Build the module
go build ./...

# Run examples
go run examples/main.go
```

### Code Quality
```bash
# Format code
go fmt ./...

# Vet for potential issues
go vet ./...
```

## Key Configuration

### Log Levels
- `LevelDebug` (1)
- `LevelInfo` (2)
- `LevelWarn` (3)
- `LevelError` (4)

### Format Flags
- `FlagDate`: Show date (1998-10-29)
- `FlagTime`: Show time (12:24:00)
- `FlagLongFile`: Full file path (i/am/the/path/file.go:12)
- `FlagShortFile`: Short file name (file.go:34)
- `FlagLevel`: Show level indicators [D], [I], [W], [E]

### Logger Flags
- `LogAsyncWrite`: Enable async writing for performance (requires explicit Flush/Close)

## Usage Patterns

### Basic Usage
```go
import "github.com/lixianmin/logo"

// Use global logger
logo.Info("Hello, %s", "world")
logo.Error("Something went wrong: %v", err)

// JSON formatted logging
logo.JsonI("key1", value1, "key2", value2)
```

### Advanced Usage
```go
// Create custom logger
logger := logo.NewLogger()
logger.AddFlag(logo.LogAsyncWrite)
logger.SetFilterLevel(logo.LevelInfo)

// Add hooks
console := logo.NewConsoleHook(logo.ConsoleHookArgs{Flag: flag})
file := logo.NewRollingFileHook(logo.RollingFileHookArgs{Flag: flag})
logger.AddHook(console)
logger.AddHook(file)

// Remember to close for cleanup
defer logger.Close()
```

## Important Notes

- The framework depends on the author's `got` utility library for concurrent primitives and utilities
- Async logging requires proper cleanup via `Close()` or `Flush()` to avoid data loss
- Hook implementations should be thread-safe as they may be called concurrently
- Message platform hooks (DingTalk, Lark) include rate limiting and message splitting logic