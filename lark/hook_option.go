package lark

import (
	"github.com/lixianmin/logo"
	"github.com/lixianmin/logo/ding"
	"time"
)

/********************************************************************
created:    2021-12-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type hookOptions struct {
	FilterLevel int
	Recoverable *ding.RecoverableError
}

type HookOption func(*hookOptions)

func WithFilterLevel(level int) HookOption {
	return func(options *hookOptions) {
		if level > logo.LevelNone {
			options.FilterLevel = level
		}
	}
}

// WithRecoverableError 对可自动恢复的错误消息：如果只是interval之内收到，则不会发送；如果持续时间超过interval，则会发送
func WithRecoverableError(interval time.Duration, checker ding.RecoverableChecker) HookOption {
	return func(options *hookOptions) {
		options.Recoverable = ding.NewRecoverableError(interval, checker)
	}
}
