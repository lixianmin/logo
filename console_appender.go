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
	LevelFilter int
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

func (my *ConsoleAppender) Write(message Message) {
	var level = message.GetLevel()
	var args = my.args
	if level < args.LevelFilter {
		return
	}

	switch level {
	case LevelInfo:
		my.writeMessage(os.Stdout, message)
	case LevelWarn:
		my.writeMessage(os.Stderr, message)
	case LevelError:
		my.writeMessage(os.Stderr, message)
	}
}

func (my *ConsoleAppender) Close() error {
	return nil
}

func (my *ConsoleAppender) writeMessage(fout *os.File, message Message) {
	var buffer = my.formatter.format(message)
	_, _ = fout.Write(buffer)
}

func checkConsoleAppenderArgs(args *ConsoleAppenderArgs) {
	if args.LevelFilter <= 0 {
		args.LevelFilter = LevelInfo
	}
}
