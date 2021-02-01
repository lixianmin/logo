package ding

import (
	"github.com/lixianmin/got/loom"
)

/********************************************************************
created:    2020-02-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type MessageQueue struct {
	mutex loom.Mutex
	buf   []Message
}

func (my *MessageQueue) Push(msg Message) {
	my.mutex.Lock()
	my.buf = append(my.buf, msg)
	my.mutex.Unlock()
}

func (my *MessageQueue) PopBatchMessage() (Message, int) {
	my.mutex.Lock()
	var buf = my.buf
	var first = buf[0]
	var batch = 1
	for i := 1; i < len(buf); i++ {
		var msg = buf[i]
		if msg.Text != first.Text || msg.Level != first.Level && msg.Title != first.Title {
			break
		}

		batch++
	}

	var newSize = len(buf) - batch
	for i := 0; i < newSize; i++ {
		buf[i] = buf[i+batch]
	}

	my.buf = buf[:newSize]
	my.mutex.Unlock()
	return first, batch
}

func (my *MessageQueue) Size() int {
	my.mutex.Lock()
	var size = len(my.buf)
	my.mutex.Unlock()

	return size
}
