package logo

import "testing"

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestConsoleAppender(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(2)

	var console = NewConsoleAppender(flag, LevelInfo)
	l.AddAppender(console)

	l.Info("This is info")
	l.Warn("I am a warning")
	l.Error("Error occurred")

	_ = l.Close()
}

func TestRollingFileAppender(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(2)

	var console = NewRollingFileAppender(flag, LevelInfo, "logs", "log_", 16)
	l.AddAppender(console)

	l.Info("This is info")
	l.Warn("I am a warning")
	l.Error("Error occurred")

	_ = l.Close()
}
