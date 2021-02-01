package ding

import "github.com/lixianmin/logo"

/********************************************************************
created:    2020-02-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type talkHookOptions struct {
	FilterLevel int
}

type TalkHookOption func(*talkHookOptions)

func WithFilterLevel(level int) TalkHookOption {
	return func(options *talkHookOptions) {
		if level > logo.LevelNone {
			options.FilterLevel = level
		}
	}
}
