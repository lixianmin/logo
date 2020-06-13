package logo

import (
	"context"
	"fmt"
	"github.com/lixianmin/got/loom"
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
	cancel    context.CancelFunc
}

func NewLogger() *Logger {
	const chanLen = 128
	var ctx, cancel = context.WithCancel(context.Background())
	var logger = &Logger{
		funcCallDepth: -1,
		messageChan:   make(chan Message, chanLen),
		cancel:        cancel,
	}

	go logger.goLoop(ctx)
	return logger
}

func (my *Logger) goLoop(ctx context.Context) {
	defer loom.DumpIfPanic()
	for {
		select {
		case message := <-my.messageChan:
			my.writeMessage(message)
		case <-ctx.Done():
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
	// Flush()需要放到wc.Close()的前面。
	// 否则如果先调用wc.Close()的话，一旦goLoop()的协程先于Flush()退出，则Flush()方法可能死锁
	my.Flush()
	my.cancel()

	for i, appender := range my.appenderList {
		my.appenderList[i] = nil
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
func (my *Logger) Debug(first interface{}, args ...interface{}) {
	var text = formatLog(first, args...)
	my.pushMessage(Message{text: text, level: LevelDebug})
}

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

	// 如果未开启异步写模式
	if !my.HasFlag(LogAsyncWrite) {
		my.Flush()
	}
}

func (my *Logger) writeMessage(message Message) {
	for _, appender := range my.appenderList {
		if appender != nil {
			appender.Write(message)
		}
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
