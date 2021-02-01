package ding

/********************************************************************
created:    2020-02-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type MarkdownParams struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type MarkdownMessage struct {
	MsgType  string         `json:"msgtype"`
	Markdown MarkdownParams `json:"markdown"`
}