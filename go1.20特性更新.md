Go 1.20 已发布！性能和构建速度上有很大提升！

哈喽，大家好。

就在今天凌晨 Go 团队发布了 Go 1.20，这是有史以来 2 月份发布最早的一次。

Image
图片来源：https://golang.design/history/#timeline
Go 1.20 amd64 版本 95MB，相比 Go 1.19 144MB，以及 Go 1.18.10 138MB 来说，降幅达到了 49MB。

简直爱了！赶紧升级 Go 1.20！

Image

下图是 Go 1.20 的部分变更内容：

Image
来源：https://github.com/golang/go/milestone/250?closed=1




本次更新涉及很多内容，特别是



Runtime: 对垃圾收集器的一些内部数据结构进行了重组，以提高空间和 CPU 效率。此更改减少了内存开销，并将整体 CPU 性能提高了2% 。

Compiler：Go 1.20增加了对概要文件引导的优化(PGO)的预览支持。

Build：Go 1.20 improves build speeds by up to 10%。



我们之前已经分析过的一些特性：

带你提前看看 Go 1.20 包括哪些重大变更和性能提升

Go 1.20 新特性之 time.Compare



https://go.dev/doc/go1.20
