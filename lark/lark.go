package lark

import (
	"bytes"
	"context"
	"fmt"
	"github.com/lixianmin/got/convert"
	"github.com/lixianmin/got/timex"
	"github.com/lixianmin/logo/ding"
	"github.com/lixianmin/logo/lark/internal"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"
)

/********************************************************************
created:    2021-12-09
author:     lixianmin

飞书机器人文档: https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN


Copyright (C) - All Rights Reserved
*********************************************************************/

type Lark struct {
	titlePrefix  string
	token        string
	cancel       context.CancelFunc
	sendingCount int32
	messageQueue ding.MessageQueue
}

func NewLark(titlePrefix string, token string) *Lark {
	if token == "" {
		panic("token should not be empty")
	}

	var ctx, cancel = context.WithCancel(context.Background())
	var talk = &Lark{
		titlePrefix: titlePrefix,
		token:       token,
		cancel:      cancel,
	}

	go talk.goLoop(ctx)
	return talk
}

func (talk *Lark) goLoop(ctx context.Context) {
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
	var block = ding.NewMessageBlock()

	// 预先准备一个bucket
	bucketChan <- struct{}{}

	defer func() {
		producerTicker.Stop()
		checkTicker.Stop()
		close(bucketChan)
	}()

	// 格式化并直接发送消息
	var sendDirect = func(msg ding.Message, batch int) {
		var title1 = fmt.Sprintf("[%s(%d) %s] %s", msg.Level, batch, talk.titlePrefix, msg.Title)
		var key = title1 + msg.Text
		if !block.CheckBlocked(key) {
			atomic.AddInt32(&talk.sendingCount, int32(-batch))
			var text = msg.Text + "\n" + msg.Timestamp.Format(timex.Layout)

			if _, err := SendPost(title1, text, msg.Token); err != nil {
				fmt.Printf("err=%q\n", err)
			}
		}
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

func (talk *Lark) Close() error {
	// 本来想给一点点flush的时间，但是进程结束时并不允许这样的等待时间
	//const timeout = 1 * time.Second
	//const step = 50 * time.Millisecond
	//for i := 0; talk.sendingCount > 0 && i < int(timeout/step); i++ {
	//	time.Sleep(step)
	//}

	talk.cancel()
	return nil
}

func (talk *Lark) PostMessage(level int, title string, text string) {
	atomic.AddInt32(&talk.sendingCount, 1)

	var msg = ding.Message{
		Level:     ding.GetLevelName(level),
		Title:     title,
		Text:      text,
		Timestamp: time.Now(),
		Token:     talk.token,
	}

	talk.messageQueue.Push(msg)
}

func (talk *Lark) SendMessage(level int, title string, text string) {
	var text1 = text + "\n" + time.Now().Format(timex.Layout)

	const batch = 1
	var levelName = ding.GetLevelName(level)
	var title1 = fmt.Sprintf("[%s(%d) %s] %s", levelName, batch, talk.titlePrefix, title)
	if _, err := SendPost(title1, text1, talk.token); err != nil {
		fmt.Printf("err=%q\n", err)
	}
}

func (talk *Lark) GetTitlePrefix() string {
	return talk.titlePrefix
}

func SendPost(title string, text string, token string) ([]byte, error) {
	var message = internal.Message{
		MsgType: "post",
		Content: internal.Content{Post: internal.Post{
			ZhCN: internal.ZhCN{
				Title: title,
				Content: [][]internal.Item{
					{{Tag: "text", Text: text}},
				}},
		},
		},
	}
	var data, err = convert.ToJsonE(message)
	if err != nil {
		return nil, err
	}

	const webHook = "https://open.feishu.cn/open-apis/bot/v2/hook/"
	var url = webHook + token

	// 裁剪待发送消息体的最大长度
	const cutLength = 1024
	if len(data) > cutLength {
		data = append(data[:cutLength], "..."...)
	}

	// 发送
	var sending = bytes.NewBuffer(data)
	response, err := http.Post(url, "application/json; charset=utf-8", sending)
	if err != nil {
		return nil, err
	}

	var body = response.Body
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	return bodyBytes, err
}
