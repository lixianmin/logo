package logo

import "testing"

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestConsoleHook(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(4)

	var console = NewConsoleHook(ConsoleHookArgs{Flag: flag, FilterLevel: LevelDebug})
	l.AddHook(console)

	l.Debug("Debug", "Message")
	l.Info("This is info")
	l.Warn("I am a warning")
	l.Error("Error occurred")

	l.Close()
}

func TestRollingFileHook(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(4)

	var fileHook = NewRollingFileHook(RollingFileHookArgs{
		FilterLevel:          LevelWarn,
		Flag:                 flag,
		MaxFileSize:          16,
		CheckRollingInterval: 10,
	})

	l.AddHook(fileHook)

	for i := 0; i < 10; i++ {
		l.Info("This is info")
		l.Warn("I am a warning")
		l.Error("Error occurred")
	}

	l.Close()
}

func TestFileHookFilterLevel(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(4)

	var fileHook = NewRollingFileHook(RollingFileHookArgs{
		FilterLevel: LevelWarn,
		Flag:        flag,
		MaxFileSize: 16,
	})

	l.AddHook(fileHook)

	l.Info("1 warn level: This is info")
	l.Warn("2 warn level: I am a warning")
	l.Error("3 warn level: Error occurred")

	if fileHook.args.FilterLevel != LevelWarn {
		t.Fatal()
	}

	fileHook.SetFilterLevel(LevelInfo)
	if fileHook.args.FilterLevel != LevelInfo {
		t.Fatal()
	}

	l.Info("4 info level: This is info")
	l.Warn("5 info level: I am a warning")
	l.Error("6 info level: Error occurred")

	_ = l.Close()
}

func TestLogAnyObject(t *testing.T) {
	Info(123.45678)
	Info(t)
}

func TestAutoFlush(t *testing.T) {
	var logger = GetLogger().(*Logger)
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

	logger.Flush()
}

func TestClose(t *testing.T) {
	var logger = GetLogger().(*Logger)
	logger.AddFlag(LogAsyncWrite)

	for i := 0; i < 50; i++ {
		logger.Info(i)
	}

	_ = logger.Close()

	logger.Info("closed")
}

func TestJsonI(t *testing.T) {
	JsonI("age", 10, "name", "lixianmin", "male", true, "height", 1.78)
}
