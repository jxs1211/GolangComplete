# 自动发现 Go 项目 Bug 的神器

[Go招聘](javascript:void(0);) *2022-05-24 11:50* *Posted on 北京*

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

```
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

```
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

```
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

```
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

```
str := "Ö"
runeArray := []rune(str)
fmt.Println("Str len:", len(str), "Rune array len:", len(runeArray))
```

将看到以下输出：

```
Str len: 2 Rune array len: 1
```

有了关于 Go 的字符串实现的重要信息，再次重写 if 语句，从 `if n >= len(str)` 改为 `if n >= utf8.RuneCountInString(str)`。因此我们想要比较的是字符数而不是字节数：

```
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

```
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

```
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

```
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