# 自动发现 Go 项目 Bug 的神器


Go1.18 新特性中有一个神器：Fuzzing，对于发现 Go 项目中的 Bug 很有帮助。本文通过一个具体的例子来介绍它的基本使用，希望你能掌握并应用。

以下这个函数，你能找到几个 bug？它的功能看起来很简单——对于一个字符串，用一个新的用户定义字符覆盖它的第一个字符 `n` 次。例如，如果我们运行`OverwriteString("Hello, World!", "A", 5)`，正确的输出是：`"AAAAA, World!"`。

```go
// overwrite_string.go

// OverwriteString overwrites the first 'n' characters in a string with
// the rune 'value'
func OverwriteString(str string, value rune, n int) string {
 // If asked to overwrite more than the entire string then no need to loop,
 // just return string length * the rune
 if n > len(str) {
  return strings.Repeat(string(value), len(str))
 }

 result := []rune(str)
 for i := 0; i <= n; i++ {
  result[i] = value
 }
 return string(result)
}
```

在为代码提供一次快速可视化所需的时间内，Go 的新模糊测试工具可以通过该函数运行超过 500 万个程序生成的输入，并在这种情况下在*一秒钟*内找到导致越界数组访问的输入。

比如，使用这组参数运行函数：`OverwriteString("0000", rune('A'), 4)`会导致 panic：

```
--- FAIL: FuzzBasicOverwriteString (0.05s)
    --- FAIL: FuzzBasicOverwriteString (0.00s)
        testing.go:1349: panic: runtime error: index out of range [4] with length 4
            goroutine 96 [running]:
            runtime/debug.Stack()
             /home/everest/sdk/gotip/src/runtime/debug/stack.go:24 +0x90
            testing.tRunner.func1()
           ...<snip> 
    
    Failing input written to testdata/fuzz/FuzzBasicOverwriteString/2bac7bdf139ad0b2de37275db2a606ecb335bd344500173b451e9dfc3658c12f
    To re-run:
    go test -run=FuzzBasicOverwriteString/2bac7bdf139ad0b2de37275db2a606ecb335bd344500173b451e9dfc3658c12f
```

模糊测试（Fuzzing）是一种强大的测试技术，它非常擅长发现开发人员通常会遗漏的 Bug 和漏洞，并且在发现开源 Go 代码中的数百个关键错误方面有着良好的记录。

将我们的小示例问题扩展到关键应用程序中的千行代码路径，通过数十亿输入的模糊器只需几分钟即可发现细微的错误，否则这些错误在生产中需要数天才能解决。下面首先介绍如何使用 Go 的最新测试工具并尽快开始发现自己的错误。

## 入门

这是 Go 1.18 的新特性：模糊测试功能，因此在开始之前，请确保您`$ go version`的版本至少为 1.18。如果你的版本低于 1.18，请升级。

如果想跟着代码做，可以在 github.com/fuzzbuzz/go-fuzzing-tutorial 找到这篇文章的代码。对于本教程的其余部分，所有命令都是从`introduction`目录中运行的。

这是模糊测试的基本写法：

```go
// overwrite_string_test.go

func FuzzBasicOverwriteString(f *testing.F) {
 f.Fuzz(func(t *testing.T, str string, value rune, n int) {
  OverwriteString(str, value, n)
 })
}
```

与期望来自固定输入的特定行为的单元测试相反，模糊测试通过其测试的功能运行数千个程序生成的输入，而无需开发人员手动提供输入。在这种特定情况下，我们希望将测试的函数传递给`f.Fuzz`，因此模糊器将生成一个新`string`、`rune` 和 `int` 来填充每次测试迭代的参数。

默认情况下，模糊测试将检测崩溃、挂起和极端内存使用情况，因此即使不编写任何断言，我们也已经为我们的函数构建了一个有用的健壮性测试。

要运行此测试，执行如下命令：

```
go test -fuzz FuzzBasicOverwriteString
```

在大约一秒钟内，你应该会看到带有类似于以下错误信息的测试退出：（你的运行结果不会完全一样）

```
--- FAIL: FuzzBasicOverwriteString (0.05s)
    --- FAIL: FuzzBasicOverwriteString (0.00s)
        testing.go:1349: panic: runtime error: index out of range [4] with length 4
...SNIP
    Failing input written to testdata/fuzz/FuzzBasicOverwriteString/2bac7bdf139ad0b2de37275db2a606ecb335bd344500173b451e9dfc3658c12f
    To re-run:
    go test -run=FuzzBasicOverwriteString/2bac7bdf139ad0b2de37275db2a606ecb335bd344500173b451e9dfc3658c12f
```

