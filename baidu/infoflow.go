package baidu

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
created:    2020-05-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type InfoFlow struct {
	titlePrefix string
	token       string
	messageChan chan InfoFlowMessage
	wc          *logo.WaitClose
}

func NewInfoFlow(titlePrefix string, token string) *InfoFlow {
	if token == "" {
		panic("token should not be empty")
	}

	var talk = &InfoFlow{
		titlePrefix: titlePrefix,
		token:       token,
		messageChan: make(chan InfoFlowMessage, 32),
		wc:          logo.NewWaitClose(),
	}

	go talk.goLoop()
	return talk
}

func (talk *InfoFlow) goLoop() {
	// https://static.im.baidu.com/robotSetDoc/index.html
	// 消息发送频率限制
	// 为了保障群成员使用体验，以防收到大量消息的打扰，机器人发消息限制200条/分钟，超出后，将限流10分钟。
	// 限流期间的消息将会被丢弃。
	// 限流针对机器人和群聊两个元素。即：机器人在群A限流，仍可以向B群发送消息；群里R1机器人限流，R2机器人仍然可以发消息到本群。

	// 令牌桶发生器
	const tokenFrequency = 300 * time.Millisecond // 每300毫秒发一个令牌
	const maxBucket = 20

	var producerTicker = time.NewTicker(tokenFrequency)
	var checkTicker = time.NewTicker(500 * time.Millisecond)
	var bucketChan = make(chan struct{}, maxBucket)

	defer func() {
		producerTicker.Stop()
		checkTicker.Stop()
		close(bucketChan)
	}()

	// 格式化并直接发送消息
	var sendDirect = func(msg InfoFlowMessage) {
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

func (talk *InfoFlow) Close() error {
	return talk.wc.Close()
}

func (talk *InfoFlow) SendInfo(title string, text string) {
	talk.sendMessage(title, text, "Info")
}

func (talk *InfoFlow) SendWarn(title string, text string) {
	talk.sendMessage(title, text, "Warn")
}

func (talk *InfoFlow) SendError(title string, text string) {
	talk.sendMessage(title, text, "Error")
}

func (talk *InfoFlow) sendMessage(title string, text string, level string) {
	talk.messageChan <- InfoFlowMessage{
		Level:     level,
		Title:     title,
		Text:      text,
		Timestamp: time.Now(),
		Token:     talk.token,
	}
}

func SendMarkdown(title string, text string, token string) ([]byte, error) {
	var content = "####" + title + "\n" + text
	var message = MarkdownMessage{Type: "MD", Content: content}
	var data, err = json.Marshal(message)
	if err != nil {
		return nil, err
	}

	const webHook = "http://apiin.im.baidu.com/api/msg/groupmsgsend?access_token="
	var url = webHook + token
	response, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var body = response.Body
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	return bodyBytes, err
}
