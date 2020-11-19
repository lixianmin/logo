package logo

/********************************************************************
created:    2020-11-19
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type loggerFetus struct {
	hooks []IHook
}

func (fetus *loggerFetus) WriteMessage(message Message) {
	for _, hook := range fetus.hooks {
		if hook != nil {
			hook.Write(message)
		}
	}
}

func (fetus *loggerFetus) FlushMessage(messageChan chan Message) {
	var count = len(messageChan)
	for i := 0; i < count; i++ {
		var message = <-messageChan
		fetus.WriteMessage(message)
	}
}
