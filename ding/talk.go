package ding

import (
	"bytes"
	"context"
	"fmt"
	"github.com/lixianmin/got/convert"
	"github.com/lixianmin/got/timex"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"
)

/********************************************************************
created:    2019-10-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Talk struct {
	titlePrefix  string
	token        string
	cancel       context.CancelFunc
	sendingCount int32
	messageQueue MessageQueue
}

func NewTalk(titlePrefix string, token string) *Talk {
	if token == "" {
		panic("token should not be empty")
	}

	var ctx, cancel = context.WithCancel(context.Background())
	var talk = &Talk{
		titlePrefix: titlePrefix,
		token:       token,
		cancel:      cancel,
	}

	go talk.goLoop(ctx)
	return talk
}

func (talk *Talk) goLoop(ctx context.Context) {
	// 令牌桶发生器
	// 钉钉机器人发送频率限制是 20条/每分钟，如果超过限制，会返回 send too fast 错误信息，
	// 再发，就返回302错误，并限制发送10分钟

	const tokenFrequency = 3 * time.Second // 每3秒发一个令牌
	const maxBucket = 10

	var producerTicker = time.NewTicker(tokenFrequency)
	var checkTicker = time.NewTicker(1 * time.Second)
	var bucketChan = make(chan struct{}, maxBucket)

	// 预先准备一个bucket
	bucketChan <- struct{}{}

	defer func() {
		producerTicker.Stop()
		checkTicker.Stop()
		close(bucketChan)
	}()

	// 格式化并直接发送消息
	var sendDirect = func(msg Message, batch int) {
		atomic.AddInt32(&talk.sendingCount, int32(-batch))
		const layout = "2006-01-02 15:04:05"
		var text = msg.Text + "  \n  " + msg.Timestamp.Format(layout)

		var title1 = fmt.Sprintf("[%s(%d) %s] %s", msg.Level, batch, talk.titlePrefix, msg.Title)
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
			for len(bucketChan) > 0 && talk.messageQueue.Size() > 0 {
				<-bucketChan

				var msg, batch = talk.messageQueue.PopBatchMessage()
				sendDirect(msg, batch)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (talk *Talk) Close() error {
	// 给一点点flush的时间：dingTalk是3s一个令牌，发的很慢的，因此在单元测试里不一定能flush成功
	// 本来想给一点点flush的时间，但是进程结束时并不允许这样的等待时间
	//const timeout = 1 * time.Second
	//const step = 50 * time.Millisecond
	//for i := 0; talk.sendingCount > 0 && i < int(timeout/step); i++ {
	//	time.Sleep(step)
	//}

	talk.cancel()
	return nil
}

func (talk *Talk) PostDebug(title string, text string) {
	talk.postMessage(title, text, Debug)
}

func (talk *Talk) PostInfo(title string, text string) {
	talk.postMessage(title, text, Info)
}

func (talk *Talk) PostWarn(title string, text string) {
	talk.postMessage(title, text, Warn)
}

func (talk *Talk) PostError(title string, text string) {
	talk.postMessage(title, text, Error)
}

func (talk *Talk) postMessage(title string, text string, level string) {
	atomic.AddInt32(&talk.sendingCount, 1)

	var msg = Message{
		Level:     level,
		Title:     title,
		Text:      text,
		Timestamp: time.Now(),
		Token:     talk.token,
	}

	talk.messageQueue.Push(msg)
}

func (talk *Talk) SendMessage(title string, text string, level string) {
	var text1 = text + "  \n  " + time.Now().Format(timex.Layout)

	const batch = 1
	var title1 = fmt.Sprintf("[%s(%d) %s] %s", level, batch, talk.titlePrefix, title)
	var text2 = fmt.Sprintf("### %s  \n  %s", title1, text1)
	_, _ = SendMarkdown(title1, text2, talk.token)
}

func (talk *Talk) GetTitlePrefix() string {
	return talk.titlePrefix
}

func SendMarkdown(title string, text string, token string) ([]byte, error) {
	var message = MarkdownMessage{MsgType: "markdown", Markdown: MarkdownParams{Title: title, Text: text}}
	var data, err = convert.ToJsonE(message)
	if err != nil {
		return nil, err
	}

	const webHook = "https://oapi.dingtalk.com/robot/send?access_token="
	var url = webHook + token

	// 裁剪待发送消息体的最大长度
	const cutLength = 1024
	if len(data) > cutLength {
		data = append(data[:cutLength], "..."...)
	}

	// 发送
	var sending = bytes.NewBuffer(data)
	response, err := http.Post(url, "application/json", sending)
	if err != nil {
		return nil, err
	}

	var body = response.Body
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	return bodyBytes, err
}
