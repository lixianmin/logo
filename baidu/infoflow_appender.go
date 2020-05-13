package baidu

import (
	"fmt"
	"github.com/lixianmin/logo"
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
	LevelFilter int
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

func (my *InfoFlowAppender) Write(message logo.Message) {
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

func checkInfoFlowAppenderArgs(args *InfoFlowAppenderArgs) {
	if args.Talker == nil {
		panic("Talker should not be null")
	}

	if args.LevelFilter <= 0 {
		args.LevelFilter = logo.LevelInfo
	}
}
