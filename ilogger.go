package logo

/********************************************************************
created:    2020-06-13
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

type ILogger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}
