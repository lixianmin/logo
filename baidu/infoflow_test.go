package baidu

import (
	"github.com/lixianmin/logo"
	"testing"
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

func TestInfoFlow(t *testing.T) {
	var talk = createTalk()
	defer talk.Close()

	talk.SendInfo("Info title", "This is an info")
	talk.SendWarn("Warn title", "This is a warning")
	talk.SendError("Error title", "This is an error")
}

func TestInfoFlowHook(t *testing.T) {
	var talk = createTalk()

	var l = logo.NewLogger()
	defer l.Close()
	l.SetFuncCallDepth(4)

	var talkHook = NewInfoFlowHook(talk, WithFilterLevel(logo.LevelWarn))

	l.AddHook(talkHook)

	l.Info("This is info, but will not appear in ding talk.")
	l.Warn("This warning will appear in ding talk.")
	l.Error("This is an %q.", "error")
}
