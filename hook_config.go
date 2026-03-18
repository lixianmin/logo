package logo

/********************************************************************
created:    2026-03-18
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type HookConfig struct {
	Flag        int
	FilterLevel int
}

type HookOption func(*HookConfig)

func WithFlag(flag int) HookOption {
	return func(config *HookConfig) {
		config.Flag = flag
	}
}

func WithFilterLevel(level int) HookOption {
	return func(config *HookConfig) {
		if level > LevelNone {
			config.FilterLevel = level
		}
	}
}
