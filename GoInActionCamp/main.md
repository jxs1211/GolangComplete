# Go实战训练营

![总链接](https://drive.weixin.qq.com/s?k=ACQADQdwAAohB4R0DS#/)

## 预习：Golang 基础语法和 Web 框架起步

- 1.Go 基本语法和 Web 框架起步  
- 2.type 定义与 Server 抽象
- 3.错误处理与简单路由树实现
- 4.并发编程、文件操作与泛型

代码地址：flycash/toy-web: 用于极客时间 go 基础课程 (github.com)

## 模块一：Go 进阶语法

- 第一讲：并发编程·context包
    多个key，value创建需要多次调用context.WithValue()
- 第二讲：并发编程·sync包与channel
    <!-- 如果有一个资源是并发安全的，对他的所有操作加上锁后，才能确保万无一失 -->
    <!-- 太多读写是否导致性能问题 -->
    <!-- 记住锁的状态转换流程图，只有cas和enqueue两个部分 -->
    <!-- 什么情况下cpu会100%，goroutine自旋（spin）的情况下就会100% -->
    ```go
        type SafeMap struct {
            m  map[string]string
            mu sync.RWMutex
        }

        func (m *SafeMap) Get(key string) string {
            // 只读需要枷锁，因为有其他goroutine在同时写
            m.mu.RLock()
            defer m.mu.RUnlock()
            return m.m[key]
        }

        func (m *SafeMap) Set(key, val string) {
            m.mu.Lock()
            defer m.mu.Unlock()
            m.m[key] = val
        }
    ```
- 第一周作业：服务器优雅退出
- 第三讲-1：并发编程·channel与反射
- 第三讲-2：并发编程·channel与反射
- 第二周作业：生成 INSERT 语句
 https://gitee.com/geektime-geekbang/geektime-go/tree/master/advance/reflect/homework
 
