# Go：服务发现原理分析与源码解读

https://mp.weixin.qq.com/s/CPi_sPaTUR1JYl7NEybFhQ

zhoushuguang [Go语言中文网](javascript:void(0);) *2022-08-18 08:52* *发表于北京*

在微服务架构中，有许多绕不开的技术话题。比如服务发现、负载均衡、指标监控、链路追踪，以及服务治理相关的超时控制、熔断、降级、限流等，还有RPC框架。这些都是微服务架构的基础，只有打牢这些基础，才敢说对微服务是有了一点理解，出门也好意思和别人打招呼了，被人提问的时候也能侃侃而谈了，线上出了问题往往也能寻根溯源内心不慌了，旁边的女同事小芳看着你的时候也是满眼的小可爱了。

在《微服务实践》公众号，之前写了《go-zero微服务实战系列》的系列文章，这个系列的文章更多的是偏向业务功能和高并发下的服务优化等。本人水平有限，难免有写的不足的地方，但也依然得到了大家的支持与鼓励，倍感荣幸，所以决定趁热打铁，乘胜追击，继续给大家输出干货。

《彻底搞懂系列》会基于 go-zero v1.3.5 和 grpc-go v1.47.0 和大家一起学习微服务架构的方方面面，主要形式是理论+源码+案例，如果时间允许也可能会加上配套视频。

本篇文章作为该系列的第一篇，会先介绍相对比较简单的服务发现相关内容。

撸袖子开搞，奥利给！！！

## **服务发现**

为什么在微服务架构中，需要引入服务发现呢？本质上，服务发现的目的是解耦程序对服务具体位置的依赖，对于微服务架构来说，服务发现不是可选的，而是必须的。因为在生产环境中服务提供方都是以集群的方式对外提供服务，集群中服务的IP随时都可能发生变化，比如服务重启，发布，扩缩容等，因此我们需要用一本“通讯录”及时获取到对应的服务节点，这个获取的过程其实就是“服务发现”。

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg17gVbgX0Cks9MeibYicExbmgfn3VAjDNAFcelms7GHDfNeQjO0t0yLpJsuvPHqLDkP68jZFLqsa8IA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

要理解服务发现，需要知道服务发现解决了如下三个问题：

- 服务的注册（Service Registration）

  当服务启动的时候，应该通过某种形式（比如调用API、产生上线事件消息、在Etcd中记录、存数据库等等）把自己（服务）的信息通知给服务注册中心，这个过程一般是由微服务框架来完成，业务代码无感知。

- 服务的维护（Service Maintaining）

  尽管在微服务框架中通常都提供下线机制，但并没有办法保证每次服务都能优雅下线（Graceful Shutdown），而不是由于宕机、断网等原因突然失联，所以，在微服务框架中就必须要尽可能的保证维护的服务列表的正确性，以避免访问不可用服务节点的尴尬。

- 服务的发现（Service Discovery）

  这里所说的发现是狭义的，它特指消费者从微服务框架（服务发现模块）中，把一个服务标识（一般是服务名）转换为服务实际位置（一般是ip地址）的过程。这个过程（可能是调用API，监听Etcd，查询数据库等）业务代码无感知。

服务发现有两种模式，分别是服务端服务发现和客户端服务发现，下面分别进行介绍。

### 服务端服务发现

对于服务端服务发现来说，服务调用方无需关注服务发现的具体细节，只需要知道服务的DNS域名即可，支持不同语言的接入，对基础设施来说，需要专门支持负载均衡器，对于请求链路来说多了一次网络跳转，可能会有性能损耗。也可以把咱们比较熟悉的 nginx 反向代理理解为服务端服务发现。

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg17gVbgX0Cks9MeibYicExbmgCu9sKMiboHvD9DWnhewcQ42rwmZ7ibz0d1V4CUfvGpjguDMH5IaZnjVA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

### 客户端服务发现

对于客户端服务发现来说，由于客户端和服务端采用了直连的方式，比服务端服务发现少了一次网络跳转，对于服务调用方来说需要内置负载均衡器，不同的语言需要各自实现。

