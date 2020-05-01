package baidu

import "time"

/********************************************************************
created:    2020-05-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type InfoFlowMessage struct {
	Level     string
	Title     string
	Text      string
	Timestamp time.Time
	Token     string
}

type MarkdownBody struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type MarkdownMessage struct {
	Body []MarkdownBody `json:"body"`
}

type Markdown struct {
	Message MarkdownMessage `json:"message"`
}
