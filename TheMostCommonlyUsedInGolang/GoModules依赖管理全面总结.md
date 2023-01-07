

# [Go Modules 依赖管理，这篇总结的挺全](https://mp.weixin.qq.com/s/7HGJaHaBjStKuVxecNi-QA)

[网管叨bi叨](javascript:void(0);) *2023-01-06 08:45* *Posted on 北京

## 前言

> 哈喽，大家好啊。一般编程语言都会提供依赖库管理工具，例如python的pip、node.js的npm，java的maven，rust的cargo，Go语言也有提供自己的依赖库管理工具，Go语言在1.11提出了Go mod，每次版本或多或少都会对go.mod进行改进优化，go mod也越来越好，当前大多数公司都使用go mod来管理依赖库，所以本文我们一起来入门go mod（参考资料文末给出）；

## Go Modules发展史

### go get阶段

起初Go语言在1.5之前没有依赖管理工具，若想引入依赖库，需要执行go get命令将代码拉取放入GOPATH/src目录下，作为GOPATH下的全局依赖，这也就意味着没有版本控制及隔离项目的包依赖；

### vendor阶段

为了解决隔离项目的包依赖问题，Go1.5版本推出了vendor机制，环境变量中有一个GO15VENDOREXPERIMENT需要设置为1，该环境变量在Go1.6版本时变成默认开启，目前已经退出了历史舞台；

vendor其实就是将原来放在GOPATH/src的依赖包放到工程的vendor目录中进行管理，不同工程独立地管理自己的依赖包，相互之间互不影响，原来是包共享的模式，通过vendor这种机制进行隔离，在项目编译的时候会先去vendor目录查找依赖，如果没有找到才会再去GOPATH目录下查找；

优点：保证了功能项目的完整性，减少了下载依赖包，直接使用vendor就可以编译

缺点：仍然没有解决版本控制问题，go get仍然是拉取最新版本代码；

### 社区管理工具

很多优秀的开发者在这期间也都实现了不错的包依赖管理工具，例如：

godep：https://github.com/tools/godep

govendor：https://github.com/kardianos/govendor

glide：https://github.com/Masterminds/glide

dep：https://github.com/golang/dep

dep应该是其中最成功的，得到了Go语言官方的支持，该项目也被放到了https://github.com/golang/dep，但是为什么dep没有称为官宣的依赖工具呢？

其实因为随着Russ Cox 与 Go 团队中的其他成员不断深入地讨论，发现 dep 的一些细节似乎越来越不适合 Go，因此官方采取了另起 proposal 的方式来推进，其方案的结果一开始先是释出 vgo，最终演变为我们现在所见到的 Go modules;

### go modules

go modules是Russ Cox推出来的，发布于Go1.11，成长于Go1.12，丰富于Go1.13，正式于Go1.14推荐在生产上使用，几乎后续的每个版本都或多或少的有一些优化，在Go1.16引入go mod retract、在Go1.18引入go work工作区的概念，这些我们在本文都会介绍到；

## Go Modules知识点

### GO111MODULE环境变量

这个环境变量是Go Modules的开关，主要有以下参数：

- auto：只在项目包含了go.mod文件时启动go modules，在Go1.13版本中是默认值
- on：无脑启动Go Modules，推荐设置，Go1.14版本以后的默认值
- off：禁用Go Modules，一般没有使用go modules的工程使用；

我现在使用的Go版本是1.19.3，默认GO111MODULE=on，感觉该变量也会像GO15VENDOREXPERIMENT最终推出系统环境变量的舞台；

### GOPROXY

该环境变量用于设置Go模块代理，Go后续在拉取模块版本时能够脱离传统的VCS方式从镜像站点快速拉取，GOPROXY的值要以英文逗号分割，默认值是`https://proxy.golang.org,direct`，但是该地址在国内无法访问，所以可以使用`goproxy.cn`来代替(七牛云配置)，设置命令：