对于微服务架构来说，我们期望的是去中心化依赖，中心化的依赖会让架构变得复杂，当出现问题的时候也会让整个排查链路变得繁琐，所以在 go-zero 中采用的是客户端服务发现的模式。

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg17gVbgX0Cks9MeibYicExbmg3iaVsVAeicpBFdCETR6MDibAMeIxJpByxIGicY4opiaGlFWXrD66gyYWJlw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

## **gRPC的服务发现**

gRPC提供了自定义Resolver的能力来实现服务发现，通过 **Register**方法来进行注册自定义的Resolver，自定义的Resolver需要实现Builder接口，定义如下：

**grpc-go/resolver/resolver.go:261**

```
type Builder interface {
    Build(target Target, cc ClientConn, opts BuildOptions) (Resolver, error)
    Scheme() string
}
```

先说下 `Scheme()` 方法的作用，该方法返回一个stirng。注册的 `Resolver` 会被保存在一个全局的变量m中，m是一个map，这个map的key即为 `Scheme()` 方法返回的字符串。也就是多个Resolver是通过Scheme来进行区分的，所以我们定义 `Resolver` 的时候 `Scheme` 不要重复，否则 `Resolver` 就会被覆盖。

**grpc-go/resolver/resolver.go:49**

```
func Register(b Builder) {
    m[b.Scheme()] = b
}
```

再来看下Build方法，Build方法有三个参数，还有Resolver返回值，乍一看不知道这些参数是干嘛的，遇到这种情况该怎么办呢？其实也很简单，去源码里看一下Build方法在哪里被调用的，就知道传入的参数是哪里来的，是什么含义了。

使用gRPC进行服务调用前，需要先创建一个 **ClientConn** 对象，最终发起调用的时候，其实是调用了 **ClientConn** 的 **Invoke** 方法，可以看下如下代码，其中 **ClientConn** 是通过调用 **NewGreeterClient** 传入的，**NewGreeterClient** 为 **protoc** 自动生成的代码，并赋值给 **cc** 属性，示例代码中创建 **ClientConn** 调用的是 **Dial** 方法，底层也会调用 **DialContext**：

**grpc-go/clientconn.go:104**

```
func Dial(target string, opts ...DialOption) (*ClientConn, error) {
    return DialContext(context.Background(), target, opts...)
}
```

创建 **ClientConn** 对象，并传递给自动生成的 **greeterClient**

**grpc-go/examples/helloworld/greeter_client/main.go:42**

```
func main() {
    flag.Parse()

    // Set up a connection to the server.
    conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    c := pb.NewGreeterClient(conn)
    // Contact the server and print out its response.
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("Greeting: %s", r.GetMessage())
}
```

最终通过 **Invoke** 方法真正发起调用请求。

**grpc-go/examples/helloworld/helloworld/helloworld_grpc.pb.go:39**

```
func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
    out := new(HelloReply)
    err := c.cc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, opts...)
    if err != nil {
        return nil, err
    }
    return out, nil
}
```

在了解了客户端调用发起的流程后，我们重点看下 **ClientConn** 方法，该方法巨长，只看我们关注的Resolver部分。**ClientConn** 第二个参数 **Target** 的语法可以参考 **https://github.com/grpc/grpc/blob/master/doc/naming.md** ，采用了URI的格式，其中第一部分表示Resolver的名称，即自定义Builder方法Scheme的返回值。格式如下：

```
dns:[//authority/]host[:port] -- DNS（默认）
```

继续往下看，通过调用 **parseTargetAndFindResolver** 方法来获取Resolver

**grpc-go/clientconn.go:251**

```
resolverBuilder, err := cc.parseTargetAndFindResolver()
```

在 **parseTargetAndFindResolver** 方法中，主要就是把 **target** 中的resolver name解析出来，然后根据resolver name去上面我们提到的保存Resolver的全局变量m中去找对应的Resolver。

**grpc-go/clientconn.go:1574**

