## Go In Action

<https://pepa.holla.cz/wp-content/uploads/2016/10/Go-in-Action.pdf>

1. INTRODUCING GO
1.1. Solving modern programming challenges with Go
1.1.1. Development Speed
1.1.2. Concurrency
1.1.3. Go’s type system
1.1.4. Memory management
1.2. Hello, Go
1.2.1. Introducing the Go Playground
1.3. Summary
2. GO QUICK START
2.1. Program architecture
2.2. Main package
2.3. Search package
2.3.1. search.go
2.3.2. feed.go
2.3.3. match.go/default.go
2.4. RSS matcher
2.5. Summary
3. PACKAGING AND TOOLING
3.1. Packages
3.1.1. Package-naming conventions
3.1.2. Package main
3.2. Imports
3.2.1. Remote imports
3.2.2. Named imports
3.3. Init
3.4. Using Go tools
3.5. Going further with Go developer tools
3.5.1. Go vet
3.5.2. Go format
3.5.3. Go documentation
3.6. Collaborating with other Go developers
3.6.1. Creating repositories for sharing
3.7. Dependency management
3.7.1. Vendoring dependencies
3.7.2. Introducing gb
3.8. Summary
4. ARRAYS, SLICES AND MAPS
4.1. Array internals and fundamentals
4.1.1. Internals
4.1.2. Declaring and initializing
4.1.3. Working with arrays
4.1.4. Multidimensional arrays
4.1.5. Passing arrays between functions
4.2. Slice internals and fundamentals
4.2.1. Internals
4.2.2. Creating and initializing
4.2.3. Working with slices
4.2.4. Multidimensional slices
4.2.5. Passing slices between functions
4.3. Map internals and fundamentals
4.3.1. Internals
4.3.2. Creating and initializing
4.3.3. Working with maps
4.3.4. Passing maps between functions
4.4. Summary
5. GO’S TYPE SYSTEM
5.1. User-defined types
5.2. Methods
5.3. The nature of types
5.3.1. Built-in types
5.3.2. Reference types
5.3.3. Struct types
5.4. Interfaces
5.4.1. Standard library
5.4.2. Implementation
5.4.3. Method sets
5.4.4. Polymorphism
5.5. Type embedding
5.6. Exporting and unexporting identifiers
5.7. Summary
6. CONCURRENCY
6.1. Concurrency versus parallelism
6.2. Goroutines
6.3. Race conditions
6.4. Locking shared resources
6.4.1. Atomic functions
6.4.2. Mutexes
6.5. Channels
6.5.1. Unbuffered channels
6.5.2. Buffered channels
6.6. Summary
7. CONCURRENCY PATTERNS
7.1. Runner
7.2. Pooling
7.3. Work
7.4. Summary
8. STANDARD LIBRARY
8.1. Documentation and Source Code
8.2. Logging
8.2.1. Log Package
8.2.2. Customized Loggers
8.2.3. Conclusion
8.3. Encoding/Decoding
8.3.1. Decoding JSON
8.3.2. Encoding JSON
8.3.3. Conclusion
8.4. Input and Output
8.4.1. Writer and Reader Interfaces
8.4.2. Working Together
8.4.3. Simple Curl
8.4.4. Conclusion
8.5. Summary
9. TESTING AND BENCHMARKING
9.1. Unit testing
9.1.1. Basic unit test
9.1.2. Table tests
9.1.3. Mocking calls
9.1.4. Testing endpoints
9.2. Examples
9.3. Benchmarking
9.4. Summary

## Go语言设计与实现

<https://draveness.me/golang/>

第一部分 预备知识
第一章 准备工作
1.1 调试源代码
第二章 编译原理
2.1 编译过程
2.2 词法分析和语法分析
2.3 类型检查
2.4 中间代码生成
2.5 机器码生成
第二部分 基础知识
第三章 数据结构
3.1 数组
3.2 切片
3.3 哈希表
3.4 字符串
第四章 语言基础
4.1 函数调用
4.2 接口
4.3 反射
第五章 常用关键字
5.1 for 和 range
5.2 select
5.3 defer
5.4 panic 和 recover
5.5 make 和 new
第三部分 运行时
第六章 并发编程
6.1 上下文 Context
6.2 同步原语与锁
6.3 定时器
6.4 Channel
6.5 调度器
6.6 网络轮询器
6.7 系统监控
第七章 内存管理
7.1 内存分配器
7.2 垃圾收集器
7.3 栈内存管理
第四部分 进阶内容
第八章 元编程
8.1 插件系统
8.2 代码生成
第九章 标准库
9.1 JSON
9.2 HTTP
9.3 数据库

