# 使用go test框架驱动的自动化测试

Original 白明的赞赏账户 [TonyBai](javascript:void(0);) *2023-03-31 08:01* *Posted on 辽宁*

收录于合集

\#自动化测试1个

\#go语言102个

\#golang176个

\#测试3个

\#subtest2个

## 一. 背景

团队的测试人员稀缺，无奈只能“自己动手，丰衣足食”，针对我们开发的系统进行自动化测试，这样**既节省的人力，又提高了效率，还增强了对系统质量保证的信心**。

我们的目标是让自动化测试覆盖三个环境，如下图所示：

![Image](https://mmbiz.qpic.cn/mmbiz_png/cH6WzfQ94mZZ5254h5hzDggXia1JlId11rgRkMMq43HiapHLFcmMO60tjFXXiaKWxhgyicKaVbHLFiaBmI2LdojUX7g/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

我们看到这三个环境分别是：

- CI/CD流水线上的自动化测试
- 发版后在各个stage环境中的自动化冒烟/**验收测试**[1]
- 发版后在生产环境的自动化冒烟/验收测试

我们会建立统一的用例库或针对不同环境建立不同用例库，但这些都不重要，重要的是我们**用什么语言来编写这些用例、用什么工具来驱动这些用例**。

下面看看方案的诞生过程。

## 二. 方案

最初组内童鞋使用了**YAML文件**[2]来描述测试用例，并用Go编写了一个独立的工具读取这些用例并执行。这个工具运作起来也很正常。但这样的方案存在一些问题：

- 编写复杂

编写一个最简单的connect连接成功的用例，我们要配置近80行yaml。一个稍微复杂的测试场景，则要150行左右的配置。

- 难于扩展

由于最初的YAML结构设计不足，缺少了扩展性，使得扩展用例时，只能重新建立一个用例文件。

- 表达能力不足

我们的系统是消息网关，有些用例会依赖一定的时序，但基于YAML编写的用例无法清晰地表达出这种用例。

- 可维护性差

如果换一个人来编写新用例或维护用例，这个人不仅要看明白一个个百十来行的用例描述，还要翻看一下驱动执行用例的工具，看看其执行逻辑。很难快速cover这个工具。

为此我们想重新设计一个工具，测试开发人员可以利用该工具支持的**外部DSL文法**[3]来编写用例，然后该工具读取这些用例并执行。

> 注：根据Martin Fowler的**《领域特定语言》**[4]一书对DSL的分类，DSL有三种选型：通用配置文件(xml, json, yaml, toml)、自定义领域语言，这两个合起来称为外部DSL。如：正则表达式、awk, sql、xml等。利用通用编程语言片段/子集作为DSL则称为内部dsl，像ruby等。

后来基于待测试的场景数量和用例复杂度粗略评估了一下DSL文法(甚至借助ChatGPT生成过几版DSL文法)，发现这个“小语言”那也是“麻雀虽小五脏俱全”。如果用这样的DSL编写用例，和利用通用语言(比如Python)编写的用例在代码量级上估计也不相上下了。

既然如此，自己设计外部DSL意义也就不大了。还不如用Python来整。但转念一想，既然用通用语言的子集了，团队成员对Python又不甚熟悉，那为啥不回到Go呢^_^。

让我们进行一个大胆的设定：将Go testing框架作为“内部DSL”来编写用例，用go test命令作为执行这些用例的测试驱动工具。此外，有了GPT-4加持，生成TestXxx、补充用例啥的应该也不是大问题。

下面我们来看看如何组织和编写用例并使用go test驱动进行自动化测试。

## 三. 实现

### 1. 测试用例组织

我的**《Go语言精进之路vol2》**[5]书中的**第41条“有层次地组织测试代码”**[6]中对基于go test的测试用例组织做过系统的论述。结合Go test提供的**TestMain**[7]、TestXxx与**sub test**[8]，我们完全可以基于go test建立起一个层次清晰的测试用例结构。这里就以一个对开源mqtt broker的自动化测试为例来说明一下。

> 注：你可以在本地搭建一个单机版的开源mqtt broker服务作为被测对象，比如使用**Eclipse的mosquitto**[9]。

在组织用例之前，我先问了一下ChatGPT对一个mqtt broker测试都应该包含哪些方面的用例，ChatGPT给了我一个简单的表：

![Image](https://mmbiz.qpic.cn/mmbiz_png/cH6WzfQ94mZZ5254h5hzDggXia1JlId113picwuCcygSnyJnYFNeQqHO4Zibv74DELj9ZAkLcyibwX6MgUhBYXXU3g/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

如果你对**MQTT协议**[10]有所了解，那么你应该觉得ChatGPT给出的答案还是很不错的。

这里我们就以connection、subscribe和publish三个场景(scenario)来组织用例：

```
$tree -F .
.
├── Makefile
├── go.mod
├── go.sum
├── scenarios/
│   ├── connection/              // 场景：connection
│   │   ├── connect_test.go      // test suites
│   │   └── scenario_test.go
│   ├── publish/                 // 场景：publish
│   │   ├── publish_test.go      // test suites
│   │   └── scenario_test.go
│   ├── scenarios.go             // 场景中测试所需的一些公共函数
│   └── subscribe/               // 场景：subscribe
│       ├── scenario_test.go     
│       └── subscribe_test.go    // test suites
└── test_report.html             // 生成的默认测试报告
```

简单说明一下这个测试用例组织布局：

- 我们将测试用例分为多个场景(scenario)，这里包括connection、subscribe和publish；
- 由于是由go test驱动，所以每个存放test源文件的目录中都要遵循Go对Test的要求，比如：源文件以_test.go结尾等。
- 每个场景目录下存放着测试用例文件，一个场景可以有多个_test.go文件。这里设定_test.go文件中的每个TestXxx为一个test suite，而TestXxx中再基于subtest编写用例，这里每个subtest case为一个最小的test case；
- 每个场景目录下的scenario_test.go，都是这个目录下包的TestMain入口，主要是考虑为所有包传入统一的命令行标志与参数值，同时你也针对该场景设置在TestMain中设置setup和teardown。该文件的典型代码如下：

```
// github.com/bigwhite/experiments/automated-testing/scenarios/subscribe/scenario_test.go

package subscribe
  
import (
    "flag"
    "log"
    "os"
    "testing"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

var addr string

func init() {
    flag.StringVar(&addr, "addr", "", "the broker address(ip:port)")
}

func TestMain(m *testing.M) {
    flag.Parse()

    // setup for this scenario
    mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)

    // run this scenario test
    r := m.Run()

    // teardown for this scenario
    // tbd if teardown is needed

    os.Exit(r)
}
```

接下来我们再来看看具体测试case的实现。

### 2. 测试用例实现

我们以稍复杂一些的subscribe场景的测试为例，我们看一下subscribe目录下的subscribe_test.go中的测试suite和cases：

```
// github.com/bigwhite/experiments/automated-testing/scenarios/subscribe/subscribe_test.go

package subscribe

import (
 scenarios "bigwhite/autotester/scenarios"
 "testing"
)

func Test_Subscribe_S0001_SubscribeOK(t *testing.T) {
 t.Parallel() // indicate the case can be ran in parallel mode

 tests := []struct {
  name  string
  topic string
  qos   byte
 }{
  {
   name:  "Case_001: Subscribe with QoS 0",
   topic: "a/b/c",
   qos:   0,
  },
  {
   name:  "Case_002: Subscribe with QoS 1",
   topic: "a/b/c",
   qos:   1,
  },
  {
   name:  "Case_003: Subscribe with QoS 2",
   topic: "a/b/c",
   qos:   2,
  },
 }

 for _, tt := range tests {
  tt := tt
  t.Run(tt.name, func(t *testing.T) {
   t.Parallel() // indicate the case can be ran in parallel mode
   client, testCaseTeardown, err := scenarios.TestCaseSetup(addr, nil)
   if err != nil {
    t.Errorf("want ok, got %v", err)
    return
   }
   defer testCaseTeardown()

   token := client.Subscribe(tt.topic, tt.qos, nil)
   token.Wait()

   // Check if subscription was successful
   if token.Error() != nil {
    t.Errorf("want ok, got %v", token.Error())
   }

   token = client.Unsubscribe(tt.topic)
   token.Wait()
   if token.Error() != nil {
    t.Errorf("want ok, got %v", token.Error())
   }
  })
 }
}

func Test_Subscribe_S0002_SubscribeFail(t *testing.T) {
}
```

这个测试文件中的测试用例与我们日常编写单测并没有什么区别！有一些需要注意的地方是：

- Test函数命名

这里使用了Test_Subscribe_S0001_SubscribeOK、Test_Subscribe_S0002_SubscribeFail命名两个Test suite。命名格式为：

```
Test_场景_suite编号_测试内容缩略
```

之所以这么命名，一来是测试用例组织的需要，二来也是为了后续在生成的Test report中区分不同用例使用。

- testcase通过subtest呈现

每个TestXxx是一个test suite，而基于表驱动的每个sub test则对应一个test case。

- test suite和test case都可单独标识为是否可并行执行

通过testing.T的Parallel方法可以标识某个TestXxx或test case(subtest)是否是可以并行执行的。

- 针对每个test case，我们都调用setup和teardown

这样可以保证test case间都相互独立，互不影响。

### 3. 测试执行与报告生成

设计完布局，编写完用例后，接下来就是执行这些用例。那么怎么执行这些用例呢？

前面说过，我们的方案是基于go test驱动的，我们的执行也要使用go test。

在顶层目录automated-testing下，执行如下命令：

```
$go test ./... -addr localhost:30083 
```

go test会遍历执行automated-testing下面每个包的测试，在执行每个包的测试时会将-addr这个flag传入。如果localhost:30083端口并没有mqtt broker服务监听，那么上面的命令将输出如下信息：

```
$go test ./... -addr localhost:30083
?    bigwhite/autotester/scenarios [no test files]
[ERROR] [client]   dial tcp [::1]:30083: connect: connection refused
[ERROR] [client]   Failed to connect to a broker
--- FAIL: Test_Connection_S0001_ConnectOKWithoutAuth (0.00s)
    connect_test.go:20: want ok, got network Error : dial tcp [::1]:30083: connect: connection refused
FAIL
FAIL bigwhite/autotester/scenarios/connection 0.015s
[ERROR] [client]   dial tcp [::1]:30083: connect: connection refused
[ERROR] [client]   Failed to connect to a broker
--- FAIL: Test_Publish_S0001_PublishOK (0.00s)
    publish_test.go:11: want ok, got network Error : dial tcp [::1]:30083: connect: connection refused
FAIL
FAIL bigwhite/autotester/scenarios/publish 0.016s
[ERROR] [client]   dial tcp [::1]:30083: connect: connection refused
[ERROR] [client]   dial tcp [::1]:30083: connect: connection refused
[ERROR] [client]   Failed to connect to a broker
[ERROR] [client]   Failed to connect to a broker
[ERROR] [client]   dial tcp [::1]:30083: connect: connection refused
[ERROR] [client]   Failed to connect to a broker
--- FAIL: Test_Subscribe_S0001_SubscribeOK (0.00s)
    --- FAIL: Test_Subscribe_S0001_SubscribeOK/Case_002:_Subscribe_with_QoS_1 (0.00s)
        subscribe_test.go:39: want ok, got network Error : dial tcp [::1]:30083: connect: connection refused
    --- FAIL: Test_Subscribe_S0001_SubscribeOK/Case_003:_Subscribe_with_QoS_2 (0.00s)
        subscribe_test.go:39: want ok, got network Error : dial tcp [::1]:30083: connect: connection refused
    --- FAIL: Test_Subscribe_S0001_SubscribeOK/Case_001:_Subscribe_with_QoS_0 (0.00s)
        subscribe_test.go:39: want ok, got network Error : dial tcp [::1]:30083: connect: connection refused
FAIL
FAIL bigwhite/autotester/scenarios/subscribe 0.016s
FAIL
```

这也是一种测试失败的情况。

在自动化测试时，我们一般会把错误或成功的信息保存到一个测试报告文件(多是html)中，那么我们如何基于上面的测试结果内容生成我们的测试报告文件呢？

首先go test支持将输出结果以结构化的形式展现，即传入-json这个flag。这样我们仅需基于这些json输出将各个字段读出并写入html中即可。好在有现成的开源工具可以做到这点，那就是**go-test-report**[11]。下面是通过命令行管道让go test与go-test-report配合工作生成测试报告的命令行：

> 注：go-test-report工具的安装方法：go install github.com/vakenbolt/go-test-report@latest

```
$go test ./... -addr localhost:30083 -json|go-test-report
[go-test-report] finished in 1.375540542s
```

执行结束后，就会在当前目录下生成一个test_report.html文件，使用浏览器打开该文件就能看到测试执行结果：

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

通过测试报告的输出，我们可以很清楚看到哪些用例通过，哪些用例失败了。并且通过Test suite的名字或Test case的名字可以快速定位是哪个scenario下的哪个suite的哪个case报的错误！我们也可以点击某个test suite的名字，比如：Test_Connection_S0001_ConnectOKWithoutAuth，打开错误详情查看错误对应的源文件与具体的行号：

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

为了方便快速敲入上述命令，我们可以将其放入Makefile中方便输入执行，即在顶层目录下，执行make即可执行测试：

```
$make
go test ./... -addr localhost:30083 -json|go-test-report
[go-test-report] finished in 2.011443636s
```

如果要传入自定义的mqtt broker的服务地址，可以用：

```
$make broker_addr=192.168.10.10:10083
```

## 四. 小结

在这篇文章中，我们介绍了如何实现基于go test驱动的自动化测试，介绍了这样的测试的结构布局、用例编写方法、执行与报告生成等。

这个方案的不足是**要求测试用例所在环境需要部署go与go-test-report**。

go test支持将test编译为一个可执行文件，不过不支持将多个包的测试编译为一个可执行文件：

```
$go test -c ./...
cannot use -c flag with multiple packages
```

此外由于go test编译出的可执行文件**不支持将输出内容转换为JSON格式**[12]，因此也无法对接go-test-report将测试结果保存在文件中供后续查看。

本文涉及的源码可以在**这里**[13]下载 - https://github.com/bigwhite/experiments/tree/master/automated-testing

------

**“Gopher部落”知识星球**[14]旨在打造一个精品Go学习和进阶社群！高品质首发Go技术文章，“三天”首发阅读权，每年两期Go语言发展现状分析，每天提前1小时阅读到新鲜的Gopher日报，网课、技术专栏、图书内容前瞻，六小时内必答保证等满足你关于Go语言生态的所有需求！2023年，Gopher部落将进一步聚焦于如何编写雅、地道、可读、可测试的Go代码，关注代码质量并深入理解Go核心技术，并继续加强与星友的互动。欢迎大家加入！

![Image](https://mmbiz.qpic.cn/mmbiz_png/cH6WzfQ94mYbN4SR0aJeoKt82pr7ibmCk1icF8xqVslY1JfrDvW4fJKB5RIWtClXGPn5Y0qsJvSibnQd6Bb9EtYWg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

![Image](https://mmbiz.qpic.cn/mmbiz_png/cH6WzfQ94mb54jsFJZ3Knmz8obUsf3PBShthmdSw5E01TcYmUReGkj0BWpxHak1HlnlzHvLmKax53YSGr7aNlA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

![Image](https://mmbiz.qpic.cn/mmbiz_png/cH6WzfQ94mYKSeNd014VMtNhYulia0OHrHVoyrVYb2JvBa5ycFaeDfscQdubicnZkxB6je42bo3J4cZcx0FticLmQ/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

![Image](https://mmbiz.qpic.cn/mmbiz_jpg/cH6WzfQ94mb54jsFJZ3Knmz8obUsf3PBDKyzaL44T9g1YiaYeujWa3QRrVC21SnO9h9qc2ia6ibyicc6LUdnD0ibymw/640?wx_fmt=jpeg&wxfrom=5&wx_lazy=1&wx_co=1)



Gopher Daily(Gopher每日新闻)归档仓库 - https://github.com/bigwhite/gopherdaily

我的联系方式：

- 微博(暂不可用)：https://weibo.com/bigwhite20xx
- 微博2：https://weibo.com/u/6484441286
- 博客：tonybai.com
- github: https://github.com/bigwhite

![Image](https://mmbiz.qpic.cn/mmbiz_png/cH6WzfQ94mb54jsFJZ3Knmz8obUsf3PBrSoqeMvoWCticN2cpU64fJ0FYQdXJhP7ia7WRh8628uOAsQYeE2NibRRw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

商务合作方式：撰稿、出书、培训、在线课程、合伙创业、咨询、广告合作。

### 参考资料

[1] 验收测试: *http://en.wikipedia.org/wiki/Acceptance_testing*

[2] YAML文件: *https://tonybai.com/2019/02/25/introduction-to-yaml-creating-a-kubernetes-deployment/*

[3] 外部DSL文法: *https://tonybai.com/2022/05/10/introduction-of-implement-dsl-using-antlr-and-go*

[4] 《领域特定语言》: *https://book.douban.com/subject/21964984/*

[5] 《Go语言精进之路vol2》: *https://item.jd.com/13694000.html*

[6] 第41条“有层次地组织测试代码”: *https://book.douban.com/subject/35720729/*

[7] TestMain: *https://pkg.go.dev/testing#Main*

[8] sub test: *https://tonybai.com/2023/03/15/an-intro-of-go-subtest/*

[9] Eclipse的mosquitto: *https://github.com/eclipse/mosquitto*

[10] MQTT协议: *https://mqtt.org/mqtt-specification/*

[11] go-test-report: *https://github.com/vakenbolt/go-test-report*

[12] 不支持将输出内容转换为JSON格式: *https://github.com/golang/go/issues/22996*

[13] 这里: *https://github.com/bigwhite/experiments/tree/master/automated-testing*

[14] “Gopher部落”知识星球: *https://wx.zsxq.com/dweb2/index/group/51284458844544*