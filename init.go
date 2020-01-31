package logo

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var defaultLogger = NewLogger()

func init() {
	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	defaultLogger.SetFuncCallDepth(3)

	var console = NewConsoleAppender(ConsoleAppenderArgs{Flag: flag})
	defaultLogger.AddAppender(console)
}

func GetDefaultLogger() *Logger {
	return defaultLogger
}

func Info(first interface{}, args ...interface{}) {
	defaultLogger.Info(first, args...)
}

func Warn(first interface{}, args ...interface{}) {
	defaultLogger.Warn(first, args...)
}

func Error(first interface{}, args ...interface{}) {
	defaultLogger.Error(first, args...)
}
