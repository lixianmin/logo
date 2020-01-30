package logo

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

const (
	LevelInfo  = 1
	LevelWarn  = 2
	LevelError = 3

	FlagDate      = 0x0001
	FlagTime      = 0x0002
	FlagLongFile  = 0x0004
	FlagShortFile = 0x0008
	FlagLevel     = 0x0010
)

var levelHints = []string{"", "[I]", "[W]", "[E]"}
var levelHintsConsole = []string{"1;37", "1;34", "1;33", "1;31"}

func init() {
	pre := "\033["
	reset := "\033[0m"

	for i := 0; i < len(levelHintsConsole); i++ {
		var color = levelHintsConsole[i]
		levelHintsConsole[i] = pre + color + "m" + levelHints[i] + reset
	}
}