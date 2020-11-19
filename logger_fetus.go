package logo

import (
	"fmt"
	"io"
)

/********************************************************************
created:    2020-11-19
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type loggerFetus struct {
	hooks []IHook
}

func newLoggerFetus() *loggerFetus {
	return &loggerFetus{
		hooks: make([]IHook, 0, 4),
	}
}

// 外面保证传入的hook不是nil
func (fetus *loggerFetus) AddHook(hook IHook) {
	fetus.hooks = append(fetus.hooks, hook)
}

func (fetus *loggerFetus) WriteMessage(message Message) {
	for _, hook := range fetus.hooks {
		hook.Write(message)
	}
}

func (fetus *loggerFetus) FlushMessage(messageChan chan Message) {
	var count = len(messageChan)
	for i := 0; i < count; i++ {
		var message = <-messageChan
		fetus.WriteMessage(message)
	}
}

func (fetus *loggerFetus) SetFilterLevel(level int) {
	for _, hook := range fetus.hooks {
		hook.SetFilterLevel(level)
	}
}

func (fetus *loggerFetus) Close() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	for _, hook := range fetus.hooks {
		if closer, ok := hook.(io.Closer); ok {
			var err = closer.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	fetus.hooks = nil
}
