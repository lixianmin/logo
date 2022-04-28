package tools

import (
	"runtime"
	"strings"
)

/********************************************************************
created:    2020-06-13
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func CallersFrames(skip int, fullStack bool) []runtime.Frame {
	const depth = 16
	var pcs [depth]uintptr // 程序计算器
	var total = runtime.Callers(skip, pcs[:])

	var fetch = total
	if !fullStack && total > 1 {
		fetch = 1
	}

	var frames = runtime.CallersFrames(pcs[:fetch])
	var results = make([]runtime.Frame, 0, fetch)
	for {
		var frame, more = frames.Next()
		results = append(results, frame)
		if !more {
			break
		}
	}

	return results
}

func AppendFrameInfo(buffer []byte, frame runtime.Frame) []byte {
	buffer = append(buffer, frame.File...)
	buffer = append(buffer, ':')
	Itoa(&buffer, frame.Line, -1)

	if frame.Function != "" {
		buffer = append(buffer, ' ')
		buffer = append(buffer, trimUrlPath(frame.Function)...)
		buffer = append(buffer, '(', ')')
	}

	return buffer
}

// frame.Function可能会是 github.com/lixianmin/logo.TestConsoleHook, 然后提取完成后会变成logo.TestConsoleHook, 包含package+object+function名
func trimUrlPath(function string) string {
	if function != "" {
		var lastIndex = strings.LastIndexByte(function, '/')
		if lastIndex > 0 {
			var s = function[lastIndex+1:]
			//fmt.Printf("function=%s, s = %s\n\n", function, s)
			return s
		}
	}

	return function
}

func Itoa(buf *[]byte, i int, wid int) {
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
