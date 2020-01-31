package logo

import "testing"

/********************************************************************
created:    2020-01-31
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/


func TestAddFlag(t *testing.T) {
	var f = Flag{}
	f.AddFlag(1)
	f.AddFlag(1)

	if !f.HasFlag(1) {
		t.Failed()
	}

	f.RemoveFlag(1)
	if f.HasFlag(1) {
		t.Failed()
	}
}