```
func (cc *ClientConn) parseTargetAndFindResolver() (resolver.Builder, error) {
    // 非关键代码省略 ...
  
    var rb resolver.Builder
    parsedTarget, err := parseTarget(cc.target)
  
    // 非关键代码省略 ...
  
    rb = cc.getResolver(parsedTarget.Scheme)
    if rb == nil {
        return nil, fmt.Errorf("could not get resolver for default scheme: %q", parsedTarget.Scheme)
    }
    cc.parsedTarget = parsedTarget
    return rb, nil
}
```

接着往下看，找到我们自己注册的Resolver之后，又调用了 **newCCResolverWrapper** 方法，把我们自己的Resolver也传了进去

**grpc-go/clientconn.go:292**

```
rWrapper, err := newCCResolverWrapper(cc, resolverBuilder)
```

进入到 **newCCResolverWrapper** 方法中，在这个方法中终于找到了我们自定义的 **Builder** 的 **Build** 方法在哪里被调用了，在 **grpc-go/resolver_conn_wrapper.go:72** 调用了我们自定义的Build方法，其中第一参数target传入的为cc.parseTarget，cc为 **newCCResolverWrapper** 第一个参数，即 **ClientConn** 对象。cc.parseTarget是在上面提到的获取自定义Resolver方法 **parseTargetAndFindResolver** 中最后赋值的，其中Scheme、Authority、Endpoint分别对应Target语法中定义的三部分，这几个属性即将被废弃，只保留URL属性，定义如下：

**grpc-go/resolver/resolver.go:245**

```
type Target struct {
    // Deprecated: use URL.Scheme instead.
    Scheme string
    // Deprecated: use URL.Host instead.
    Authority string
    // Deprecated: use URL.Path or URL.Opaque instead. The latter is set when
    // the former is empty.
    Endpoint string
    // URL contains the parsed dial target with an optional default scheme added
    // to it if the original dial target contained no scheme or contained an
    // unregistered scheme. Any query params specified in the original dial
    // target can be accessed from here.
    URL url.URL
}
```

URL的Scheme对应Target的Scheme，URL的Host对应Target的Authority，URL的Path对应Target的Endpoint

**/usr/local/go/src/net/url/url.go:358**

```
type URL struct {
    Scheme      string
    Opaque      string    // encoded opaque data
    User        *Userinfo // username and password information
    Host        string    // host or host:port
    Path        string    // path (relative paths may omit leading slash)
    RawPath     string    // encoded path hint (see EscapedPath method)
    ForceQuery  bool      // append a query ('?') even if RawQuery is empty
    RawQuery    string    // encoded query values, without '?'
    Fragment    string    // fragment for references, without '#'
    RawFragment string    // encoded fragment hint (see EscapedFragment method)
}
```

继续看传入自定义Build方法的第二个参数cc，这个cc参数是一个接口 **ClientConn**，不要和我们之前讲的创建客户端调用用的 ClientConn混淆，这个 **ClientConn**定义如下：

**grpc-go/resolver/resolver.go:203**

```
type ClientConn interface {
    UpdateState(State) error
    ReportError(error)
    NewAddress(addresses []Address)
    NewServiceConfig(serviceConfig string)
    ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult
}
```

**ccResolverWrapper** 实现了这个接口，并作为自定义Build方法的第二个参数传入

**grpc-go/resolver_conn_wrapper.go:36**

```
type ccResolverWrapper struct {
    cc         *ClientConn
    resolverMu sync.Mutex
    resolver   resolver.Resolver
    done       *grpcsync.Event
    curState   resolver.State

    incomingMu sync.Mutex // Synchronizes all the incoming calls.
}
```

自定义Build方法的第三个参数为一些配置项，newCCResolverWrapper实现如下：

**grpc-go/resolver_conn_wrapper.go:48**

