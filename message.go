package logo

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Message struct {
	text     string
	filePath string
	lineNum  int
	level    int
	stack    string
}

func (message Message) GetText() string {
	return message.text
}

func (message Message) GetFilePath() string {
	return message.filePath
}

func (message Message) GetLineNum() int {
	return message.lineNum
}

func (message Message) GetLevel() int {
	return message.level
}

func (message Message) GetStack() string {
	return message.stack
}