模糊器在 `testdata/fuzz/FuzzBasicOverwriteString`目录内存放导致问题的特定输入的文件。打开这个文件，你可以看到导致我们的函数 panic 的实际值：（你的值可能不一样）

```
go test fuzz v1
string("00")
rune('A')
int(2)
```

现在我们已经发现了一个错误，可以进入我们的代码修复问题。查看实际导致 panic ( `overwrite_string.go:16`) 的代码行，该代码似乎试图访问长度为 4 的字符串的索引 4，这导致了数组索引越界错误。你可以通过更改检查 `if n > len(str)` 以测试大于或等于来修复错误：

```go
// overwrite_string.go

// OverwriteString overwrites the first 'n' characters in a string with
// the rune 'value'
func OverwriteString(str string, value rune, n int) string {
 // If asked to overwrite more than the entire string then no need to loop,
 // just return string length * the rune
 if n >= len(str) {
  return strings.Repeat(string(value), len(str))
 }

 result := []rune(str)
 for i := 0; i <= n; i++ {
  result[i] = value
 }
 return string(result)
}
```

这将确保仅当 `n` 至少小于字符串长度 1 时才输入循环。我们也可以修复 for 循环的边界，但这隐藏了一个更有趣的错误，所以现在我们忽略它。

通过使用输出中提供的 fuzzer 命令重新运行崩溃的测试用例，确认修复了 Bug：

```sh
$ go test -v -count=1 -run=FuzzBasicOverwriteString/2bac7bdf139ad0b2de37275db2a606ecb335bd344500173b451e9dfc3658c12f

=== RUN   FuzzBasicOverwriteString
=== RUN   FuzzBasicOverwriteString/2bac7bdf139ad0b2de37275db2a606ecb335bd344500173b451e9dfc3658c12f
--- PASS: FuzzBasicOverwriteString (0.00s)
    --- PASS: FuzzBasicOverwriteString/2bac7bdf139ad0b2de37275db2a606ecb335bd344500173b451e9dfc3658c12f (0.00s)
PASS
ok   github.com/fuzzbuzz/go-fuzzing-tutorial/introduction 0.001s
```

任何时候执行 `go test`（这些输入统称为“种子”），Go 的 fuzzer 将自动运行 `testdata` 目录中的每个输入作为单元测试。将目录 testdata 提交 到版本控制会将此输入保存为永久回归测试，以确保永远不会重新引入该错误。

## 一次意外之旅

现在，我完全承认，在我写这篇文章的时候，我希望这个改变能够满足基本的模糊测试，但是如果你在这个改变之后重新运行模糊器，你会注意到一个全新的错误出现：

