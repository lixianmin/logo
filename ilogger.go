package logo

/********************************************************************
created:    2020-06-13
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

type ILogger interface {
	Debug(first interface{}, args ...interface{})
	Info(first interface{}, args ...interface{})
	Warn(first interface{}, args ...interface{})
	Error(first interface{}, args ...interface{})
}
