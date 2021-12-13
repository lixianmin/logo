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
	startTime time.Time
	touchTime time.Time
	counter   int
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
			startTime: now,
			counter:   0,
		}

		my.cache[key] = item
	}

	item.touchTime = now
	var banTime = my.getBlockTime(item)
	var isBlocked = now.Before(banTime)

	var canPass = !isBlocked
	if canPass {
		item.counter += 1
	}

	// 每分钟移除一次过期的数据
	if now.After(my.nextRemoveTime) {
		my.checkRemoveExpired()
		my.nextRemoveTime = now.Add(time.Minute)
	}

	//fmt.Printf("canPass=%t, counter=%d, key=%s\n", canPass, item.counter, key)
	return isBlocked
}

func (my *MessageBlock) getBlockTime(item *blockItem) time.Time {
	var totalBlocked = my.step * time.Duration(math.Pow(2, float64(item.counter))-1)
	// 最长禁用时间
	const maxBlock = 60 * time.Minute
	if totalBlocked > maxBlock {
		totalBlocked = maxBlock
	}

	var blockTime = item.startTime.Add(totalBlocked)
	return blockTime
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
