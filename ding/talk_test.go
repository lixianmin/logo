package ding

import (
	"testing"
	"time"
)

/********************************************************************
created:    2020-05-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestDingTalk(t *testing.T) {
	var talk = NewTalk("title前缀", "ed23007fe73228e7d16c99504d1c39bf04c8cf8946d6d3a752ccf50158394650")
	talk.SendInfo("info", "这里是正文")
	talk.SendWarn("warn", "这里是正文")
	talk.SendError("error", "这里是正文")

	for len(talk.messageChan) > 0 {
		time.Sleep(time.Second)
	}

	talk.Close()
}
