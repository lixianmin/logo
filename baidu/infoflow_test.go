package baidu

import (
	"github.com/lixianmin/logo"
	"testing"
	"time"
)

/********************************************************************
created:    2020-05-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func createTalk() *InfoFlow {
	var talk = NewInfoFlow("Title Prefix", "d2127efaae4cfd31e53e47d919c27ad5c")
	return talk
}

func destroyTalk(talk *InfoFlow) {
	for talk.sendingCount > 0 {
		time.Sleep(time.Millisecond * 100)
	}

	talk.Close()
}

func TestInfoFlow(t *testing.T) {
	var talk = createTalk()
	defer destroyTalk(talk)

	talk.SendInfo("Info title", "This is an info")
	talk.SendWarn("Warn title", "This is a warning")
	talk.SendError("Error title", "This is an error")
}

func TestInfoFlowAppender(t *testing.T) {
	var talk = createTalk()
	defer destroyTalk(talk)

	var l = logo.NewLogger()
	l.SetFuncCallDepth(2)

	var talkAppender = NewInfoFlowAppender(InfoFlowAppenderArgs{
		Talker:      talk,
		LevelFilter: logo.LevelWarn,
	})

	l.AddAppender(talkAppender)

	l.Info("This is info, but will not appear in ding talk.")
	l.Warn("This warning will appear in ding talk.")
	l.Error("This is an %q.", "error")
}