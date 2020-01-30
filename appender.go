package logo

import (
	"io"
)

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Appender interface {
	io.Closer
	Write(message Message)
}