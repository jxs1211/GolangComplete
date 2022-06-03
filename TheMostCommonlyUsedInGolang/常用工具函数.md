# Golang — Utility functions you always missed

以下是一些平时开发常用的通用函数，赶紧收藏起来，一定可以用得上。

1.检查 slice 中某元素是否存在。

Go 语言没有预定义的函数用来检测 slice 中某元素是否存在，下面的代码可以帮你实现。

```go
import "fmt"

func main() {
  // create an array of strings
  slice := []string{"apple", "grapes", "mango"}
  // You can check by changing element to "orange"
  if Contains(slice, "mango") {
    fmt.Println("Slice contains element")
  } else {
    fmt.Println("Slice doesn't contain element")
  }
}

func Contains(slice []string, element string) bool {
   for _, i := range slice {
     if i == element {
       return true
     }
   }
   return false
}


output:
Slice contains element
```

循环遍历 slice 每个元素，如果能匹配上说明存在，返回 true；否则不存在，返回 false。

2.检查给定的时间是否处于某一时间区间内。

```
func main() {
  currentTime := time.Now()
  // Time after 18 hours of currentTime
  futureTime := time.Now().Add(time.Hour * 18)
  // Time after 10 hours of currentTime
  intermediateTime := time.Now().Add(time.Hour * 10)
  if intermediateTime.After(currentTime) &&    intermediateTime.Before(futureTime) {
    fmt.Println("intermediateTime is between currentTime and  futureTime")
  } else {
    fmt.Println("intermediateTime is not inbetween currentTime and futureTime")
  }
}

output:
intermediateTime is between currentTime and futureTime
```

给定两个时间戳，用于检查当前时间是否位于两者之间。Go 标准库的 time 包提供了相应的内置函数 After() 和 Before()，如果两者都返回 true，那么就可以认为当前时间在给定时间之间。

3.计算特定时区的当前时间戳。

```
import (
 "time"
 "fmt"
)

func main() {
  timeZone := "Asia/Kolkata" // timezone value
  loc, _ := time.LoadLocation(timeZone)
  currentTime = time.Now().In(loc)
  fmt.Println("currentTime : ", currentTime)
}


output:
// for timezone = "Asia/Kolkata"
currentTime :  2022-02-09 10:42:39.164079505 +0530 IST
// for timezone = "Asia/Shanghai"
currentTime :  2022-02-09 13:14:33.986953939 +0800 CST
```

首先使用 time.LoadLocation() 加载给定时区，然后将其与 time.Now.In() 一起使用，以获取给定时区的当前时间。可以将变量 timeZone 的值修改成我们想要的时区。

4.将较小的数除以较大的数

如果将较小的整数除以较大的整数，则结果为 0，可以使用下面的方案保留小数。

```
import "fmt"

func main() {
  smallerNo := 5
  largerNo := 25
  result := float32(smallerNo) / float32(largerNo)
  fmt.Println("result : ", result)
}


output:
result : 0.2
```

将两个数字都转换为浮点数，然后将它们相除，结果将保留小数部分。

5.去重

通过下面的方案可以删除切片中重复项。

```
import "fmt"

func main() {
  // define array of strings
  fruits := []string{"Mango", "Grapes", "Kiwi", "Apple", "Grapes"}
  fmt.Println("Array before removing duplicates : ", fruits)
  // Array after duplicates removal
  dulicatesRemovedArray := RemoveDuplicatesFromSlice(fruits)
  fmt.Println("Array after removing duplicates : ",  dulicatesRemovedArray)
}

func RemoveDuplicatesFromSlice(intSlice []string) []string {
  keys := make(map[string]bool)
  list := []string{}
  for _, entry := range intSlice {
    if _, value := keys[entry]; !value {
      keys[entry] = true
      list = append(list, entry)
    }
  }
 return list
}


output:
Array before removing duplicates :  [Mango Grapes Kiwi Apple Grapes]
Array after removing duplicates :  [Mango Grapes Kiwi Apple]
```

我们定义了另一个切片，循环遍历原切片，如果元素在新切片中已经存在就忽略不处理，并返回元素不会重复的切片。

6.随机打乱

Go 语言没有相关的内置函数，可以通过下面代码实现。

```
import "fmt"

func main() {
  // shuffle array
  array := []string{"India", "US", "Canada", "UK"}
  Shuffle(array)
}

func Shuffle(array []string) {
  // seed random for changing order of elements
  random := rand.New(rand.NewSource(time.Now().UnixNano()))
  for i := len(array) - 1; i > 0; i-- {
     j := random.Intn(i + 1)
     array[i], array[j] = array[j], array[i]
  }
  fmt.Println("Shuffled array : ", array)
}


output:
Shuffled array :  [UK India Canada US]
```

想要打乱元素，我们必须引入随机数，然后再交换元素。

7.反转

可以通过下面函数来反转切片。

```
import "fmt"


func main() {
  a := []int{1, 2, 3, 4, 5, 6} // input int array
  reverseArray := ReverseSlice(a)
  fmt.Println("Reverted array : ", reverseArray) // print output
}

func ReverseSlice(a []int) []int {
  for i := len(a)/2 - 1; i >= 0; i-- {
   pos := len(a) - 1 - i
   a[i], a[pos] = a[pos], a[i]
}
 return a
}


output:
Reverted array :  [6 5 4 3 2 1]
```

8.元素求和

```
import "fmt"


func main() {
  s := []int{10, 20, 30}
  sum := sumSlice(s)
  fmt.Println("Sum of slice elements : ", sum)
}


func sumSlice(array []int) int {
  sum := 0
  for _, item := range array {
      sum += item
  }
  return sum
}


output:
Sum of slice elements :  60
```

通过循环遍历 slice 实现求和。

9.将 slice 转换为逗号分隔的字符串

```
import (
   "fmt"
   "strings"
   "strconv"
)


func main() {
   result := ConvertSliceToString([]int{10, 20, 30, 40})
   fmt.Println("Slice converted string : ", result)
}


func ConvertSliceToString(input []int) string {
   var output []string
   for _, i := range input {
      output = append(output, strconv.Itoa(i))
   }
   return strings.Join(output, ",")
}


output:
Slice converted string :  10,20,30,40
```

10.将字符串以下划线分割

下面的代码会将给定的字符串以下划线分割。

```go
import (
  "fmt"
  "strings"
  "regexp"
)


func main() {
   snakeCase := ConvertToSnakeCase("ILikeProgrammingINGo123")
   fmt.Println("String in snake case : ", snakeCase)
}

func ConvertToSnakeCase(input string) string {
  var matchChars = regexp.MustCompile("(.)([A-Z][a-z]+)")
  var matchAlpha = regexp.MustCompile("([a-z0-9])([A-Z])")

  snake := matchChars.ReplaceAllString(input, "${1}_${2}")
  snake = matchAlpha.ReplaceAllString(snake, "${1}_${2}")
  return strings.ToLower(snake)
}


output:
String in snake case :  i_like_programming_in_go123
```



via: https://blog.canopas.com/golang-utility-functions-you-always-missed-ebeabae6b276 >

作者：Nidhi D