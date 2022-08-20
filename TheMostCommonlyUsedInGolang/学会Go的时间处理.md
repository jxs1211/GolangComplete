https://mp.weixin.qq.com/s/4W1IDfYq56zBZuk24wD5ew

# 学会 Go 中的时间处理

[码农桃花源](javascript:void(0);) *2022-09-07 08:44* *Posted on 香港*

Editor's Note

这篇文章浅显易懂，值得一读。

作为程序员，我们经常需要对时间进行处理。在 Go 中，标准库 time 提供了对应的能力。

本文将介绍 time 库中一些重要的函数和方法，希望能帮助到那些一遇到 Go 时间处理问题就需要百度的童鞋。

## 应对时区问题

在编程中，我们经常会遭遇**八小时时间差**问题。这是由时区差异引起的，为了能更好地解决它们，我们需要理解几个时间定义标准。

GMT（Greenwich Mean Time），格林威治平时。GMT 根据地球的自转和公转来计算时间，它规定太阳每天经过位于英国伦敦郊区的皇家格林威治天文台的时间为中午12点。GMT 是前世界标准时。

UTC（Coordinated Universal Time），协调世界时。UTC 比 GMT 更精准，它根据原子钟来计算时间。在不需要精确到秒的情况下，可以认为 UTC=GMT。UTC 是现世界标准时。

从格林威治本初子午线起，往东为正，往西为负，全球共划分为 24 个标准时区，相邻时区相差一个小时。

```
package main

import (
 "fmt"
 "time"
)

func main() {
 fmt.Println(time.Now())
}
```

中国大陆使用的是东八时区的标准时，即北京时间 CST，China Standard Time。

```
$ go run main.go 
2022-07-17 16:37:31.186043 +0800 CST m=+0.000066647
```

这是默认时区下的结果，`time.Now()`的打印中会标注`+0800 CST`。

假设我们是在美国洛杉矶时区下，那得到的结果是什么呢？

```
$ TZ="America/Los_Angeles" go run main.go
2022-07-17 01:39:12.391505 -0700 PDT m=+0.000069514
```

可以看到，此时的结果是`-0700 PDT` 时间，即 PDT（Pacific Daylight Time）太平洋夏季时间。由于时区差异，两次执行的时间结果相差了 15 小时。

注意，在使用 Docker 容器时，系统默认的时区就是 UTC 时间（0 时区），和我们实际需要的北京时间相差八个小时，这是导致**八小时时间差**问题的经典场景。

时区问题的应对策略，可以详细查看 src/time/zoneinfo_unix.go 中 initLocal() 函数的加载逻辑。例如，可以通过指定环境变量 TZ，修改/etc/localtime文件等方式来解决。

因为时区问题非常重要，所以放在了文章第一部分讲述。下面开始介绍 time 库的使用。

## 时间瞬间 time.Time

time 库，最核心的对象是 time.Time 结构体。它的定义如下，用以表示某个瞬间的时间。

```
type Time struct {
  // wall and ext encode the wall time seconds, wall time nanoseconds,
 // and optional monotonic clock reading in nanoseconds.
   wall uint64
   ext  int64
   loc *Location
}
```

计算机在时间处理上，主要涉及到两种时钟。

- 墙上时钟（wall time），又称为钟表时间，用于表示具体的日期与时间。
- 单调时钟（monotonic clocks），总是保证时间是向前的，不会出现墙上时钟的回拨问题，因此它很适合用于测量持续时间段。

wall 和 ext 字段就是用于记录墙上时钟和单调时钟，精度为纳秒。字段的对应位数上关联着用于确定时间的具体年、月、日、小时、分钟、秒等信息。

loc 字段记录时区位置，当 loc 为 nil 时，默认为 UTC 时间。

因为 time.Time 用于表示具有纳秒精度的时间瞬间，在程序中通常应该将它作为值存储和传递，而不是指针。

即在时间变量或者结构体字段中，我们应该使用 time.Time，而非 *time.Time。

##### 获取 time.Time

我们可以通过 Now 函数获取当前本地时间

```
func Now() Time {}
```

也可以通过 Date 函数，根据年、月、日等时间和时区参数获取指定时间

```
func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time {}
```

##### 转换时间戳

计算机世界中，将 UTC 时间 1970 年1月1日 0 时 0 分 0 秒作为 Unix 时间 0。所谓的时间瞬间转换为 Unix 时间戳，即计算的是从 Unix 时间 0 到指定瞬间所经过的秒数、微秒数等。