```sh
$ go test -fuzz FuzzBasicOverwriteString
fuzz: elapsed: 0s, gathering baseline coverage: 0/17 completed
fuzz: elapsed: 0s, gathering baseline coverage: 17/17 completed, now fuzzing with 8 workers
fuzz: minimizing 177-byte failing input file
fuzz: elapsed: 0s, minimizing
--- FAIL: FuzzBasicOverwriteString (0.17s)
    --- FAIL: FuzzBasicOverwriteString (0.00s)
        testing.go:1349: panic: runtime error: index out of range [60] with length 60
            goroutine 2911 [running]:
            runtime/debug.Stack()
             /home/everest/sdk/gotip/src/runtime/debug/stack.go:24 +0x90
            testing.tRunner.func1()
             /home/everest/sdk/gotip/src/testing/testing.go:1349 +0x1f2
            panic({0x5b3700, 0xc00289c798})
             /home/everest/sdk/gotip/src/runtime/panic.go:838 +0x207
            github.com/fuzzbuzz/go-fuzzing-tutorial/introduction.OverwriteString({0xc00288ef00, 0x3d}, 0x83, 0x3c)
             /home/everest/src/fuzzbuzz/go-fuzzing-tutorial/introduction/overwrite_string.go:20 +0x270
            github.com/fuzzbuzz/go-fuzzing-tutorial/introduction.FuzzBasicOverwriteString.func1(0x5?, {0xc00288ef00?, 0x0?}, 0x0?, 0x0?)
             /home/everest/src/fuzzbuzz/go-fuzzing-tutorial/introduction/overwrite_string_test.go:24 +0x38
            reflect.Value.call({0x598d60?, 0x5cfb58?, 0x13?}, {0x5c179f, 0x4}, {0xc0028c2de0, 0x4, 0x4?})
             /home/everest/sdk/gotip/src/reflect/value.go:556 +0x845
            reflect.Value.Call({0x598d60?, 0x5cfb58?, 0x514?}, {0xc0028c2de0, 0x4, 0x4})
             /home/everest/sdk/gotip/src/reflect/value.go:339 +0xbf
            testing.(*F).Fuzz.func1.1(0x0?)
             /home/everest/sdk/gotip/src/testing/fuzz.go:337 +0x231
            testing.tRunner(0xc0028e7380, 0xc0028ec5a0)
             /home/everest/sdk/gotip/src/testing/testing.go:1439 +0x102
            created by testing.(*F).Fuzz.func1
             /home/everest/sdk/gotip/src/testing/fuzz.go:324 +0x5b8
            
    
    Failing input written to testdata/fuzz/FuzzBasicOverwriteString/2ee896e38866e089811eeece13f9919795072e6cc05ee9f782d68d1663d204c7
    To re-run:
    go test -run=FuzzBasicOverwriteString/2ee896e38866e089811eeece13f9919795072e6cc05ee9f782d68d1663d204c7
FAIL
exit status 1
FAIL github.com/fuzzbuzz/go-fuzzing-tutorial/introduction 0.174s
```

你看到的情况下，模糊器生成的实际输入可能看起来不同，以下是我看到的测试用例：

```
go test fuzz v1
string("000000000000000000000000000000Ö00000000000000000000000000000")
rune('\u0083')
int(60)
```

乍一看，这个错误看起来几乎与你刚刚修复的错误一模一样。尝试访问 60 个字符长字符串的索引 60 应该是不可能的，因为该函数将在初始 if 语句处返回。但是这就是模糊测试的力量——它揭示了开发人员没有考虑过的边缘情况，这实际上是一个完全独立的错误。

如果你检查 panic 的输入，你可能会像我一样注意到其中一个字符是 Unicode 字符。也就是说，它由一个以上的字节表示。在我的情况下，它是 `Ö`。当然，这个输入字符串有 60 个字符长，但它有*61 个字节长*。在 Go 中，通过 `len` 是获取字符串中的字节数，而不是字符数（或 rune）。

这很容易自己检查。如果你运行以下 Go 代码片段：

```go
str := "Ö"
runeArray := []rune(str)
fmt.Println("Str len:", len(str), "Rune array len:", len(runeArray))
```

将看到以下输出：

```
Str len: 2 Rune array len: 1
```

有了关于 Go 的字符串实现的重要信息，再次重写 if 语句，从 `if n >= len(str)` 改为 `if n >= utf8.RuneCountInString(str)`。因此我们想要比较的是字符数而不是字节数：

```go
// overwrite_string.go

// OverwriteString overwrites the first 'n' characters in a string with
// the rune 'value'
func OverwriteString(str string, value rune, n int) string {
 // If asked to overwrite more than the entire string then no need to loop,
 // just return string length * the rune
 if n >= utf8.RuneCountInString(str) {
  return strings.Repeat(string(value), len(str))
 }

 result := []rune(str)
 for i := 0; i <= n; i++ {
  result[i] = value
 }
 return string(result)
}
```

再次运行 fuzz 测试，观察它的变化，试图找到另一个输入来使我们的函数 panic：

```
go test -fuzz FuzzBasicOverwriteString
```

你*应该*让它运行一段时间以确保没有任何其他错误潜伏的错误，至少可以确信该函数不会在最基本的输入上崩溃。你可以按`ctrl/cmmand-C`停止模糊器。

## 功能性 Bug

到目前为止，我们已经发现了导致崩溃的错误。拒绝服务是一件大事，但我们所知道的是，这个函数不会因意外输入而崩溃。但测试函数的*正确性*也很重要。有很多方法可以解决这个问题，但是通过模糊测试，最好尝试考虑一个始终适用于你的代码*的不变量（或属性） 。*

