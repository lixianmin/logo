package logo

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var theLogger = NewLogger()

func init() {
	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	theLogger.SetFuncCallDepth(5)

	var console = NewConsoleAppender(ConsoleAppenderArgs{Flag: flag})
	theLogger.AddAppender(console)
}

func GetLogger() *Logger {
	return theLogger
}

func Debug(first interface{}, args ...interface{}) {
	theLogger.Debug(first, args...)
}

func Info(first interface{}, args ...interface{}) {
	theLogger.Info(first, args...)
}

func Warn(first interface{}, args ...interface{}) {
	theLogger.Warn(first, args...)
}

func Error(first interface{}, args ...interface{}) {
	theLogger.Error(first, args...)
}
