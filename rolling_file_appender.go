package logo

import (
	"fmt"
	"os"
	"path"
	"time"
)

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var levelNames = []string{"", "info", "warn", "error"}

type RollingFileAppender struct {
	formatter   *MessageFormatter
	levelFilter int

	files             [4]*os.File
	dirName           string
	fileNamePrefix    string
	maxFileSize       int64
	checkRollingCount int
}

func NewRollingFileAppender(flag int, levelFilter int, dirName string, fileNamePrefix string, maxFileSize int64) *RollingFileAppender {
	if maxFileSize <= 0 {
		panic("maxFileSize <= 0")
	}

	var my = &RollingFileAppender{
		formatter:   newMessageFormatter(flag, levelHints),
		levelFilter: levelFilter,

		dirName:        dirName,
		fileNamePrefix: fileNamePrefix,
		maxFileSize:    maxFileSize,
	}

	var err = EnsureDir(dirName, 0777)
	err = my.openLogFile(LevelInfo)
	err = my.openLogFile(LevelWarn)
	err = my.openLogFile(LevelError)

	if err != nil {
		fmt.Println(err)
	}

	return my
}

func (my *RollingFileAppender) Write(message Message) {
	var level = message.GetLevel()
	if level < my.levelFilter {
		return
	}

	switch message.level {
	case LevelError:
		my.writeMessage(message, LevelError)
		fallthrough
	case LevelWarn:
		my.writeMessage(message, LevelWarn)
		fallthrough
	case LevelInfo:
		my.writeMessage(message, LevelInfo)
	}
}

func (my *RollingFileAppender) Close() error {
	for i := 0; i < len(my.files); i++ {
		var fout = my.files[i]
		if fout != nil {
			_ = fout.Close()
			my.files[i] = nil
		}
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
	if err != nil {
		fmt.Println(err)
	}
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

	var size = info.Size()
	if size <= my.maxFileSize {
		return nil
	}

	var levelName = levelNames[level]
	var dirName = path.Join(my.dirName, levelName)
	err = EnsureDir(dirName, 0777)

	var lastPath = fout.Name()
	err = fout.Close()
	my.files[level] = nil

	var now = time.Now()
	year, month, day := now.Date()

	for i := 1; true; i++ {
		var name = fmt.Sprintf("%s%s-%d-%d-%d_%d.log", my.fileNamePrefix, levelName, year, month, day, i)
		var nextPath = path.Join(my.dirName, levelName, name)
		if IsPathExist(nextPath) {
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
	const fileFlag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	const fileMode = 0666

	var filepath = path.Join(my.dirName, my.fileNamePrefix+levelNames[level]+".log")
	var err error
	my.files[level], err = os.OpenFile(filepath, fileFlag, fileMode)

	return err
}