```
func newCCResolverWrapper(cc *ClientConn, rb resolver.Builder) (*ccResolverWrapper, error) {
    ccr := &ccResolverWrapper{
        cc:   cc,
        done: grpcsync.NewEvent(),
    }

    var credsClone credentials.TransportCredentials
    if creds := cc.dopts.copts.TransportCredentials; creds != nil {
        credsClone = creds.Clone()
    }
    rbo := resolver.BuildOptions{
        DisableServiceConfig: cc.dopts.disableServiceConfig,
        DialCreds:            credsClone,
        CredsBundle:          cc.dopts.copts.CredsBundle,
        Dialer:               cc.dopts.copts.Dialer,
    }

    var err error
    ccr.resolverMu.Lock()
    defer ccr.resolverMu.Unlock()
    ccr.resolver, err = rb.Build(cc.parsedTarget, ccr, rbo)
    if err != nil {
        return nil, err
    }
    return ccr, nil
}
```

好了，到这里我们已经知道了自定Resolver的Build方法在哪里被调用，以及传入的参数的由来以及含义，如果你是第一次看gRPC源码的话可能现在已经有点懵了，可以多读几遍，为大家提供了时序图配合代码阅读效果更佳：

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg17gVbgX0Cks9MeibYicExbmgGOsuGkgLzKAncNsYlgXYNWCIfHXXicCN9nSg8813ZcsMPVqjAV3LibUA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

## **go-zero中如何实现的服务发现**

通过对gRPC服务发现相关内容的学习，我们大概已经知道了服务发现是怎么回事了，有了理论，接下来我们就一起看下go-zero是如何基于gRPC做服务发现的。

通过上面的时序图可以看到第一步是需要自定义Resolver，第二步注册自定义的Resolver。

go-zero的服务发现是在客户端实现的。在创建zRPC客户端的时候，通过init方法进行了自定义Resolver的注册。

**go-zero/zrpc/internal/client.go:23**

```
func init() {
    resolver.Register()
}
```

在go-zero中默认注册了四个自定义的Resolver。

**go-zero/zrpc/resolver/internal/resolver.go:35**

```
func RegisterResolver() {
    resolver.Register(&directResolverBuilder)
    resolver.Register(&discovResolverBuilder)
    resolver.Register(&etcdResolverBuilder)
    resolver.Register(&k8sResolverBuilder)
}
```

通过 **goctl** 自动生成的rpc代码默认使用的是etcd作为服务注册与发现组件的，因此我们重点来看下go-zero是如何基于etcd实现服务注册与发现的。

etcdBuilder返回的Scheme值为etcd

**go-zero/zrpc/resolver/internal/etcdbuilder.go:7**

```
func (b *etcdBuilder) Scheme() string {
    return EtcdScheme
}
```

**go-zero/zrpc/resolver/internal/resolver.go:15**

```
EtcdScheme = "etcd"
```

还记得我们上面讲过的吗？在时序图的第五步和第六步，会通过scheme去全局的m中寻找自定义的Resolver，而scheme是从DialContext第二个参数target中解析出来的，那我们看下go-zero调用DialContext的时候，传入的target值是什么。target是通过 **BuildTarget** 方法获取来的，定义如下：

**go-zero/zrpc/config.go:72**

```
func (cc RpcClientConf) BuildTarget() (string, error) {
    if len(cc.Endpoints) > 0 {
        return resolver.BuildDirectTarget(cc.Endpoints), nil
    } else if len(cc.Target) > 0 {
        return cc.Target, nil
    }

    if err := cc.Etcd.Validate(); err != nil {
        return "", err
    }

    if cc.Etcd.HasAccount() {
        discov.RegisterAccount(cc.Etcd.Hosts, cc.Etcd.User, cc.Etcd.Pass)
    }
    if cc.Etcd.HasTLS() {
        if err := discov.RegisterTLS(cc.Etcd.Hosts, cc.Etcd.CertFile, cc.Etcd.CertKeyFile,
            cc.Etcd.CACertFile, cc.Etcd.InsecureSkipVerify); err != nil {
            return "", err
        }
    }

    return resolver.BuildDiscovTarget(cc.Etcd.Hosts, cc.Etcd.Key), nil
}
```

最终生成target结果的方法如下，也就是对于etcd来说，最终生成的target格式为：

