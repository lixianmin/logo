package internal

/********************************************************************
created:    2021-12-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Message struct {
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

type Content struct {
	Post Post `json:"post"`
}

type Post struct {
	ZhCN ZhCN `json:"zh_cn"`
}

type ZhCN struct {
	Title   string   `json:"title"`
	Content [][]Item `json:"content"`
}

type Item struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
	Href string `json:"href,omitempty"`
}
