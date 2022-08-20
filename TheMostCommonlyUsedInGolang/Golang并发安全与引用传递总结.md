导语 | 因为现在服务上云的趋势，业务代码都纷纷转向golang的技术栈。在迁移或使用的过程中，由于对golang特性的生疏经常会遇到一些问题，本文总结了golang并发安全和参数引用传值时的一些知识。



**一、Map类型并发读写引发Fatal Error**

先看一个在Go中关于Map类型并发读写的经典例子：

```go
var testMap  = map[string]string{}
func main() {
   go func() {
      for{
         _ = testMap["bar"]
      }
   }()
   go func() {
      for  {
         testMap["bar"] = "foo"
      }
   }()
   select{}
}
```

以上例子会引发一个Fatal error：

**fatal error: concurrent map read and map write**

产生这个错误的原因就是在Go中Map类型并不是并发安全的，出于安全的考虑，此时会引发一个致命错误以保证程序不出现数据的混乱。

**二、Go如何检测Map并发异常**

在Go源码map.go中，可以看到以下flags：

```go
// flags
iterator     = 1 // there may be an iterator using buckets
oldIterator  = 2 // there may be an iterator using oldbuckets
hashWriting  = 4 // a goroutine is writing to the map
sameSizeGrow = 8 // the current map growth is to a new map of the same size
```

在源码中mapaccess1、mapaccess2都用于查询mapassign和mapdelete用于修改。

对于查询操作，大致检查并发错误的流程如下：在查询前检查并发flag是否存在，如果存在就抛出异常。

```go
if h.flags&hashWriting != 0 {
    throw("concurrent map read and map write")
}
```

对于修改操作则如下：

- 写入前检查一次标记位，通过后打上标记。

- 写入完成再次检查标记位，通过后还原标记。

```go
 //各类前置操作
   ....
   if h.flags&amp;hashWriting != 0 {
      //检查是否存在并发
      throw("concurrent map writes")
   }

   //赋值标记位
   h.flags ^= hashWriting
   ....
   //后续操作
  done:
   //完成修改后，再次检查标记位
   if h.flags&hashWriting == 0 {
      throw("concurrent map writes")
   }
   //还原标记位取消hashWriting标记
   h.flags &^= hashWriting
```

**三、如何避免Map的并发问题**

go官方认为因为Map并发的问题在实际开发中并不常见，如果把Map原生设计成并发安全的会带来巨大的性能开销。因此需要使用额外方式来实现。

**（一）自行使用锁和map来解决并发问题**

参考如下：

```go

type cocurrentMap = struct {
   sync.RWMutex
   m map[string]string
}

func main() {
   var testMap = &cocurrentMap{m:make(map[string]string)}
   //写
   testMap.Lock()
   testMap.m["a"] = "foo"
   testMap.Unlock()
   //读
   testMap.RLock()
   fmt.Println(testMap.m["a"])
   testMap.RUnlock()
}
```

这个方法存在问题就是并发量巨大的时候，锁的竞争也会带来巨量消耗，性能一般。

**（二）使用sync.Map**

sync.Map通过巧妙的设计来提高并发安全下Map的性能，其设计思路是通过空间换时间来实现的，同时维护2份数据，read&dirty。read主要用来避免读写冲突。

其数据结构如下：

```go

type Map struct {
   mu Mutex //锁
   read atomic.Value //readOnly
   dirty map[interface{}]*entry //*entry
   misses int
}

type readOnly struct {
   m       map[interface{}]*entry
   amended bool // true if the dirty map contains some key not in m.
}

type entry struct {
   p unsafe.Pointer // *interface{}
}
```

使用示例如下：

```go
var m sync.Map
// 写
m.Store("test", 1)
m.Store(1, true)

// 读
val1, _ := m.Load("test")
val2, _ := m.Load(1)
fmt.Println(val1.(int))
fmt.Println(val2.(bool))

//遍历
m.Range(func(key, value interface{}) bool {
   //....
   return true
})

//删除
m.Delete("test")

//读取或写入
m.LoadOrStore("test", 1)
```

这里对sync.Map的原理不做深入展开，只提几点特性：

- read和dirty是共享内存的，尽量减少冗余内存的开销。

- read是原子性的，可以并发读，写需要加锁。

- 读的时候先read中取，如果没有则会尝试去dirty中读取（需要有标记位readOnly.amended配合）

- dirty就是原生Map类型，需要配合各类锁读写。

