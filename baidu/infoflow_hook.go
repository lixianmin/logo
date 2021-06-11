package baidu

import (
	"fmt"
	"github.com/lixianmin/logo"
	"github.com/lixianmin/logo/ding"
	"github.com/lixianmin/logo/tools"
	"path"
	"strings"
)

/********************************************************************
created:    2020-05-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type InfoFlowHook struct {
	talker      *InfoFlow
	filterLevel int
	recoverable *ding.RecoverableError
}

func NewInfoFlowHook(talker *InfoFlow, opts ...InfoFlowHookOption) *InfoFlowHook {
	if talker == nil {
		panic("Talker should not be null")
	}

	// 默认值
	var options = infoFlowHookOptions{
		FilterLevel: logo.LevelInfo,
	}

	// 初始化
	for _, opt := range opts {
		opt(&options)
	}

	var my = &InfoFlowHook{
		talker:      talker,
		filterLevel: options.FilterLevel,
		recoverable: options.Recoverable,
	}

	return my
}

func (my *InfoFlowHook) Close() error {
	return my.talker.Close()
}

func (my *InfoFlowHook) SetFilterLevel(level int) {
	if level > logo.LevelNone && level < logo.LevelMax {
		my.filterLevel = level
	}
}

func (my *InfoFlowHook) Write(message logo.Message) {
	var level = message.GetLevel()
	if level < my.filterLevel {
		return
	}

	var text = message.GetText()
	// 如果属于recoverable error，并且判断不需要发送消息，则返回
	if my.recoverable != nil && !my.recoverable.NeedPostMessage(text) {
		return
	}

	var frames = message.GetFrames()
	if len(frames) > 0 {
		var buffer = make([]byte, 0, 128)
		for i := 1; i < len(frames); i++ {
			buffer = append(buffer, '\n')
			buffer = tools.AppendFrameInfo(buffer, frames[i])
		}

		var first = frames[0]
		text = fmt.Sprintf("%s:%d [%s()] %s %s", path.Base(first.File), first.Line, getFunctionName(first.Function), text, buffer)
	}

	var talker = my.talker
	switch level {
	case logo.LevelDebug:
		talker.PostDebug("", text)
	case logo.LevelInfo:
		talker.PostInfo("", text)
	case logo.LevelWarn:
		talker.PostWarn("", text)
	case logo.LevelError:
		talker.PostError("", text)
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
