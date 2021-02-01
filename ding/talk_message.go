package ding

import "time"

/********************************************************************
created:    2019-10-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type TalkMessage struct {
	Level     string
	Title     string
	Text      string
	Timestamp time.Time
	Token     string
}