## Go101

<https://go101.org/article/101.html>

<https://github.com/jxs1211/go101>

<https://go101.org/>

### fundermental

Become Familiar With Go Code
Introduction of Source Code Elements
Keywords and Identifiers
Basic Types and Their Value Literals
Constants and Variables - also introduces untyped values and type deductions.
Common Operators - also introduces more type deduction rules.
Function Declarations and Calls
Code Packages and Package Imports
Expressions, Statements and Simple Statements
Basic Control Flows
Goroutines, Deferred Function Calls and Panic/Recover
Go Type System
Go Type System Overview - a must read to master Go programming.
Pointers
Structs
Value Parts - to gain a deeper understanding into Go values.
Arrays, Slices and Maps - first-class citizen container types.
Strings
Functions - function types and values, including variadic functions.
Channels - the Go way to do concurrency synchronizations.
Methods
Interfaces - value boxes used to do reflection and polymorphism.
Type Embedding - type extension in the Go way.
Type-Unsafe Pointers
Generics - use and read composite types
Reflections - the reflect standard package.
Some Special Topics
Line Break Rules
More About Deferred Function Calls
Some Panic/Recover Use Cases
Explain Panic/Recover Mechanism in Detail - also explains exiting phases of function calls.
Code Blocks and Identifier Scopes
Expression Evaluation Orders
Value Copy Costs in Go
Bounds Check Elimination
Concurrent Programming
Concurrency Synchronization Overview
Channel Use Cases
How to Gracefully Close Channels
Other Concurrency Synchronization Techniques - the sync standard package.
Atomic Operations - the sync/atomic standard package.
Memory Order Guarantees in Go
Common Concurrent Programming Mistakes
Memory Related
Memory Blocks
Memory Layouts
Memory Leaking Scenarios
Some Summaries
Some Simple Summaries
nil in Go
Value Conversion, Assignment and Comparison Rules
Syntax/Semantics Exceptions

### Go Details 101

### Go FAQ 101

### Go Tips 101

## practical-go-lessons

<https://www.practical-go-lessons.com/>

Chap. 1: Programming A Computer
Chap. 2: The Go Language
Chap. 3: The terminal
Chap. 4: Setup your dev environment
Chap. 5: First Go Application
Chap. 6: Binary and Decimal
Chap. 7: Hexadecimal, octal, ASCII, UTF8, Unicode, Runes
Chap. 8: Variables, constants and basic types
Chap. 9: Control Statements
Chap. 10: Functions
Chap. 11: Packages and imports
Chap. 12: Package Initialization
Chap. 13: Types
Chap. 14: Methods
Chap. 15: Pointer type
Chap. 16: Interfaces
Chap. 17: Go modules
Chap. 18: Go Module Proxies
Chap. 19: Unit Tests
Chap. 20: Arrays
Chap. 21: Slices
Chap. 22: Maps
Chap. 23: Errors
Chap. 24: Anonymous functions & closures
Chap. 25: JSON and XML
Chap. 26: Basic HTTP Server
Chap. 27: Enum, Iota & Bitmask
Chap. 28: Dates and time
Chap. 29: Data storage : files and databases
Chap. 30: Concurrency
Chap. 31: Logging
Chap. 32: Templates
Chap. 33: Application Configuration
Chap. 34: Benchmarks
Chap. 35: Build an HTTP Client
Chap. 36: Program Profiling
Chap. 37: Context
Chap. 38: Generics
Chap. 39: An object oriented programming language ?
Chap. 40: Upgrading or Downgrading Go
Chap. 41: Design Recommendations
Chap. 42: Cheatsheet

## gobyexample

<https://gobyexample.com/>

## GoLangBooks

git@github.com:diptomondal007/GoLangBooks.git

## go tool

<https://pkg.go.dev/cmd/go>