```
go env -w GOPROXY=GOPROXY=https://goproxy.cn,direct
```

也可以使用其他配置，例如阿里配置：

```
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/
```

该环境变量也可以关闭，可以设置为"off"，禁止Go在后续操作中使用任何Go module proxy；

上面的配置中我们用逗号分割后面的值是`direct`，它是什么意思呢？

`direct`为特殊指示符，因为我们指定了镜像地址，默认是从镜像站点拉取，但是有些库可能不存在镜像站点中，`direct`可以指示Go回源到模块版本的源地址去抓取，比如github，当go module proxy返回404、410这类错误时，其会自动尝试列表中的下一个，遇见`direct`时回源地址抓取；

### GOSUMDB

该环境变量的值是一个Go checksum database，用于保证Go在拉取模块版本时拉取到的模块版本数据未经篡改，若发现不一致会中止，也可以将值设置为`off`即可以禁止Go在后续操作中校验模块版本；

什么是Go checksum database?

Go checksum database主要用于保护Go不会从任何拉到被篡改过的非法Go模块版本，详细算法机制可以看一下：https://go.googlesource.com/proposal/+/master/design/25530-sumdb.md#proxying-a-checksum-database

GOSUMDB的默认值是`sum.golang.org`，默认值与自定义值的格式不一样，默认值在国内是无法访问，这个值我们一般不用动，因为我们一般已经设置好了GOPROXY，goproxy.cn支持代理`sum.golang.org`;

GOSUMDB的值自定义格式如下：

- 格式 1：`<SUMDB_NAME>+<PUBLIC_KEY>`。
- 格式 2：`<SUMDB_NAME>+<PUBLIC_KEY> <SUMDB_URL>`。

### GONOPROXY/GONOSUMDB/GOPRIVATE

这三个环境变量放在一起说，一般在项目中不经常使用，这三个环境变量主要用于私有模块的拉取，在GOPROXY、GOSUMDB中无法访问到模块的场景中，例如拉取git上的私有仓库；

GONOPROXY、GONOSUMDB的默认值是GOPRIVATE的值，所以我们一般直接使用GOPRIVATE即可，其值也是可以设置多个，以英文逗号进行分割；例如：

```
$ go env -w GOPRIVATE="github.com/asong2020/go-localcache,git.xxxx.com"
```

也可以使用通配符的方式进行设置，对域名设置通配符号，这样子域名就都不经过Go module proxy和Go checksum database;

### 全局缓存

go mod download会将依赖缓存到本地，缓存的目录是`GOPATH/pkg/mod/cache`、`GOPATH/pkg/sum`，这些缓存依赖可以被多个项目使用，未来可能会迁移到`$GOCACHE`下面；

可以使用`go clean -modcache`清理所有已缓存的模块版本数据；

### Go Modules命令

我们可以使用`go help mod`查看可以使用的命令：

```
go help mod
Go mod provides access to operations on modules.

Note that support for modules is built into all the go commands,
not just 'go mod'. For example, day-to-day adding, removing, upgrading,
and downgrading of dependencies should be done using 'go get'.
See 'go help modules' for an overview of module functionality.

Usage:

        go mod <command> [arguments]

The commands are:

        download    download modules to local cache
        edit        edit go.mod from tools or scripts
        graph       print module requirement graph
        init        initialize new module in current directory
        tidy        add missing and remove unused modules
        vendor      make vendored copy of dependencies
        verify      verify dependencies have expected content
        why         explain why packages or modules are needed

Use "go help mod <command>" for more information about a command.
```

| 命令            | 作用                                          |
| :-------------- | :-------------------------------------------- |
| go mod init     | 生成go.mod文件                                |
| go mod download | 下载go.mod文件中指明的所有依赖放到全局缓存    |
| go mod tidy     | 整理现有的依赖，添加缺失或移除不使用的modules |
| go mod graph    | 查看现有的依赖结构                            |
| go mod edit     | 编辑go.mod文件                                |
| go mod vendor   | 导出项目所有的依赖到vendor目录                |
| go mod verify   | 校验一个模块是否被篡改过                      |
| go mod why      | 解释为什么需要依赖某个模块                    |

