怎么把自己造的轮子发布到 Go Module上
网管叨bi叨 2023-02-10 08:45 Posted on 北京
The following article is from Golang技术分享 Author 机器铃砍菜刀


我们在开发 Go 项目时，经常会使用到一些外部依赖包。它们一般都是通过形如go get example.com/xxx的命令形式获取到本地使用。

本着开源精神，如果我们想将自己开发的包共享出去，让其他人也能使用go get命令获取到它们，应该如何做呢？

本文将以开源至 Github 平台为例，展示其流程。

建立仓库
首先，在 github 平台建立仓库，设置想要开源的 Go 库名。这里将项目名取为 publishdemo，状态为 Public。

Image

开发模块代码
将创建好的仓库，通过 git clone 命令拉取至本地开发。

$ git clone git@github.com:golangShare/publishdemo.git
此时，在项目根目录publishdemo/下，仅包含了 LICENSE 和 README 文件，不包含任何 Go 代码。

$ ls
LICENSE   README.md
为项目初始化 mod 文件

$ go mod init github.com/golangShare/publishdemo
假设，我们想开源的 Go 工具类的库，此时准备先提供的是对字符串相关的操作。因此在publishdemo/目录下，我们新增stringutil/子目录，在子目录中增加 reverse.go 文件，其内容如下。

package stringutil

// Reverse returns its argument string reversed rune-wise left to right.
func Reverse(s string) string {
 r := []rune(s)
 for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
  r[i], r[j] = r[j], r[i]
 }
 return string(r)
}
很简单，我们想提供的是一个翻转字符串功能函数。

当然，我们也应该提供测试代码。在同级目录增加 reverse_test.go 文件，其内容如下。

package stringutil

import "testing"

func TestReverse(t *testing.T) {
 for _, c := range []struct {
  in, want string
 }{
  {"Hello, world", "dlrow ,olleH"},
  {"Hello, 世界", "界世 ,olleH"},
  {"", ""},
 } {
  got := Reverse(c.in)
  if got != c.want {
   t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
  }
 }
}
回到项目根目录，此时，代码结构如下

.
├── LICENSE
├── README.md
├── go.mod
└── stringutil
    ├── reverse.go
    └── reverse_test.go

1 directory, 5 files
测试代码也都通过

$ go test ./...
ok   github.com/golangShare/publishdemo/stringutil 0.005s
开发完成，我们就可以将工具库共享出去了。

发布
为了避免模块中还记录了不再需要的依赖项，执行 go mod tidy 移除它们。

$ go mod tidy
使用 git tag 命令标记版本信息

$ git commit -m "add Reverse: for v0.1.0"
$ git tag v0.1.0
将其 push 至远程仓库

$ git push origin v0.1.0
使用
发布之后，其他项目就可以通过以下命令获取我们开源的 Go 包了。

$ go get github.com/golangShare/publishdemo@v0.1.0
此时项目 go.mod 文件中，将会增加以下一行记录

require github.com/golangShare/publishdemo v0.1.0
和其他三方库一样的方式使用即可

package main

import (
 "fmt"
 "github.com/golangShare/publishdemo/stringutil"
)

func main() {
 s := stringutil.Reverse("hello gopher")
 fmt.Println(s)
}
总结
看完了上述流程，可以发现：开源自己的 Go 库，其实非常简单。

但还是有些地方需要注意一下：

不要在开源的 mod 文件中记录不需要的依赖。
即使发现 bug，也不要更改已发布版本的代码，而应该发布新版本。
遵循版本命名规范，推荐参考 Module version numbering 一文：https://go.dev/doc/modules/version-numbers 一文
Image

- END -
