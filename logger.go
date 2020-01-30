package logo

import (
	"fmt"
	"runtime"
)

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Logger struct {
	appenderList  []Appender
	funcCallDepth int
	messageChan   chan Message
	wc            *WaitClose
}

func NewLogger() *Logger {
	var logger = &Logger{
		funcCallDepth: -1,
		messageChan:   make(chan Message, 128),
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

func (my *Logger) AddAppender(appender Appender) {
	if appender != nil {
		my.appenderList = append(my.appenderList, appender)
	}
}

func (my *Logger) Close() error {
	var count = len(my.messageChan)
	for i := 0; i < count; i++ {
		var message = <-my.messageChan
		my.writeMessage(message)
	}

	_ = my.wc.Close()

	for _, appender := range my.appenderList {
		var err = appender.Close()
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (my *Logger) Info(format string, args ...interface{}) {
	var text = fmt.Sprintf(format, args...)
	my.pushMessage(Message{text: text, level: LevelInfo})
}

func (my *Logger) Warn(format string, args ...interface{}) {
	var text = fmt.Sprintf(format, args...)
	my.pushMessage(Message{text: text, level: LevelWarn})
}

func (my *Logger) Error(format string, args ...interface{}) {
	var text = fmt.Sprintf(format, args...)
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

	my.messageChan <- message
}

func (my *Logger) writeMessage(message Message) {
	for _, appender := range my.appenderList {
		appender.Write(message)
	}
}
