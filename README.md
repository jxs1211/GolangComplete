# GolangComplete

## go-questions

## GolangProjectDevelopmentPractice

## GoProgrammingFromBeginnerToMaster

第一部分 熟知Go语言的一切

第1条 了解Go语言的诞生与演进2
1.1 Go语言的诞生2
1.2 Go语言的早期团队和演进历程4
1.3 Go语言正式发布并开源4

第2条 选择适当的Go语言版本6
2.1 Go语言的先祖6
2.2 Go语言的版本发布历史7
2.3 Go语言的版本选择建议11

第3条 理解Go语言的设计哲学12
3.1 追求简单，少即是多12
3.2 偏好组合，正交解耦15
3.3 原生并发，轻量高效17
3.4 面向工程，“自带电池”21

第4条 使用Go语言原生编程思维来写Go代码26
4.1 语言与思维—来自大师的观点26

4.2 现实中的“投影”27
- https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter1/sources/sieve.go#L61

4.3 Go语言原生编程思维29


第二部分 项目结构、代码风格与标识符命名

第5条 使用得到公认且广泛使用的项目结构32

5.1 Go项目的项目结构32
- https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter2/sources/GoProj/go.mod

[going@dev GoProj]$ tree
.
├── cmd
│   ├── app1
│   │   └── main.go
│   └── app2
│       └── main.go
├── go.mod
├── LICENSE
├── Makefile
├── pkg
│   ├── lib1
│   │   └── lib1.go
│   └── lib2
│       └── lib2.go
└── README.md

- https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter2/sources/GoLibProj/go.mod

[going@dev GoLibProj]$ tree
.
├── go.mod
├── internal
│   ├── ilib1
│   │   └── ilib1.go
│   └── ilib2
│       └── ilib2.go
├── lib1
│   └── lib1.go
├── lib2
│   └── lib2.go
├── lib.go
├── LICENSE
├── Makefile
└── README.md

5.2 Go语言典型项目结构35

第6条 提交前使用gofmt格式化源码40
6.1 gofmt：Go语言在解决规模化问题上的最佳实践40
6.2 使用gofmt41
6.3 使用goimports43
6.4 将gofmt/goimports与IDE或编辑器工具集成44

第7条 使用Go命名惯例对标识符进行命名47
7.1 简单且一致48
7.2 利用上下文环境，让最短的名字携带足够多的信息53

第三部分 声明、类型、语句与控制结构

第8条 使用一致的变量声明形式56
8.1 包级变量的声明形式56
8.2 局部变量的声明形式59

第9条 使用无类型常量简化代码63
9.1 Go常量溯源63
9.2 有类型常量带来的烦恼64
9.3 无类型常量消除烦恼，简化代码65

第10条 使用iota实现枚举常量68

第11条 尽量定义零值可用的类型73
11.1 Go类型的零值73
11.2 零值可用75

第12条 使用复合字面值作为初值构造器78
12.1 结构体复合字面值79
12.2 数组/切片复合字面值80
12.3 map复合字面值81

第13条 了解切片实现原理并高效使用83
13.1 切片究竟是什么83

13.2 切片的高级特性：动态扩容87
https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/slice_append.go

13.3 尽量使用cap参数创建切片90
https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/slice_benchmark_test.go

第14条 了解map实现原理并高效使用92
14.1 什么是map92
https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/map_var_as_param.go

14.2 map的基本操作93
https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/map_iterate.go

https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/map_multiple_iterate.go

https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/map_stable_iterate.go

14.3 map的内部实现97
https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/map_concurrent_read_and_write.go

14.4 尽量使用cap参数创建map103
https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/map_test.go

第15条 了解string实现原理并高效使用105
15.1 Go语言的字符串类型105
https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/string_type.go

15.2 字符串的内部表示110
https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/string_as_param_benchmark_test.go

15.3 字符串的高效构造112
https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/string_concat_benchmark_test.go

15.4 字符串相关的高效转换115
https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/string_slice_to_string.go

https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/string_mallocs_in_convert.go

https://github.com/jxs1211/GolangComplete/blob/a7f196d7c5fa7e5d96294bba884acdce0c86698a/GoProgrammingFromBeginnerToMaster/chapter3/sources/string_for_range_covert_optimize.go


第16条 理解Go语言的包导入120
16.1 Go程序构建过程121
16.2 究竟是路径名还是包名127
16.3 包名冲突问题130

第17条 理解Go语言表达式的求值顺序132
17.1 包级别变量声明语句中的表达式求值顺序133
17.2 普通求值顺序136
17.3 赋值语句的求值139
17.4 switch/select语句中的表达式求值140

