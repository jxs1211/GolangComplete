# Go的HTTP tracing

brantou [Go生态](javascript:void(0);) *2022-05-25 15:00* *Posted on 北京*

收录于合集#golang9个

**原文链接：**https://brantou.github.io/2017/05/24/go-http-trace/

**官方博文：**https://go.dev/blog/http-tracing

在Go 1.7中，引入了 *HTTP tracing* ，这是在HTTP客户端请求的整个生命周期中收集细粒度信息的工具。由 **net/http/httptrace**[1] 包提供 *HTTP tracing* 的支持。收集的信息可用于调试延迟问题，服务监控，编写自适应系统等。

## 1 HTTP事件

httptrace包提供了许多钩子，用于在HTTP往返期间收集各种事件的信息。这些事件包括：

- 连接创建
- 连接复用
- DNS 查询
- 将请求写入网路
- 读取响应

## 2 跟踪事件

可以通过将包含钩子函数的 ***httptrace.ClientTrace**[2] 放在请求的 **context.Context**[3] 中来启用 *HTTP tracing* 。**http.RoundTripper**[4] 通过查找 *context* 的 **httptrace.ClientTrace* , 并调用相关的钩子函数报告内部事件。

追踪范围限于请求的 *context* ，用户应在 *context* 上下文之前放置一个 **httptrace.ClientTrace* ， 然后才能启动请求。

```go
func main() {
  req, _ := http.NewRequest("GET", "https://google.com", nil)
  trace := &httptrace.ClientTrace{
    GotConn: func(connInfo httptrace.GotConnInfo) {
      fmt.Printf("Got Conn: %+v\n", connInfo)
    },
    DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
      fmt.Printf("DNS Info: %+v\n", dnsInfo)
    },
  }
  req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
  _, err := http.DefaultTransport.RoundTrip(req)
  if err != nil {
    log.Fatal(err)
  }
}
DNS Info: {Addrs:[{IP:192.168.83.230 Zone:}] Err:<nil> Coalesced:false}
Got Conn: {Conn:0xc42001ce00 Reused:false WasIdle:false IdleTime:0s}
```

在 *round trip* 中，/http.DefaultTransport/ 会在事件发生时调用每个钩子。一旦DNS查找完成，将打印DNS信息。当与请求的主机建立连接时，它将类似地打印连接信息。

## 3 跟踪http.Client

跟踪机制旨在跟踪单个http.Transport.RoundTrip的生命周期中的事件。但是，客户端可以进行多次往返，以完成HTTP请求。例如，在URL重定向的情况下，注册的钩子将被调用多次，客户端遵循HTTP重定向，进行多个请求。用户有职责在 *http.Client* 级别识别这些事件。下面的示例使用 *http.RoundTripper wrapper* 来标识当前的请求。

```go
package main

import (
  "fmt"
  "log"
  "net/http"
  "net/http/httptrace"
)

// transport is an http.RoundTripper that keeps track of the in-flight
// request and implements hooks to report HTTP tracing events.
type transport struct {
  current *http.Request
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to keep track
// of the current request.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
  t.current = req
  return http.DefaultTransport.RoundTrip(req)
}

// GotConn prints whether the connection has been used previously
// for the current request.
func (t *transport) GotConn(info httptrace.GotConnInfo) {
  fmt.Printf("Connection reused for %v? %v\n", t.current.URL, info.Reused)
}

func main() {
  t := &transport{}

  req, _ := http.NewRequest("GET", "https://google.com", nil)
  trace := &httptrace.ClientTrace{
    GotConn: t.GotConn,
  }
  req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

  client := &http.Client{Transport: t}
  if _, err := client.Do(req); err != nil {
    log.Fatal(err)
  }
}
```

上面示例从 google.com 重定向到 www.google.com， 输出如下：

```sh
Connection reused for https://google.com? false
Connection reused for https://www.google.com.hk/?gfe_rd=cr&ei=olwkWd3BAa-M8Qfjs73IBA? false
```

net/http包中的 *Transport* 支持跟踪 HTTP/1 和 HTTP/2 的 request。

如果你是自定义 *http.RoundTripper* 实现的作者，则可以通过检查 **httptest.ClientTrace* 的请求 *context* 来支持跟踪，并在事件发生时调用相关的钩子。

## 4 总结

对于那些有兴趣调试HTTP请求延迟和编写工具来进行出站流量的网络调试的人来说， *HTTP tracing* 是一个有价值的补充。通过启用这个新工具，希望看到来自社区的HTTP调试，基准测试和可视化工具，如**httpstat**[5]。

### 参考资料

[1]net/http/httptrace: *https://golang.org/pkg/net/http/httptrace/*[2]*httptrace.ClientTrace: *https://golang.org/pkg/net/http/httptrace/#ClientTrace*[3]context.Context: *https://golang.org/pkg/context/#Context*[4]http.RoundTripper: *https://golang.org/pkg/net/http/#RoundTripper*[5]httpstat: *https://github.com/davecheney/httpstat*

