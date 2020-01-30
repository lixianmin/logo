package logo

import (
	"os"
)

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type ConsoleAppender struct {
	levelFilter int
	formatter   *MessageFormatter
}

func NewConsoleAppender(levelFilter int, flag int) *ConsoleAppender {
	var my = &ConsoleAppender{
		levelFilter: levelFilter,
		formatter:   newMessageFormatter(flag, levelHintsConsole),
	}

	return my
}

func (my *ConsoleAppender) Write(message Message) {
	var level = message.GetLevel()
	if level < my.levelFilter {
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
