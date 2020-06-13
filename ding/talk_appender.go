package ding

import (
	"fmt"
	"github.com/lixianmin/logo"
	"path"
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

	var filePath = message.GetFilePath()
	var lineNum = message.GetLineNum()
	var text = fmt.Sprintf("[%s:%d] %s\n%s", path.Base(filePath), lineNum, message.GetText(), message.GetTrace())

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