`OverwriteString` 不变量的一个例子是，该函数永远不应填充比它被要求的数字*更多的字符*。更具体地说，如果被要求覆盖“Hello, world!” 使用 5 个 “A” 字符，应该可以检查字符串中剩余的字符是否仍然是“, world!”。（通常成为语料或种子）

这可以通过以下测试进行一般化：

```go
// overwrite_string_test.go

func FuzzOverwriteStringSuffix(f *testing.F) {
 f.Add("Hello, world!", 'A', 15)

 f.Fuzz(func(t *testing.T, str string, value rune, n int) {
  result := OverwriteString(str, value, n)
  if n > 0 && n < utf8.RuneCountInString(str) {
   // If we modified characters [0:n], then characters [n:] should stay the same
   resultSuffix := string([]rune(result)[n:])
   strSuffix := string([]rune(str)[:])
   if resultSuffix != strSuffix {
    t.Fatalf("OverwriteString modified too many characters! Expected %s, got %s.", strSuffix, resultSuffix)
   }
  }
 })
}
```

运行测试会发现另一个 Bug：

```sh
$ go test -fuzz FuzzOverwriteStringSuffix

fuzz: elapsed: 0s, gathering baseline coverage: 0/54 completed
fuzz: minimizing 66-byte failing input file
fuzz: elapsed: 0s, gathering baseline coverage: 8/54 completed
--- FAIL: FuzzOverwriteStringSuffix (0.03s)
    --- FAIL: FuzzOverwriteStringSuffix (0.00s)
        overwrite_string_test.go:38: OverwriteString modified too many characters! Expected 0, got A.
    
    Failing input written to testdata/fuzz/FuzzOverwriteStringSuffix/148139e8febb077401421c031a9bd3c3315179c5a66c90349102d223b451ec02
    To re-run:
    go test -run=FuzzOverwriteStringSuffix/148139e8febb077401421c031a9bd3c3315179c5a66c90349102d223b451ec02
FAIL
exit status 1
FAIL github.com/fuzzbuzz/go-fuzzing-tutorial/introduction 0.031s
```

请注意，这一次不是 panic，而是一条看起来非常类似于单元测试失败的消息。实际上，这段代码一直存在一个功能性错误。

在第 20 行它检查循环索引用的是 `<=`，所以它一直填充多余的字符。将循环条件从 `i <= n` 更改为 `i < n` 解决此问题。

最终`OverwriteString`函数应该如下所示：

```go
// overwrite_string.go

// OverwriteString overwrites the first 'n' characters in a string with
// the rune 'value'
func OverwriteString(str string, value rune, n int) string {
 // If asked for more no need to loop, just return
 // string length * the rune
 if n >= utf8.RuneCountInString(str) {
  return strings.Repeat(string(value), len(str))
 }

 result := []rune(str)
 for i := 0; i < n; i++ {
  result[i] = value
 }
 return string(result)
}
```

如果再次运行 fuzzer，你应该会看到 fuzzer 每秒可靠地运行数千个输入，同时没有发现其他错误。

理想情况下，应该让这个模糊测试运行至少几分钟，以增加对该代码正确性的信心（特别是如果它正在测试超过 10 行的函数）。

这篇文章中的错误是在几秒钟内发现的，但有些错误可能需要数小时或数天的时间来进行模糊测试，因为模糊器需要时间来探索被测软件的整个状态空间。后续文章将介绍大规模连续模糊测试的艺术。

## 结语

这只是对 Go 模糊测试的简要介绍。今天讨论的示例是你开始向自己的项目添加模糊测试所需的全部内容。

后续文章将深入研究你可以找到的错误类型，调查一些真实世界的模糊测试错误，并讨论如何自动化你的模糊测试，以便 CI 在你睡着时发现错误。

> 作者：Everest Munro-Zeisberger，原文链接：https://blog.fuzzbuzz.io/go-fuzzing-basics



# Go 1.18 让写测试的代码量骤减，你会开始写测试吗？

Original 网管 [网管叨bi叨](javascript:void(0);) *2022-06-20 08:45* *Posted on 北京*

收录于合集#Go语言实战技巧44个

模糊测试是一种向程序提供随机意外的输入以测试可能的崩溃或者边缘情况的方法。通过模糊测试可以揭示一些逻辑错误或者性能问题，因此使用模糊测试可以让程序的稳定性和性能都更有保证。

