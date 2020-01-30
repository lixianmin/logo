package ding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lixianmin/logo"
	"io/ioutil"
	"net/http"
	"time"
)

/********************************************************************
created:    2019-10-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Talk struct {
	titlePrefix string
	token       string
	messageChan chan TalkMessage
	wc          *logo.WaitClose
}

func NewTalk(titlePrefix string, token string) *Talk {
	var talk = &Talk{
		titlePrefix: titlePrefix,
		token:       token,
		messageChan: make(chan TalkMessage, 32),
		wc:          logo.NewWaitClose(),
	}

	go talk.goLoop()
	return talk
}

func (talk *Talk) goLoop() {
	// 令牌桶发生器
	// 钉钉机器人发送频率限制是 20条/每分钟，如果超过限制，会返回 send too fast 错误信息，
	// 再发，就返回302错误，并限制发送10分钟

	const tokenFrequency = 3 * time.Second // 每3秒发一个令牌
	const maxBucket = 10

	var producerTicker = time.NewTicker(tokenFrequency)
	var checkTicker = time.NewTicker(500 * time.Millisecond)
	var bucketChan = make(chan struct{}, maxBucket)

	defer func() {
		producerTicker.Stop()
		checkTicker.Stop()
		close(bucketChan)
	}()

	// 格式化并直接发送消息
	var sendDirect = func(msg TalkMessage) {
		var text = msg.Text + "  \n  " + msg.Timestamp.Format(time.RFC3339)

		var title1 = fmt.Sprintf("[%s: %s] %s", msg.Level, talk.titlePrefix, msg.Title)
		var text1 = fmt.Sprintf("### %s  \n  %s", title1, text)
		_, _ = SendMarkdown(title1, text1, msg.Token)
	}

	for {
		select {
		case <-producerTicker.C:
			if len(bucketChan) < maxBucket {
				bucketChan <- struct{}{}
			}
		case <-checkTicker.C:
			for len(bucketChan) > 0 && len(talk.messageChan) > 0 {
				<-bucketChan
				sendDirect(<-talk.messageChan)
			}
		case <-talk.wc.CloseChan:
			return
		}
	}
}

func (talk *Talk) Close() error {
	return talk.wc.Close()
}

func (talk *Talk) SendInfo(title string, text string) {
	talk.sendMessage(title, text, "Info")
}

func (talk *Talk) SendWarn(title string, text string) {
	talk.sendMessage(title, text, "Warn")
}

func (talk *Talk) SendError(title string, text string) {
	talk.sendMessage(title, text, "Error")
}

func (talk *Talk) sendMessage(title string, text string, level string) {
	talk.messageChan <- TalkMessage{
		Level:     level,
		Title:     title,
		Text:      text,
		Timestamp: time.Now(),
		Token:     talk.token,
	}
}

func SendMarkdown(title string, text string, token string) ([]byte, error) {
	var message = MarkdownMessage{MsgType: "markdown", Markdown: MarkdownParams{Title: title, Text: text}}
	var content, err = json.Marshal(message)
	if err != nil {
		return nil, err
	}

	const webHook = "https://oapi.dingtalk.com/robot/send?access_token="
	var url = webHook + token
	response, err := http.Post(url, "application/json", bytes.NewBuffer(content))
	if err != nil {
		return nil, err
	}

	var body = response.Body
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	return bodyBytes, err
}
