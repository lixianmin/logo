package ding

import (
	"fmt"
	"github.com/lixianmin/got/loom"
	"github.com/lixianmin/logo"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

/********************************************************************
created:    2020-05-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestMain(m *testing.M) {
	fmt.Println("setup: ding talk test")
	var retCode = m.Run()
	fmt.Println("teardown: ding talk test")
	os.Exit(retCode)
}

func createTalk() *Talk {
	var talk = NewTalk("Title Prefix", "ed23007fe73228e7d16c99504d1c39bf04c8cf8946d6d3a752ccf50158394650")
	return talk
}

func TestDingTalk(t *testing.T) {
	var talk = createTalk()
	defer talk.Close()

	talk.PostInfo("Info title", "This is an info")
	talk.PostWarn("Warn title", "This is a warning")
	talk.PostError("Error title", "This is an error")
}

func TestMessageBan(t *testing.T) {
	var talk = createTalk()
	defer talk.Close()

	for i := 0; i < 200; i++ {
		talk.PostInfo("Info title", "This is an info")
		time.Sleep(time.Second)
	}
}

func TestDingTalkHook(t *testing.T) {
	var talk = createTalk()

	var l = logo.NewLogger()
	defer l.Close()
	l.SetFuncCallDepth(4)

	var talkHook = NewTalkHook(talk, WithFilterLevel(logo.LevelWarn))
	l.AddHook(talkHook)

	l.Info("This is info, but will not appear in dingTalk.")

	// 这是一个batch消息，可以合并的
	for i := 0; i < 10; i++ {
		l.Warn("This warning will appear in dingTalk.")
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

func TestDumpIfPanic(t *testing.T) {
	defer loom.DumpIfPanic()
	var talk = createTalk()

	loom.Initialize(func(data []byte) {
		var message = string(data)
		talk.SendMessage("", message, Warn)
	})

	panic("hello")
}
