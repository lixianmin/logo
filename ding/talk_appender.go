package ding

import (
	"fmt"
	"github.com/lixianmin/logo"
	"github.com/lixianmin/logo/tools"
	"path"
	"strings"
)

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 这里没有选择使用TalkArgs，是为了给talk.go中的Talk类留出未来
type TalkAppenderArgs struct {
	Talker      *Talk
	FilterLevel int
}

type TalkAppender struct {
	args TalkAppenderArgs
}

func NewTalkAppender(args TalkAppenderArgs) *TalkAppender {
	checkTalkAppenderArgs(&args)

	var my = &TalkAppender{
		args: args,
	}

	return my
}

func (my *TalkAppender) Close() error {
	return my.args.Talker.Close()
}

func (my *TalkAppender) SetFilterLevel(level int) {
	if level > logo.LevelNone && level < logo.LevelMax {
		my.args.FilterLevel = level
	}
}

func (my *TalkAppender) Write(message logo.Message) {
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
			buffer = append(buffer, "  \n  "...)
			buffer = tools.AppendFrameInfo(buffer, frames[i])
		}

		var first = frames[0]
		text = fmt.Sprintf("%s:%d [%s()] %s %s", path.Base(first.File), first.Line, getFunctionName(first.Function), text, buffer)
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

func checkTalkAppenderArgs(args *TalkAppenderArgs) {
	if args.Talker == nil {
		panic("Talker should not be null")
	}

	if args.FilterLevel <= 0 {
		args.FilterLevel = logo.LevelInfo
	}
}

func getFunctionName(function string) string {
	if function != "" {
		var lastIndex = strings.LastIndexByte(function, '.')
		if lastIndex > 0 {
			var s = function[lastIndex+1:]
			return s
		}
	}

	return function
}
