package logo

import (
	"fmt"
	"github.com/lixianmin/got/loom"
	"github.com/lixianmin/got/osx"
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

const archiveDirectory = "archive"

type RollingFileHookArgs struct {
	Flag                 int
	FilterLevel          int
	DirName              string
	FileNamePrefix       string
	MaxFileSize          int64         // 当文件达到MaxFileSize后自动分隔成小文件
	ExpireTime           time.Duration // 文件最后修改时间超过ExpireTime后自动删除
	CheckRollingInterval int64         // 每间隔多少行检查rolling一个文件
}

type RollingFileHook struct {
	wc        loom.WaitClose
	args      RollingFileHookArgs
	formatter *MessageFormatter

	files [LevelMax] struct {
		*os.File
		checkRollingCount int64
	}
}

func NewRollingFileHook(args RollingFileHookArgs) *RollingFileHook {
	checkRollingFileHookArgs(&args)

	var my = &RollingFileHook{
		args:      args,
		formatter: newMessageFormatter(args.Flag, levelHints),
	}

	_ = os.MkdirAll(args.DirName, os.ModePerm)

	for level := args.FilterLevel; level < LevelMax; level++ {
		var err = my.openLogFile(level)
		checkPrintError(err)
	}

	loom.Go(my.goLoop)
	return my
}

func (my *RollingFileHook) goLoop(later loom.Later) {
	var removeTicker = later.NewTicker(6 * time.Hour)
	var closeChan = my.wc.C()

	for {
		select {
		case <-removeTicker.C:
			my.checkRemoveExpiredLogFiles()
		case <-closeChan:
			return
		}
	}
}

func (my *RollingFileHook) checkRemoveExpiredLogFiles() {
	// 遍历并删除过期的文件
	var args = my.args
	var removeTime = time.Now().Add(-args.ExpireTime)
	var dirName = path.Join(args.DirName, archiveDirectory)
	_ = filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() && info.ModTime().Before(removeTime) {
			_ = os.Remove(path)
		}
		return nil
	})
}

func (my *RollingFileHook) Write(message Message) {
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

func (my *RollingFileHook) Close() error {
	return my.wc.Close(func() error {
		for level := LevelNone + 1; level < LevelMax; level++ {
			_ = my.closeLogFile(level)
		}

		return nil
	})
}

func (my *RollingFileHook) writeMessage(message Message, level int) {
	var fout = my.files[level]
	if fout.File == nil {
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

func (my *RollingFileHook) checkRollFile(level int) (err error) {
	my.files[level].checkRollingCount++

	if my.files[level].checkRollingCount%my.args.CheckRollingInterval != 0 {
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

	var dirName = path.Join(args.DirName, archiveDirectory)
	err = os.MkdirAll(dirName, os.ModePerm)

	var lastPath = fout.Name()
	my.files[level].File = nil
	err = fout.Close()

	var levelName = levelNames[level]
	var now = time.Now()
	year, month, day := now.Date()

	for i := 1; true; i++ {
		var name = fmt.Sprintf("%s%s-%d-%d-%d_%d.log", args.FileNamePrefix, levelName, year, month, day, i)
		// 原来的归档是按levelName分类的，但是这样的话，当前使用中的debug.log, info.log, warn.log, error.log将会被分散到多处，
		// 因此，现在统一归集到archive目录下
		var nextPath = path.Join(args.DirName, archiveDirectory, name)
		if osx.IsPathExist(nextPath) {
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

func (my *RollingFileHook) openLogFile(level int) error {
	if my.files[level].File != nil {
		return nil
	}

	const fileFlag = os.O_WRONLY | os.O_CREATE | os.O_APPEND

	var args = my.args
	var fullPath = path.Join(args.DirName, args.FileNamePrefix+levelNames[level]+".log")
	var err error
	my.files[level].File, err = os.OpenFile(fullPath, fileFlag, 0666) // 在docker中创建的文件必须让外面的人可以读

	return err
}

func (my *RollingFileHook) closeLogFile(level int) error {
	var files = my.files
	if level > LevelNone && level < LevelMax && files[level].File != nil {
		var file = files[level]
		files[level].File = nil
		var err = file.Close()
		return err
	}

	return nil
}

func (my *RollingFileHook) SetFilterLevel(level int) {
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

func checkRollingFileHookArgs(args *RollingFileHookArgs) {
	if args.FilterLevel <= LevelNone || args.FilterLevel >= LevelMax {
		args.FilterLevel = LevelInfo
	}

	if args.DirName == "" {
		args.DirName = "logs"
	}

	if args.MaxFileSize <= 0 {
		args.MaxFileSize = 10 * 1024 * 1024 // 默认大小为10M
	}

	if args.ExpireTime <= 0 {
		args.ExpireTime = 7 * 24 * time.Hour // 默认7天后删除
	}

	if args.CheckRollingInterval <= 0 {
		args.CheckRollingInterval = 1024
	}
}
