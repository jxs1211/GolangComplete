# Go 语言自带设计模式

[GoCN](javascript:void(0);) *2023-03-21 08:15* *Posted on 湖北*

The following article is from 洋芋编程 Author 蛮荆

[![img](http://wx.qlogo.cn/mmhead/Q3auHgzwzM4uG0Oe9TfmsLHuKzrBwHubd2KvMoZSOH6W1oRRTZBibyw/0)**洋芋编程**.Go 语言、云原生](https://mp.weixin.qq.com/s/oKNlRbbCaCWgtz5iE7lpkA#)

## 👇我在这儿 

![img](http://mmbiz.qpic.cn/mmbiz_png/vbERicIdYZbC53rq7PQziczkzCA4pIPx8Xdx2r0a6FgQWEicTOgFNY1KlJD3jSKQKlTiavwTUgHS1HMXza2RYApkDw/0?wx_fmt=png)

**GoCN**

最具规模和生命力的 Go 开发者社区

607篇原创内容



公众号



## 概述

> 在软件工程中，设计模式（design pattern）是对软件设计中普遍存在（反复出现）的各种问题，所提出的解决方案。 -- 维基百科

和传统的 `GOF`, `Java`, `C#` 教科书式的 `设计模式` 不同，Go 语言设计从一开始就力求简洁，有其他编程语言基础的读者在学习和使用 Go 语言时， 万万不可按图索骥、生搬硬套，简单的事情复杂化。

本文带领大家一起看一下，Go 语言标准库中自带的 `编程设计模式`。

## 单例模式

**确保一个类只有一个实例，并提供对该实例的全局访问**。

通过使用标准库中的 `sync.Once` 对业务对象进行简单封装，即可实现 `单例模式`，简单安全高效。

```
package main

import "sync"

var (
    once     sync.Once
    instance Singleton
)

// Singleton 业务对象
type Singleton struct {
}

// NewInstance 单例模式方法
func NewInstance() Singleton {
    once.Do(func() {
        instance = Singleton{}
    })
    return instance
}

func main() {
    // 调用方代码
    s1 := NewInstance()
    s2 := NewInstance()
    s3 := NewInstance() 
}
```

![Image](https://mmbiz.qpic.cn/mmbiz_png/2mstuQELkbtwzRwt8icDiaaokwxp7XIQbJeyzXFCiaFybG0gTS9icBUQXSGDySHnKsxIRpsR25riab3GgYBx9P5IJTA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

Go 标准库单例模式

## 简单工厂模式

Go 语言本身没有 `构造方法` 特性，工程实践中一般使用 `NewXXX` 创建新的对象 (XXX 为对象名称)，比如标准库中的:

```
// errors/errors.go

func New(text string) error {
    return &errorString{text}
}

// sync/cond.go
func NewCond(l Locker) *Cond {
    return &Cond{L: l}
}
```

在这个基础上，如果方法返回的是 `interface` 的时候，其实就等于是 `简单工厂模式`，然后再加一层抽象的话，就接近于 `抽象工厂模式`。

```
package main

// ConfigParser 配置解析接口
type ConfigParser interface {
    Parse(p []byte)
}

// JsonParser Json 文件解析器
type JsonParser struct {
}

func (j *JsonParser) Parse(p []byte) {

}

func newJsonParser() *JsonParser {
    return &JsonParser{}
}

// YamlParser Yaml 文件解析器
type YamlParser struct {
}

func (y *YamlParser) Parse(p []byte) {

}

func newYamlParser() *YamlParser {
    return &YamlParser{}
}

type ConfigType uint8

const (
    JsonType ConfigType = 1 << iota
    YamlType
)

// NewConfig 根据不同的类型创建对应的解析器
func NewConfig(t ConfigType) ConfigParser {
    switch t {
    case JsonType:
        return newJsonParser()
    case YamlType:
        return newYamlParser()
    default:
        return nil
    }
}

func main() {
    // 调用方代码
    jsonParser := NewConfig(JsonType)
    yamlParser := NewConfig(YamlType)
}
```

![Image](https://mmbiz.qpic.cn/mmbiz_png/2mstuQELkbtwzRwt8icDiaaokwxp7XIQbJhLrEC8My3eoBTicUIa1fZd7tGOnH2rQwdpibZuGNPNAFaS9Bt93n8IlA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

Go 实现简单工厂模式

## 对象池模式

通过回收利用对象避免获取和释放资源所需的昂贵成本，我们可以直接使用 `sync.Pool` 对象来实现功能。

```
package main

import (
    "net/http"
    "sync"
)

var (
    // HTTP Request 对象池
    reqPool = sync.Pool{
        New: func() any {
            return http.Request{}
        },
    }
)

func main() {
    // 调用方代码
    r1 := reqPool.Get()
    r2 := reqPool.Get()
    r3 := reqPool.Get()

    reqPool.Put(r1)
    reqPool.Put(r2)
    reqPool.Put(r3)
}
```

## 构建模式 (Builder)

将一个复杂对象的构建与它的表示分离，使得同样的构建过程可以创建不同的表示。

如果用传统的方法实现 `构建模式`，对应的 Go 语言代码大致是下面这个样子:

```
package main

type QueryBuilder interface {
    Select(table string, columns []string) QueryBuilder
    Where(conditions ...string) QueryBuilder
    GetRawSQL() string
}

type MySQLQueryBuilder struct {
}

func (m *MySQLQueryBuilder) Select(table string, columns ...string) QueryBuilder {
    // 具体实现代码跳过
    return nil
}

func (m *MySQLQueryBuilder) Where(conditions ...string) QueryBuilder {
    // 具体实现代码跳过
    return nil
}

func (m *MySQLQueryBuilder) GetRawSQL() string {
    // 具体实现代码跳过
    return ""
}

func main() {
    // 调用方代码
    m := &MySQLQueryBuilder{}

    sql := m.Select("users", "username", "password").
        Where("id = 100").
        GetRawSQL()

    println(sql)
}
```

![Image](https://mmbiz.qpic.cn/mmbiz_png/2mstuQELkbtwzRwt8icDiaaokwxp7XIQbJeib59iaoMjR8ZRv5SQgzYxTBW8XsB24PhibaRhbXeiaWpJGWS0RSribB70A/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

Go 实现构建模式

上面的代码中，通过经典的链式调用来构造出具体的 SQL 语句，但是在 Go 语言中，我们一般使用另外一种模式来实现同样的功能 `FUNCTIONAL OPTIONS`， 这似乎也是 Go 语言中最流行的模式之一。

```
package main

type SQL struct {
    Table   string
    Columns []string
    Where   []string
}

type Option func(s *SQL)

func Table(t string) Option {
    // 注意返回值类型
    return func(s *SQL) {
        s.Table = t
    }
}

func Columns(cs ...string) Option {
    // 注意返回值类型
    return func(s *SQL) {
        s.Columns = cs
    }
}

func Where(conditions ...string) Option {
    // 注意返回值类型
    return func(s *SQL) {
        s.Where = conditions
    }
}

func NewSQL(options ...Option) *SQL {
    sql := &SQL{}

    for _, option := range options {
        option(sql)
    }

    return sql
}

func main() {
    // 调用方代码
    sql := NewSQL(Table("users"),
        Columns("username", "password"),
        Where("id = 100"),
    )

    println(sql)
}
```

![Image](https://mmbiz.qpic.cn/mmbiz_png/2mstuQELkbtwzRwt8icDiaaokwxp7XIQbJP0ESJ3Ohenl2DTGyCTNv9buptl1XkBG9IBjgEMvGoWJCIlickZWOiaLw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

Go FUNCTIONAL OPTIONS 模式

## 观察者模式

在对象间定义一个一对多的联系性，由此当一个对象改变了状态，所有其他相关的对象会被通知并且自动刷新。

如果用传统的方法实现 `观察者模式`，对应的 Go 语言代码大致是下面这个样子:

```
package main

import "math"

// Observer 观察者接口
type Observer interface {
    OnNotify(Event)
}

// Notifier 订阅接口
type Notifier interface {
    Register(Observer)
    Deregister(Observer)
    Notify(Event)
}

type (
    Event struct {
        Data int64
    }

    eventObserver struct {
        id int
    }

    eventNotifier struct {
        observers map[Observer]struct{}
    }
)

// OnNotify 观察者收到订阅的时间回调
func (o *eventObserver) OnNotify(e Event) {
}

// Register 注册观察者
func (o *eventNotifier) Register(l Observer) {
    o.observers[l] = struct{}{}
}

// Deregister 移除观察者
func (o *eventNotifier) Deregister(l Observer) {
    delete(o.observers, l)
}

// Notify 发出通知
func (o *eventNotifier) Notify(e Event) {
    for p := range o.observers {
        p.OnNotify(e)
    }
}

func main() {
    // 调用方代码
    notifier := eventNotifier{
        observers: make(map[Observer]struct{}),
    }

    notifier.Register(&eventObserver{1})
    notifier.Register(&eventObserver{2})
    notifier.Register(&eventObserver{3})

    notifier.Notify(Event{Data: math.MaxInt64})
}
```

![Image](https://mmbiz.qpic.cn/mmbiz_png/2mstuQELkbtwzRwt8icDiaaokwxp7XIQbJQwvN36zWle19pDd4WA3hwP0cS21hdZVcsEsZ1yW5orL1icg5QnibanZw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

Go 实现观察者模式

但其实我们有更简洁的方法，直接使用标准库中的 `sync.Cond` 对象，改造之后的 `观察者模式` 代码大概是这个样子:

```
package main

import (
    "fmt"
    "sync"
    "time"
)

var done = false

func read(name string, c *sync.Cond) {
    fmt.Println(name, "starts reading")

    c.L.Lock()
    for !done {
        c.Wait() // 等待发出通知
    }
    c.L.Unlock()
}

func write(name string, c *sync.Cond) {
    fmt.Println(name, "starts writing")
    time.Sleep(100 * time.Millisecond)

    c.L.Lock()
    done = true // 设置条件变量
    c.L.Unlock()

    fmt.Println(name, "wakes all")
    c.Broadcast() // 通知所有观察者
}

func main() {
    cond := sync.NewCond(&sync.Mutex{}) // 创建时传入一个互斥锁

    // 3 个观察者
    go read("reader1", cond)
    go read("reader2", cond)
    go read("reader3", cond)

    time.Sleep(time.Second) // 模拟延时

    write("writer-1", cond) // 发出通知

    time.Sleep(time.Second) // 模拟延时
}
```

![Image](https://mmbiz.qpic.cn/mmbiz_png/2mstuQELkbtwzRwt8icDiaaokwxp7XIQbJpMWWwN3FDpzAhs4sBeQbY7cPTz8OppPrmiat9d42n06GQPhpUrot7vg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

Go 标准库观察者模式

将代码改造为 `sync.Cond` 之后，代码量更好，结构更简洁。

## ok/error 模式

在 Go 语言中，经常在一个表达式返回 `2` 个参数时使用这种模式:

- • 第 1 个参数是一个值或者 `nil`
- • 第 2 个参数是 `true/false` 或者 `error`

在一个需要赋值的 `if` 条件语句中，使用这种模式去检测第 2 个参数值会让代码显得优雅简洁。

### 在函数返回时检测错误

```
package main

func foo() (int, error){
    return 0, nil
}

func main() {
    if v, err := foo(); err != nil {
        panic(err)
    } else {
        println(v)
    }
}
```

### 检测 map 是否存在指定的 key

```
package main

func main() {
    m := make(map[int]string)

    if v, ok := m[0]; ok {
        println(v)
    }
}
```

### 类型断言

```
package main

func foo() interface{} {
    return 1024
}

func main() {
    n := foo()
    if v, ok := n.(int); ok {
        println(v)
    }
}
```

### 检测通道是否关闭

```
package main

func main() {
    ch := make(chan int)

    go func() {
        for i := 0; i < 5; i++ {
            ch <- i
        }
        close(ch)
    }()

    for {
        if v, ok := <-ch; ok {
            println(v)
        } else {
            return
        }
    }
}

// $ go run main.go
// 输出如下
// 0
// 1
// 2
// 3
// 4
```

## 附加内容

## 闭包

有时候，我们可以利用 `闭包` 实现一些短小精悍的内部函数。

### 计数器

```
package main

func main() {
    newSeqInc := func() func() int {
        seq := 0
        return func() int {
            seq++
            return seq
        }
    }

    seq := newSeqInc() // 创建一个计数器
    println(seq())     // 1
    println(seq())     // 2
    println(seq())     // 3

    seq2 := newSeqInc() // 创建另一个计数器
    println(seq2())     // 1
    println(seq2())     // 2
    println(seq2())     // 3
}
```

## 小结

下面表格列出了常用的 `设计模式`，其中 Go 标准库自带的 `模式` 已经用删除线标识，读者可以和自己常用的 `设计模式` 进行对比。

| 创建型模式 | 结构性模式 | 行为型模式 |
| ---------- | ---------- | ---------- |
| 单例       | 适配器     | 策略       |
| 简单工厂   | 装饰者     | 观察者     |
| 抽象工厂   | 代理       | 状态       |
| 对象池     |            | 责任链     |
| 构建       |            |            |

长期以来，`设计模式` 一直处于尴尬的位置：初学者被各种概念和关系搞得不知所云，有经验的程序员会觉得 “这种代码写法 (这里指设计模式)，我早就知道了啊”。 鉴于这种情况，本文中没有涉及到的 `设计模式`，笔者不打算再一一描述，感兴趣的读者可以直接跳到 仓库代码[1] 查看示例代码。

> 相比于设计模式，更重要的是理解语言本身的特性以及最佳实践。



#### 引用链接

`[1]` 仓库代码: *https://github.com/duanbiaowu/go-examples-for-beginners/tree/master/patterns*
`[2]` 设计模式 - 维基百科: *https://zh.wikipedia.org/wiki/%E8%AE%BE%E8%AE%A1%E6%A8%A1%E5%BC%8F_(%E8%AE%A1%E7%AE%97%E6%9C%BA)*
`[3]` go-examples-for-beginners/patterns: *https://github.com/duanbiaowu/go-examples-for-beginners/tree/master/patterns*
`[4]` 圣杯与银弹 · 没用的设计模式: *https://draveness.me/holy-grail-design-pattern/*
`[5]` tmrts/go-patterns: *https://github.com/tmrts/go-patterns*
`[6]` DESIGN PATTERNS in GO: *https://refactoring.guru/design-patterns/go*
`[7]` 解密“设计模式”: *http://www.yinwang.org/blog-cn/2013/03/07/design-patterns*
`[8]` Go 编程模式 - 酷壳: *https://coolshell.cn/articles/series/go%e7%bc%96%e7%a8%8b%e6%a8%a1%e5%bc%8f*