### go.mod文件

go.mod是启用Go modules的项目所必须且最重要的文件，其描述了当前项目的元信息，每个go.mod文件开头符合包含如下信息：

**module**：用于定义当前项目的模块路径（突破$GOPATH路径）

**go**：当前项目Go版本，目前只是标识作用

**require**：用设置一个特定的模块版本

**exclude**：用于从使用中排除一个特定的模块版本

**replace**：用于将一个模块版本替换为另外一个模块版本，例如chromedp使用golang.org/x/image这个package一般直连是获取不了的，但是它有一个github.com/golang/image的镜像，所以我们要用replace来用镜像替换它

**restract**：用来声明该第三方模块的某些发行版本不能被其他模块使用，在Go1.16引入

例子：

image-20230102181623314

接下来我们分模块详细介绍一下各部分；

#### module path

go.mod文件的第一行是module path，采用仓库+module name的方式定义，例如上面的项目：

```
module github.com/asong2020/go-localcache
```

因为Go module遵循语义化版本规范2.0.0，所以如果工程的版本已经大于2.0.0，按照规范需要加上major的后缀，module path改成如下：

```
module github.com/asong2020/go-localcache/v2
module github.com/asong2020/go-localcache/v3
......
```

#### go version

go.mod文件的第二行是go version，其是用来指定你的代码所需要的最低版本：

```
go 1.19.3
```

其实这一行不是必须的，目前也只是标识作用，可以不写；

#### require

require用来指定该项目所需要的各个依赖库以及他们的版本，从上面的例子中我们看到版本部分有不同的写法，还有注释，接下来我们来解释一下这部分；

##### indirect注释

```
github.com/davecgh/go-spew v1.1.0 // indirect
```

以下场景才会添加indirect注释：

- 当前项目依赖包A，但是A依赖包B，但是A的go.mod文件中缺失B，所以在当前项目go.mod中补充B并添加indirect注释
- 当前项目依赖包A，但是依赖包A没有go.mod文件，所以在当前项目go.mod中补充B并添加indirect注释
- 当前项目依赖包A，依赖包A又依赖包B，当依赖包A降级不在依赖B时，这个时候就会标记indirect注释，可以执行go mod tidy移除该依赖；

Go1.17版本对此做了优化，indirect 的 module 将被放在单独 require 块的，这样看起来更加清晰明了。

##### incompatible标记

我们在项目中会看到有一些库后面添加了incompatible标记：

```
github.com/dgrijalva/jwt-go v3.2.0+incompatible
```

jwt-go这个库就是这样的，这是因为jwt-go的版本已经大于2了，但是他们的module path仍然没有添加v2、v3这样的后缀，不符合Go的module管理规范，所以go module把他们标记为incompatible，不影响引用；

##### 版本号

go module拉取依赖包本质也是go get行为，go get主要提供了以下命令：

| 命令               | 作用                                                         |
| :----------------- | :----------------------------------------------------------- |
| go get             | 拉取依赖，会进行指定性拉取（更新），并不会更新所依赖的其它模块。 |
| go get -u          | 更新现有的依赖，会强制更新它所依赖的其它全部模块，不包括自身。 |
| go get -u -t ./... | 更新所有直接依赖和间接依赖的模块版本，包括单元测试中用到的。 |

go get拉取依赖包取决于依赖包是否有发布的tags：

1. 拉取的依赖包没有发布tags

2. - 默认取主分支最近一次的commit的commit hash，生成一个伪版本号

3. 拉取的依赖包有发布tags

4. - 如果只有单个模块，那么就取主版本号最大的那个tag
   - 如果有多个模块，则推算相应的模块路径，取主版本号最大的那个tag

