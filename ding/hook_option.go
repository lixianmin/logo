package ding

import "github.com/lixianmin/logo"

/********************************************************************
created:    2020-02-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type hookOptions struct {
	FilterLevel int
}

type HookOption func(*hookOptions)

func WithFilterLevel(level int) HookOption {
	return func(options *hookOptions) {
		if level > logo.LevelNone {
			options.FilterLevel = level
		}
	}
}
