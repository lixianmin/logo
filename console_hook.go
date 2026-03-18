package logo

import (
	"os"
)

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type ConsoleHookArgs struct {
	Flag        int
	FilterLevel int
}

type ConsoleHook struct {
	args      consoleHookOptions
	formatter *MessageFormatter
}

func NewConsoleHook(args ConsoleHookArgs) *ConsoleHook {
	return NewConsoleHookWithOptions(
		WithFlag(args.Flag),
		WithFilterLevel(args.FilterLevel),
	)
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

func (my *ConsoleHook) SetFilterLevel(level int) {
	if level > LevelNone && level < LevelMax {
		my.args.FilterLevel = level
	}
}

func (my *ConsoleHook) Write(message Message) {
	var level = message.GetLevel()
	var args = my.args
	if level < args.FilterLevel {
		return
	}

	var buffer = my.formatter.Format(message)
	switch level {
	case LevelDebug, LevelInfo:
		_, _ = os.Stdout.Write(buffer)
	case LevelWarn, LevelError:
		_, _ = os.Stderr.Write(buffer)
	}
}
