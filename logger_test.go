package logo

import (
	"fmt"
	"github.com/lixianmin/got/loom"
	"github.com/lixianmin/got/randx"
	"math"
	"strconv"
	"sync"
	"testing"
	"time"
)

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestConsoleHook(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(4)

	var console = NewConsoleHook(ConsoleHookArgs{Flag: flag, FilterLevel: LevelDebug})
	l.AddHook(console)

	l.Debug("Debug", "Message")
	l.Info("This is info", "info args")
	l.Warn("I am a warning")
	l.Error("Error occurred")

	l.Close()
}

func TestRollingFileHook(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(4)

	var fileHook = NewRollingFileHook(RollingFileHookArgs{
		FilterLevel:          LevelWarn,
		Flag:                 flag,
		MaxFileSize:          16,
		CheckRollingInterval: 10,
	})

	l.AddHook(fileHook)

	// 测试archive中生成的文件名
	for i := 0; i < 200; i++ {
		l.Info("This is info")
		l.Warn("I am a warning")
		l.Error("Error occurred")
	}

	l.Close()
}

func TestFileHookFilterLevel(t *testing.T) {
	var l = NewLogger()

	const flag = FlagDate | FlagTime | FlagShortFile | FlagLevel
	l.SetFuncCallDepth(4)

	var fileHook = NewRollingFileHook(RollingFileHookArgs{
		FilterLevel: LevelWarn,
		Flag:        flag,
		MaxFileSize: 16,
	})

	l.AddHook(fileHook)

	l.Info("1 warn level: This is info")
	l.Warn("2 warn level: I am a warning")
	l.Error("3 warn level: Error occurred")

	if fileHook.args.FilterLevel != LevelWarn {
		t.Fatal()
	}

	fileHook.SetFilterLevel(LevelInfo)
	if fileHook.args.FilterLevel != LevelInfo {
		t.Fatal()
	}

	l.Info("4 info level: This is info")
	l.Warn("5 info level: I am a warning")
	l.Error("6 info level: Error occurred")

	_ = l.Close()
}

//func TestLogAnyObject(t *testing.T) {
//	Info(123.45678)
//	Info(t)
//}

func TestAutoFlush(t *testing.T) {
	var logger = GetLogger().(*Logger)
	logger.AddFlag(LogAsyncWrite)
	var i = 0
	for i < 10 {
		Info(strconv.Itoa(i))
		i++
	}

	logger.RemoveFlag(LogAsyncWrite)
	for i < 20 {
		Warn(strconv.Itoa(i))
		i++
	}

	logger.AddFlag(LogAsyncWrite)
	for i < 30 {
		Error(strconv.Itoa(i))
		i++
	}

	logger.Flush()
}

func TestClose(t *testing.T) {
	var logger = GetLogger().(*Logger)
	logger.AddFlag(LogAsyncWrite)

	for i := 0; i < 50; i++ {
		logger.Info(strconv.Itoa(i))
	}

	_ = logger.Close()

	logger.Info("closed")
}

func TestJson(t *testing.T) {
	type Pig struct {
		Weight  int32     `json:"weight"`
		Birth   time.Time `json:"birth"`
		Name    string    `json:"name"`
		Message *Message  `json:"message"`
	}

	var pig = Pig{
		Weight: 135,
		Birth:  time.Now(),
		Name:   "panda",
	}

	theLogger.(*Logger).SetFilterLevel(LevelDebug)

	JsonD("int", math.MinInt64)
	JsonD("int8", int8(math.MinInt8))
	JsonD("int16", int16(math.MinInt16))
	JsonD("int32", int32(math.MinInt32))
	JsonD("int64", int64(math.MinInt64))

	JsonI("uint8", uint8(math.MaxUint8))
	JsonI("uint16", uint16(math.MaxUint16))
	JsonI("uint32", uint32(math.MaxUint32))
	JsonI("uint64", uint64(math.MaxUint64))

	JsonW("nil", nil)
	JsonW("bool", true, "bool", false)
	JsonW("float32", float32(1.234), "float64", 10.29)
	JsonW("string", "lixianmin\"' \t\r\n 你好啊小朋友")
	JsonW("bytes", []byte("this is a byte buffer"))
	JsonW("time", time.Now())
	JsonW("error", fmt.Errorf("this is an error: %d", 1029))

	JsonE("age", 10, "pig", pig)
	JsonE("slice", []Pig{pig, pig})
	JsonE("map", map[int]string{1: "hello", 2: "world"})

	JsonI(1, "test")
	JsonI(2, "奇数个参数", "third")
}

func TestJsonConcurrent(t *testing.T) {
	var count = 100
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(i int) {
			JsonI("key", i)
			time.Sleep(randx.Duration(100*time.Millisecond, 500*time.Millisecond))
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestPanicDump(t *testing.T) {
	defer loom.DumpIfPanic()

	loom.Initialize(func(data []byte) {
		Warn(string(data))
	})

	panic("faint")
}

func TestFormatLog(t *testing.T) {
	type S struct {
		Name string
		Age  int
	}

	fmt.Println(formatLog("only first string"))
	fmt.Println(formatLog([]byte("first is not string")))
	fmt.Println(formatLog(S{"panda", 15}))

	fmt.Println(formatLog("formatter: %s, age=%d", "pet", 23))
	fmt.Println(formatLog("only title: ", "pet", 23))
}
