package ding

import "time"

/********************************************************************
created:    2021-06-11
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type RecoverableChecker func(message string) bool

type RecoverableError struct {
	isRecoverableError RecoverableChecker
	interval           time.Duration
	lastTime           time.Time
}

func NewRecoverableError(checker RecoverableChecker, interval time.Duration) *RecoverableError {
	if checker == nil {
		checker = func(message string) bool {
			return false
		}
	}

	var my = &RecoverableError{
		isRecoverableError: checker,
		interval:           interval,
	}

	return my
}

func (my *RecoverableError) NeedPostMessage(message string) bool {
	if my.isRecoverableError(message) {
		var now = time.Now()
		var pastTime = now.Sub(my.lastTime)

		if pastTime > 2*my.interval {
			// 第一次，或长时间没有遇到这个error，不发送提醒消息，只记录时间
			my.lastTime = now
			return false
		} else if pastTime > my.interval {
			// 1分钟之前刚刚遇到过，发送提醒消息，记录时间
			my.lastTime = now
			return true
		} else {
			// 刚刚遇到的，不发送提醒消息，不修改时间
			return false
		}
	}

	return true
}
