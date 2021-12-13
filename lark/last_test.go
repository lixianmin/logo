package lark

import (
	"github.com/lixianmin/logo"
	"github.com/lixianmin/logo/ding"
	"strconv"
	"strings"
	"testing"
	"time"
)

/********************************************************************
created:    2021-12-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func createLark() *Lark {
	var talk = NewLark("Title Prefix", "5ff9b6ab-fbe3-490f-8980-71509263efe2")
	return talk
}

func TestLark(t *testing.T) {
	var talk = createLark()
	defer talk.Close()

	talk.PostMessage("Info title", "This is an info", ding.Info)
	talk.PostMessage("Warn title", "This is a warning", ding.Warn)
	talk.PostMessage("Error title", "This is an error", ding.Error)

	time.Sleep(time.Minute)
}

func TestLarkHook(t *testing.T) {
	var talk = createLark()

	var l = logo.NewLogger()
	defer l.Close()
	l.SetFuncCallDepth(4)

	var talkHook = ding.NewHook(talk, ding.WithFilterLevel(logo.LevelWarn))
	l.AddHook(talkHook)

	l.Info("This is info, but will not appear in lark.")

	// 这是一个batch消息，可以合并的
	for i := 0; i < 10; i++ {
		l.Warn("This warning will appear in lark.")
	}

	l.Error("This is an %q.", "error")

	// 测试huge message
	var content = make([]string, 0, 2048)
	for i := 0; i < 2048; i++ {
		content = append(content, strconv.Itoa(i))
	}

	l.Warn("HugeMessage: " + strings.Join(content, ","))

	time.Sleep(time.Minute)
}
