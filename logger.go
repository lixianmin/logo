package logo

import (
	"fmt"
	"github.com/lixianmin/got/loom"
	"github.com/lixianmin/got/std"
	"github.com/lixianmin/logo/tools"
	"io"
	"strings"
	"sync"
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

	wc        loom.WaitClose
	tasks     *loom.TaskQueue
	waitClose sync.WaitGroup
}

type loggerFetus struct {
	hooks []IHook
}

func NewLogger() *Logger {
	const chanLen = 128
	var my = &Logger{
		funcCallDepth: -1,
		messageChan:   make(chan Message, chanLen),
		filterLevel:   LevelInfo,
		stackLevel:    LevelError,
	}

	my.tasks = loom.NewTaskQueue(loom.WithCloseChan(my.wc.C()))
	my.waitClose.Add(1)
	loom.Go(my.goLoop)

	return my
}

func (my *Logger) goLoop(later loom.Later) {
	var closeChan = my.wc.C()
	var fetus = &loggerFetus{
		hooks: make([]IHook, 0, 4),
	}

	defer func() {
		my.closeHooks(fetus)
		my.waitClose.Done()
	}()

	for {
		select {
		case message := <-my.messageChan:
			my.writeMessage(fetus, message)
		case task := <-my.tasks.C:
			_ = task.Do(fetus)
		case <-closeChan:
			if len(my.messageChan) == 0 {
				return
			}
		}
	}
}

// 这个方法是否需要考虑设计成线程安全？
// 暂时没有必要：hook列表基本上是在工程启动最前期初始化完成，目前没遇到运行中需要改动的情况
func (my *Logger) AddHook(hook IHook) {
	if !std.IsNil(hook) {
		my.tasks.SendCallback(func(args interface{}) (interface{}, error) {
			var fetus = args.(*loggerFetus)
			fetus.hooks = append(fetus.hooks, hook)
			return nil, nil
		}).Get1()
	}
}

func (my *Logger) closeHooks(fetus *loggerFetus) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	for i, hook := range fetus.hooks {
		fetus.hooks[i] = nil
		if closer, ok := hook.(io.Closer); ok {
			var err = closer.Close()
			if err != nil {
				fmt.Println(err)
			}
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
	//my.waitFlush.Wait()
}

func (my *Logger) Close() error {
	// Flush()需要放到wc.Close()的前面。
	// 否则如果先调用wc.Close()的话，一旦goLoop()的协程先于Flush()退出，则Flush()方法可能死锁
	my.Flush()
	_ = my.wc.Close(nil)
	my.waitClose.Wait()

	return nil
}

func (my *Logger) SetFilterLevel(level int) {
	if level > LevelNone && level < LevelMax {
		atomic.StoreInt32(&my.filterLevel, int32(level))

		my.tasks.SendCallback(func(args interface{}) (interface{}, error) {
			var fetus = args.(*loggerFetus)
			for _, hook := range fetus.hooks {
				if hook != nil {
					hook.SetFilterLevel(level)
				}
			}

			return nil, nil
		})
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
	if my.wc.IsClosed() {
		return
	}

	var fullStack = message.level >= int(atomic.LoadInt32(&my.stackLevel))
	var depth = int(atomic.LoadInt32(&my.funcCallDepth))
	message.frames = tools.CallersFrames(depth, fullStack)

	// 原来的使用waitGroup来同步Flush()的方案是错误的，会报如下错误
	// sync: WaitGroup is reused before previous Wait has returned
	//my.waitFlush.Add(1)
	my.messageChan <- message

	// 如果未开启异步写模式
	if !my.HasFlag(LogAsyncWrite) {
		my.Flush()
	}
}

func (my *Logger) writeMessage(fetus *loggerFetus, message Message) {
	for _, hook := range fetus.hooks {
		if hook != nil {
			hook.Write(message)
		}
	}

	//my.waitFlush.Add(-1)
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
