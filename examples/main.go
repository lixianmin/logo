package main

import (
	"fmt"
	"github.com/lixianmin/logo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var theLogger = logo.NewLogger()
	// 开启异步写标记，提高日志输出性能
	theLogger.AddFlag(logo.LogAsyncWrite)

	// 控制台日志
	const flag = logo.FlagDate | logo.FlagTime | logo.FlagShortFile | logo.FlagLevel
	theLogger.SetFuncCallDepth(5)

	var console = logo.NewConsoleAppender(logo.ConsoleAppenderArgs{Flag: flag})
	theLogger.AddAppender(console)

	// 文件日志
	var rollingFile = logo.NewRollingFileAppender(logo.RollingFileAppenderArgs{Flag: flag})
	theLogger.AddAppender(rollingFile)

	theLogger.Info("this is info")
	theLogger.Warn("that is warn")

	var signalChan = make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-signalChan:
			fmt.Printf("exiting application by signal=%q\n", sig)

			// 程序退出时关闭logger以及所有实现了Closer接口的appenders
			println("hello")
			theLogger.Info("[main()] logger is exiting...")
			theLogger.Info("[main()] logger is exiting1...")
			theLogger.Info("[main()] logger is exiting2...")
			_ = theLogger.Close()
			println("world")

			os.Exit(0)
		}
	}()

	time.Sleep(time.Hour)
}