```
etcd://127.0.0.1:2379/product.rpc
```

**go-zero/zrpc/resolver/target.go:17**

```
func BuildDiscovTarget(endpoints []string, key string) string {
    return fmt.Sprintf("%s://%s/%s", internal.DiscovScheme,
        strings.Join(endpoints, internal.EndpointSep), key)
}
```

似乎有点不对劲，scheme不应该是etcd么？为什么是discov？其实是因为etcd和discov共用了一套Resolver逻辑，也就是gRPC通过scheme找到已经注册的discov Resolver，该Resolver对应的Build方法同样适用于etcd，discov可以认为是对服务发现的一个抽象，etcdResolver的定义如下：

**go-zero/zrpc/resolver/internal/etcdbuilder.go:3**

```
type etcdBuilder struct {
    discovBuilder
}
```

### 服务注册

在详细看基于etcd的自定义Resolver逻辑之前，我们先来看下go-zero的服务注册，即如何把服务信息注册到etcd中的，我们以 **lebron/apps/product/rpc** 这个服务为例进行说明。

在product-rpc的配置文件中配置了Etcd，包括etcd的地址和服务对应的key，如下：

**lebron/apps/product/rpc/etc/product.yaml:4**

```
ListenOn: 127.0.0.1:9002

Etcd:
  Hosts:
  - 127.0.0.1:2379
  Key: product.rpc
```

调用zrpc.MustNewServer创建gRPC server，接着会调用 **NewRpcPubServer** 方法，定义如下：

**go-zero/zrpc/internal/rpcpubserver.go:17**

```
func NewRpcPubServer(etcd discov.EtcdConf, listenOn string, opts ...ServerOption) (Server, error) {
    registerEtcd := func() error {
        pubListenOn := figureOutListenOn(listenOn)
        var pubOpts []discov.PubOption
        if etcd.HasAccount() {
            pubOpts = append(pubOpts, discov.WithPubEtcdAccount(etcd.User, etcd.Pass))
        }
        if etcd.HasTLS() {
            pubOpts = append(pubOpts, discov.WithPubEtcdTLS(etcd.CertFile, etcd.CertKeyFile,
                etcd.CACertFile, etcd.InsecureSkipVerify))
        }
        pubClient := discov.NewPublisher(etcd.Hosts, etcd.Key, pubListenOn, pubOpts...)
        return pubClient.KeepAlive()
    }
    server := keepAliveServer{
        registerEtcd: registerEtcd,
        Server:       NewRpcServer(listenOn, opts...),
    }

    return server, nil
}
```

在启动Server的时候，调用Start方法，在Start方法中会调用registerEtcd进行真正的服务注册

**go-zero/zrpc/internal/rpcpubserver.go:44**

```
func (s keepAliveServer) Start(fn RegisterFn) error {
    if err := s.registerEtcd(); err != nil {
        return err
    }

    return s.Server.Start(fn)
}
```

在KeepAlive方法中，首先创建etcd连接，然后调用register方法进行服务注册，在register首先创建租约，租约默认时间为10秒钟，最后通过Put方法进行注册。

**go-zero/core/discov/publisher.go:125**

```
func (p *Publisher) register(client internal.EtcdClient) (clientv3.LeaseID, error) {
    resp, err := client.Grant(client.Ctx(), TimeToLive)
    if err != nil {
        return clientv3.NoLease, err
    }

    lease := resp.ID
    if p.id > 0 {
        p.fullKey = makeEtcdKey(p.key, p.id)
    } else {
        p.fullKey = makeEtcdKey(p.key, int64(lease))
    }
    _, err = client.Put(client.Ctx(), p.fullKey, p.value, clientv3.WithLease(lease))

    return lease, err
}
```

key的规则定义如下，其中key为在配置文件中配置的Key，这里为product.rpc，id为租约id。value为服务的地址。

**go-zero/core/discov/clients.go:39**

```
func makeEtcdKey(key string, id int64) string {
    return fmt.Sprintf("%s%c%d", key, internal.Delimiter, id)
}
```

