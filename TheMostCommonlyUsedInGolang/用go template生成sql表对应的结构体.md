# 如何根据mysql表生成结构体｜一个小工具的探索之旅

Original 陪计算机走过漫长岁月 [陪计算机走过漫长岁月](javascript:void(0);) *2022-07-05 09:00* *Posted on 北京*

收录于合集#go1个

## 1. 目录

![Image](https://mmbiz.qpic.cn/mmbiz_png/BtqZzg1Lt0xH8wRbv1oziaOKdC6ULuyESR9xAqmdPCk7EvPaW2tuVYkoQJDu2pHS6WXZtdsk95gwnYsElrPicy2w/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

## 2. 背景

​最近在工作中会有根据mysql表在go中编写一个对应的结构体这样的coding，虽然数据表并不是复杂，字段不是很多，代码写起来也比较快，为了快速的完成工作我一开始就是按照数据表的列一个接着一个的来写。但我是个懒人，重复的工作希望可以通过代码帮我完成，因为后面也有类似的工作，如果我有对应的代码生成工具会方便很多，并且用自己做出来的工具内心中或多或少会有一些成就感。所以我心生一个想法，为什么我不搞一个简单的工具，来根据表的结构生成结构体呢？所以我就研究了一下，看了一些资料，包括mysql的information_schame数据库， sqlx， go标准库里的text/template等内容 ，特此写下文章分享给读者朋友们，希望读者朋友们有所收获。话不多说，让我们开始本次的探索之旅。

## 3. 怎么找到mysql的表结构信息

​在安装mysql的时候，会发现除了自己创建的数据库之外，还会有一些别的数据库是默认给你创建好的。我们可以登陆mysql使用下面语句查看。

```sql
mysql> SHOW DATABASES;
+--------------------+
| Database           |
+--------------------+
| elliot_test        |
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
5 rows in set (0.01 sec)
```

大家可以看下我本地的mysql，其中elliot_test是我创建的数据库，其他的information_schema, mysql, performance_schema, sys（5.7 以上的叫sys，5.6的自带的是test）,这些是mysql自带的。那么这些数据库有什么用呢，记录的都是什么信息呢？

1. **information_schema**：保存了MySQl服务所有数据库的信息。具体MySQL服务有多少个数据库，各个数据库有哪些表，各个表中的字段是什么数据类型，各个表中有哪些索引，各个数据库要什么权限才能访问。
2. **mysql**：保存MySQL的权限、参数、对象和状态信息。如哪些user可以访问这个数据，DB的参数。
3. **performance_schema**：主要用于收集数据库服务器性能参数，提供进程等待的详细信息， 保存历史的事件汇总信息，为提供MySQL服务器性能做出详细的判断。
4. **test**：5.6自带，没有什么东西。
5. **sys**：Sys库所有的数据源来自：performance_schema。目标是把performance_schema的把复杂度降低。

​在了解了mysql自带数据库的功能之后，这时候我们就知道了要查看mysql的表结构信息其实我们只需要在information_schema这里面找就好了，因为这里包含了所有数据库的信息，包括有哪些表，什么表有什么字段，这正是我们需要的。由于information_schema里面表比较多，这里就不作展示并且一一介绍了，感兴趣的读者朋友们可以登陆mysql使用use information_schema; 命令切换数据库，然后使用show tables；命令查看这下面有什么表。这里就介绍几个比较重要的，或者说我们可能会用到的。

1. **SCHEMATA表**：提供了当前mysql实例中所有数据库的信息。是show databases的结果取之此表。
2. **TABLES表**：提供了关于数据库中的表的信息（包括视图）。详细表述了某个表属于哪个schema，表类型，表引擎，创建时间等信息。是show tables from schemaname的结果取之此表。
3. **COLUMNS表**：提供了表中的列信息。详细表述了某张表的所有列以及每个列的信息。是show columns from schemaname.tablename的结果取之此表。

​我们这里关注的是COLUMNS这张表，这张表记录的是每张表的列信息，我们使用desc命令看下这张表都有些什么。

```sql
mysql> desc columns;
+--------------------------+---------------------+------+-----+---------+-------+
| Field                    | Type                | Null | Key | Default | Extra |
+--------------------------+---------------------+------+-----+---------+-------+
| TABLE_CATALOG            | varchar(512)        | NO   |     |         |       |
| TABLE_SCHEMA             | varchar(64)         | NO   |     |         |       |
| TABLE_NAME               | varchar(64)         | NO   |     |         |       |
| COLUMN_NAME              | varchar(64)         | NO   |     |         |       |
| ORDINAL_POSITION         | bigint(21) unsigned | NO   |     | 0       |       |
| COLUMN_DEFAULT           | longtext            | YES  |     | NULL    |       |
| IS_NULLABLE              | varchar(3)          | NO   |     |         |       |
| DATA_TYPE                | varchar(64)         | NO   |     |         |       |
| CHARACTER_MAXIMUM_LENGTH | bigint(21) unsigned | YES  |     | NULL    |       |
| CHARACTER_OCTET_LENGTH   | bigint(21) unsigned | YES  |     | NULL    |       |
| NUMERIC_PRECISION        | bigint(21) unsigned | YES  |     | NULL    |       |
| NUMERIC_SCALE            | bigint(21) unsigned | YES  |     | NULL    |       |
| DATETIME_PRECISION       | bigint(21) unsigned | YES  |     | NULL    |       |
| CHARACTER_SET_NAME       | varchar(32)         | YES  |     | NULL    |       |
| COLLATION_NAME           | varchar(32)         | YES  |     | NULL    |       |
| COLUMN_TYPE              | longtext            | NO   |     | NULL    |       |
| COLUMN_KEY               | varchar(3)          | NO   |     |         |       |
| EXTRA                    | varchar(30)         | NO   |     |         |       |
| PRIVILEGES               | varchar(80)         | NO   |     |         |       |
| COLUMN_COMMENT           | varchar(1024)       | NO   |     |         |       |
| GENERATION_EXPRESSION    | longtext            | NO   |     | NULL    |       |
+--------------------------+---------------------+------+-----+---------+-------+
21 rows in set (0.00 sec)
```

​可以看到里面记录的东西还是很多的，如果聚焦于我们的需求：根据表结构构建结构体。那我们的关注点只需要知道他的列名（COLUNM_NAME）,以及对应的数据类型（DATA_TYPE）就可以了。所以我们很容易写出下面这条sql查询到我们需要的信息。

```sql
mysql> select column_name, data_type from columns where table_name='benchmark' and table_schema='bitstorm'; 
+-------------------+-----------+
| column_name       | data_type |
+-------------------+-----------+
| id                | int       |
| space_id          | char      |
| name              | varchar   |
| description       | varchar   |
| benchmark_version | int       |
| is_default        | tinyint   |
| creator           | varchar   |
| created_time      | datetime  |
| last_updated_user | varchar   |
| last_updated_time | datetime  |
| is_runtime        | tinyint   |
+-------------------+-----------+
11 rows in set (0.00 sec)
```

不过我们也可以用下图这种方式去查看。

```sql
mysql> show columns from bitstorm.benchmark;
+-------------------+---------------+------+-----+-------------------+----------------+
| Field             | Type          | Null | Key | Default           | Extra          |
+-------------------+---------------+------+-----+-------------------+----------------+
| id                | int(11)       | NO   | PRI | NULL              | auto_increment |
| space_id          | char(36)      | NO   | MUL | NULL              |                |
| name              | varchar(255)  | NO   |     | NULL              |                |
| description       | varchar(1024) | NO   |     | NULL              |                |
| benchmark_version | int(11)       | NO   |     | NULL              |                |
| is_default        | tinyint(4)    | NO   |     | NULL              |                |
| creator           | varchar(255)  | NO   |     | system            |                |
| created_time      | datetime      | NO   |     | CURRENT_TIMESTAMP |                |
| last_updated_user | varchar(255)  | NO   |     | system            |                |
| last_updated_time | datetime      | NO   |     | CURRENT_TIMESTAMP |                |
| is_runtime        | tinyint(4)    | NO   |     | 0                 |                |
+-------------------+---------------+------+-----+-------------------+----------------+
11 rows in set (0.00 sec)

mysql> 
```

​获得了这些信息之后，就可以开始我们的编码工作了。

## 4. 技术实现

​在这一章节会探讨如何实现这个需求，会探讨如何简单实现一个mysql的client，也就是一个简单的mysql驱动，当然并不是说在这里我要实现这个东西，只是探讨一下如果要写一个的话大概需要怎么做。另外会讲到实现这个需求用到的一些主要技术，包括sqlx， go标准库的text/template，还有实现这个需求的核心逻辑讲解。

### 4.1 以何种方式与数据库交互

​其实这里一般是没有讨论的必要的，选择一个开源的库去和执行上面提到的获取数据表信息的sql语句就好了。不过如果我想把这个工具当作一个开源的项目去做，可能会考虑为了轻量化而尽量减少这个项目的依赖，以及开源的mysql驱动库对于这个项目来说会提供一些我们本身并不需要的功能。不过在这里我并不打算要自己动手实现一个简单的查询语句处理的工具，因为想要快速的完成这个小工具，如果花费大量的时间在实现别的东西上面，那么我工作就失去了焦点，这并不是我想要的。不过之前看过不少go-mysql这个开源库的源码，对数据库驱动库是怎么实现的还是有一定了解的，实际上实现起来也并不是很难。这里的话可以和读者朋友们分享一下如果我要手动实现的话，我要怎么做。

​这里思路可以简单的概括一下，就是我要吧自己模仿成mysql的client端，只要我们遵循mysql的通讯协议就可以了。就好像go的结构与其实现类，实现类只要实现了接口的方法就可以看作是这个接口的实现类。而我们只需要遵循mysqld的client端的一些通讯协议，mysql自然也会把我们看成是一个client端，现在主流的CDC组件一般的做法就是把自己伪装成mysql的从节点去获取mysql的数据变更，其实也是一样的道理。那么要把自己变成mysql的client端我们都需要遵循哪些协议呢？

​首先在TCP层面的客户端与服务端三次握手建立了TCP连接之后，mysql的服务端会主动发送一个握手初始化的包给客户端，这里的主要内容是服务器的一些信息，告诉客户端要遵循什么什么协议，打个比方说如果mysql 的binlog协议就有好几个版本，mysql 5.15之前是v0，5.6之前是v1，5.6之后是v2版本，握手的作用就是告诉客户端一些服务端的信息，后面的通讯要按照一些规范去进行。在握手初始化之后要进行认证登陆的过程，这里客户端会把账户和密码给服务端，让服务端去验证，最后服务器把认证结果发送回来，要是成功了，后面就可以开始执行我们的发送的sql语句了。在这个过程中涉及到的协议内容和格式我贴在下方参考资料中（connect_phase packet）

![Image](https://mmbiz.qpic.cn/mmbiz_png/BtqZzg1Lt0xH8wRbv1oziaOKdC6ULuyESCM3Zlo7lziar1cQup5hCibLTJBSlR1NKFK7RNfpL6tGZw9rQqJBat23g/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

​在执行sql语句这一段，我们这个需求其实只是执行一条query语句，所以按照COM_QUERY协议从客户端构建一个数据包，然后等服务端返回的时候按照ResultSet协议去解析返回的内容就好了。

​这一系列流程我是在go-mysql这个开源库上看到的，之前学习过这个库，详细看了不少代码的实现，我把代码的位置放在参考资料上，感兴趣的朋友可以看看～看的过程中如果有不太明白的地方可以留言联系我，我们可以交流交流。

​这里就不打算把所有的协议的格式都展开讲了，不过我会在下方参考资料那边贴上mysql的官方文档，里面会有对协议内容详细的介绍。这里我想要说的是一个实现的大致流程。可以看到其实这里并不是说很难，不过实现起来确实也是需要时间的。既然讲到这里了，说句题外话，既然实现了协议就可以被当作是mysql的客户端，那么我可不可以实现一些协议，被当作一个mysql的服务端呢？当然也是可以的，之前我也做过类似的尝试。所以说很多时候网络确实不是百分百安全的，对方只需要满足一些行为就可以获取信任，但是有时候我们并不能完全确认对方是不是真的值得信任。

### 4.2 代码实现

其实实现的流程比较简单，没用多长时间就写完了，这里主要记录一些使用到的技术和一些值得注意或者值得分享的地方。

#### 4.2.1 与数据库的交互--sqlx

​我主要用了sqlx这个库和数据库做交互，感觉这个库写的还不错，使用起来也很方便，github地址贴在这里：https://github.com/jmoiron/sqlx，感兴趣的朋友可以看看readme里面对于它使用方法的详细介绍，这里稍微贴一段代码吧。

```go
type Person struct {
    FirstName string `db:"first_name"`
    LastName  string `db:"last_name"`
    Email     string
}

 // Query the database, storing results in a []Person (wrapped in []interface{})
people := []Person{}
db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
```

不过说实在的，就我个人而言，其实感觉orm框架的用法都大同小异，用哪个其实不是问题的关键。

#### 4.2.2 如何生成代码 -- go text/template

一开始是想着用字符串拼接去做的，但是感觉这样子写出来的代码有点点丑。其实丑是一方面，另外一方面是，因为不同orm框架生成的struct可能是不一样的，最明显的就是tag里面的标签是不一样的，如果后续要拓展的话，难道判断生成什么orm框架下面的struct这样的逻辑要嵌入到字符串拼接里面吗？这显然不是一个好的办法。所以我想到了用template模版去解析的方式去做。

​go标准库中提供了text/template这个包，实现了数据驱动的用于生成文本输出的模板，其实这个很像前端的mvvm那一套，我们定义好模版，然后传入参数，他就会解析模版，将我们传入的参数替换到模版对应的位置，从而生成我们想要的文本。在模版里面还可以写逻辑呢，比如简单的逻辑判断，循环遍历之类的。话不多说，这里展示一个简单的demo。关于text/template详细的用法和讲解文章，我贴到参考资料中，感兴趣的朋友可以看看。

```go
package main 
import(
    "os"
    "text/template"
)

type Inventory struct {
    Material string
    Count    uint
}

func main() {
    sweaters := Inventory{"wool", 17}
    tmpl, err := template.New("test").Parse("{{.Count}} of {{.Material}}\n")//{{.Count}}获取的是struct对象中的Count字段的值
    if err != nil { panic(err) }
    err = tmpl.Execute(os.Stdout, sweaters)//返回 17 of wool
    if err != nil { panic(err) }
}
```

这个例子是将Inventory的对象传进模版中做解析，其中{{.Count}},{{.Material}}就是获取这个对象的属性，然后替换。

好了，讲到这里实现这个需求的所有技术，所有基础知识我们都知道了，下面我们看看核心逻辑。

#### 4.2.3 核心逻辑

核心的流程如下：

1. **获取列数据**。包括列的名字以及数据类型。
2. **数据类型转化**，因为mysql和go的类型其实不大一样的，如果我们要生成go的struct就需要做一个mysql数据类型和go数据类型的转化，我这里的做法比较粗糙，写了一个map配置，key是mysql数据类型，value是go数据类型。直接拿就完事了。但是其实这样并不是很合理，比如mysql的TINYINT类型对应的是go的int8，但是在go中如果用int32，int64，去表示可不可以呢？其实是可以的。这里后续可以作为一个优化点，或者这个配置的能力向用户开放更好。
3. **执行模版引擎**。会提前写好一个模版文件，然后用拿到的数据去解析。不过这个模版文件是可配置的，这样可以提供比较灵活的方式去生成自己想要的代码。比如我现在写的这个模版文件比较淳朴，除了生成struct之外就什么都没有了，但是在一些环境下，可能用户会想着我生成某某个接口的实现类或者有一些特定的注释，那他可以自定义模版文件去生成，不过解析模版文件需要对应的数据，这里数据的结构体还是在我的掌控之内。所以后面可以考虑把这个数据也开放出去，这样就比较完美了。
4. **格式化生成的文件**。因为模版生成出来的文件不是特别好看。所以我这里手动执行了一些go fmt去格式化代码，看起来会舒服很多。

```go
func (g *Generator) Gen(config *GenInfo) (isSuccess bool, err error) {
   if err := checkGenInfo(config); err != nil {
      return false, err
   }
   tableInfos, err := g.executeQuery(config.Schema, config.Table)
   if err != nil {
      return false, err
   }
   templateMetaDatas, err := convertTableInfoToMeta(tableInfos)
   templateData := &TemplateData{
      PackageName: config.PackageName,
      StructName:  config.StructName,
      Meta:        templateMetaDatas,
   }
   var genPath string
   if config.ExportFolder == "" {
      genPath = config.FileName
   } else {
      genPath = fmt.Sprintf("%s/%s", config.ExportFolder, config.FileName)
   }
   isSuccess, err = genCodeByTemplate(genPath, config.TemplatePath, templateData)
   if err == nil && isSuccess {
      _, _ = exec.Command("go", "fmt", genPath).Output()
   }
   return isSuccess, err
}
```

template文件：

```go
package {{.PackageName}}

type {{.StructName}} struct {
    {{- range $i, $v := .Meta }}
        {{$v.CamelName}} {{$v.DataTypeInGo}}
    {{- end }}
}
```

测试：
```go
func TestGenerator_Gen(t *testing.T) {
  config := &Config{
    Host:     "127.0.0.1",
    Port:     3306,
    Username: "root",
    Password: "", // 我才不告诉你我的密码呢。
  }
  g, err := NewGenerator(config)
  if err != nil {
    t.Errorf("have err during NewGenerator, err is %s", err)
    return
  }
  genInfo := &GenInfo{
    Schema:       "elliot_test",
    Table:        "test_table",
    ExportFolder: "",
    TemplatePath: "struct_gen_test_template",
    FileName:     "test_gen.go",
    PackageName:  "table_gen",
    StructName:   "TestGenStruct",
  }
  isSuccess, err := g.Gen(genInfo)
  if err != nil {
    t.Errorf("have err during Gen file, err is %s", err)
    return
  }
  t.Logf("does it gen file successully?  %v", isSuccess)
}
```

结果：

```go
package table_gen

type TestGenStruct struct {
  Name string
  Age  int
}
```

完美，这就达到了我们最初的梦想。这里就不对细节上的东西做过多的介绍了，感觉该讲的东西大部分都讲了。详细的code我已经开源到我的github上了。地址：https://github.com/elliotchenzichang/GenStructByTable。感兴趣的小伙伴可以看看。有问题或者发现代码中的错误，可以留言联系我，互相交流学习。

## 5. 总结

​做这个东西是一时兴起，做完之后还是有很多地方感到不足和可优化，后面有时间可以慢慢优化一下。就我个人而言有时间的话还是比较喜欢倒腾一些小东西。一方面做成一个东西的过程中探索求知这个过程是很快乐的，另一方面做成之后的喜悦就像攀登了一座座小山峰最终到达目标，那瞬间扑面而来的快乐是人生少有的。好了，今天这个探索的故事就讲到这里了，感谢各位读者朋友们赏脸读到末尾处hhh。

## 6. 参考资料

1. mysql information_schema 详解：https://zhuanlan.zhihu.com/p/88342863
2. mysql connection_phase packet :https://dev.mysql.com/doc/internals/en/connection-phase.html
3. mysql COM_QUERY packet: https://dev.mysql.com/doc/internals/en/com-query.html
4. mysql COM_QUERY Response :https://dev.mysql.com/doc/internals/en/com-query-response.html
5. go text/template :https://www.cnblogs.com/wanghui-garcia/p/10385062.html
6. go-mysql client模块代码：https://github.com/go-mysql-org/go-mysql/tree/master/client