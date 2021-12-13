package ding

import "github.com/lixianmin/logo"

/********************************************************************
created:    2021-03-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var levelNames = []string{"", "Debug", "Info", "Warn", "Error"}

func GetLevelName(level int) string {
	if level > logo.LevelNone && level < logo.LevelMax {
		return levelNames[level]
	}

	return ""
}
