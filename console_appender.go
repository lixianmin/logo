package logo

import (
	"os"
)

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type ConsoleAppenderArgs struct {
	Flag        int
	FilterLevel int
}

type ConsoleAppender struct {
	args      ConsoleAppenderArgs
	formatter *MessageFormatter
}

func NewConsoleAppender(args ConsoleAppenderArgs) *ConsoleAppender {
	checkConsoleAppenderArgs(&args)

	var my = &ConsoleAppender{
		args:      args,
		formatter: newMessageFormatter(args.Flag, levelHintsConsole),
	}

	return my
}

func (my *ConsoleAppender) SetFilterLevel(level int) {
	if level > LevelNone && level < LevelMax {
		my.args.FilterLevel = level
	}
}

func (my *ConsoleAppender) Write(message Message) {
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

func (my *ConsoleAppender) writeMessage(fout *os.File, message Message) {
	var buffer = my.formatter.format(message)
	_, _ = fout.Write(buffer)
}

func checkConsoleAppenderArgs(args *ConsoleAppenderArgs) {
	if args.FilterLevel <= 0 {
		args.FilterLevel = LevelInfo
	}
}
