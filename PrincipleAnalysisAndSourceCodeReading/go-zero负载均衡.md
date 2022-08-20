# Go：负载均衡原理分析与源码解读

https://mp.weixin.qq.com/s/5Xskm1gZmaOw0jatrJOMOA

zhoushuguang [Go语言中文网](javascript:void(0);) *2022-08-23 08:52* *发表于北京*

上一篇文章一起学习了Resolver的原理和源码分析，本篇继续和大家一起学习下和Resolver关系密切的Balancer的相关内容。这里说的负载均衡主要指数据中心内的负载均衡，即RPC间的负载均衡。

传送门 [服务发现原理分析与源码解读](http://mp.weixin.qq.com/s?__biz=MzAxMTA4Njc0OQ==&mid=2651453321&idx=1&sn=8eec649f6c16e803c8bd8c9a3dedb3e4&chksm=80bb297bb7cca06dc70ca34da83fcba0581cad6d8cb3295bfbd452e492f49bedfc8a9f099071&scene=21#wechat_redirect)

基于 go-zero v1.3.5 和 grpc-go v1.47.0

## **负载均衡**

每一个被调用服务都会有多个实例，那么服务的调用方应该将请求，发向被调用服务的哪一个服务实例，这就是负载均衡的业务场景。

负载均衡的第一个关键点是公平性，即负载均衡需要关注被调用服务实例组之间的公平性，不要出现旱的旱死，涝的涝死的情况。

负载均衡的第二个关键点是正确性，即对于有状态的服务来说，负载均衡需要关心请求的状态，将请求调度到能处理它的后端实例上，不要出现不能处理和错误处理的情况。

### 无状态的负载均衡

无状态的负载均衡是我们日常工作中接触比较多的负载均衡模型，它指的是参与负载均衡的后端实例是无状态的，所有的后端实例都是对等的，一个请求不论发向哪一个实例，都会得到相同的并且正确的处理结果，所以无状态的负载均衡策略不需要关心请求的状态。下面介绍两种无状态负载均衡算法。

#### 轮询

轮询的负载均衡策略非常简单，只需要将请求按顺序分配给多个实例，不用再做其他的处理。例如，轮询策略会将第一个请求分配给第一个实例，然后将下一个请求分配给第二个实例，这样依次分配下去，分配完一轮之后，再回到开头分配给第一个实例，再依次分配。轮询在路由时，不利用请求的状态信息，属于无状态的负载均衡策略，所以它不能用于有状态实例的负载均衡器，否则正确性会出现问题。在公平性方面，因为轮询策略只是按顺序分配请求，所以适用于请求的工作负载和实例的处理能力差异都较小的情况。

#### 权重轮询

权重轮询的负载均衡策略是将每一个后端实例分配一个权重，分配请求的数量和实例的权重成正比轮询。例如有两个实例 A，B，假设我们设置 A 的权重为 20，B 的权重为 80，那么负载均衡会将 20% 的请求数量分配给 A，80 % 的请求数量分配给 B。权重轮询在路由时，不利用请求的状态信息，属于无状态的负载均衡策略，所以它也不能用于有状态实例的负载均衡器，否则正确性会出现问题。在公平性方面，因为权重策略会按实例的权重比例来分配请求数，所以，我们可以利用它解决实例的处理能力差异的问题，认为它的公平性比轮询策略要好。

### 有状态负载均衡

有状态负载均衡是指，在负载均衡策略中会保存服务端的一些状态，然后根据这些状态按照一定的算法选择出对应的实例。

#### P2C+EWMA

在go-zero中默认使用的是P2C的负载均衡算法。该算法的原理比较简单，即随机从所有可用节点中选择两个节点，然后计算这两个节点的负载情况，选择负载较低的一个节点来服务本次请求。为了避免某些节点一直得不到选择导致不平衡，会在超过一定的时间后强制选择一次。

在该复杂均衡算法中，采用了EWMA指数移动加权平均的算法，表示是一段时间内的均值。该算法相对于算数平均来说对于突然的网络抖动没有那么敏感，突然的抖动不会体现在请求的lag中，从而可以让算法更加均衡。

**go-zero/zrpc/internal/balancer/p2c/p2c.go:133**

```
atomic.StoreUint64(&c.lag, uint64(float64(olag)*w+float64(lag)*(1-w)))
```

**go-zero/zrpc/internal/balancer/p2c/p2c.go:139**

```
atomic.StoreUint64(&c.success, uint64(float64(osucc)*w+float64(success)*(1-w)))
```

系数w是一个时间衰减值，即两次请求的间隔越大，则系数w就越小。

**go-zero/zrpc/internal/balancer/p2c/p2c.go:124**

```
w := math.Exp(float64(-td) / float64(decayTime))
```

节点的load值是通过该连接的请求延迟 **lag** 和当前请求数 **inflight** 的乘积所得，如果请求的延迟越大或者当前正在处理的请求数越多表明该节点的负载越高。

**go-zero/zrpc/internal/balancer/p2c/p2c.go:199**

```
func (c *subConn) load() int64 {
  // plus one to avoid multiply zero
  lag := int64(math.Sqrt(float64(atomic.LoadUint64(&c.lag) + 1)))
  load := lag * (atomic.LoadInt64(&c.inflight) + 1)
  if load == 0 {
    return penalty
  }

  return load
}
```

## **源码分析**

> 如下源码会涉及go-zero和gRPC，请根据给出的代码路径进行区分

在gRPC中，Balancer和Resolver一样也可以自定义，同样也是通过Register方法进行注册

**grpc-go/balancer/balancer.go:53**

```
func Register(b Builder) {
  m[strings.ToLower(b.Name())] = b
}
```

Register的参数Builder为接口，在Builder接口中，Build方法的第一个参数ClientConn也为接口，Build方法的返回值Balancer同样也是接口，定义如下：

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg00wMY19Rickys6twWMhF9PqowaFicADYYpibpnrn8sym2quDMfh0V8VhoUrhgDVpcM6p36WZIvnKJaA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

可以看出，要想实现自定义的Balancer的话，就必须要实现balancer.Builder接口。

在了解了gRPC提供的Balancer的注册方式之后，我们看一下go-zero是在什么地方进行Balancer注册的

**go-zero/zrpc/internal/balancer/p2c/p2c.go:36**

```
func init() {
  balancer.Register(newBuilder())
}
```

在go-zero中并没有实现 **balancer.Builder** 接口，而是使用gRPC提供的 **base.baseBuilder** 进行注册，**base.baseBuilder** 实现了**balancer.Builder** 接口。创建baseBuilder的时候调用了 **base.NewBalancerBuilder** 方法，需要传入 **PickerBuilder** 参数，PickerBuilder为接口，在go-zero中 **p2c.p2cPickerBuilder** 实现了该接口。

PickerBuilder接口Build方法返回值 **balancer.Picker** 也是一个接口，**p2c.p2cPicker** 实现了该接口。

**grpc-go/balancer/base/base.go:65**

```
func NewBalancerBuilder(name string, pb PickerBuilder, config Config) balancer.Builder {
  return &baseBuilder{
    name:          name,
    pickerBuilder: pb,
    config:        config,
  }
}
```

各结构之间的关系如下图所示，其中各结构模块对应的包为：

- balancer：grpc-go/balancer
- base：grpc-go/balancer/base
- p2c: go-zero/zrpc/internal/balancer/p2c

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg00wMY19Rickys6twWMhF9PqsJ3uR4uNM6CvdplHy8G8ZZTftCJS0JWQsucDSDicqu4fNfKDBXILgWg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

### 在哪里获取已注册的Balancer？

通过上面的流程步骤，已经知道了如何自定义Balancer，以及如何注册自定义的Blancer。既然注册了肯定就会获取，接下来看一下是在哪里获取已经注册的Balancer的。

我们知道Resolver是通过解析DialContext的第二个参数target，从而得到Resolver的name，然后根据name获取到对应的Resolver的。获取Balancer同样也是根据名称，Balancer的名称是在创建gRPC Client的时候通过配置项传入的，这里的p2c.Name为注册Balancer时指定的名称 p2c_ewma ，如下：

**go-zero/zrpc/internal/client.go:50**

```
func NewClient(target string, opts ...ClientOption) (Client, error) {
  var cli client

  svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, p2c.Name)
  balancerOpt := WithDialOption(grpc.WithDefaultServiceConfig(svcCfg))
  opts = append([]ClientOption{balancerOpt}, opts...)
  if err := cli.dial(target, opts...); err != nil {
    return nil, err
  }

  return &cli, nil
}
```

在上一篇文章中，我们已经知道当创建gRPC客户端的时候，会触发调用自定义Resolver的Build方法，在Build方法内部获取到服务地址列表后，通过cc.UpdateState方法进行状态更新，后面当监听到服务状态变化的时候同样也会调用cc.UpdateState进行状态的更新，而这里的cc指的就是 **ccResolverWrapper** 对象，这一部分如果忘记的话，可以再去回顾一下讲解Resolver的那篇文章，以便能丝滑接入本篇：

**go-zero/zrpc/resolver/internal/kubebuilder.go:51**

```
if err := cc.UpdateState(resolver.State{
  Addresses: addrs,
}); err != nil {
  logx.Error(err)
}
```

这里有几个重要的模块对象，如下：

- ClientConn：grpc-go/clientconn.go:464
- ccResolverWrapper：grpc-go/resolver_conn_wrapper.go:36
- ccBalancerWrapper：grpc-go/balancer_conn_wrappers.go:48
- Balancer：grpc-go/internal/balancer/gracefulswitch/gracefulswitch.go:46
- balancerWrapper：grpc-go/internal/balancer/gracefulswitch/gracefulswitch.go:247

当监听到服务状态的变更后（首次启动或者通过Watch监听变化）调用 **ccResolverWrapper.UpdateState** 触发更新状态的流程，各模块间的调用链路如下所示：

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg00wMY19Rickys6twWMhF9PqgcyRNoperRcnNic18Br5zRAjibjc7CbIynWcT0NBrpqnVgTtgaYiauOGw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

获取Balancer的动作是在 **ccBalancerWrapper.handleSwitchTo** 方法中触发的，代码如下所示：

**grpc-go/balancer_conn_wrappers.go:266**

```
builder := balancer.Get(name)
if builder == nil {
  channelz.Warningf(logger, ccb.cc.channelzID, "Channel switches to new LB policy %q, since the specified LB policy %q was not registered", PickFirstBalancerName, name)
  builder = newPickfirstBuilder()
} else {
  channelz.Infof(logger, ccb.cc.channelzID, "Channel switches to new LB policy %q", name)
}

if err := ccb.balancer.SwitchTo(builder); err != nil {
  channelz.Errorf(logger, ccb.cc.channelzID, "Channel failed to build new LB policy %q: %v", name, err)
  return
}
ccb.curBalancerName = builder.Name()
```

然后在 **Balancer.SwitchTo** 方法中，调用了自定义Balancer的Build方法：

**grpc-go/internal/balancer/gracefulswitch/gracefulswitch.go:121**

```
newBalancer := builder.Build(bw, gsb.bOpts)
```

上文有提到Build方法的第一个参数为接口 **balancer.ClientConn** ，而这里传入的为 **balancerWrapper** ，所以gracefulswitch.balancerWrapper实现了该接口：

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg00wMY19Rickys6twWMhF9Pqxwnmfniav96ULD9VgyiaQlVV9u9TpRFjasZtPk3Rib4aiaPQx9yicK42LKg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

到这里我们已经知道了获取自定义Balancer是在哪里触达的，以及在哪里获取的自定义的Balancer，和balancer.Builder的Build方法在哪里被调用。

通过上文可知这里的balancer.Builder为baseBuilder，所以调用的Build方法为baseBuilder的Build方法，Build方法的定义如下：

**grpc-go/balancer/base/balancer.go:39**

```
func (bb *baseBuilder) Build(cc balancer.ClientConn, opt balancer.BuildOptions) balancer.Balancer {
  bal := &baseBalancer{
    cc:            cc,
    pickerBuilder: bb.pickerBuilder,

    subConns: resolver.NewAddressMap(),
    scStates: make(map[balancer.SubConn]connectivity.State),
    csEvltr:  &balancer.ConnectivityStateEvaluator{},
    config:   bb.config,
  }
  bal.picker = NewErrPicker(balancer.ErrNoSubConnAvailable)
  return bal
}
```

Build方法返回了baseBalancer，可以知道baseBalancer实现了balancer.Balancer接口：

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg00wMY19Rickys6twWMhF9PqmFAIbcbTwtvmmPr0T1MPRibiaWr4RWo5LT3NGl71mYxEk1NzaSWFwwXQ/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

再来回顾下这个流程，其实主要做了如下几件事：

1. 在自定义的Resolver中监听服务状态的变更
2. 通过UpdateState来更新状态
3. 获取自定义的Balancer
4. 执行自定义Balancer的Build方法获取Balancer

### 如何创建连接？

继续回到ClientConn的updateResolverState方法，在方法的最后调用balancerWrapper.updateClientConnState方法更新客户端的连接状态：

**grpc-go/clientconn.go:664**

```
uccsErr := bw.updateClientConnState(&balancer.ClientConnState{ResolverState: s, BalancerConfig: balCfg})
if ret == nil {
ret = uccsErr // prefer ErrBadResolver state since any other error is
// currently meaningless to the caller.
}
```

后面的调用链路如下图所示：

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg00wMY19Rickys6twWMhF9Pq3N8XP7lzxvpeGvjv2xT2Wus6Gr73M2FGJFRAEMh7XicnbicBvNeFvNKA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

最终会调用baseBalancer.UpdateClientConnState方法：

**grpc-go/balancer/base/balancer.go:94**

```
func (b *baseBalancer) UpdateClientConnState(s balancer.ClientConnState) error {
  // .............
  b.resolverErr = nil
  addrsSet := resolver.NewAddressMap()
  for _, a := range s.ResolverState.Addresses {
    addrsSet.Set(a, nil)
    if _, ok := b.subConns.Get(a); !ok {
      sc, err := b.cc.NewSubConn([]resolver.Address{a}, balancer.NewSubConnOptions{HealthCheckEnabled: b.config.HealthCheck})
      if err != nil {
        logger.Warningf("base.baseBalancer: failed to create new SubConn: %v", err)
        continue
      }
      b.subConns.Set(a, sc)
      b.scStates[sc] = connectivity.Idle
      b.csEvltr.RecordTransition(connectivity.Shutdown, connectivity.Idle)
      sc.Connect()
    }
  }
  for _, a := range b.subConns.Keys() {
    sci, _ := b.subConns.Get(a)
    sc := sci.(balancer.SubConn)
    if _, ok := addrsSet.Get(a); !ok {
      b.cc.RemoveSubConn(sc)
      b.subConns.Delete(a)
    }
  }

  // ................
}
```

当第一次触发调用UpdateClientConnState的时候，如下代码中 ok 为 false：

```
_, ok := b.subConns.Get(a);
```

所以会创建新的连接：

```
sc, err := b.cc.NewSubConn([]resolver.Address{a}, balancer.NewSubConnOptions{HealthCheckEnabled: b.config.HealthCheck})
```

这里的 **b.cc** 即为 **balancerWrapper**，忘记的盆友可以往上翻看复习一下，也就是会调用 balancerWrapper.NewSubConn创建连接

**grpc-go/internal/balancer/gracefulswitch/gracefulswitch.go:328**

```
func (bw *balancerWrapper) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) {
  // .............

  sc, err := bw.gsb.cc.NewSubConn(addrs, opts)
  if err != nil {
    return nil, err
  }
  
  // .............
  
  bw.subconns[sc] = true
  
  // .............
}
```

bw.gsb.cc即为ccBalancerWrapper，所以这里会调用ccBalancerWrapper.NewSubConn创建连接：

**grpc-go/balancer_conn_wrappers.go:299**

```
func (ccb *ccBalancerWrapper) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) {
  if len(addrs) <= 0 {
    return nil, fmt.Errorf("grpc: cannot create SubConn with empty address list")
  }
  ac, err := ccb.cc.newAddrConn(addrs, opts)
  if err != nil {
    channelz.Warningf(logger, ccb.cc.channelzID, "acBalancerWrapper: NewSubConn: failed to newAddrConn: %v", err)
    return nil, err
  }
  acbw := &acBalancerWrapper{ac: ac}
  acbw.ac.mu.Lock()
  ac.acbw = acbw
  acbw.ac.mu.Unlock()
  return acbw, nil
}
```

最终返回的是acBalancerWrapper对象，acBalancerWrapper实现了balancer.SubConn接口：

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg00wMY19Rickys6twWMhF9PqVh15G3DNt4FjstD0NQiarjYrf0pSzPKPicLEA9ciaAiapU8oEsyvohx5Aw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

调用流程图如下所示：

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg00wMY19Rickys6twWMhF9Pq5pOMfO8zwVCr7N3VqXlUkZB9rVvCm79nabTpBTOcIwqE4nl5BWYmTw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

创建连接的默认状态为 **connectivity.Idle** :

**grpc-go/clientconn.go:699**

```
func (cc *ClientConn) newAddrConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (*addrConn, error) {
  ac := &addrConn{
    state:        connectivity.Idle,
    cc:           cc,
    addrs:        addrs,
    scopts:       opts,
    dopts:        cc.dopts,
    czData:       new(channelzData),
    resetBackoff: make(chan struct{}),
  }
 
  // ...........
}
```

在gRPC中为连接定义了五种状态，分别如下：

```
const (
  // Idle indicates the ClientConn is idle.
  Idle State = iota
  // Connecting indicates the ClientConn is connecting.
  Connecting
  // Ready indicates the ClientConn is ready for work.
  Ready
  // TransientFailure indicates the ClientConn has seen a failure but expects to recover.
  TransientFailure
  // Shutdown indicates the ClientConn has started shutting down.
  Shutdown
)
```

在 **baseBalancer ** 中通过b.scStates保存创建的连接，初始状态也为connectivity.Idle，之后通过sc.Connect()进行连接：

**grpc-go/balancer/base/balancer.go:112**

```
b.subConns.Set(a, sc)
b.scStates[sc] = connectivity.Idle
b.csEvltr.RecordTransition(connectivity.Shutdown, connectivity.Idle)
sc.Connect()
```

这里sc.Connetc调用的是acBalancerWrapper的Connect方法，可以看到这里创建连接是异步进行的：

**grpc-go/balancer_conn_wrappers.go:406**

```
func (acbw *acBalancerWrapper) Connect() {
  acbw.mu.Lock()
  defer acbw.mu.Unlock()
  go acbw.ac.connect()
}
```

最后会调用addrConn.connect方法：

**grpc-go/clientconn.go:786**

```
func (ac *addrConn) connect() error {
  ac.mu.Lock()
  if ac.state == connectivity.Shutdown {
    ac.mu.Unlock()
    return errConnClosing
  }
  if ac.state != connectivity.Idle {
    ac.mu.Unlock()
    return nil
  }
  ac.updateConnectivityState(connectivity.Connecting, nil)
  ac.mu.Unlock()

  ac.resetTransport()
  return nil
}
```

从connect开始的调用链路如下所示：

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg00wMY19Rickys6twWMhF9PqYfaOCbEHb00Fwd8YOXqYseRhHaaYpnZVy5YfrajmvBzxFzMyGbibfxg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

在baseBalancer的UpdateSubConnState方法的最后，更新了Picker为自定义的Picker：

**grpc-go/balancer/base/balancer.go:221**

```
b.cc.UpdateState(balancer.State{ConnectivityState: b.state, Picker: b.picker})
```

在addrConn方法的最后会调用ac.resetTransport()真正的进行连接的创建：

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg00wMY19Rickys6twWMhF9Pqhmgyt9tYBRMcibmQrceUj6qPa57YnRwmTzFROiauGzYh8bGl9vVpq3Eg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

当连接已经创建好，处于Ready状态，最后调用baseBalancer.UpdateSubConnState方法，此时s==connectivity.Ready为true，而oldS == connectivity.Ready为false，所以会调用b.regeneratePicker()方法：

```
if (s == connectivity.Ready) != (oldS == connectivity.Ready) ||
    b.state == connectivity.TransientFailure {
    b.regeneratePicker()
}
func (b *baseBalancer) regeneratePicker() {
  if b.state == connectivity.TransientFailure {
    b.picker = NewErrPicker(b.mergeErrors())
    return
  }
  readySCs := make(map[balancer.SubConn]SubConnInfo)

  // Filter out all ready SCs from full subConn map.
  for _, addr := range b.subConns.Keys() {
    sci, _ := b.subConns.Get(addr)
    sc := sci.(balancer.SubConn)
    if st, ok := b.scStates[sc]; ok && st == connectivity.Ready {
      readySCs[sc] = SubConnInfo{Address: addr}
    }
  }
  b.picker = b.pickerBuilder.Build(PickerBuildInfo{ReadySCs: readySCs})
}
```

在regeneratePicker中获取了处于connectivity.Ready状态可用的连接，同时更新了picker。还记得b.pickerBuilder吗？b.b.pickerBuilder为在go-zero中自定义实现的base.PickerBuilder接口。

**go-zero/zrpc/internal/balancer/p2c/p2c.go:42**

```
func (b *p2cPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
  readySCs := info.ReadySCs
  if len(readySCs) == 0 {
    return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
  }

  var conns []*subConn
  for conn, connInfo := range readySCs {
    conns = append(conns, &subConn{
      addr:    connInfo.Address,
      conn:    conn,
      success: initSuccess,
    })
  }

  return &p2cPicker{
    conns: conns,
    r:     rand.New(rand.NewSource(time.Now().UnixNano())),
    stamp: syncx.NewAtomicDuration(),
  }
}
```

最后把自定义的Picker赋值为 ClientConn.blockingpicker.picker属性。

**grpc-go/balancer_conn_wrappers.go:347**

```
func (ccb *ccBalancerWrapper) UpdateState(s balancer.State) {
  ccb.cc.blockingpicker.updatePicker(s.Picker)
  ccb.cc.csMgr.updateState(s.ConnectivityState)
}
```

### 如何选择已创建的连接？

现在已经知道了如何创建连接，以及连接其实是在 **baseBalancer.scStates** 中管理，当连接的状态发生变化，则会更新 **baseBalancer.scStates ** 。那么接下来我们来看一下gRPC是如何选择一个连接进行请求的发送的。

当gRPC客户端发起调用的时候，会调用ClientConn的Invoke方法，一般不会主动使用该方法进行调用，该方法的调用一般是自动生成：

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

如下为发起请求的调用链路，最终会调用p2cPicker.Pick方法获取连接，我们自定义的负载均衡算法一般都在Pick方法中实现，获取到连接之后，通过sendMsg发送请求。

![图片](https://mmbiz.qpic.cn/mmbiz_png/UyIojWicPOg00wMY19Rickys6twWMhF9PqvwLiavmXCz2icEfPTgoyxlL69VOqCRV7aqrN9CfH5bUsYZxfcGW6rcWg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

**grpc-go/stream.go:945**

```
func (a *csAttempt) sendMsg(m interface{}, hdr, payld, data []byte) error {
  cs := a.cs
  if a.trInfo != nil {
    a.mu.Lock()
    if a.trInfo.tr != nil {
      a.trInfo.tr.LazyLog(&payload{sent: true, msg: m}, true)
    }
    a.mu.Unlock()
  }
  if err := a.t.Write(a.s, hdr, payld, &transport.Options{Last: !cs.desc.ClientStreams}); err != nil {
    if !cs.desc.ClientStreams {
      return nil
    }
    return io.EOF
  }
  if a.statsHandler != nil {
    a.statsHandler.HandleRPC(a.ctx, outPayload(true, m, data, payld, time.Now()))
  }
  if channelz.IsOn() {
    a.t.IncrMsgSent()
  }
  return nil
}
```

源码分析到此就结束了，由于篇幅有限没法做到面面俱到，所以本文只列出了源码中的主要路径。

## **结束语**

Balancer相关的源码还是有点复杂的，笔者也是读了好几遍才理清脉络，所以如果读了一两遍感觉没有头绪也不用着急，对照文章的脉络多读几遍就一定能搞懂。

如果有疑问可以随时找我讨论，在社区群中可以搜索dawn_zhou找到我。

希望本篇文章对你有所帮助，你的点赞是作者持续输出的最大动力。