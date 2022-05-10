package logo

import (
	"github.com/lixianmin/logo/tools"
)

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var theLogger ILogger

func init() {
	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	var log = NewLogger()
	log.SetFuncCallDepth(5)

	var console = NewConsoleHook(ConsoleHookArgs{Flag: flag})
	log.AddHook(console)

	theLogger = log
}

func SetLogger(log ILogger) {
	if log != nil {
		theLogger = log
	}
}

func GetLogger() ILogger {
	return theLogger
}

func Debug(format string, args ...interface{}) {
	theLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	theLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	theLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	theLogger.Error(format, args...)
}

func JsonD(args ...interface{}) {
	theLogger.Debug(tools.FormatJson(args...))
}

func JsonI(args ...interface{}) {
	theLogger.Info(tools.FormatJson(args...))
}

func JsonW(args ...interface{}) {
	theLogger.Warn(tools.FormatJson(args...))
}

func JsonE(args ...interface{}) {
	theLogger.Error(tools.FormatJson(args...))
}
