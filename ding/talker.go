package ding

import "io"

/********************************************************************
created:    2021-12-13
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Talker interface {
	io.Closer
	PostMessage(level int, title string, text string)
}
