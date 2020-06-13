package baidu

import (
	"fmt"
	"github.com/lixianmin/logo"
	"github.com/lixianmin/logo/tools"
	"path"
)

/********************************************************************
created:    2020-05-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 这里没有选择使用InfoFlowArgs，是为了给infoflow.go中的InfoFlow类留出未来
type InfoFlowAppenderArgs struct {
	Talker      *InfoFlow
	FilterLevel int
}

type InfoFlowAppender struct {
	args InfoFlowAppenderArgs
}

func NewInfoFlowAppender(args InfoFlowAppenderArgs) *InfoFlowAppender {
	checkInfoFlowAppenderArgs(&args)

	var my = &InfoFlowAppender{
		args: args,
	}

	return my
}

func (my *InfoFlowAppender) Close() error {
	return my.args.Talker.Close()
}

func (my *InfoFlowAppender) SetFilterLevel(level int) {
	if level > logo.LevelNone && level < logo.LevelMax {
		my.args.FilterLevel = level
	}
}

func (my *InfoFlowAppender) Write(message logo.Message) {
	var level = message.GetLevel()
	var args = my.args
	if level < args.FilterLevel {
		return
	}

	var text = message.GetText()
	var frames = message.GetFrames()
	if len(frames) > 0 {
		var buffer = make([]byte, 0, 128)
		for i := 1; i < len(frames); i++ {
			buffer = append(buffer, '\n')
			buffer = tools.AppendFrameInfo(buffer, frames[i])
		}

		var first = frames[0]
		text = fmt.Sprintf("[%s:%d] %s\n%s", path.Base(first.File), first.Line, text, buffer)
	}

	var talker = args.Talker
	switch level {
	case logo.LevelDebug:
		talker.SendDebug("", text)
	case logo.LevelInfo:
		talker.SendInfo("", text)
	case logo.LevelWarn:
		talker.SendWarn("", text)
	case logo.LevelError:
		talker.SendError("", text)
	}
}

func checkInfoFlowAppenderArgs(args *InfoFlowAppenderArgs) {
	if args.Talker == nil {
		panic("Talker should not be null")
	}

	if args.FilterLevel <= 0 {
		args.FilterLevel = logo.LevelInfo
	}
}
