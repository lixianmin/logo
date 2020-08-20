package logo

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type IHook interface {
	SetFilterLevel(level int)
	Write(message Message)
}
