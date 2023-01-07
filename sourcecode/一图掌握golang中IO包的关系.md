# 一图掌握golang中IO包的关系

Original 渔夫子 [Go学堂](javascript:void(0);) *2023-05-20 11:34* *Posted on 北京*

收录于合集

\#go13个

\#个人成长5个

\#io1个

大家好，我是渔夫子。

今天在知乎上看到这样一个问题：Golang的IO库那么多，我该怎么选。今天就跟大家聊聊这个问题。

首先，我们要知道，golang中有哪些IO包。我整理了一下，大概有io包、bufio包、ioutil、os、net等。

其次，要知道这些io包的各自的定位。我整理了一张图供大家参考：![Image](https://mmbiz.qpic.cn/mmbiz_png/l6hSQtEH2593FWibsiaEuZOT0d3MAZQ9nT198popuuxwTKGwRJCZKw0aZgRVLVhlukmibDBtBeNibsHJlF8f3wDSMQ/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

我们大致讲解下上图：

- **io**：基础的IO库，提供了Reader和Writer接口。其中的os包、net包、string包、bytes包以及bufio包都实现了io中的Reader或Writer接口。
- **os**：提供了访问底层操作系统资源的能力，如文件读写、进程控制等。
- **net**：提供了网络相关的IO功能，如TCP、UDP通信、HTTP请求等。
- **string.Reader**：提供了string的读取。因为string不能写，所以只有Reader。
- **bytes.Buffer和Reader**：提供了对字节内容的读写。
- **bufio**：提供带缓存的I/O操作，解决频繁、少量读取场景下的性能问题。这里利用了计算机的局部性原理。
- **ioutil**：提供了一些方便的文件读写函数，如ReadFile和WriteFile。

我们以iotuil包为例，看下ReadDir函数的实现。ReadDir函数的功能就是从一个目录中读取所有的文件列表。这个操作其实包含两步：打开文件、读取目录下的文件。ReadDir函数就把这两步做了封装，供客户端调用，是不是就更方便了。如下代码：

```
func ReadDir(dirname string) ([]fs.FileInfo, error) {
 f, err := os.Open(dirname)
 if err != nil {
  return nil, err
 }
 list, err := f.Readdir(-1)
 f.Close()
 if err != nil {
  return nil, err
 }
 sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
 return list, nil
}
```

所以，选择哪个库主要取决于你要进行什么样的操作。如果只是简单的文件读写，可以使用ioutil库；如果需要处理大量数据，则应该选择bufio库；如果需要访问底层操作系统功能，则可以使用os库；如果涉及到网络通信，则选择net库。

好了，今天就分享到这里。

你的关注，是我写下去的最大动力。点击下方公众号卡片，直接关注。关注送《100个go常见的错误》pdf文档、经典go学习资料。