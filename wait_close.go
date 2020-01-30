package logo

import (
	"sync"
	"sync/atomic"
)

/********************************************************************
created:    2018-10-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type WaitClose struct {
	CloseChan chan struct{}
	m         sync.Mutex
	done      int32
}

func NewWaitClose() *WaitClose {
	var wd = &WaitClose{
		CloseChan: make(chan struct{}),
	}

	return wd
}

func (wc *WaitClose) Close() error {
	if 0 == atomic.LoadInt32(&wc.done) {
		wc.m.Lock()
		if 0 == wc.done {
			atomic.StoreInt32(&wc.done, 1)
			close(wc.CloseChan)
		}
		wc.m.Unlock()
	}

	return nil
}

func (wc *WaitClose) IsClosed() bool {
	return atomic.LoadInt32(&wc.done) == 1
}
