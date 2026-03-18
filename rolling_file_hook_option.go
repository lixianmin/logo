package logo

import (
	"time"
)

/********************************************************************
created:    2026-03-18
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func WithHookFlag(flag int) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		WithFlag(flag)(&options.HookConfig)
	}
}

func WithHookFilterLevel(level int) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		WithFilterLevel(level)(&options.HookConfig)
	}
}

func WithDirName(dirName string) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		if dirName != "" {
			options.DirName = dirName
		}
	}
}

func WithFileNamePrefix(prefix string) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		options.FileNamePrefix = prefix
	}
}

func WithMaxFileSize(size int64) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		if size > 0 {
			options.MaxFileSize = size
		}
	}
}

func WithExpireTime(duration time.Duration) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		if duration > 0 {
			options.ExpireTime = duration
		}
	}
}

func WithCheckRollingInterval(interval int64) RollingFileHookOption {
	return func(options *rollingFileHookOptions) {
		if interval > 0 {
			options.CheckRollingInterval = interval
		}
	}
}
