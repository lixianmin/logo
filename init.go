package logo

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
