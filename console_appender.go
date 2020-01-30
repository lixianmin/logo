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
	formatter   *MessageFormatter
	levelFilter int
}

func NewConsoleAppender(flag int, levelFilter int) *ConsoleAppender {
	var my = &ConsoleAppender{
		formatter:   newMessageFormatter(flag, levelHintsConsole),
		levelFilter: levelFilter,
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
