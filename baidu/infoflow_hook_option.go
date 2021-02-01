package baidu

import "github.com/lixianmin/logo"

/********************************************************************
created:    2020-02-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type infoFlowHookOptions struct {
	FilterLevel int
}

type InfoFlowHookOption func(*infoFlowHookOptions)

func WithFilterLevel(level int) InfoFlowHookOption {
	return func(options *infoFlowHookOptions) {
		if level > logo.LevelNone {
			options.FilterLevel = level
		}
	}
}
