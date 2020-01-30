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

type TalkArgs struct {
	Talker      *Talk
	LevelFilter int
}

type TalkAppender struct {
	args TalkArgs
}

func NewTalkAppender(args TalkArgs) *TalkAppender {
	checkTalkArgs(&args)

	var my = &TalkAppender{
		args: args,
	}

	return my
}

func (my *TalkAppender) Write(message logo.Message) {
	var level = message.GetLevel()
	var args = my.args
	if level < args.LevelFilter {
		return
	}

	var filePath = message.GetFilePath()
	var lineNum = message.GetLineNum()
	var text = fmt.Sprintf("[%s:%d] %s", path.Base(filePath), lineNum, message.GetText())

	var talker = args.Talker
	switch level {
	case logo.LevelInfo:
		talker.SendInfo("", text)
	case logo.LevelWarn:
		talker.SendWarn("", text)
	case logo.LevelError:
		talker.SendError("", text)
	}
}

func (my *TalkAppender) Close() error {
	return nil
}

func checkTalkArgs(args *TalkArgs) {
	if args.Talker == nil {
		panic("Talker should not be null")
	}

	if args.LevelFilter <= 0 {
		args.LevelFilter = logo.LevelInfo
	}
}
