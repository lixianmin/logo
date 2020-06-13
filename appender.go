package logo

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Appender interface {
	SetFilterLevel(level int)
	Write(message Message)
}