```
func (t Time) Unix() int64 {}      // 从 Unix 时间 0 经过的秒数
func (t Time) UnixMicro() int64 {} // 从 Unix 时间 0 经过的微秒数
func (t Time) UnixMilli() int64 {} // 从 Unix 时间 0 经过的毫秒数
func (t Time) UnixNano() int64 {}  // 从 Unix 时间 0 经过的纳秒数
```

##### 获取基本字段

```
 t := time.Now()
 fmt.Println(t.Date())      // 2022 July 17
 fmt.Println(t.Year())      // 2022
 fmt.Println(t.Month())     // July
 fmt.Println(t.ISOWeek())   // 2022 28
 fmt.Println(t.Clock())     // 22 21 56
 fmt.Println(t.Day())       // 17
 fmt.Println(t.Weekday())   // Sunday
 fmt.Println(t.Hour())      // 22
 fmt.Println(t.Minute())    // 21
 fmt.Println(t.Second())    // 56
 fmt.Println(t.Nanosecond())// 494313000
 fmt.Println(t.YearDay())   // 198
```

## 持续时间 time.Duration

持续时间 time.Duration 用于表示两个时间瞬间 time.Time 之间所经过的时间。它通过 int64 表示纳秒计数，能表示的极限大约为 290 年。

```
// A Duration represents the elapsed time between two instants
// as an int64 nanosecond count. The representation limits the
// largest representable duration to approximately 290 years.
type Duration int64
```

在 Go 中，持续时间只是一个以纳秒为单位的数字而已。如果持续时间等于 1000000000，则它代表的含义是 1 秒或 1000 毫秒或 1000000 微秒或 1000000000 纳秒。

例如，相隔 1 小时的两个时间瞬间 time.Time 值，它们之间的持续时间 time.Duration 值为

```
1*60*60*1000*1000*1000
```

Go 的 time 包中定义了这些持续时间常量值

```
const (
 Nanosecond  Duration = 1
 Microsecond          = 1000 * Nanosecond
 Millisecond          = 1000 * Microsecond
 Second               = 1000 * Millisecond
 Minute               = 60 * Second
 Hour                 = 60 * Minute
)
```

同时，time.Duration 提供了能获取各时间粒度数值的方法

```
func (d Duration) Nanoseconds() int64 {}   // 纳秒
func (d Duration) Microseconds() int64 {}  // 微秒
func (d Duration) Milliseconds() int64 {}  // 毫秒
func (d Duration) Seconds() float64 {}     // 秒
func (d Duration) Minutes() float64 {}     // 分钟
func (d Duration) Hours() float64 {}       // 小时
```

## 时间计算

在学习了时间瞬间和持续时间之后，我们来看如何做时间计算。

```
func (t Time) Add(d Duration) Time {}
```

- Add 函数用于增加/减少（ d 的正值表示增加、负值表示减少） time.Time 的持续时间。我们可以对某瞬时时间，增加或减少指定纳秒级以上的时间。

```
func (t Time) Sub(u Time) Duration {}
```

- Sub  函数可以得出两个时间瞬间之间的持续时间。

```
func (t Time) AddDate(years int, months int, days int) Time {}
```

- AddDate 函数基于年、月和日的维度增加/减少 time.Time 的值。

当然，基于当前时间瞬间 time.Now() 的计算是最普遍的需求。因此，time 包还提供了以下便捷的时间计算函数。

```
func Since(t Time) Duration {}
```

Since 函数是 time.Now().Sub(t) 的快捷方法。

```
func Until(t Time) Duration {}
```

Until 函数是 t.Sub(time.Now()) 的快捷方法。

##### 使用示例