在了解了服务注册的流程后，我们启动product-rpc服务，然后通过如下命令查看服务注册的地址：

```
$ etcdctl get product.rpc --prefix
product.rpc/7587864068988009477
127.0.0.1:9002
```

在 **KeepAlive** 方法中，服务注册完后，最后会调用 **keepAliveAsync** 进行租约的续期，以保证服务一直是存活的状态，如果服务异常退出了，那么也就无法进行续期，服务发现也就能自动识别到该服务异常下线了。

### 服务发现

现在已经把服务注册到etcd中了，继续来看如何发现这些服务地址。我们回到 **etcdBuilder** 的Build方法的实现。

还记得第一个参数target是什么吗？如果不记得了可以往上翻再复习一下，首先从target中解析出etcd的地址，和服务对应的key。然后创建etcd连接，接着执行update方法，在update方法中，通过调用cc.UpdateState方法进行服务状态的更新。

**go-zero/zrpc/resolver/internal/discovbuilder.go:14**

```
func (b *discovBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (
    resolver.Resolver, error) {
    hosts := strings.FieldsFunc(targets.GetAuthority(target), func(r rune) bool {
        return r == EndpointSepChar
    })
    sub, err := discov.NewSubscriber(hosts, targets.GetEndpoints(target))
    if err != nil {
        return nil, err
    }

    update := func() {
        var addrs []resolver.Address
        for _, val := range subset(sub.Values(), subsetSize) {
            addrs = append(addrs, resolver.Address{
                Addr: val,
            })
        }
        if err := cc.UpdateState(resolver.State{
            Addresses: addrs,
        }); err != nil {
            logx.Error(err)
        }
    }
    sub.AddListener(update)
    update()

    return &nopResolver{cc: cc}, nil
}
```

如果忘记了Build方法第二个参数cc的话，可以往上翻翻再复习一下，cc.UpdateState方法定义如下，最终会调用 **ClientConn** 的 **updateResolverState** 方法：

**grpc-go/resolver_conn_wrapper.go:94**

```
func (ccr *ccResolverWrapper) UpdateState(s resolver.State) error {
    ccr.incomingMu.Lock()
    defer ccr.incomingMu.Unlock()
    if ccr.done.HasFired() {
        return nil
    }
    ccr.addChannelzTraceEvent(s)
    ccr.curState = s
    if err := ccr.cc.updateResolverState(ccr.curState, nil); err == balancer.ErrBadResolverState {
        return balancer.ErrBadResolverState
    }
    return nil
}
```

继续看 **Build** 方法，update方法会被添加到事件监听中，当有PUT和DELETE事件触发，都会调用update方法进行服务状态的更新，事件监听是通过etcd的Watch机制实现，代码如下：

**go-zero/core/discov/internal/registry.go:295**

```
func (c *cluster) watchStream(cli EtcdClient, key string) bool {
    rch := cli.Watch(clientv3.WithRequireLeader(c.context(cli)), makeKeyPrefix(key), clientv3.WithPrefix())
    for {
        select {
        case wresp, ok := <-rch:
            if !ok {
                logx.Error("etcd monitor chan has been closed")
                return false
            }
            if wresp.Canceled {
                logx.Errorf("etcd monitor chan has been canceled, error: %v", wresp.Err())
                return false
            }
            if wresp.Err() != nil {
                logx.Error(fmt.Sprintf("etcd monitor chan error: %v", wresp.Err()))
                return false
            }

            c.handleWatchEvents(key, wresp.Events)
        case <-c.done:
            return true
        }
    }
}
```

当有事件触发的时候，会调用事件处理函数 **handleWatchEvents** ，最终会调用 **Build** 方法中定义的update进行服务状态的更新：

**go-zero/core/discov/internal/registry.go:172**

