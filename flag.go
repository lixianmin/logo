package logo

/********************************************************************
created:    2020-01-31
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Flag struct {
	flags int
}

func (my *Flag) AddFlag(flag int) {
	my.flags |= flag
}

func (my *Flag) RemoveFlag(flag int) {
	my.flags &= ^flag
}

func (my *Flag) HasFlag(flag int) bool {
	return (my.flags & flag) != 0
}
