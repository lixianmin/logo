package logo

/********************************************************************
created:    2021-01-18
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type loggerOptions struct {
	BufferSize int // 数据是否压缩
}

type LoggerOption func(*loggerOptions)

func WithBufferSize(size int) LoggerOption {
	return func(options *loggerOptions) {
		if size > 0 {
			options.BufferSize = size
		}
	}
}
