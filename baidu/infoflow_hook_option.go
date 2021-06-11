package baidu

import (
	"github.com/lixianmin/logo"
	"github.com/lixianmin/logo/ding"
	"time"
)

/********************************************************************
created:    2020-02-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type infoFlowHookOptions struct {
	FilterLevel int
	Recoverable *ding.RecoverableError
}

type InfoFlowHookOption func(*infoFlowHookOptions)

func WithFilterLevel(level int) InfoFlowHookOption {
	return func(options *infoFlowHookOptions) {
		if level > logo.LevelNone {
			options.FilterLevel = level
		}
	}
}

// WithRecoverableError 对可自动恢复的错误消息：如果只是interval之内收到，则不会发送；如果持续时间超过interval，则会发送
func WithRecoverableError(checker ding.RecoverableChecker, interval time.Duration) InfoFlowHookOption {
	return func(options *infoFlowHookOptions) {
		options.Recoverable = ding.NewRecoverableError(checker, interval)
	}
}