- 当read中miss次数等于dirty长度时，dirty会提升为read，并且清理已经删除的k-v（延迟更新，具体如何清理需要enrty中的p标记位配合）

- 双检查（在加锁后会再次对值检查一遍是否依然符合条件）

- sync.Map适用于读多写少的场景。

- sync.Map没有提供获取长度size的方法，需要通过遍历来计算。



**四、切片类型Slice是并发安全的吗**

与Map一样，Slice也不是并发安全的：

```go
var testSlice []int
func main() {
   for i:=0; i<1000; i++ {
      go func() {
         testSlice = append(testSlice, i)
      }()
   }
   for idx, val := range testSlice {
      fmt.Printf("idx:%d val:%d\n", idx, val)
   }
}
```

可以看到输出如下：

........

**idx:901 val:999**

**idx:902 val:999**

.........

但是在切片中并不会引发panic，如果程序无意中对切片使用了并发读写，严重的话会导致获取的数据和之后存储的数据错乱，所以这里要格外小心，可以通过加锁来避免。

**五、Map、Slice作为参数传递的问题**

切片除了并发有问题外，当他作为参数传递的时候，也会导致意料之外的问题，Go官方说明在Go中所有的传递都是值传递，没有引用传递的问题，但是在实际使用时，切片偶尔会引起一些疑惑，例如以下情况：

```go

func changeVal(testSlice []string, idx int, val string){
   testSlice[idx] = val
}

func main() {
   var testSlice []string
   testSlice = make([]string, 5)
   testSlice[0] = "foo"
   changeVal(testSlice, 0, "bar")
   fmt.Println(testSlice[0])
}
```

以上代码执行后可以看到打印出的值为：

**bar**

这里就奇怪了，如果按照Go官方说明在该语言中传递都是值传递的话，为什么函数内修改切片会导致原切片也一起修改呢？这里要分2个问题来看：

- Go只会对基础值类型在传参中使用深拷贝，实际上对于Slice和Map类型，使用的是浅拷贝，Slice作为传参，其指向的内存地址依然是原数据。

- Slice扩容机制的影响：向Slice中添加元素超出容量的时候，我们知道会触发扩容机制，而扩容机制会创建一份新的【原数据】此时，它与浅拷贝获取到的变量是没有任何关联的。

可以通过以下代码验证，我们故意构造触发扩容的场景：

```go
func appendVal(testSlice []string, val string){
   fmt.Printf("testSlice:%p\n", testSlice)
   testSlice = append(testSlice, "addCap") //触发了扩容机制
   fmt.Printf("after append testSlice:%p\n", testSlice)
   testSlice[0] = val
}

func main() {
   var testSlice []string
   testSlice = make([]string, 5)
   testSlice[0] = "foo"
   appendVal(testSlice,"bar")
   fmt.Println(testSlice[0]) //此时打印出的值为foo
}
```

可以看到控制台打印如下：

**testSlice:0xc00005a050**

**after append testSlice:0xc0000700a0**

**foo**

此时因为扩容的影响导致原切片和传递后的切片不再有关联，因此打印值回到了最初的原数据foo

除了扩容机制外，我们也可以利用go中的copy函数来强制深拷贝：

```go
var newTestSlice []string
newTestSlice = make([]string, len(testSlice))
copy(newTestSlice, testSlice)
fmt.Printf("testSlice:%p\n", testSlice)
fmt.Printf("newTestSlice:%p\n", newTestSlice)
```

**testSlice:0xc0000d6000**

**newTestSlice:0xc0000d6050**

另外对于数组类型，如果无意中转换为切片时，也极容易导致这种不确定性发生。切片作为参数传递时，在函数内对切片进行修改，需要时刻注意。

回过头再来看Map就一目了然了，因为Map的操作对象一直是引用，其即使扩容后，引用的地址不会改变，所以不会出现时而可以修改，时而不能修改的情况：

```go
func changeMap(testMap map[string]string, k string, v string){
   testMap[k] = v
}

func main() {
   var testMap map[string]string
   testMap = make(map[string]string)
   testMap["foo"] = "bar"
   changeMap(testMap, "foo", "rab")
   fmt.Println(testMap)
}
```

**输出：map[foo:rab]**

可以看到函数内修改了原参数的值。



**六、总结**

Go因为其简洁的语法和高效的性能在当今微服务领域笑傲江湖，但是其本身语言特性在使用时，也会带来不少坑，本文总结了并发场景和参数传递时容易引发的问题，从而注意避免这些情况的发生。