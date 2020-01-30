package logo

import (
	"path"
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

	var text = message.text
	my.buffer = append(my.buffer, text...)
	if len(text) == 0 || text[len(text)-1] != '\n' {
		my.buffer = append(my.buffer, '\n')
	}

	return my.buffer
}

// formatHeader writes log header to buffer in following order:
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided).
func (my *MessageFormatter) formatHeader(message Message, flag int) {
	var buffer = &my.buffer
	var t = time.Now()

	if flag&FlagDate != 0 {
		year, month, day := t.Date()
		itoa(buffer, year, 4)
		*buffer = append(*buffer, '-')
		itoa(buffer, int(month), 2)
		*buffer = append(*buffer, '-')
		itoa(buffer, day, 2)
		*buffer = append(*buffer, ' ')
	}

	if flag&FlagTime != 0 {
		hour, min, sec := t.Clock()
		itoa(buffer, hour, 2)
		*buffer = append(*buffer, ':')
		itoa(buffer, min, 2)
		*buffer = append(*buffer, ':')
		itoa(buffer, sec, 2)
		*buffer = append(*buffer, ' ')
	}

	// levelHints放到前面，[I], [W], [E]，因为它们的长度比较统一，容易对齐
	if flag&FlagLevel != 0 && my.levelHints != nil {
		var name = my.levelHints[message.level]
		*buffer = append(*buffer, name...)
		*buffer = append(*buffer, ' ')
	}

	if flag&(FlagLongFile|FlagShortFile) != 0 {
		var filePath = message.filePath
		if flag&FlagShortFile != 0 {
			filePath = path.Base(filePath)
		}

		*buffer = append(*buffer, filePath...)
		*buffer = append(*buffer, ':')
		itoa(buffer, message.lineNum, -1)
		*buffer = append(*buffer, ' ')
	}
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
