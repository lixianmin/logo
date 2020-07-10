package logo

import (
	"fmt"
	"github.com/lixianmin/got/loom"
	"github.com/lixianmin/logo/tools"
	"os"
	"path"
	"path/filepath"
	"time"
)

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var levelNames = []string{"", "debug", "info", "warn", "error"}

type RollingFileAppenderArgs struct {
	Flag           int
	FilterLevel    int
	DirName        string
	FileNamePrefix string
	MaxFileSize    int64         // 当文件达到MaxFileSize后自动分隔成小文件
	ExpireTime     time.Duration // 文件最后修改时间超过ExpireTime后自动删除
}

type RollingFileAppender struct {
	wc        *loom.WaitClose
	args      RollingFileAppenderArgs
	formatter *MessageFormatter

	files             [LevelMax]*os.File
	checkRollingCount int
}

func NewRollingFileAppender(args RollingFileAppenderArgs) *RollingFileAppender {
	checkRollingFileAppenderArgs(&args)

	var my = &RollingFileAppender{
		wc:        loom.NewWaitClose(),
		args:      args,
		formatter: newMessageFormatter(args.Flag, levelHints),
	}

	_ = os.MkdirAll(args.DirName, 0666)

	for level := args.FilterLevel; level < LevelMax; level++ {
		var err = my.openLogFile(level)
		checkPrintError(err)
	}

	go my.goLoop()
	return my
}

func (my *RollingFileAppender) goLoop() {
	defer loom.DumpIfPanic()

	var args = my.args
	var ticker = time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 遍历并删除过期的文件
			var removeTime = time.Now().Add(-args.ExpireTime)
			for level := LevelNone + 1; level < LevelMax; level++ {
				var levelName = levelNames[level]
				var dirName = path.Join(args.DirName, levelName)
				_ = filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
					if info != nil && !info.IsDir() && info.ModTime().Before(removeTime) {
						_ = os.Remove(path)
					}
					return nil
				})
			}
		case <-my.wc.CloseChan:
			return
		}
	}
}

func (my *RollingFileAppender) Write(message Message) {
	var level = message.GetLevel()
	var filterLevel = my.args.FilterLevel
	if level < filterLevel {
		return
	}

	switch level {
	case LevelError:
		my.writeMessage(message, LevelError)
		fallthrough
	case LevelWarn:
		my.writeMessage(message, LevelWarn)
		fallthrough
	case LevelInfo:
		my.writeMessage(message, LevelInfo)
		fallthrough
	case LevelDebug:
		my.writeMessage(message, LevelDebug)
	}
}

func (my *RollingFileAppender) Close() error {
	_ = my.wc.Close()
	for level := LevelNone + 1; level < LevelMax; level++ {
		_ = my.closeLogFile(level)
	}

	return nil
}

func (my *RollingFileAppender) writeMessage(message Message, level int) {
	var fout = my.files[level]
	if fout == nil {
		return
	}

	var buffer = my.formatter.format(message)
	_, err := fout.Write(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = my.checkRollFile(level)
	checkPrintError(err)
}

func (my *RollingFileAppender) checkRollFile(level int) (err error) {
	my.checkRollingCount++

	const checkInterval = 1024
	if my.checkRollingCount%checkInterval != 0 {
		return nil
	}

	// 检测文件大小是否超过maxFileSize
	var fout = my.files[level]
	var info os.FileInfo
	info, err = fout.Stat()

	if err != nil {
		return err
	}

	var args = my.args
	var size = info.Size()
	if size <= args.MaxFileSize {
		return nil
	}

	var levelName = levelNames[level]
	var dirName = path.Join(args.DirName, levelName)
	err = os.MkdirAll(dirName, 0666)

	var lastPath = fout.Name()
	my.files[level] = nil
	err = fout.Close()

	var now = time.Now()
	year, month, day := now.Date()

	for i := 1; true; i++ {
		var name = fmt.Sprintf("%s%s-%d-%d-%d_%d.log", args.FileNamePrefix, levelName, year, month, day, i)
		var nextPath = path.Join(args.DirName, levelName, name)
		if tools.IsPathExist(nextPath) {
			continue
		}

		err = os.Rename(lastPath, nextPath)
		if err != nil {
			return err
		}

		err = my.openLogFile(level)
		return err
	}

	return nil
}

func (my *RollingFileAppender) openLogFile(level int) error {
	if my.files[level] != nil {
		return nil
	}

	const fileFlag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	const fileMode = 0666

	var args = my.args
	var filepath = path.Join(args.DirName, args.FileNamePrefix+levelNames[level]+".log")
	var err error
	my.files[level], err = os.OpenFile(filepath, fileFlag, fileMode)

	return err
}

func (my *RollingFileAppender) closeLogFile(level int) error {
	var files = my.files
	if level > LevelNone && level < LevelMax && files[level] != nil {
		var file = files[level]
		files[level] = nil
		var err = file.Close()
		return err
	}

	return nil
}

func (my *RollingFileAppender) SetFilterLevel(level int) {
	if level > LevelNone && level < LevelMax {
		my.args.FilterLevel = level

		for i := LevelNone + 1; i < level; i++ {
			var err = my.closeLogFile(i)
			checkPrintError(err)
		}

		for i := level; i < LevelMax; i++ {
			var err = my.openLogFile(i)
			checkPrintError(err)
		}
	}
}

func checkPrintError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func checkRollingFileAppenderArgs(args *RollingFileAppenderArgs) {
	if args.FilterLevel <= LevelNone || args.FilterLevel >= LevelMax {
		args.FilterLevel = LevelInfo
	}

	if args.DirName == "" {
		args.DirName = "logs"
	}

	if args.FileNamePrefix == "" {
		args.FileNamePrefix = "log_"
	}

	if args.MaxFileSize <= 0 {
		args.MaxFileSize = 10 * 1024 * 1024 // 默认大小为10M
	}

	if args.ExpireTime <= 0 {
		args.ExpireTime = 7 * 24 * time.Hour // 默认7天后删除
	}
}
