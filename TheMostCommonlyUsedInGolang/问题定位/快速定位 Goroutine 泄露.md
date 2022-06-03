# 快速定位 Goroutine 泄露

[Golang技术分享](javascript:void(0);) *2022-05-07 08:01*

Editor's Note

煎鱼大佬，Go 领域专家，出版过畅销书《Go 语言编程之旅》，获得 GOP（Go 领域最有观点专家）的荣誉称号。关注一波！

The following article is from 脑子进煎鱼了 Author 陈煎鱼

Go 语言能够广受大家喜欢，其中一个原因就是他的协程做的非常非常简单，初学的入门者都可以使用。

![Image](https://mmbiz.qpic.cn/mmbiz_png/KVl0giak5ib4gTxOG4qAu7VefHGf1hzDSa5Bl5Uib6PRUBzk0OwpzB8E01Wp467zicabiay4pPA0PpMwOgtZExfmA3w/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

平时只需 go 关键字一下，成千上万个 goroutines 就出现了：

```
for ...
  go func(){}
```

起协程就跟下饺子的。这时就个大问题，因为协程用起来简单，出问题出起来也很快。

也就是常常会出现 goroutine 泄露，查起来很费劲。

## goleak

今天给大家推荐一个好物，定位为纯介绍。他来自 Uber 的 **Goroutine leak detector**[1]，他能够结合单元测试去快速的检测 goroutine 泄露，达到避免和排查的目的。

![Image](data:image/gif;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQImWNgYGBgAAAABQABh6FO1AAAAABJRU5ErkJggg==)

在命令行执行第三方库安装：

```
go get -u go.uber.org/goleak
```

channel 会一直阻塞，导致泄露的方法：

```go
func leak() {
 ch := make(chan struct{})
 go func() {
  ch <- struct{}{}
 }()
}
```

如果我们直接在常规的测试方法中调用：

```go
func TestLeak(t *testing.T) {
 leak()
}
```

输出结果：

```
=== RUN   TestLeak
--- PASS: TestLeak (0.00s)
PASS
```

是不会有任何变化的，正常通过测试。

这时候我们需要在代码中添加 `goleak.VerifyNone` 方法，如下：

```
import (
 "testing"
 "go.uber.org/goleak"
)

func TestLeak(t *testing.T) {
 defer goleak.VerifyNone(t)
 leak()
}
```

再进行验证：

```shell
=== RUN   TestLeak
    leaks.go:78: found unexpected goroutines:
        [Goroutine 7 in state chan send, with github.com/eddycjy/awesome-project/tools.leak.func1 on top of the stack:
        goroutine 7 [chan send]:
        github.com/eddycjy/awesome-project/tools.leak.func1(0xc0000562a0)
         /Users/eddycjy/go-application/awesomeProject/tools/leak.go:6 +0x35
        created by github.com/eddycjy/awesome-project/tools.leak
         /Users/eddycjy/go-application/awesomeProject/tools/leak.go:5 +0x4e
        ]
--- FAIL: TestLeak (0.46s)
FAIL
```

可以从报告中看到，运行结构会明确的告诉你发现泄露的 goroutine 的代码堆栈和泄露类型，非常的省心。

另外在内部的 CI/CD 流程里，把 goleak 结合上，对于平时的 CR 也是一个不错的辅助。

## 总结

今天我们的好物分享介绍了一个小工具，能够解决大家团队中时长遇到的 goroutine 泄露问题，除了本文提到的 uber-go/goleak，也还有 ysmood/gotrace 等同类型的库能够达到类似的效果。

希望对你排查 goroutine 泄露有所帮助：）

### 参考资料

[1]Goroutine leak detector: *github.com/uber-go/goleak*