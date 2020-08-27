package logo

import (
	"github.com/lixianmin/logo/tools"
	"path"
	"strings"
	"time"
)

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type MessageFormatter struct {
	flag       int
	levelHints []string
	buffer     []byte // for accumulating text to write
}

func newMessageFormatter(flag int, levelHints []string) *MessageFormatter {
	var formatter = &MessageFormatter{
		flag:       flag,
		levelHints: levelHints,
		buffer:     nil,
	}

	return formatter
}

func (my *MessageFormatter) format(message Message) []byte {
	my.buffer = my.buffer[:0]
	my.formatHeader(message, my.flag)

	var buffer = append(my.buffer, message.text...)

	var frames = message.GetFrames()
	for i := 1; i < len(frames); i++ {
		buffer = append(buffer, '\n')
		buffer = tools.AppendFrameInfo(buffer, frames[i])
	}

	if len(buffer) == 0 || buffer[len(buffer)-1] != '\n' {
		buffer = append(buffer, '\n')
	}

	var withStack = len(frames) > 1
	if withStack {
		buffer = append(buffer, '\n')
	}

	my.buffer = buffer
	return buffer
}

// formatHeader writes log header to buffer in following order:
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided).
func (my *MessageFormatter) formatHeader(message Message, flag int) {
	var buffer = &my.buffer
	var t = time.Now()

	var hasTimeFlag = flag&FlagTime != 0
	if flag&FlagDate != 0 {
		year, month, day := t.Date()
		tools.Itoa(buffer, year, 4)
		*buffer = append(*buffer, '-')
		tools.Itoa(buffer, int(month), 2)
		*buffer = append(*buffer, '-')
		tools.Itoa(buffer, day, 2)

		// 之所以使用'T'而不是空格，原因是为了方便日志搜索
		var c byte = 'T'
		if !hasTimeFlag {
			c = ' '
		}

		*buffer = append(*buffer, c)
	}

	if hasTimeFlag {
		hour, min, sec := t.Clock()
		tools.Itoa(buffer, hour, 2)
		*buffer = append(*buffer, ':')
		tools.Itoa(buffer, min, 2)
		*buffer = append(*buffer, ':')
		tools.Itoa(buffer, sec, 2)
		*buffer = append(*buffer, ' ')
	}

	// levelHints放到前面，[I], [W], [E]，因为它们的长度比较统一，容易对齐
	if flag&FlagLevel != 0 && my.levelHints != nil {
		var name = my.levelHints[message.level]
		*buffer = append(*buffer, name...)
		*buffer = append(*buffer, ' ')
	}

	if flag&(FlagLongFile|FlagShortFile) != 0 && len(message.frames) > 0 {
		var first = message.frames[0]
		var filePath = first.File
		if flag&FlagShortFile != 0 {
			filePath = path.Base(filePath)
		}

		*buffer = append(*buffer, filePath...)
		*buffer = append(*buffer, ':')
		tools.Itoa(buffer, first.Line, -1)
		*buffer = append(*buffer, ' ')

		if first.Function != "" {
			*buffer = append(*buffer, '[')
			*buffer = append(*buffer, getFunctionName(first.Function)...)
			*buffer = append(*buffer, '(', ')', ']', ' ')
		}
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
