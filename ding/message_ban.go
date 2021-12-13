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

type BanItem struct {
	startTime time.Time
	touchTime time.Time
	counter   int
}

type MessageBan struct {
	cache map[string]*BanItem
	step  time.Duration
}

func NewMessageBan() *MessageBan {
	var my = &MessageBan{
		cache: make(map[string]*BanItem, 8),
		step:  time.Second,
	}

	return my
}

func (my *MessageBan) CheckBanned(key string) bool {
	var now = time.Now()
	var item, ok = my.cache[key]
	if !ok {
		item = &BanItem{
			startTime: now,
			counter:   0,
		}

		my.cache[key] = item
	}

	item.touchTime = now
	var banTime = my.getBanTime(item)
	var isBanned = now.Before(banTime)

	var canSpeak = !isBanned
	if canSpeak {
		item.counter += 1
	}

	//fmt.Printf("canSpeak=%t, counter=%d, key=%s\n", canSpeak, item.counter, key)
	return isBanned
}

func (my *MessageBan) getBanTime(item *BanItem) time.Time {
	var totalBanned = my.step * time.Duration(math.Pow(2, float64(item.counter))-1)
	// 最长禁用时间
	const maxBanned = 60 * time.Minute
	if totalBanned > maxBanned {
		totalBanned = maxBanned
	}

	var speakTime = item.startTime.Add(totalBanned)
	return speakTime
}

func (my *MessageBan) CheckRemoveExpired() {
	const despairTime = 10 * time.Minute
	var expireTime = time.Now().Add(-despairTime)
	for message, item := range my.cache {
		if expireTime.After(item.touchTime) {
			delete(my.cache, message)
		}
	}
}
