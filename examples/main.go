package main

import (
	"fmt"
	"github.com/lixianmin/logo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Error struct {
	message string
}

func (my *Error) Error() string {
	return my.message
}

func main() {
	var theLogger = logo.NewLogger()
	// 开启异步写标记，提高日志输出性能
	theLogger.AddFlag(logo.LogAsyncWrite)

	// 控制台日志
	const flag = logo.FlagDate | logo.FlagTime | logo.FlagShortFile | logo.FlagLevel
	theLogger.SetFuncCallDepth(5)

	var console = logo.NewConsoleHook(logo.ConsoleHookArgs{Flag: flag})
	theLogger.AddHook(console)

	// 文件日志
	var rollingFile = logo.NewRollingFileHook(logo.RollingFileHookArgs{Flag: flag})
	theLogger.AddHook(rollingFile)

	theLogger.Info("this is info")
	year, month, day := time.Now().Date()
	var name = fmt.Sprintf("warn of : %4d-%02d-%02d", year, month, day)
	theLogger.Warn(name)

	var err = func() error {
		var err1 *Error
		return err1
	}()

	var b []byte = nil

	// err是个空接口指针, 但必须能正确处理
	logo.JsonI("err", err, "err2", &Error{message: "world"}, "b", b)

	var signalChan = make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-signalChan:
			fmt.Printf("exiting application by signal=%q\n", sig)

			// 程序退出时关闭logger以及所有实现了Closer接口的hooks
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
