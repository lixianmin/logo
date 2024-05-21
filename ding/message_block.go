package ding

import (
	"math"
	"time"
)

/********************************************************************
created:    2021-12-13
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type blockItem struct {
	blockTime time.Time
	touchTime time.Time
	counter   int
}

func (item *blockItem) incCounter(step time.Duration) {
	item.counter += 1

	var n = float64(item.counter)
	var candidate = step * time.Duration(math.Pow(2, n))

	// 最长禁用时间
	const maxStep = 60 * time.Minute
	if candidate > maxStep {
		candidate = maxStep
	}

	item.blockTime = item.blockTime.Add(candidate)
}

type MessageBlock struct {
	cache          map[string]*blockItem
	step           time.Duration
	nextRemoveTime time.Time
}

func NewMessageBlock() *MessageBlock {
	var my = &MessageBlock{
		cache:          make(map[string]*blockItem, 8),
		step:           time.Second,
		nextRemoveTime: time.Now(),
	}

	return my
}

func (my *MessageBlock) CheckBlocked(key string) bool {
	var now = time.Now()
	var item, ok = my.cache[key]
	if !ok {
		item = &blockItem{
			blockTime: now,
			counter:   0,
		}

		my.cache[key] = item
	}

	item.touchTime = now
	var blockTime = item.blockTime
	var isBlocked = now.Before(blockTime)

	var canPass = !isBlocked
	if canPass {
		item.incCounter(my.step)
	}

	// 每分钟移除一次过期的数据
	if now.After(my.nextRemoveTime) {
		my.checkRemoveExpired()
		my.nextRemoveTime = now.Add(time.Minute)
	}

	// 如果发送太频繁, 则会被block一段时间, 这时先不打印日志了
	//fmt.Printf("canPass=%t=[blockTime(%s) < now(%s)], delta=%q, counter=%d, key=%q\n", canPass, timex.FormatTime(blockTime), timex.FormatTime(now),
	//	timex.FormatDuration(now.Sub(blockTime)), item.counter, key)
	return isBlocked
}

func (my *MessageBlock) checkRemoveExpired() {
	const despairTime = 10 * time.Minute
	var expireTime = time.Now().Add(-despairTime)
	for message, item := range my.cache {
		if expireTime.After(item.touchTime) {
			delete(my.cache, message)
		}
	}
}
