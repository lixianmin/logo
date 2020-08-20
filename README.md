

-----

#### logo简介

一个仅具备基本功能的golang日志框架。



目前，golang日志框架的状态如下：

1. 内置的log过于简单，只能输出到stderr，可惜了这么好的一个package name
2. github上已有的开源项目，翻来翻去，没能找到一个心仪可用的
3. 各大开源项目都选择自己攢一个，往往跟项目本身耦合，且扩展性受限



我很懒，写这么个框架都是被逼的。



话说logo这名字，log + go的组合，其实是实在找不到合适的名字，都被用过了



---

#### 功能简介



1. 不依赖任何第三方代码库
2. 仅具备简单的输出到console与file的功能
3. 部分设计参考log4j，可以自定义hook，从而支持输出到其它设备
4. ConsoleHook输出到控制台
5. RollingFileHook输出到文件，支持自定义输出目录和文件名，支持按文件大小的归档
6. **赠送**：ding.TalkHook输出日志到钉钉，baidu.InfoflowHook输出日志到百度如流



---

#### 项目实例

```go
func main() {
	// main()方法退出时关闭logger以及所有实现了Closer接口的hooks
	defer logo.GetDefaultLogger().Close()

	// 开启异步写标记，提高日志输出性能
	var logger = logo.GetDefaultLogger()
	logger.AddFlag(logo.LogAsyncWrite)

	// 开启文件日志
	const flag = logo.FlagDate | logo.FlagTime | logo.FlagShortFile | logo.FlagLevel
	var rollingFile = logo.NewRollingFileHook(logo.RollingFileHookArgs{Flag: flag})
	logger.AddHook(rollingFile)

	// 下面是业务代码
	logo.Info(1234.5678)
	logo.Warn("ding talk")
	logo.Error("hello, %s", "logo")
}
```