```
func (c *cluster) handleWhandleWatchEventsatchEvents(key string, events []*clientv3.Event) {
    c.lock.Lock()
    listeners := append([]UpdateListener(nil), c.listeners[key]...)
    c.lock.Unlock()

    for _, ev := range events {
        switch ev.Type {
        case clientv3.EventTypePut:
            c.lock.Lock()
            if vals, ok := c.values[key]; ok {
                vals[string(ev.Kv.Key)] = string(ev.Kv.Value)
            } else {
                c.values[key] = map[string]string{string(ev.Kv.Key): string(ev.Kv.Value)}
            }
            c.lock.Unlock()
            for _, l := range listeners {
                l.OnAdd(KV{
                    Key: string(ev.Kv.Key),
                    Val: string(ev.Kv.Value),
                })
            }
        case clientv3.EventTypeDelete:
            c.lock.Lock()
            if vals, ok := c.values[key]; ok {
                delete(vals, string(ev.Kv.Key))
            }
            c.lock.Unlock()
            for _, l := range listeners {
                l.OnDelete(KV{
                    Key: string(ev.Kv.Key),
                    Val: string(ev.Kv.Value),
                })
            }
        default:
            logx.Errorf("Unknown event type: %v", ev.Type)
        }
    }
}
```

第一次会调用 **load** 方法，获取key对应的服务列表，通过etcd前缀匹配的方式获取，获取方式如下：

```
func (c *cluster) load(cli EtcdClient, key string) {
    var resp *clientv3.GetResponse
    for {
        var err error
        ctx, cancel := context.WithTimeout(c.context(cli), RequestTimeout)
        resp, err = cli.Get(ctx, makeKeyPrefix(key), clientv3.WithPrefix())
        cancel()
        if err == nil {
            break
        }

        logx.Error(err)
        time.Sleep(coolDownInterval)
    }

    var kvs []KV
    for _, ev := range resp.Kvs {
        kvs = append(kvs, KV{
            Key: string(ev.Key),
            Val: string(ev.Value),
        })
    }

    c.handleChanges(key, kvs)
}
```

获取的服务地址列表，通过map存储在本地，当有事件触发的时候通过操作map进行服务列表的更新，这里有个隐藏的设计考虑是当 etcd 连不上或者出现故障时，内存里的服务地址列表不会被更新，保障了当 etcd 有问题时，服务发现依然可以工作，保障服务继续正常运行。逻辑相对比较直观，这里就不再赘述，代码逻辑在 **go-zero/core/discov/subscriber.go:76** ，下面是go-zero服务发现的时序图

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg17gVbgX0Cks9MeibYicExbmgk2IGp9HmicAqIXWicSqsaYpR0CLbMIK2tliaTWt0bOojZqYOn5Z2iazoow/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

## **结束语**

到这里服务发现相关的内容已经讲完了，内容还是有点多的，特别是代码部分需要反复仔细阅读才能加深理解。

我们一起来简单回顾下本篇的内容：

- 首先介绍了服务发现的概念，以及服务发现需要解决哪些问题

- 服务发现的两种模式，分别是服务端发现模式和客户端发现模式

- 接着一起学习了gRPC提供的注册Resolver的能力，通过注册Resolver来实现自定义的服务发现功能，以及gRPC内部是如何寻找到自定义的Resolver和触发调用自定义Resolver的逻辑

- 最后学习了go-zero中服务发现的实现原理，

- - 先是介绍了go-zero的服务注册流程，演示了最终注册的效果
  - 接着从自定义Resolver的Build方法出发，了解到先是通过前缀匹配的方式获取对应的服务列表存在本地，然后调用UpdateState方法更新服务状态
  - 通过Watch的方式监听服务状态的变化，监听到变化后会触发调用update方法更新本地的服务列表和调用UpdateState更新服务的状态。

服务发现是理解微服务架构的基础，希望大家能仔细的阅读本文，如果有疑问可以随时找我讨论，在社区群中可以搜索dawn_zhou找到我。

通过服务发现获取到服务列表后，接着就会通过Invoke方法进行服务调用，在服务调用的时候就涉及到负载均衡，通过负载均衡选择一个合适的节点发起请求。负载均衡是下一篇文章要讲的内容，敬请期待。

希望本篇文章对你有所帮助，你的点赞是作者持续输出的最大动力。