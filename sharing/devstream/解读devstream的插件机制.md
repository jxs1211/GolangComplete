本文已参与「开源摘星计划」，欢迎正在阅读的你加入。活动链接：https://github.com/weopenprojects/WeOpen-Star

devstream的插件机制基于go的插件系统实现

简单介绍下go的插件系统

https://img.draveness.me/2020-04-03-15859025269151-plugin-system.png

通过在主程序和共享库直接定义一系列的约定或者接口，再通过以下的代码动态加载其他人编译的 Go 语言共享对象。
主要步骤：

- 1 通过Open调用.so文件
- 2 再通过Lookup加载Symbol实现共享对象的调用

```go
type Driver interface {
    Name() string
}

func main() {
    p, err := plugin.Open("driver.so")
    if err != nil {
	   panic(err)
    }

    newDriverSymbol, err := p.Lookup("NewDriver")
    if err != nil {
        panic(err)
    }

    newDriverFunc := newDriverSymbol.(func() Driver)
    newDriver := newDriverFunc()
    fmt.Println(newDriver.Name())
}
```

devstream使用go的插件机制，实现plugin的模板的自动生成
