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
	args      ConsoleHookArgs
	formatter *MessageFormatter
}

func NewConsoleHook(args ConsoleHookArgs) *ConsoleHook {
	checkConsoleHookArgs(&args)

	var my = &ConsoleHook{
		args:      args,
		formatter: newMessageFormatter(args.Flag, levelHintsConsole),
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

	switch level {
	case LevelDebug, LevelInfo:
		my.writeMessage(os.Stdout, message)
	case LevelWarn, LevelError:
		my.writeMessage(os.Stderr, message)
	}
}

func (my *ConsoleHook) writeMessage(fout *os.File, message Message) {
	var buffer = my.formatter.format(message)
	_, _ = fout.Write(buffer)
}

func checkConsoleHookArgs(args *ConsoleHookArgs) {
	if args.FilterLevel <= 0 {
		args.FilterLevel = LevelInfo
	}
}