没有发布的tags：

```
github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
```

**v0.0.0**：根据commit的base version生成的：

- 如果没有base version，那么就是vx.0.0的形式
- 如果base version是一个预发版本，那么就是vx.y.z-pre.0的形式
- 如果base version是一个正式发布的版本，那么它就patch号加1，就是vx.y.(z+1)-0的形式

**20190718012654**：是这次提交的时间，格式是`yyyyMMddhhmmss`

**fb15b899a751**：是这个版本的commit id，通过这个可以确定这个库的特定的版本

```
github.com/beego/bee v1.12.0
```

#### replace

replace用于解决一些错误的依赖库的引用或者调试依赖库；

场景举例：

举例1：

日常开发离不开第三方库，大部分场景都可以满足我们的需要，但是有些时候我们需要对依赖库做一些定制修改，依赖库修改后，我们想引起最小的改动，就可以使用replace命令进行重新引用，调试也可以使用replace进行替换，Go1.18引入了工作区的概念，调试可以使用work进行代替，后面会介绍；

举例2：

golang.org/x/crypto库一般我们下载不下来，可以使用replace引用到github.com/golang/crypto：

```
go mod edit -replace golang.org/x/crypto=github.com/golang/crypto@v0.0.0-20160511215533-1f3b11f56072
```

#### exclude

用于跳过某个依赖库的版本，使用场景一般是我们知道某个版本有bug或者不兼容，为了安全起可以使用exclude跳过此版本；

```
exclude (
 go.etcd.io/etcd/client/v2 v2.305.0-rc.0
)
```

#### retract

这个特性是在Go1.16版本中引入，用来声明该第三方模块的某些发行版本不能被其他模块使用；

使用场景：发生严重问题或者无意发布某些版本后，模块的维护者可以撤回该版本，支持撤回单个或多个版本；

这种场景以前的解决办法：

维护者删除有问题版本的tag，重新打一个新版本的tag；

使用者发现有问题的版本tag丢失，手动介入升级，并且不明真因；

引入retract后，维护者可以使用retract在go.mod中添加有问题的版本：

```
// 严重bug...
retract (
  v0.1.0
  v0.2.0
)
```

重新发布新版本后，在引用该依赖库的使用执行go list可以看到 版本和"严重bug..."的提醒；

该特性的主要目的是将问题更直观的反馈到开发者的手中；

### go.sum文件

go.sun文件也是在go mod init阶段创建，go.sum的介绍文档偏少，我们一般也很少关注go.sum文件，go.sum主要是记录了所有依赖的module的校验信息，内容如下：

image-20230102193717816

从上面我们可以看到主要是有两种形式：

- h1:
- /go.mod h1:

其中module是依赖的路径，version是依赖的版本号。hash是以`h1:`开头的字符串，hash 是 Go modules 将目标模块版本的 zip 文件开包后，针对所有包内文件依次进行 hash，然后再把它们的 hash 结果按照固定格式和算法组成总的 hash 值。

h1 hash 和 go.mod hash两者要不同时存在，要不就是只存在go.mod hash，当Go认为肯定用不到某个版本的时候就会省略它的h1 hash，就只有go.mod hash；

## Go Modules在项目中使用

使用go modules的一个前置条件是Go语言版本大于等于Go1.11；

然后我们要检查环境变量GO111MODULE是否开启，执行`go env`查看：

```
$ go env
GO111MODULE=off
```

执行如下命令打开go mod：

```
$ go env -w GO111MODULE=on
```

接下来我们随意创建一个项目：

```
$ mkdir -p asong/demo
$ cd asong/demo
```

执行go mod init初始化该项目：

```
$ go mod init github.com/asong/demo
go: creating new go.mod: module github.com/asong/demo
```

接下来我们在demo目录下创建main.go文件，写下如下代码：

