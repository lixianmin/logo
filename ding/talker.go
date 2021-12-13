package ding

import "io"

/********************************************************************
created:    2021-12-13
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Talker interface {
	io.Closer
	PostMessage(title string, text string, level string)
}
