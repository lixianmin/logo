package logo

import "testing"

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestConsoleAppender(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(2)

	var console = NewConsoleAppender(ConsoleAppenderArgs{Flag: flag})
	l.AddAppender(console)

	l.Info("This is info")
	l.Warn("I am a warning")
	l.Error("Error occurred")
}

func TestRollingFileAppender(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(2)

	var console = NewRollingFileAppender(RollingFileAppenderArgs{
		LevelFilter: LevelWarn,
		Flag:        flag,
		MaxFileSize: 16,
	})

	l.AddAppender(console)

	l.Info("This is info")
	l.Warn("I am a warning")
	l.Error("Error occurred")
}

func TestLogAnyObject(t *testing.T) {
	Info(123.45678)
	Info(t)
}

func TestAutoFlush(t *testing.T) {
	logger.AddFlag(LogAsyncWrite)
	var i = 0
	for i < 10 {
		Info(i)
		i++
	}

	logger.RemoveFlag(LogAsyncWrite)
	for i < 20 {
		Warn(i)
		i++
	}

	logger.AddFlag(LogAsyncWrite)
	for i < 30 {
		Error(i)
		i++
	}

	//logger.Flush()
}