```
package main

import (
 "fmt"
 cache "github.com/asong2020/go-localcache"
)

func main() {
 c, err := cache.NewCache()
 if err != nil {
  return
 }
 key := "asong"
 value := []byte("公众号：Golang梦工厂")
 err = c.Set(key, value)
 if err != nil {
  return
 }
 entry, err := c.Get(key)
 if err != nil {
  return
 }
 fmt.Printf("get value is %s\n", string(entry))

 err = c.Delete(key)
 if err != nil {
  return
 }
}
```

然后执行**go mod tidy**命令：

image-20230102202521167

自动根据main.go文件更新依赖，我们再看一下go.mod文件：

image-20230102202600643

以上就是在项目对go.mod的简单使用；

## go1.18新特性：工作区

工作区用来解决什么问题？

场景1：我们有时在本地会对一些三方依赖库进行特制修改，然后想在项目修改依赖库引用到本地进行调试，这时我们可以使用replace做替换，这样就可以在本地进行开发联调，这样虽然可以解决问题，但是会存在问题，因为是在项目的go.mod文件直接修改的，如果误传到远程仓库，会影响到其他开发同学；

场景2：我们在本地开发了一些依赖库，这时想在本地测试一下，还未发到远程仓库，那么我们在其他项目中引入该依赖库后，执行`go mod tidy`就会报远程库没有找到的问题，所以就必须要把依赖库先推送到远程，在引用调试；

正是这些问题，Go语言在Go1.18正式增加了`go work`工作区的概念，其实就是将N个Go Module组成一个Go Work，工作区的读取优先级是最高的，执行`go help work`可以查看go work提供的功能：

```
$ go help work
Usage:

        go work <command> [arguments]

The commands are:

        edit        edit go.work from tools or scripts
        init        initialize workspace file
        sync        sync workspace build list to modules
        use         add modules to workspace file

Use "go help work <command>" for more information about a command.
```

执行`go work init`命令初始化一个新的工作区，在项目中生成一个go.work文件：

```
go 1.18
// 多模块添加
use (...)

replace XXXXX => XXXX v1.4.5
```

go.work文件与go.mod文件语法一致，go.work支持三个指令：

- go：声明go版本号
- use：声明应用所依赖模块的具体文件路径，路径可以是绝对路径或相对路径，即使路径是当前应用目录外也可
- replace：声明替换某个模块依赖的导入路径，优先级高于 go.mod 中的 replace 指令；

所以针对上述场景，我们使用go work init命令在项目中对本地依赖库进行关联即可解决，后续我们只需要在git配置文件中添加go.work文件不推送到远程即可；

我们也可以在编译时通过`-workfile=off`指令禁用工作区模式：

```
$ go build -workfile=offf .
```

go.work的推出主要是用于在本地调试，不会因为修改go.mod引入问题；

### 参考文献

- [Go1.18 新特性：多 Module 工作区模式](https://mp.weixin.qq.com/s?__biz=MzUxMDI4MDc1NA==&mid=2247493859&idx=1&sn=248e6a4ee3c62e8a32512fbd42f069ea&scene=21#wechat_redirect)
- [Go Modules 终极入门](https://mp.weixin.qq.com/s?__biz=MzUxMDI4MDc1NA==&mid=2247483713&idx=1&sn=817ffef56f8bc5ca09a325c9744e00c7&scene=21#wechat_redirect)
- [Go mod 七宗罪](https://mp.weixin.qq.com/s?__biz=MjM5MDUwNTQwMQ==&mid=2257485369&idx=1&sn=79e1cb77d41411a42a5386a6ce1e5e65&scene=21#wechat_redirect)
- 深入Go Module之go.mod文件解析

## 总结

现在大小公司的项目应该都已经在使用`Go Modules`进行依赖包管理了，虽然`Go Modules`相比于Maven、npm还不是很完善，但也在不断地进行优化，变得越来越好，如果你现在项目还没有使用`go modules`，可以准备将项目迁移到go mod了，推荐你使用；