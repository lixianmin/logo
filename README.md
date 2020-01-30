

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



1. 仅具备简单的输出到console与file的功能
2. 部分设计参考log4j，可以自定义appender，从而支持输出到其它设备
3. ConsoleAppender输出到控制台
4. RollingFileAppender输出到文件，支持自定义输出目录和文件名，支持按文件大小的归档
5. 附送：ding.TalkAppender输出日志到钉钉





