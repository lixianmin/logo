package baidu

import (
	"testing"
	"time"
)

/********************************************************************
created:    2020-05-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestInfoFlow(t *testing.T) {
	var talk = NewInfoFlow("title前缀", "d2127efaae4cfd31e53e47d919c27ad5c")
	talk.SendWarn("warn", "这里是正文")

	time.Sleep(10 * time.Second)
}
