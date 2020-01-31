package logo

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
)

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Logger struct {
	Flag
	appenderList  []Appender
	funcCallDepth int
	messageChan   chan Message

	waitFlush sync.WaitGroup
	wc        *WaitClose
}

func NewLogger() *Logger {
	const chanLen = 128
	var logger = &Logger{
		funcCallDepth: -1,
		messageChan:   make(chan Message, chanLen),
		Flag:          Flag{flags: LogAutoFlush},
		wc:            NewWaitClose(),
	}

	go logger.goLoop()
	return logger
}

func (my *Logger) goLoop() {
	defer DumpIfPanic()
	for {
		select {
		case message := <-my.messageChan:
			my.writeMessage(message)
		case <-my.wc.CloseChan:
			return
		}
	}
}

func (my *Logger) SetFuncCallDepth(depth int) {
	my.funcCallDepth = depth
}

// 这个方法是否需要考虑设计成线程安全？
// 暂时没有必要：appender列表基本上是在工程启动最前期初始化完成，目前没遇到运行中需要改动的情况
func (my *Logger) AddAppender(appender Appender) {
	if appender != nil {
		my.appenderList = append(my.appenderList, appender)
	}
}

func (my *Logger) Flush() {
	my.waitFlush.Wait()
}

func (my *Logger) Close() error {
	_ = my.wc.Close()
	my.Flush()

	for _, appender := range my.appenderList {
		if closer, ok := appender.(io.Closer); ok {
			var err = closer.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

// 第一个参数有可能是format，也有可能是任意其它类型的对象
func (my *Logger) Info(first interface{}, args ...interface{}) {
	var text = formatLog(first, args...)
	my.pushMessage(Message{text: text, level: LevelInfo})
}

func (my *Logger) Warn(first interface{}, args ...interface{}) {
	var text = formatLog(first, args...)
	my.pushMessage(Message{text: text, level: LevelWarn})
}

func (my *Logger) Error(first interface{}, args ...interface{}) {
	var text = formatLog(first, args...)
	my.pushMessage(Message{text: text, level: LevelError})
}

func (my *Logger) pushMessage(message Message) {
	if my.funcCallDepth > 0 {
		var ok bool
		_, file, line, ok := runtime.Caller(my.funcCallDepth)
		if !ok {
			file = "???"
			line = 0
		}

		message.filePath = file
		message.lineNum = line
	}

	my.waitFlush.Add(1)
	my.messageChan <- message

	// 如果开启了autoFlush
	if my.HasFlag(LogAutoFlush) {
		my.Flush()
	}
}

func (my *Logger) writeMessage(message Message) {
	for _, appender := range my.appenderList {
		appender.Write(message)
	}

	my.waitFlush.Add(-1)
}

func formatLog(first interface{}, args ...interface{}) string {
	var msg string
	switch first.(type) {
	case string:
		msg = first.(string)
		if len(args) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(args))
		}
	default:
		msg = fmt.Sprint(first)
		if len(args) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(args))
	}
	return fmt.Sprintf(msg, args...)
}
