package logo

import (
	"testing"
)

/********************************************************************
created:    2026-03-18
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestHookConfigDefaults(t *testing.T) {
	var config HookConfig

	if config.Flag != 0 {
		t.Errorf("expected default Flag=0, got %d", config.Flag)
	}

	if config.FilterLevel != 0 {
		t.Errorf("expected default FilterLevel=0, got %d", config.FilterLevel)
	}
}
