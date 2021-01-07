package logo

import (
	"github.com/lixianmin/got/convert"
	"strconv"
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

// Info() for json
func I(args ...interface{}) {
	theLogger.Info(formatJson(args...))
}

func formatJson(args ...interface{}) string {
	var count = len(args)
	var halfCount = count >> 1

	var results = make([]byte, 128)
	results = append(results, '{')
	for i := 0; i < halfCount; i++ {
		var index = i << 1
		var key, _ = args[index].(string)
		var value = args[index+1]

		results = strconv.AppendQuote(results, key)
		results = append(results, ':')
		results = convert.AppendArg(results, value)

		if i+1 < halfCount {
			results = append(results, ',')
		}
	}

	results = append(results, '}')
	var text = string(results)
	return text
}
