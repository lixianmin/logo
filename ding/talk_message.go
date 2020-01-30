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

type MarkdownParams struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type MarkdownMessage struct {
	MsgType  string         `json:"msgtype"`
	Markdown MarkdownParams `json:"markdown"`
}
