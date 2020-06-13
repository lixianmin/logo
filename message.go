package logo

import "runtime"

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Message struct {
	text   string
	level  int
	frames []runtime.Frame
}

func (message Message) GetText() string {
	return message.text
}

func (message Message) GetLevel() int {
	return message.level
}

func (message Message) GetFrames() []runtime.Frame {
	return message.frames
}
