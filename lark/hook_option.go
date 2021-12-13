package lark

import (
	"github.com/lixianmin/logo"
)

/********************************************************************
created:    2021-12-09
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