```
 t := time.Now()
 fmt.Println(t)                      // 2022-07-17 22:41:06.001567 +0800 CST m=+0.000057466

 //时间增加 1小时
 fmt.Println(t.Add(time.Hour * 1))   // 2022-07-17 23:41:06.001567 +0800 CST m=+3600.000057466
 //时间增加 15 分钟
 fmt.Println(t.Add(time.Minute * 15))// 2022-07-17 22:56:06.001567 +0800 CST m=+900.000057466
 //时间增加 10 秒钟
 fmt.Println(t.Add(time.Second * 10))// 2022-07-17 22:41:16.001567 +0800 CST m=+10.000057466

 //时间减少 1 小时
 fmt.Println(t.Add(-time.Hour * 1))  // 2022-07-17 21:41:06.001567 +0800 CST m=-3599.999942534
 //时间减少 15 分钟
 fmt.Println(t.Add(-time.Minute * 15))// 2022-07-17 22:26:06.001567 +0800 CST m=-899.999942534
 //时间减少 10 秒钟
 fmt.Println(t.Add(-time.Second * 10))// 2022-07-17 22:40:56.001567 +0800 CST m=-9.999942534

 time.Sleep(time.Second * 5)
 t2 := time.Now()
 // 计算 t 到 t2 的持续时间
 fmt.Println(t2.Sub(t))              // 5.004318874s
 // 1 年之后的时间
 t3 := t2.AddDate(1, 0, 0)
 // 计算从 t 到当前的持续时间
 fmt.Println(time.Since(t))          // 5.004442316s
 // 计算现在到明年的持续时间
 fmt.Println(time.Until(t3))         // 8759h59m59.999864s
```

## 格式化时间

在其他语言中，一般会使用通用的时间模板来格式化时间。例如 Python，它使用 %Y 代表年、%m 代表月、%d 代表日等。

但是，Go 不一样，它使用**固定的时间**（需要注意，使用其他的时间是不可以的）作为布局模板，而这个固定时间是 Go 语言的诞生时间。

```
Mon Jan 2 15:04:05 MST 2006
```

格式化时间涉及到两个转换函数

```
func Parse(layout, value string) (Time, error) {}
```

- Parse 函数用于将时间字符串根据它所能对应的布局转换为 time.Time 对象。

```
func (t Time) Format(layout string) string {}
```

- Formate 函数用于将 time.Time 对象根据给定的布局转换为时间字符串。

##### 示例

```
const (
   layoutISO = "2006-01-02"
   layoutUS  = "January 2, 2006"
)
date := "2012-08-09"
t, _ := time.Parse(layoutISO, date)
fmt.Println(t)                  // 2012-08-09 00:00:00 +0000 UTC
fmt.Println(t.Format(layoutUS)) // August 9, 2012
```

在 time 库中，Go 提供了一些预定义的布局模板常量，这些可以直接拿来使用。

```
const (
 Layout      = "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
 ANSIC       = "Mon Jan _2 15:04:05 2006"
 UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
 RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
 RFC822      = "02 Jan 06 15:04 MST"
 RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
 RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
 RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
 RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
 RFC3339     = "2006-01-02T15:04:05Z07:00"
 RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
 Kitchen     = "3:04PM"
 // Handy time stamps.
 Stamp      = "Jan _2 15:04:05"
 StampMilli = "Jan _2 15:04:05.000"
 StampMicro = "Jan _2 15:04:05.000000"
 StampNano  = "Jan _2 15:04:05.000000000"
)
```

下面是我们可选的布局参数对照表：

```
年         06/2006
月         01/1/Jan/January
日         02/2/_2
星期       Mon/Monday
小时       03/3/15
分         04/4
秒         05/5
毫秒       .000/.999
微秒       .000000/.999999
纳秒       .000000000/.999999999
am/pm     PM/pm
时区       MST
时区小时数差-0700/-07/-07:00/Z0700/Z07:00
```

## 时区转换

在文章开头，我们介绍了时区问题。如果在代码中，需要获取同一个 time.Time 在不同时区下的结果，我们可以使用它的 In 方法。

```
func (t Time) In(loc *Location) Time {}
```

它的使用非常简单，直接看示例代码

```
now := time.Now()
fmt.Println(now)          // 2022-07-18 21:19:59.9636 +0800 CST m=+0.000069242

loc, _ := time.LoadLocation("UTC")
fmt.Println(now.In(loc)) // 2022-07-18 13:19:59.9636 +0000 UTC

loc, _ = time.LoadLocation("Europe/Berlin")
fmt.Println(now.In(loc)) // 2022-07-18 15:19:59.9636 +0200 CEST

loc, _ = time.LoadLocation("America/New_York")
fmt.Println(now.In(loc)) // 2022-07-18 09:19:59.9636 -0400 EDT

loc, _ = time.LoadLocation("Asia/Dubai")
fmt.Println(now.In(loc)) // 2022-07-18 17:19:59.9636 +0400 +04
```

## 总结

整体而言，time 库提供的时间处理函数和方法，基本满足我们的使用需求。

有意思的是，Go 时间格式化转换必须采用 Go 诞生时间，确实有够自恋。