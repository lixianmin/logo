package logo

import (
	"fmt"
	"github.com/lixianmin/got/convert"
	"github.com/lixianmin/got/loom"
	"github.com/lixianmin/got/std"
	"github.com/lixianmin/got/taskx"
	"github.com/lixianmin/logo/tools"
	"strings"
	"sync/atomic"
	"unsafe"
)

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Logger struct {
	Flag
	funcCallDepth int32
	messageChan   chan Message
	filterLevel   int32
	stackLevel    int32

	wc    loom.WaitClose
	tasks *taskx.Queue
}

func NewLogger(opts ...LoggerOption) *Logger {
	// 默认值
	var options = loggerOptions{
		BufferSize: 4096,
	}

	// 初始化
	for _, opt := range opts {
		opt(&options)
	}

	var my = &Logger{
		funcCallDepth: -1,
		messageChan:   make(chan Message, options.BufferSize),
		filterLevel:   LevelInfo,
		stackLevel:    LevelError,
	}

	my.tasks = taskx.NewQueue(taskx.WithCloseChan(my.wc.C()))
	loom.Go(my.goLoop)

	return my
}

func (my *Logger) goLoop(later loom.Later) {
	var closeChan = my.wc.C()
	var fetus = newLoggerFetus()

	defer fetus.Close()

	for {
		select {
		case message := <-my.messageChan:
			fetus.WriteMessage(message)
		case task := <-my.tasks.C:
			_ = task.Do(fetus)
		case <-closeChan:
			fetus.FlushMessage(my.messageChan)
			return
		}
	}
}

func (my *Logger) AddHook(hook IHook) {
	if !std.IsNil(hook) {
		select {
		case <-my.wc.C():
		default:
			my.tasks.SendCallback(func(args interface{}) (interface{}, error) {
				var fetus = args.(*loggerFetus)
				fetus.AddHook(hook)
				return nil, nil
			}).Get1()
		}
	}
}

func (my *Logger) Write(p []byte) (n int, err error) {
	if p != nil {
		my.pushMessage(Message{
			text:  *(*string)(unsafe.Pointer(&p)),
			level: int(atomic.LoadInt32(&my.filterLevel)),
		})
	}

	return len(p), nil
}

func (my *Logger) Flush() {
	select {
	case <-my.wc.C():
	default:
		my.tasks.SendCallback(func(args interface{}) (interface{}, error) {
			var fetus = args.(*loggerFetus)
			fetus.FlushMessage(my.messageChan)
			return nil, nil
		}).Get1()
	}
}

func (my *Logger) Close() error {
	return my.wc.Close(func() error {
		my.Flush()
		return nil
	})
}

func (my *Logger) SetFilterLevel(level int) {
	select {
	case <-my.wc.C():
	default:
		if level > LevelNone && level < LevelMax {
			atomic.StoreInt32(&my.filterLevel, int32(level))

			my.tasks.SendCallback(func(args interface{}) (interface{}, error) {
				var fetus = args.(*loggerFetus)
				fetus.SetFilterLevel(level)
				return nil, nil
			})
		}
	}
}

func (my *Logger) SetFuncCallDepth(depth int32) {
	atomic.StoreInt32(&my.funcCallDepth, depth)
}

func (my *Logger) SetStackLevel(level int32) {
	if level > LevelNone && level < LevelMax {
		atomic.StoreInt32(&my.stackLevel, level)
	}
}

// Debug 第一个参数有可能是format，也有可能是任意其它类型的对象
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
	if my.wc.IsClosed() {
		return
	}

	var fullStack = message.level >= int(atomic.LoadInt32(&my.stackLevel))
	var depth = int(atomic.LoadInt32(&my.funcCallDepth))
	message.frames = tools.CallersFrames(depth, fullStack)

	// 原来的使用waitGroup来同步Flush()的方案是错误的，会报如下错误
	// sync: WaitGroup is reused before previous Wait has returned
	//my.waitFlush.Add(1)
	select {
	case my.messageChan <- message:
	case <-my.wc.C():
	}

	// 如果未开启异步写模式
	if !my.HasFlag(LogAsyncWrite) {
		my.Flush()
	}
}

func formatLog(first interface{}, args ...interface{}) string {
	var message string
	switch first := first.(type) {
	case string:
		message = first
	case []byte:
		message = convert.String(first)
	default:
		message = fmt.Sprint(first)
	}

	if len(args) == 0 {
		return message
	}

	var format = message
	var isFormat = strings.Contains(message, "%") && !strings.Contains(message, "%%")
	if !isFormat {
		// 如果不是format的, 则需要增加一些格式化的参数
		format = message + strings.Repeat(" %v", len(args))
	}

	return fmt.Sprintf(format, args...)
}