Go 从1.18 版本开始正式把模糊测试（Go Fuzz）加入到了其工具集中，不再依靠三方库就能在程序代码中进行模糊测试。那么为什么要引入模糊测试呢，引入后我们在写单元测试的时候要有哪些调整呢？

首先我们来聊聊为什么引入模糊测试。

## 为什么引入模糊测试

大家看文章开头第一段的解释，那就是Go官方要引入模糊测试的原因。估计各位看了想要打人，哈，那我就结合个简单的例子再把上面那段话要表达的意思，用代码再解释一遍。

大家先不考虑什么模糊测试的事儿，就单纯给下面这个工具函数写一个单测，我们该怎么写。

```go
func Equal(a []byte, b []byte) bool {
 for i := range a {
  if a[i] != b[i] {
   return false
  }
 }
 return true
}
```

这个工具函数将接收两个字节切片，比较他们的值是否相等。那么为了通过单测测试这个工具函数是否能如预期那样完成任务，我们就需要提供一些样本数据，来测试函数的知识结果。

### **单元测试怎么写**

我们在之前Go 单元测试入门中，给大家介绍过表格测试，就是为单测的执行提供样本数据的，那么这个表格测试该怎么写呢？这里直接放代码了，如果对表格测试和各种Go单测知识不了解的可以回看之前的文章：Go单元测试基础，文末会给出链接。

```go
func TestEqualWithTable(t *testing.T) {
 tests := []struct {
  name   string
  inputA []byte
  inputB []byte
  want   bool
 }{
  {"right case", []byte{'f', 'u', 'z', 'z'}, []byte{'f', 'u', 'z', 'z'}, true},
  {"right case", []byte{'a', 'b', 'c'}, []byte{'b', 'c', 'd'}, false},
 }

 for _, tt := range tests {
  tt := tt
  t.Run(tt.name, func(t *testing.T) {
   if got := Equal(tt.inputA, tt.inputB); got != tt.want {
    t.Error("expected " + strconv.FormatBool(tt.want) + ",  got " +
        strconv.FormatBool(got))
   }
  })
 }
}
```

上面这个单元测试使用的两个样本数据能让测试通过，但不代表我们的工具函数就完美无缺了，毕竟这里的两个样本都太典型了，如果你把输入的两个切片搞的不一样，工具函数直接就`index out of range`，程序直接挂掉了。

如果没有模糊测试呢，我们就需要在表格测试里尽量多的提供样本，才能测出各种边界情况下程序是否符合预期。

不过让自己提供样本测试，主观性太强，有的人能想到很多边界条件有的就不行，再加上我国互联网公司程序员糟糕的职场生存环境，又要保证BUG少稳定，又要快，这个时候模糊测试确实能帮助我们节省很多想样本的工作量。

### **用模糊测试简化**

现在我们换用Go 1.18 的 Fuzz 模糊测试，来测试下我们的工具函数。

```go
func FuzzEqual(f *testing.F) {
 //f.Add([]byte{'a', 'b', 'c'}, []byte{'a', 'b', 'c'})
 f.Fuzz(func(t *testing.T, a []byte, b []byte) {
  Equal(a, b)
 })
}
```

虽然模糊测试是1.18 新引入的，但只是节省了我们写表格测试提供样本的流程，其他流程和以前的单元测试并不差别，所用到的知识也没有变化。

可以看到使用模糊测试后，代码量明显减少了很多。模糊测试会帮我们生产随机的输入，来供要测试的目标来使用。上面两个参数的输入是随机产生的（也有规则，模糊测试会先测各种空输入，这个规则我们可以不用管）

也可以通过`f.Add()`方法添加语料，注意这里语料设置的个数和顺序要和目标函数里的输入参数保持一致（就是除了 testing.T之外的参数）

此外还有点明显的差异大家一定要注意，使用模糊测试后，测试函数的声明跟普通单测的不一样

```go
// 普通单元测试
TestXXX(t *testing.T){}
// 使用模糊测试的测试函数，必须以Fuzz开头，形参类型为*testing.F
FuzzXXX(f *testing.F) {}
```

### **执行模糊测试**

模糊测试执行的时候需要给 `go test`加上`-fuzz`这个标记。

```sh
➜  go test -fuzz .
warning: starting with empty corpus
fuzz: elapsed: 0s, execs: 0 (0/sec), new interesting: 0 (total: 0)
fuzz: minimizing 57-byte failing input file
fuzz: elapsed: 0s, minimizing
--- FAIL: FuzzEqual (0.04s)
   --- FAIL: FuzzEqual (0.00s)
       testing.go:1349: panic: runtime error: index out of range [0] with length 0
```

