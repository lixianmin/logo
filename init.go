package logo

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var logger = NewLogger()

func init() {
	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	logger.SetFuncCallDepth(3)

	var console = NewConsoleAppender(ConsoleAppenderArgs{Flag: flag})
	logger.AddAppender(console)
}

func Close() error {
	return logger.Close()
}

func AddAppender(appender Appender) {
	logger.AddAppender(appender)
}

func Info(format string, args ...interface{}) {
	logger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	logger.Error(format, args...)
}
