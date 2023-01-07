# Go error 源码解读、错误处理的优化与最佳实践

作者：烤冷面

- 2023-03-19

  广东

- 本文字数：6813 字

  阅读完需：约 22 分钟

![Go error 源码解读、错误处理的优化与最佳实践](https://static001.geekbang.org/infoq/3c/3cc216e4f87233c822f931f4ffcf1829.png)

Go 语言自从诞生起，它的错误处理机制一直被喷出翔🙂。

没错，Go 语言在诞生初期确实简陋得不行，但在多个版本迭代以及各位前辈的摸索下还是找到了 Go 语言「错误处理」的最佳实践。

下面我们深入了解下 Go 的 error 包，并讨论如何让我们的 Go 项目拥有清爽的错误处理。

## Go 的 errors 包

Go 中的 error 是一个简单的内置接口类型。只要实现了这个接口，就可以将其视为一种 error。

```
type error interface {    Error() string}
```

复制代码

与此同时，Go 的 errors 包实现了这个接口：调用 `errors.New()` 就会返回`error接口`的实现类`errorString`，通过源码我们看到`errorString`的底层就是一字符串，可真是"省事"啊🙃。

```go
package errors

func New(text string) error { return &errorString{text} }

type errorString struct{ s string }

func (e *errorString) Error() string { return e.s }
```

复制代码

> `errors.New()函数`返回的是`errorString`的指针类型，这样做的目的是为了防止字符串产生碰撞。
>
> 我们可以做个小测试：`error1`和`error2`的 text 都是`"error"`，但是二者并不相等。

```go
func TestErrString(t *testing.T) {
	var error1 = errors.New("error")
	var error2 = errors.New("error")
	if error1 != error2 {
		log.Println("error1 != error2")
	}
}
---------------------代码运行结果--------------------------
=== RUN   TestXXXX2022/03/25 22:05:40 error1 != error2
```

复制代码

这种创建 error 的方式很常见，在 Go 源码以及三方包源码中大量出现。

```go
var EOF = errors.New("EOF")
var ErrUnexpectedEOF = errors.New("unexpected EOF")
var ErrNoProgress = errors.New("multiple Read calls return no data or error")
```

复制代码

然而很可惜的是，Go 的 error 设计并不能满足所有场景。

## Go error 的设计缺陷

### error 具有二义性

发生`error != nil`时不再意味着一定发生了错误，比如 io.Reader 返回 io.EOF 来告知调用者数据已经读取完毕，而这并不算是一个错误。

### 在两个包之间创建了依赖

比如我们使用了 io.EOF 来检查数据是否读取完毕，那么代码里一定会导入 io 包。

### 错误信息太单薄

只有一个字符串表达错误，过于单薄。

## 改进 Go error

现在我们知道 error 底层其实就是一字符串，它很简洁，但反过来也意味着"简陋"，无法携带更多错误信息。

### 自定义错误类型

所以程序员们决定自己封装一个 error 结构体，比如 Go 源码中的 os.PathError。

```go
type PathError struct {
	Op   string
	Path string
	Err  error
}
```

复制代码

### 封装 error 堆栈信息

将 error 封装后确实能表达更多的错误信息，但是它还有一个致命问题：没有堆栈信息。

比如这种日志，鬼知道代码哪一行报了错，Debug 时简直要命。

```go
SERVICE ERROR 2022-03-25T16:32:10.687+0800!!!       
Error 1406: Data too long for column 'content' at row 1
```

复制代码

我们可以使用`github.com/pkg/error包`解决这个问题，这个包提供了`errors.withStack()方法`将堆栈信息封装进 error：

```go
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}

type withStack struct{
	error *stack 
}
```

复制代码

### 防止 error 被覆盖

上层 error 想附带更多日志信息时，往往会使用`fmt.Errorf()`，`fmt.Errorf()`会创建一个新的 error 覆盖掉原本的 error 类型，我们写一个 demo 测试一下:

```
var errNoRows = errors.New("no rows")

// 模仿sql库返回一个errNoRowsfunc sqlExec() error {    return errNoRows}
func serviceNoErrWrap() error {
	err := sqlExec()
	if err != nil {
		// fmt.Errorf() 吞掉了原本的 errNoRows 类型错误。
		return fmt.Errorf("sqlExec failed.Err:%v", err)
	}
	return nil
}

func TestErrWrap(t *testing.T) {
	// 使用fmt.Errorf创建了一个新的err，丢失了底层err
	err := serviceNoErrWrap()
	if err != errNoRows {
		log.Println("===== errType don't equal errNoRows =====")
	}
}

-------------------------------代码运行结果----------------------------------
=== RUN   TestErrWrap2022/03/26 17:19:43 ===== errType don't equal errNoRows =====
```

复制代码

同样，使用`github.com/pkg/error包`的`errors.Wrap()函数`可以帮助我们为 error 添加自定义的文本信息。

```go
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   message,
	}
	return &withStack{
		err,
		callers(),
	}
}
```

复制代码

> `github.com/pkg/error包` 内容很多，这里不展开聊了，后面单独讲。

到此为止，我们深入认识了 Go 的 error，现在我们谈谈如何在大型项目中做好错误处理。

## error 处理最佳实践

### 优先处理 error

当一个函数返回 error 时，应该优先处理 error，忽略其他返回值。

### 只处理 error 一次

在 Go 中，每个 err 只应该被处理一次。 如果一个函数返回了 err，那么我们有两个选择：

1. 选择一：立即处理 err（包括记日志等行为），然后 return nil（把错误吞掉）。

> 这个行为可以被认为是对 error 做降级处理，所以一定要小心处理函数返回值。

1. 选择二：直接 return err，把 err 抛给调用者。

如果我们违反了这个原则会导致什么后果？请看反例：

```go
// 试想如果writeAll函数出错，会打印两遍日志
// 如果整个项目都这么做，最后会惊奇的发现我们在处处打日志，项目中存在大量没有价值的垃圾日志
// unable to write:io.EOF
// could not write config:io.EOF
type config struct{}

func writeAll(w io.Writer, buf []byte) error {
	_, err := w.Write(buf)
	if err != nil {
		log.Println("unable to write:", err)
		return err
	}
	return nil
}

func writeConfig(w io.Writer, conf *config) error {
	buf, err := json.Marshal(conf)
	if err != nil {
		log.Printf("could not marshal config:%v", err)
	}
	if err := writeAll(w, buf); err != nil {
		log.Println("count not write config: %v", err)
		return err
	}
	return nil
}
```

复制代码

### 不要反复包装 error

我们应该包装 error，但仅包装一次，否则会造成日志重复打印。

上层业务代码建议`Wrap error`，但是底层`基础Kit库`则不建议这样做。比如 Go 的`sql库`会返回`sql.ErrNoRows`这种预定义错误，然后我们的业务代码将其包装后 return。

### 不透明的错误处理

在大型项目中，推荐使用`不透明的错误处理(Opaque errors)`：不关心错误类型，只关心 error 是否为 nil。



![img](https://static001.geekbang.org/infoq/32/32c6528b0f01e15cd8292084cb6bc50f.png)

Go error 源码解读、错误处理的优化与最佳实践



这种方式代码耦合小，不需要判断特定错误类型，也就不需要导入相关包的依赖。

当然了，在这种情况下，只要我们调用函数，就一定跟着一组`if err != nil{}`，这也是大家经常吐槽 Go 项目`if err != nil{}`满天飞的原因😂。

> 目前我们只讨论在调用 Go 内置库和第三方库时产生的 error 的最佳处理实践，业务层面的错误处理是一个单独的话题，以后单独写一篇聊。

## 优化错误处理流程

Go 因为代码中无数的`if err != nil`被诟病，现在我教大家一个优化技巧：

我们先看看 bufio.scan() 是如何简化 error 处理的：

```go
// CountLines() 实现了"读取内容的行数"功能
func CountLines(r io.Reader) (int, error) {
	var (
		br    = bufio.NewReader(r)
		lines int
		err   error
	)
	for {
		_, err := br.ReadString('\n')
		lines++
		if err != nil {
			break
		}
	}
	if err != io.EOF {
		return 0, nil
		// sadwawa
	}
	return lines, nil
}

// 利用 bufio.scan() 简化 error 的处理：
func CountLinesGracefulErr(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	lines := 0
	for sc.Scan() {
		lines++
	}
	return lines, sc.Err()
}
```

复制代码

源码解读：`bufio.NewScanner()` 返回一个 `Scanner` 对象，结构体内部包含了 error 类型，调用`Err()`方法即可返回封装好的 error。

```go
type Scanner struct {
    r            io.Reader // The reader provided by the client.
    split        SplitFunc // The function to split the tokens.
    maxTokenSize int       // Maximum size of a token; modified by tests.
    token        []byte    // Last token returned by split.
    buf          []byte    // Buffer used as argument to split.
    start        int       // First non-processed byte in buf.
    end          int       // End of data in buf.
    err          error     // Sticky error.
    empties      int       // Count of successive empty tokens.
    scanCalled   bool      // Scan has been called; buffer is in use.
    done         bool      // Scan has finished.
}

func (s *Scanner) Err() error {
    if s.err == io.EOF {
        return nil
    }
    return s.err
}
```

复制代码

利用上面学到的思路，我们可以自己实现一个`errWriter`对象，简化对 error 的处理：

```go
type Header struct {
    Key, Value string
}

type Status struct {
    Code   int
    Reason string
}

// WriteResponse()函数实现了"构建HttpResponse"功能
func WriteResponse(w io.Writer, st Status, headers []Header, body io.Reader) error {
    _, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", st.Code, st.Reason)
    if err != nil {
        return err
    }
    
    for _, h := range headers {
        _, err := fmt.Fprintf(w, "%s: %s\r\n", h.Key, h.Value)
        if err != nil {
            return err
        }
    }
    
    if _, err := fmt.Fprintf(w, "\r\n"); err != nil {
        return err
    }
    
    _, err = io.Copy(w, body)
    return err
}

// 优化错误处理
type errWriter struct {
    io.Writer
    err error
}

func (e *errWriter) Write(buf []byte) (n int, err error) {
    if e.err != nil {
        return 0, e.err
    }
    
    n, e.err = e.Writer.Write(buf)
    
    return n, nil
}

func WriteResponseGracefulErr(w io.Writer, st Status, headers []Header, body io.Reader) error {
    ew := &errWriter{w, nil}
    
    fmt.Fprintf(ew, "HTTP/1.1 %d %s\r\n", st.Code, st.Reason)
    
    for _, h := range headers {
        fmt.Fprintf(ew, "%s: %s\r\n", h.Key, h.Value)
    }
    
    fmt.Fprintf(w, "\r\n")
    
    io.Copy(ew, body)
    
    return ew.err
}
```

复制代码

## Go.1.13 版本 error 的新特性

Go 1.13 版本借鉴了`github.com/pkg/error包`，大幅增强了 Golang 语言判断 error 类型的能力，这些函数平时还是用得到的，我们深入学习下：

### errors.UnWrap()

```go
// 与errors.Wrap()行为相反
// 获取err链中的底层err
func Unwrap(err error) error {
    u, ok := err.(interface {
        Unwrap() error
    })
    if !ok {
        return nil
    }
    return u.Unwrap()
}
```

复制代码

### errors.Is()

在 1.13 版本之前，我们可以用`err == targetErr`判断 err 类型

`errors.Is()`是其增强版：error 链上的`任一err == targetErr`，即`return true`，我们写个 demo 跑一下：

```go
var errNoRows = errors.New("no rows")

// 模仿sql库返回一个errNoRows
func sqlExec() error {
    return errNoRows
}

func service() error {
    err := sqlExec()
    if err != nil {
        return errors.WithStack(err)    // 包装errNoRows
    }
    
    return nil
}

func TestErrIs(t *testing.T) {
    err := service()
    
    // errors.Is递归调用errors.UnWrap，命中err链上的任意err即返回true
    if errors.Is(err, errNoRows) {
        log.Println("===== errors.Is() succeeded =====")
    }
    
    //err经errors.WithStack包装，不能通过 == 判断err类型
    if err == errNoRows {
        log.Println("err == errNoRows")
    }
}
-------------------------------代码运行结果----------------------------------
=== RUN   TestErrIs
2022/03/25 18:35:00 ===== errors.Is() succeeded =====
```

复制代码

例子解读：因为使用`errors.WithStack`包装了`sqlError`，`sqlError`位于 error 链的底层，上层的 error 已经不再是`sqlError`类型，所以使用`==`无法判断出底层的`sqlError`

源码解读：

- 内部调用了`err = Unwrap(err)`方法来获取 error 链中每一个 error。
- 兼容自定义 error 类型。

```go
func Is(err, target error) bool {
    if target == nil {
        return err == target
    }
    
    isComparable := reflectlite.TypeOf(target).Comparable()
    for {
        if isComparable && err == target {
            return true
        }
        // 自定义的 error 可以实现`Is接口`自定义 error 类型判断逻辑
        if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
            return true
        }
        if err = Unwrap(err); err == nil {
            return false
        }
    }
}
```

复制代码

下面我们尝试使用`erros.Is()`识别自定义 error 类型：

```go
type errNoRows struct {
    Desc string
}

func (e errNoRows) Unwrap() error { return e }

func (e errNoRows) Error() string { return e.Desc }

func (e errNoRows) Is(err error) bool {
    return reflect.TypeOf(err).Name() == reflect.TypeOf(e).Name()
}

// 模仿sql库返回一个errNoRows
func sqlExec() error {
    return &errNoRows{"Kaolengmian NB"}
}

func service() error {
    err := sqlExec()
    if err != nil {
        return errors.WithStack(err)
    }
    
    return nil
}

func serviceNoErrWrap() error {
    err := sqlExec()
    if err != nil {
        return fmt.Errorf("sqlExec failed.Err:%v", err)
    }

    return nil
}

func TestErrIs(t *testing.T) {
    err := service()

    if errors.Is(err, errNoRows{}) {
        log.Println("===== errors.Is() succeeded =====")
    }
}
-------------------------------代码运行结果----------------------------------
=== RUN   TestErrIs
2022/03/25 18:35:00 ===== errors.Is() succeeded =====
```

复制代码

### errors.As()

在 1.13 版本之前，我们可以用`if _,ok := err.(targetErr)`判断 err 类型，现在`errors.As()`是其增强版：error 链上的`任一err与targetErr类型相同`，即`return true`，我们写个 demo 跑一下：

```go
// errors.WithStack 包装了 sqlError 
// sqlError 位于 error 链的底层，上层的 error 已经不再是 sqlError 类型
// 使用类型断言无法判断出底层的 sqlError，而使用 errors.As() 函数可以判断出底层的 sqlError
type sqlError struct {
    error
}

func (e *sqlError) IsNoRows() bool {
    t, ok := e.error.(ErrNoRows)
    return ok && t.IsNoRows()
}

type ErrNoRows interface {
    IsNoRows() bool
}

// 返回一个sqlError
func sqlExec() error {
    return sqlError{}
}

// errors.WithStack包装sqlError
func service() error {
    err := sqlExec()
    if err != nil {
        return errors.WithStack(err)
    }

    return nil
}

func TestErrAs(t *testing.T) {
    err := service()

    // 递归使用errors.UnWrap，只要Err链上有一种Err满足类型断言，即返回true
    sr := &sqlError{}
    if errors.As(err, sr) {
        log.Println("===== errors.As() succeeded =====")
    }

    // 经errors.WithStack包装后，不能通过类型断言将当前Err转换成底层Err
    if _, ok := err.(sqlError); ok {
        log.Println("===== type assert succeeded =====")
    }
}
----------------------------------代码运行结果--------------------------------------------
=== RUN   TestErrAs
2022/03/25 18:09:02 ===== errors.As() succeeded =====

```

复制代码

## 总结

这篇文章我们认识了 Go 的 error，研究了`error包`、`github.com/pkg/error包`的源码，也聊了聊针对 Go 项目错误处理的优化与最佳实践，文中有大量 Demo 代码，建议 copy 代码跑上两遍，有助于理解我单薄的文字，有助于快速掌握 Go 的 error 处理。

------

参考：

1. 《Effective GO》
2. Go 程序设计语言》
3. [https://dave.cheney.net/practical-go/presentations/qcon-china.html#_error_handling](https://xie.infoq.cn/link?target=https%3A%2F%2Fdave.cheney.net%2Fpractical-go%2Fpresentations%2Fqcon-china.html%23_error_handling)

------

文章归档：[Go源码解读](https://xie.infoq.cn/link?target=https%3A%2F%2Frustyscript.com%2Fzh-cn%2Ftags%2Fgo%E6%BA%90%E7%A0%81%E8%A7%A3%E8%AF%BB%2F)

转载声明：本文允许转载，原文地址：[Go error 源码解读、错误处理的优化与最佳实践](https://xie.infoq.cn/link?target=https%3A%2F%2Frustyscript.com%2Fzh-cn%2Fgo-error%2F)