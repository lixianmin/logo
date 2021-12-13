package ding

import (
	"github.com/lixianmin/got/loom"
	"math"
	"sync"
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
	m     sync.Mutex
	step  time.Duration
}

func NewMessageBan() *MessageBan {
	var my = &MessageBan{
		cache: make(map[string]*BanItem, 8),
		step:  time.Second,
	}

	loom.Go(my.goLoop)
	return my
}

func (my *MessageBan) goLoop(later loom.Later) {
	var removeTicker = later.NewTicker(10 * time.Second)
	for {
		select {
		case <-removeTicker.C:
			my.checkRemoveExpired()
		}
	}
}

func (my *MessageBan) CheckBanned(message string) bool {
	my.m.Lock()
	defer my.m.Unlock()

	var now = time.Now()
	var item, ok = my.cache[message]
	if !ok {
		item = &BanItem{
			startTime: now,
			counter:   0,
		}

		my.cache[message] = item
	}

	item.touchTime = now
	var banTime = my.getBanTime(item)
	var isBanned = now.Before(banTime)

	var canSpeak = !isBanned
	if canSpeak {
		item.counter += 1
	}

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

func (my *MessageBan) checkRemoveExpired() {
	my.m.Lock()
	defer my.m.Unlock()

	const despairTime = 10 * time.Minute
	var expireTime = time.Now().Add(-despairTime)
	for message, item := range my.cache {
		if expireTime.After(item.touchTime) {
			delete(my.cache, message)
		}
	}
}
