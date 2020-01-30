package ding

import (
	"fmt"
	"github.com/lixianmin/btc-trade/app/core/logo"
	"path"
)

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type TalkAppender struct {
	levelFilter int

	talk *Talk
}

func NewTalkAppender(levelFilter int, talk *Talk) *TalkAppender {
	if talk == nil {
		panic("talk should not be null")
	}

	var my = &TalkAppender{
		levelFilter: levelFilter,
		talk:        talk,
	}

	return my
}

func (my *TalkAppender) Write(message logo.Message) {
	var level = message.GetLevel()
	if level < my.levelFilter {
		return
	}

	var filePath = message.GetFilePath()
	var lineNum = message.GetLineNum()
	var text = fmt.Sprintf("[%s:%d] %s", path.Base(filePath), lineNum, message.GetText())

	switch level {
	case logo.LevelInfo:
		my.talk.SendInfo("", text)
	case logo.LevelWarn:
		my.talk.SendWarn("", text)
	case logo.LevelError:
		my.talk.SendError("", text)
	}
}

func (my *TalkAppender) Close() error {
	return nil
}
