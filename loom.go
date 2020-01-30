package logo

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

/********************************************************************
created:    2020-01-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func DumpIfPanic() {
	var panicData = recover()
	if panicData == nil {
		return
	}

	var exeName = filepath.Base(os.Args[0]) // 获取程序名称
	var now = time.Now()                    // 获取当前时间
	var pid = os.Getpid()                   // 获取进程ID

	// 设定时间格式
	var timestamp = now.Format(time.RFC3339)
	// 保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）
	var logDir = "logs"
	var logFilePath = fmt.Sprintf("%s/dump.%s.%d.%s.log", logDir, exeName, pid, timestamp)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.MkdirAll(logDir, os.ModePerm)
	}

	fmt.Println("dump to file ", logFilePath)

	f, err := os.Create(logFilePath)
	if err != nil {
		return
	}
	defer f.Close()

	// 输出panic信息
	//writeOneMessage(f, "------------------------------------\r\n")
	//writeOneMessage(f, message)
	writeOneMessage(f, "------------------------------------\r\n")
	writeOneMessage(f, fmt.Sprintf("%v\r\n", panicData))
	writeOneMessage(f, "------------------------------------\r\n")

	// 输出堆栈信息
	writeOneMessage(f, string(debug.Stack()))

	// 直接退出？
	os.Exit(1)
}

func writeOneMessage(fout *os.File, message string) {
	_, _ = fout.WriteString(message)
	_, _ = os.Stderr.WriteString(message)
}
