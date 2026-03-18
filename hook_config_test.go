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

func TestWithFlag(t *testing.T) {
	var tests = []struct {
		name     string
		flag     int
		wantFlag int
	}{
		{"zero flag", 0, 0},
		{"single flag", FlagDate, FlagDate},
		{"combined flags", FlagDate | FlagTime | FlagShortFile, FlagDate | FlagTime | FlagShortFile},
		{"all flags", FlagDate | FlagTime | FlagLongFile | FlagShortFile | FlagLevel, FlagDate | FlagTime | FlagLongFile | FlagShortFile | FlagLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config HookConfig
			var opt = WithFlag(tt.flag)
			opt(&config)

			if config.Flag != tt.wantFlag {
				t.Errorf("WithFlag(%d): got Flag=%d, want %d", tt.flag, config.Flag, tt.wantFlag)
			}
		})
	}
}