执行模糊测试后，就能测出我上面说的索引越界的问题，这个时候我们就可以回去完善我们的工具函数，然后再进行模糊测试了，通过几轮执行，会让被测试的函数足够健壮。

我们示例的工具函数足够简单，所以修复起来也超简单，价格长度判断就可以了。

```go
func Equal(a []byte, b []byte) bool {
 if len(a) != len(b) {
  return false
 }

 for i := range a {
  if a[i] != b[i] {
   return false
  }
 }
 return true
}
```

再度执行模糊测试后程序不再会报错，不过这个时候你应该发现，测试程序会一直执行，除非主动停下来，或者发现了测试失败的情况才能让模糊测试终止。

这就是模糊测试和普通单测的另一个大区别了，普通单测执行完我们提供的 Case 后就会停止，而模糊测试是会不停的跑样本，直到发生测试失败的情况才会停止。这个时候我们就可以用命令指定一个测试时长，让模糊测试到时自动停止。

```sh
go test -fuzz=Fuzz -fuzztime=10s  .
```

这里我们通过 fuzztime 这个标志，给模糊测试指定了 10 s的测试时长，到时模糊测试就会自动停止。

```sh
➜  go test -fuzz=Fuzz -fuzztime=10s  .
fuzz: elapsed: 0s, gathering baseline coverage: 0/10 completed
fuzz: elapsed: 0s, gathering baseline coverage: 10/10 completed, now fuzzing with 8 workers
fuzz: elapsed: 3s, execs: 852282 (284056/sec), new interesting: 0 (total: 10)
fuzz: elapsed: 6s, execs: 1748599 (298745/sec), new interesting: 0 (total: 10)
fuzz: elapsed: 9s, execs: 2653073 (301474/sec), new interesting: 0 (total: 10)
fuzz: elapsed: 10s, execs: 2955038 (274558/sec), new interesting: 0 (total: 10)
PASS
ok      golang-unit-test-demo/fuzz_test_demo    11.783s
➜  fuzz_test_demo git:(master) ✗ 
```

## 怎么写好一个模糊测试

相信通过上面的例子，其实大家已经看到模糊测试该怎么编写。为了让内容更吸引人，文章并没有上来就给大家罗列一堆名称概念，这里我们再把理论上的一些东西补充一下，这样未来自己编写模糊测试的时候自己心里就更有谱啦。

### **模糊测试的结构**

下面是官方文档里，给出的一张 "模糊测试构成元素" 的图

![Image](https://mmbiz.qpic.cn/mmbiz_png/z4pQ0O5h0f5VoaTYCN25RU1HB6qljlEj7kPwuTe6b4ldAicwIXUwTfDOw1v6jwVbzpDgAJTqna6bSn4P2z0q6nQ/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)



这张图里能看出来这几点：

- Fuzz test: 即整个模糊测试函数，它的函数签名要求，函数名必须以关键字 Fuzz 开头，只有一个类型为`*testing.F`的参数，且没有返回值。
- 模糊测试也是测试，所以跟单测一样，必须位于`_test.go`文件中，可以语单测在统一文件。
- Fuzz Target：模糊测试中，由 f.Fuzz 指定的要执行的测试函数叫 fuzz target，一个模糊测试中只能包含一个 fuzz target，且它的第一个参数必须是`*testing.T`类型的，后面跟至少一个模糊参数，这个也好理解，如果没有这个参数，那随机输入该往哪输入呢。
- Fuzz argument：这个一条说过了，就是fuzz target 中第一个参数以后的参数都叫模糊参数，用来接收模糊测试随机生成的样本，这个参数数量应该是要跟我们的被测函数的形参数一致的。
- Seed corpus：语料，这个单词儿我也没见过，大家记住就是提供了它后，生产的随机参数都跟这个语料有相关性，不是瞎随机的，且用 f.Add 设置的语料个数，要跟模糊参数的个数、顺序、类型上保持一致。

更详细的解释，请参考官方文档：https://go.dev/doc/fuzz/
模糊测试的结构，来自：https://go.dev/doc/fuzz/
https://www.youtube.com/watch?v=7KWPiRq3ZYI
