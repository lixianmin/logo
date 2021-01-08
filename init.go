package logo

import (
	"github.com/lixianmin/logo/tools"
	"strconv"
	"sync"
)

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var theLogger ILogger
var bufferPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 256)
	},
}

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

func JsonD(args ...interface{}) {
	theLogger.Debug(formatJson(args...))
}

func JsonI(args ...interface{}) {
	theLogger.Info(formatJson(args...))
}

func JsonW(args ...interface{}) {
	theLogger.Warn(formatJson(args...))
}

func JsonE(args ...interface{}) {
	theLogger.Error(formatJson(args...))
}

func formatJson(args ...interface{}) string {
	var results = bufferPool.Get().([]byte)
	results = append(results, '{')
	{
		var count = len(args)
		var halfCount = (count + 1) >> 1
		for i := 0; i < halfCount; i++ {
			var index = i << 1
			var key, _ = args[index].(string)

			// 如果只有奇数个参数，则输出默认值null
			index++
			var value interface{} = nil
			if index < count {
				value = args[index]
			}

			results = strconv.AppendQuote(results, key)
			results = append(results, ':')
			results = tools.AppendJson(results, value)

			if i+1 < halfCount {
				results = append(results, ',')
			}
		}
	}
	results = append(results, '}')

	// 这里results马上就要还回去了，不要使用unsafe的[]byte转string了
	var text = string(results)
	bufferPool.Put(results[:0])
	return text
}