第18条 理解Go语言代码块与作用域143
18.1 Go代码块与作用域简介143
18.2 if条件控制语句的代码块145
18.3 其他控制语句的代码块规则简介148

第19条 了解Go语言控制语句惯用法及使用注意事项154
19.1 使用if控制语句时应遵循“快乐路径”原则154
19.2 for range的避“坑”指南156
19.3 break跳到哪里去了165
19.4 尽量用case表达式列表替代fallthrough167

第四部分 函数与方法

第20条 在init函数中检查包级变量的初始状态170
20.1 认识init函数170
20.2 程序初始化顺序171
20.3 使用init函数检查包级变量的初始状态174

第21条 让自己习惯于函数是“一等公民”179
21.1 什么是“一等公民”179
21.2 函数作为“一等公民”的特殊运用183

第22条 使用defer让函数更简洁、更健壮192
22.1 defer的运作机制193
22.2 defer的常见用法194
22.3 关于defer的几个关键问题199

第23条 理解方法的本质以选择
正确的receiver类型206
23.1 方法的本质207
23.2 选择正确的receiver类型208
23.3 基于对Go方法本质的理解巧解难题210

第24条 方法集合决定接口实现214
24.1 方法集合215
24.2 类型嵌入与方法集合216
24.3 defined类型的方法集合226
24.4 类型别名的方法集合227

第25条 了解变长参数函数的妙用230
25.1 什么是变长参数函数230
25.2 模拟函数重载233
25.3 模拟实现函数的可选参数与默认参数236
25.4 实现功能选项模式238

第五部分 接口

第26条 了解接口类型变量的内部表示246
26.1 nil error值 != nil247
26.2 接口类型变量的内部表示248
26.3 输出接口类型变量内部表示的详细信息254
26.4 接口类型的装箱原理258

第27条 尽量定义小接口263
27.1 Go推荐定义小接口263
27.2 小接口的优势265
27.3 定义小接口可以遵循的一些点267

第28条 尽量避免使用空接口作为函数参数类型270

第29条 使用接口作为程序水平组合的连接点274
29.1 一切皆组合274
29.2 垂直组合回顾275
29.3 以接口为连接点的水平组合276

第30条 使用接口提高代码的可测试性281
30.1 实现一个附加免责声明的电子邮件发送函数282
30.2 使用接口来降低耦合283

第六部分 并发编程

第31条 优先考虑并发设计288
31.1 并发与并行288
31.2 Go并发设计实例290

第32条 了解goroutine的调度原理299
32.1 goroutine调度器299
32.2 goroutine调度模型与演进过程300
32.3 对goroutine调度器原理的进一步理解302
32.4 调度器状态的查看方法305
32.5 goroutine调度实例简要分析307

第33条 掌握Go并发模型和常见并发模式315
33.1 Go并发模型315
33.2 Go常见的并发模式317

第34条 了解channel的妙用340
34.1 无缓冲channel341
34.2 带缓冲channel347
34.3 nil channel的妙用354
34.4 与select结合使用的一些惯用法357

第35条 了解sync包的正确用法359
35.1 sync包还是channel359
35.2 使用sync包的注意事项360
35.3 互斥锁还是读写锁362
35.4 条件变量365
35.5 使用sync.Once实现单例模式 368
35.6 使用sync.Pool减轻垃圾回收压力370

第36条 使用atomic包实现伸缩性更好的并发读取374
36.1 atomic包与原子操作374
36.2 对共享整型变量的无锁读写375
36.3 对共享自定义类型变量的无锁读写377

第七部分 错误处理
第37条 了解错误处理的4种策略382
37.1 构造错误值383
37.2 透明错误处理策略385
37.3 “哨兵”错误处理策略385
37.4 错误值类型检视策略388
37.5 错误行为特征检视策略390

第38条 尽量优化反复出现的if err != nil392
38.1 两种观点393
38.2 尽量优化395
38.3 优化思路395

第39条 不要使用panic进行正常的错误处理405
39.1 Go的panic不是Java的checked exception405
39.2 panic的典型应用408
39.3 理解panic的输出信息412
（以上为本书内容，以下为第2册内容。）

第八部分 测试、性能剖析与调试

第40条 理解包内测试与包外测试的差别
40.1 官方文档的“自相矛盾”
40.2 包内测试与包外测试
第41条 有层次地组织测试代码

## The-Golang-Standard-Library-by-Example

## TheMostCommonlyUsedInGolang

## toy-web-1

## go语言经典库，保持关注
