package lark

import (
	"testing"
)

/********************************************************************
created:    2021-12-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestNewLark(t *testing.T) {
	var token = "5ff9b6ab-fbe3-490f-8980-71509263efe2"
	var lark = NewLark("hello", token)
	lark.SendMessage("title", "text", "info")
}
