# 

利用反射，将一个结构体实例转  `sertStmt(User{})` 有的单元测试。在这个作业里面，你只需要考虑生成一个能在 MySQL 上执行的语句。我们已经预先定义好了方法：

```go
func InsertStmt(entity interface{}) (string, []interface{}, error) { }
```
为了减轻一些工作量，也为了便于维护，所以输入你只需要考虑有限的几种情况：
* nil
* 结构体
* 一级指针，但是指针指向的是结构体
* 组合，但是不能是指针形态的组合
实际上，在考虑到执行 INSERT 语句，数据库驱动只支持特定的一些类型，包括：

* 基本类型
* string, []byte 和 time.Time
* 实现了 driver.Valuer 和 sql.Scanner 两个接口的类型
具体可以参考测试用例。

从实现上来说，就是简单的遍历所有的字段，并且检查每个字段：

* 如果字段是普通的字段，那么直接进行处理，也就是认为它是一个列。它的名字作为列名，而对应的值作为插入时候的参数；
* 如果字段其实代表的是组合。注意，在 Go 里面，所谓组合其实就是一个比较特殊的字段，在课堂上演示过，所谓的匿名字段。如果字段类型此时是结构体，那么应该递归进去解析这个结构体的每一个字段；
>以下文档是一个很全面的考虑，但是我们作业是极简版，所以简单定义一个 InsertStmt 方法就可以，主要还是以语法练习为主。
>实际上在我们后面构造 ORM 框架的时候，采用的才是下面文档的设计。下面设计的核心要素是使用 Builder 模式，和我们作业的单一方法比起来要复杂很多，大概看一下就可以，ORM 模块的时候我们会详细分析这个文档。
>另外你们也可以注意到，实际上在技术文档里面，把调研（也就是需求分析那一块）做好之后，其实设计是比较简单的事情。大多数同学写不出代码或者不知道怎么解决问题，其实都是卡在了搞清楚需求这一步。
# 背景

当我们希望把一个类型实例的数据存储进去数据库的时候，就必然面临着一个问题：如何将实例转化为一个 INSERT 语句，并且生成插入语句的参数。

在实际业务中，有很多种情况是需要考虑进去的：

* 批量插入
* 指定列，以及指定列的表达式
* 组合：一般的公司会考虑定义一些公共的结构体，要求所有的业务方必须组合这些结构体
* upsert：即所谓的 INSERT or UPDATE 语义，常见的场景则是数据不存在则插入，数据存在则删除
## 名词解释

* 多重组合：一个结构体同时组合了多个结构体
* 深层组合：一个结构体 A 组合了另一个结构体 B，而 B 本身也组合了 C
* 多级指针：指向指针的指针，如 **int
* 方言：SQL虽然有一个标准，但是往往数据库支持的语义可能有所差别，并且有些数据库会有自己独特的语法，因此不同数据库支持的 SQL 我们称之为方言
* upsert：即所谓的 INSERT or UPDATE语义，如果数据不存在，则插入；如果数据存在，则更新。存在与否的判断，一般是通过主键或者唯一索引来进行
# 需求分析

## 场景分析

这里基本上就是考虑三个因素：

* 开发者会如何构建输入。核心就是他们会使用指针还是非指针；
* 开发者如何定义模型。关键是他们是否会使用组合，以及如何使用组合；
* 开发者使用什么类型来定义字段。关键在于开发者如果用到了复杂结构体来作为字段类型，该如何处理；
* 开发者使用批量插入
* 开发者如何指定列，以及列的表达式
* 开发者希望执行 upsert。并且需要额外考虑以下情况
    * 指定冲突的列或者索引（在 PostgreSQL 或者 SQLite3）上
    * 指定更新的列，以及对应的值
|**场景**|**例子**|**说明**|
|:----|:----|:----|
|使用结构体作为输入|InsertStmt(User{})|能够正常处理|
|使用指针作为输入|InsertStmt(&User{})|能够正常处理|
|使用多级指针作为输入|InsertStmt(&&User{})|返回错误，不需要支持。很显然，正常情况下，开发者没有道理使用这种作为输入|
|使用 nil 作为输入|InsertStmt(nil)|不需要支持，这只能是开发者偶发性错误，没有在调用之前检测是否为 nil。返回错误会比较合适|
|使用基本类型或者内置类型作为输入|InsertStmt([]int)|不需要支持。正常的开发者都不会尝试这样使用，返回错误就可以|
|使用了组合|type Buyer struct {<br>    User<br>}|支持。这是常见用法。尤其是在一些公司有内部规范的情况下，例如部分公司可能会要求所有的模型都必须强制组合一个 BaseEntity，这个结构体里面包含基本的 Id, CreateTime, UpdateTime 等常见字段。在这种情况下，User 里面的字段会被认为是 Buyer 的一部分。因此生成的 SQL 里面也要包含这些字段|
|使用了组合，但是组合是指针形态|type Buyer struct {<br>    *User<br>}|不需要支持。大多数情况下，组合都会采用非指针的形态，尤其是在数据库模型定义这个层面上，这种定义方式并没有明显的优势，也找不出一个非它不可的场景。<br><br>另外一个理由是，大多数时候，如果追求性能，我们会尝试用 unsafe 来取代反射操作，但是 unsafe 对指针类型的组合优势不大|
|使用了复杂结构体作为字段|type Buyer struct {<br>    User User<br>}|不做校验。严格来讲，这种定义方式和组合定义方式，在语义上是有区别的。组合意味着同一张表，这种更加接近关联关系。User 整体会被认为是一个参数，但是实际上这个参数在 driver 执行的时候会出错。不过我们并不会对这个执行校验。<br>理由非常简单：driver 已经提供了这种保障。开发者应该知道 driver 支持什么，不支持什么。|
|使用了内置类型作为字段类型|    |同上|
|使用了 time.Time 作为字段|    |同上|
|使用了 driver.Valuer 作为字段类型|    |    |
|使用了组合，但是有同名字段|type User struct {<br>    Id int64<br>}<br><br>type Buyer struct {<br>    Id int64<br>    User<br>}|User 和 Buyer 都定义了自己的 Id 字段，这种情况下只需要取一个，可以按照谁先就用谁的值的原则来构建 SQL。|
|批量插入|InsertBatch(u1, u2...)|也就是意味着用户可以在一个 SQL 语句里面插入多行。<br>插入多行需要注意一点，就是说有的插入的数据需要有同样的类型。<br>|
|批量插入，但是不同类型|u := &User{}<br>o := &Order{}<br>InsertBatch(u, o)|返回错误|
|批量插入，不同类型但是字段完全相同|u1 := &UserV1{}<br>u2 := &UserV2{}<br>InsertBatch(u1, u2)|假定说 UserV1 和 UserV2 两个类型的字段一模一样，在这种场景下，依旧返回错误。<br>原则上来说，这种是可以正确生成 SQL 的，但是在用户层面上，他们不应该这么使用|
|批量插入，但是多批次|InsertBatch(users, 10)|假如说 users 传入了 1000 个实例，而且用户要求 10 个一批，那么 1000 个实例。<br><br>类似这种需求，其实不属于 ORM 层面上的需求，应该在应用层面上处理。<br><br>对于 ORM 来说，如果要支持该场景，那么需要解决：<br><ul><li>要不要开事务？</li><li>部分批次失败，部分批次成功，怎么办？</li></ul><br>所以ORM 框架不应该处理这种情况|
|插入，同时指定插入的列|Insert(user, "age", "first_name")|在一些场景之下，用户不希望插入所有的列，而是希望能够只插入部分的列。<br>如果用户指定的列不存在，那么应该返回错误|
|插入复杂的列表达式|Insert(user, "age", "create_time=now()")|在插入的时候，用户希望使用 MySQL 函数 now() 来作为插入的值。其生成的语句类似于：<br>INSERT xxx(col, col2) VALUES(val1, now());<br><br>类似于 now() 这种表达式，是跟使用的数据库相关的，所以不需要对表达式的正确性进行校验，用户需要对此负责。|
|自增主键|    |用户在插入数据的时候，如果主键有值，那么应该使用主键的值，如果主键没有值，那么应该自增生成一个主键。<br><br>|
|自增主键，但是 0 值|type User struct {<br>    ID int<br>}<br><br>Insert(&User{})|在这种情况下，用户使用基本类型来作为主键类型，那么用户在没有设置值的情况下，它的默认值是 0, 0 应该被看做没有设置值，从而触发自增主键。|
|单个插入，获得自增主键|    |在单个插入的情况下，我们可以确定无疑获得自增主键，而且是必然正确的主键|
|批量插入，获得自增主键|    |用户可能期望，如果他一次性插入 100 条数据，假如说第一条的 ID 是 201，那么下一条是 202,203,204...<br>实际上，有些数据库它同一批次插入的 ID，都不是连续的<br>因此实际上在批量插入的时候不需要返回所有的 ID，只需要返回 last_insert_id，用户根据自己的数据库的配置，来计算其它 ID。<br><br>关于自增主键 ID 是否一定连续，[可以参考 stackoverflow 的讨论](https://stackoverflow.com/questions/34200805/when-i-insert-multiple-rows-into-a-mysql-table-will-the-ids-be-increment-by-1-e)|
|新增或者更新（Upsert）|InsertOrUpdate(users)<br><br>|用户希望，如果要是数据冲突了（可能是主键冲突，也可能是唯一索引冲突），那么就执行更新。<br><br>在不同的方言下，upsert 的写法是不同的。所以需要考虑兼容不同的方言；<br><br>具体参考方言分析部分。|
|INSERT...SELECT|    |这种用法就是用户插入一个查询返回的数据。<br>这部分将在子查询部分进一步考虑。|

## 方言分析

此处我们忽略了 Oracle，是因为 Oracle 缺乏开源免费版本，因此对测试非常不友好。

### MySQL

MySQL 的 INSERT 总体上有三种形态，在[它的文档](https://dev.mysql.com/doc/refman/8.0/en/insert.html)里面有详细描述。

第一种是：

```go
INSERT [LOW_PRIORITY | DELAYED | HIGH_PRIORITY] [IGNORE]
    [INTO] tbl_name
    [PARTITION (partition_name [, partition_name] ...)]
    [(col_name [, col_name] ...)]
    { {VALUES | VALUE} (value_list) [, (value_list)] ... }
    [AS row_alias[(col_alias [, col_alias] ...)]]
    [ON DUPLICATE KEY UPDATE assignment_list]
```
比较值得注意的就是它采用了 `ON DUPLICATE KEY UPDATE` 来解决 upsert 的问题。这种形态也是我们最常见的形态。
第二种是：

```plain
INSERT [LOW_PRIORITY | DELAYED | HIGH_PRIORITY] [IGNORE]
    [INTO] tbl_name
    [PARTITION (partition_name [, partition_name] ...)]
    SET assignment_list
    [AS row_alias[(col_alias [, col_alias] ...)]]
    [ON DUPLICATE KEY UPDATE assignment_list]
```
和第一种比起来，这里用了  `SET assignment_list` 的语法。
第三种形态是：

```plain
INSERT [LOW_PRIORITY | HIGH_PRIORITY] [IGNORE]
    [INTO] tbl_name
    [PARTITION (partition_name [, partition_name] ...)]
    [(col_name [, col_name] ...)]
    { SELECT ... 
      | TABLE table_name 
      | VALUES row_constructor_list
    }
    [ON DUPLICATE KEY UPDATE assignment_list]
```
这种形态使用了 SELECT 子句。
#### ON DUPLICATE KEY UPDATE

在 MySQL 的  `ON DUPLICATE KEY UPDATE` 部分，它后面可以跟着一个  `assinment_list` ，而  `assignment_list` 的定义是：

```plain
assignment:
    col_name = 
          value
        | [row_alias.]col_name
        | [tbl_name.]col_name
        | [row_alias.]col_alias
assignment_list:
    assignment [, assignment] ...
```
关键是 assignment 有四种：
1. value：一个纯粹的值
2. row_alias.col_name：在使用了行别名的时候才会有的形态
3. tbl_name.col_name：指定更新为插入的值
4. row_alias.col_alias：这个和第二种比较像，所不同的是这里使用的是别名
在 ORM 框架中，我们不需要支持 2 和 4，因为在插入部分，我们就不会使用行的别名，所谓的 row_alias。

实际上 MySQL 的这个规范写得还是漏了一部分，也就是我们其实可以在 assignment 里面使用复杂的表达式，例如  `ON DUPLICATE KEY UPDATE update_time=now()` 又或者  `ON DUPLICATE KEY UPDATE epoch = epoch +1` 。ORM 框架需要将这种用法纳入考虑范围。

另外一个值得注意的点是：ON DUPLICATE KEY 是无法指定 KEY的，也就是说，假如我们的表上面定义了很多个唯一索引，那么任何一个唯一索引冲突（包含主键）都会引起执行更新。这和下面讨论的 PostgreSQL，SQLite3 非常不一样。

### PostgreSQL

[PostgreSQL 的语法](https://www.postgresql.org/docs/14/sql-insert.html)在 INSERT 部分和 MySQL 的第一种形态接近。但是它的 upsert 部分采用的也是 ON CONFLICT 语法：

```plain
[ WITH [ RECURSIVE ] with_query [, ...] ]
INSERT INTO table_name [ AS alias ] [ ( column_name [, ...] ) ]
    [ OVERRIDING { SYSTEM | USER } VALUE ]
    { DEFAULT VALUES | VALUES ( { expression | DEFAULT } [, ...] ) [, ...] | query }
    [ ON CONFLICT [ conflict_target ] conflict_action ]
    [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]
where conflict_target can be one of:
    ( { index_column_name | ( index_expression ) } [ COLLATE collation ] [ opclass ] [, ...] ) [ WHERE index_predicate ]
    ON CONSTRAINT constraint_name
and conflict_action is one of:
    DO NOTHING
    DO UPDATE SET { column_name = { expression | DEFAULT } |
                    ( column_name [, ...] ) = [ ROW ] ( { expression | DEFAULT } [, ...] ) |
                    ( column_name [, ...] ) = ( sub-SELECT )
                  } [, ...]
              [ WHERE condition ]
```
ON CONFLICT 部分简单概括可以看做是：
* ON CONFLICT(col1, col2) DO NOTHING
* ON CONFLICT(co1, col2) DO UPDATE SET col1=xxx, ...
举例来说：

```sql
-- Do nothing 的例子
INSERT INTO distributors (did, dname) VALUES (7, 'Redline GmbH')
    ON CONFLICT (did) DO NOTHING;

-- 这个是指定了索引的例子，索引必须是唯一索引    
INSERT INTO distributors (did, dname) VALUES (9, 'Antwerp Design')
    ON CONFLICT ON CONSTRAINT distributors_pkey DO NOTHING;

-- update 的例子，还在 Update 里面带了 where    
INSERT INTO distributors AS d (did, dname) VALUES (8, 'Anvil Distribution')
    ON CONFLICT (did) DO UPDATE
    SET dname = EXCLUDED.dname || ' (formerly ' || d.dname || ')'
    WHERE d.zipcode <> '21201';
```

### SQLite3

SQLite3 的语法整体上要简单很多。在 INSERT 部分类似于 MySQL 的第一种形态，也就是我们所熟知的那种形态。而在 UPSERT 部分则是采用的是 ON CONFLICT 语法。完整的语法结构参考 [INSERT 语句](https://www.sqlite.org/lang_insert.html)。

换句话说，SQLite3 的语法和 PostgreSQL 的语法更加接近。

## 框架分析

### GORM 分析

GORM 的跟 Insert 有关的方法有：

* Create: 可以插入单个，也可以插入批量
* CreateInBatches：分批次插入，例如可以指定 10 个一批，那么 100 个就会拆成 10 批
* Save：接近于 INSERT or Update 的语义
```go
// Create insert the value into database
func (db *DB) Create(value interface{}) (tx *DB) {
   if db.CreateBatchSize > 0 {
      return db.CreateInBatches(value, db.CreateBatchSize)
   }
   tx = db.getInstance()
   tx.Statement.Dest = value
   return tx.callbacks.Create().Execute(tx)
}

// CreateInBatches insert the value in batches into database
func (db *DB) CreateInBatches(value interface{}, batchSize int) (tx *DB) {
   // ...
   return
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (db *DB) Save(value interface{}) (tx *DB) {
   // ... 
}
```
代码位于：[gorm/gorm.go at master · go-gorm/gorm (github.com)](https://github.com/go-gorm/gorm/blob/master/gorm.go)
这里比较有特色的地方是 CreateInBatches，即 ORM 帮助用户解决了分批次插入的问题。但是整体上来说，我（们）认为，这个应该是应用层面上用户自己解决的问题，我们没有必要帮助用户解决。

>注：GORM 最终拼出来一个 SQL 的过程是很复杂的，不必细究，因为后面我们 ORM 的思路和它不一样，我在 ORM 的小课上大概提到了 GORM 的设计思路，可以去看看
### Beego ORM

Beego ORM 的接口定义是放在 [beego/types.go at develop · beego/beego (github.com)](https://github.com/beego/beego/blob/develop/client/orm/types.go)

里面提供了方法：

* Insert 和 InsertWithCtx：插入单个
* InsertOrUpdate 和 InsertOrUpdateWithCtx：插入或者更新，同时支持指定冲突列。在实际过程中，用户会困惑与 MySQL 怎么指定冲突列，因为 MySQL 语法根本不支持冲突列
* InsertMulti 和 InsertMultiWithCtx：类似于 GORM 的 CreateInBatches，分批次插入
```go
// Data Manipulation Language
	type DML interface {
		// insert model data to database
		// for example:
		//  user := new(User)
		//  id, err = Ormer.Insert(user)
		//  user must be a pointer and Insert will set user's pk field
		Insert(md interface{}) (int64, error)
		InsertWithCtx(ctx context.Context, md interface{}) (int64, error)
		// mysql:InsertOrUpdate(model) or InsertOrUpdate(model,"colu=colu+value")
		// if colu type is integer : can use(+-*/), string : convert(colu,"value")
		// postgres: InsertOrUpdate(model,"conflictColumnName") or InsertOrUpdate(model,"conflictColumnName","colu=colu+value")
		// if colu type is integer : can use(+-*/), string : colu || "value"
		InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error)
		InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error)
		// insert some models to database
		InsertMulti(bulk int, mds interface{}) (int64, error)
		InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error)
}
```
从实现的角度来说， Beego ORM 的代码和 GORM 的代码一样很复杂，相比之下，GORM 的代码属于设计上的复杂——即职责被切分给了很多接口，而 Beego 的复杂，则在于所有的代码都混在了一起，难以搞懂每个步骤具体是干什么的。
## 功能需求

### 生成查询

将一个结构体转化为一个对应的 INSERT 查询，查询包括：

* SQL
* 参数
要支持以下选项：

* 单个或者批量：但是我们不会把一次提交的数据拆分成过个批次，也就是不管用户提交了多少数据，我们都是一批次插入进去
* 指定列：
    * 普通列
    * 复杂表达式在初期阶段不支持，在下一阶段支持
* 返回主键：只需要返回 last_insert_id
* upsert：支持 MySQL，PostgreSQL 和  SQlite3 的语法，并且在后两者中支持 UPDATE 和 DO NOTHING 两种动作。ON CONFLICT 部分可以只支持传入列名，而不支持指定冲突索引
### 方言支持

在 upsert 处涉及到了方言的处理，也就是说需要有依据方言来构建不同 SQL 的能力。因此需要设计一个综合的方言解决方案，该方案要求：

* 可以方便扩展新的方言
* 不同方言之间的实现相互隔离，互不影响
目前仅仅需要考虑支持 MySQL, PostgreSQL 和 SQLite3。

## 非功能需求

* 扩展性。要求我们在将来支持更加复杂的 upsert 语句的情况下，变更极小
* 耦合性。方言之间不存耦合
# 设计

## 总体设计

总体设计上，我们采用 Builder 模式，即定义一个新的  `Inserter` :

```go
type Inserter[T any] struct {
    
}

// 构建 SQL
func (i *Inserter[T]) Build() (Query, error) {}



type Query struct {
  SQL string
  Args []any
}
```
在 Builder 模式之下，任何复杂部分都可以拆成几个单一的方法。
与此同时，我们在 Inserter 上使用了泛型了，用来约束用户所能传入的类型。因此插入不同类型这种情况，是不会出现的，因为用户会得到编译错误。

### 方言

于此同时，为了解决方言的问题，我们需要进一步引入一个新的抽象  `Dialect` 。这个 `Dialect` 用于屏蔽不同方言之间的不同。

同时，除了为 Dialect 提供各种不同的实现以外，还会有一个基于标准 SQL 的实现，standardSQL。在引入了 standardSQL 之后，整体的方言抽象就变成了：


![图片](data:image/png;base64,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f3XJcHJOxVshENoNab1kXS246zacnlRUgNkN99dXcLqnbZCXqOva5Rk5OeCxM/+enHBInv0p5l3EIXpn/pXsXcQBBxxwwAEnJu/Mn2Q5eYw44IADDjhET9n8kZwccMABBxxwwCF6Zy5PqOvakYNR17X/EnJywGNn/iLWeYo44IADDnhMnZt/zryLOOCAAw54TJ2bf8m65BEccEie/SPLnmJ0wAGH6DFlt5wc8JT9QU4OeMr+NOs8RRxwwAEHHKJ35h9l7yIOOOCAAw44MXln/iTLyWPEAQcccMAhesrmj+TkgAMOOOCAAw7RO/ODV9e1IwejrmuXb8neRRxwwAEHHHBi8s78SZaTx4gDDjjggEP0lM0fyckBBxxwwAEHHKJ35h/k5ICn7F+wLnkEBxySZ//IOk8RBxxwwAGH6J25fEdd1847dV27/GzmXcSJnZt/i7mZuZn5jzAzzwmH6J25/I3qunbkYNR17f91ZuZm5i9h/hLmL2FddIjemf8pZuZm5t9mbmZuZv4jzMxzwiF6Z/6fVde1IwejrmuXlzI3Mzcz/xFm5jnhEL0z/8nMn2NmbmYuL1fXtf+G/CPK9jXrEVJuCXxLIAR+WAiBsDO624p1tWLmOxpE5L8ohMBLBV4i8F1l4HI9AolZ4E8JIfB9gRD4YSEEws7obivW1YqZ72gQkV9HIAR+WAiBsDO624p1tWLmOxp+lsBzQgjIj5siP1XZrlgsJlTrkZgyu4afKNC+6YjxlovVQEFE5CcoW1aLFavVitVqxWqxYFIt6Yl0tqPh3yrQvumI8ZaL1UBBROSdQPumI8ZbLlYDBTkUvyE/ld3dAokun9M2gZ8utOz3LSIiP8+M+fEVV7d8dEzqNpyfNITAv1to2e9bRES+EFr2+xY5LL8hP1Wz29MgIvILCQ3trqFFRETk/88UERERERGRAzNFRERERETkwEwRERERERE5MFNEREREREQOzBQREREREZEDM0VEREREROTATBERERERETkwU0RERERERA7MFBERERERkQMzRURERERE5MBMEREREREROTBTREREREREDswUERERERGRAzNFRERERETkwEwRERERERE5MFNEREREREQOzBQREREREZEDM0VEREREROTATBERERERETkwU0RERERERA7MFBERERERkQMzRURERERE5MBMEREREREROTBTREREREREDswUERERERGRAzNFRERERETkwEwRERERERE5MFNEREREREQOzBQREREREZEDM0VEREREROTATBERERERETkwU0RERERERA7Mb7x1dHT0+6tXr/6HiDxydHT0O3Iwjo6Ofn/16tX/EJFHjo6OfkcOxtHR0e+vXr36HyLyyNHR0e//B2WvQwKUZoRHAAAAAElFTkSuQmCC)


总结：

* 如果 SQL 的某一部分，不同方言之间有差异，那么就在 Dialect 里面新增一个方法
* 如果方言遵守 SQL 标准，那么我们不需要对它进行特殊处理
* 如果方言不遵守 SQL 标准，那么具体方言的实现就需要提供自己的实现
## 详细设计

### Insert 

在 INSERT 部分，需要考虑的就是两个问题： 如何指定插入的行，以及如何指定插入的列。

#### Values 方法

```go
func (i *Inserter[T]) Values(vals...*T) *Inserter {
  panic("implement me")
}
```
在这个方法里面，如果用户传入单个值，例如  `Values(&User)` ，那么就是插入一行；如果用户传入了多个值，例如  `Values(&user{}, &User{})` 那么就是批量插入。如果用户没有调用，或者调用了但是没有传值，例如  `Values()` 那么我们将在构建 SQL 的时候返回错误。
在这种设计之下，我们并没有区别单个插入还是批量插入。

>实际上这里 vals 可以不用指针
#### Columns 方法

```go
func (i *Inserter[T]) Columns(cols...string) *Inserter[T] {
  panic("implement me")
}
```
目前这种设计，我们放弃了支持复杂表达式，所以用户只能传入具体的列，而不能指定列的表达式，例如 now() 这种数据库方法调用。
#### Exec 方法

```go
func (i *Inserter[T]) Exec(ctx context.Context) (sql.Result, error) {
  panic("implement me")
}
```
将会发起查询，并且返回结果。用户可以通过  `sql.Result` 拿到  `last_insert_id` ,也可以拿到受影响行数，即插入数量。
### Upsert 

在 Insert 部分，我们只是构建了 INSERT 的主要部分，但是如果用户想要使用 upsert 特性，那么就需要调用额外的方法。

为了支持 upsert，首先需要定义额外的结构体：

```go
type Upsert struct {
  doNothing bool
  updateColumns []string
  conflictColumns []string
}
type UpsertBuilder[T any] struct {
  i *Inserter[T]
  conflictColumns []string
}

// Update 指定在冲突的时候要执行 update 的列
func (u *UpsertBuilder[T]) Update(cols...string) *Inserter[T] {
   i.upsert = Upsert{conflictColumns: u.conflictColumns, updateColumns: cols}
   return i
}

func (u *UpsertBuilder[T]) DoNothing() *Inserter[T] {
  i.upsert = Upsert{conflictColumns: u.conflictColumns, doNothing: true}
}

func (u *UpsertBuilder[T]) ConflictColumns(cols...string) *UpsertBuilder[T] {
  u.conflictColumns = cols
}
```
这里面明显分成两类方法：
* 中间方法：ConflictColumns，用户指定冲突的列，在将来我们可以增加对索引的支持
* 终结方法：DoNothing 和 Update，调用这两个方法后，重新回到 Inserter
很显然，在 Inserter 里面就是要暴露一个方法允许用户跳过去 UpsertBuilder 里面：

```go
func (i *Inserter[T]) Upsert() *UpsertBuilder[T] {
    return &UpsertBuilder {
      i: i
    }
}
```
使用起来的效果是：
![图片](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAABa4AAAEqCAIAAAByb1CBAAAgAElEQVR4AezBDVCUhaI//u8+u6y6LGvqSAEiKygiYBroCTYVFbWwoqgr6h6nQ2FHratCk8c8+kvseH2pScUc7cU8esvj0fwrpic1s1BzNZRcIeTFq7IgIuKmILuwsLvPf2ZnnhkYWF9Twf1+PjJRFEFERERERERE5BkEEBERERERERF5DAXuisPhMBqNBoPhwoULpaWl169ft1gsdrsdRA+bQqFQufj5+QUGBoaEhDz11FNhYWFyuRxEHs/hcBiNRoPBcOHChdLS0uvXr1ssFrvdDiKih0ehUKhc/Pz8AgMDQ0JCnnrqqbCwMLlcDiKP53A4jEajwWC4cOFCaWnp9evXLRaL3W4H0cOmUChULn5+foGBgSEhIU899VRYWJhcLke7JxNFEbfNZrMdPXo0Ozv7yJEjtbW1IOogvL294+Linn322djYWLlcDiIPY7PZjh49mp2dfeTIkdraWhARtXve3t5xcXHPPvtsbGysXC4HkYex2WxHjx7Nzs4+cuRIbW0tiDoIb2/vuLi4Z599NjY2Vi6Xo72SiaKI22C327OystavX19dXQ2XiIiIYcOGhYaGBgUFde/eXaVSeXl5gehha2pqqq+vt1gsly5dKi8vz8/PNxgMV65cgUtQUNBbb701evRoQRBA5AHsdntWVtb69eurq6vhEhERMWzYsNDQ0KCgoO7du6tUKi8vLxARPTxNTU319fUWi+XSpUvl5eX5+fkGg+HKlStwCQoKeuutt0aPHi0IAog8gN1uz8rKWr9+fXV1NVwiIiKGDRsWGhoaFBTUvXt3lUrl5eUFooetqampvr7eYrFcunSpvLw8Pz/fYDBcuXIFLkFBQW+99dbo0aMFQUD7IxNFETfldDr379+/bt26iooKADqdLj4+fvjw4T169ABRRyCK4vnz5w0Gw4EDBwoKCgD079//7bff1ul0MpkMRI8op9O5f//+devWVVRUANDpdPHx8cOHD+/RoweIiNo3URTPnz9vMBgOHDhQUFAAoH///m+//bZOp5PJZCB6RDmdzv37969bt66iogKATqeLj48fPnx4jx49QNQRiKJ4/vx5g8Fw4MCBgoICAP3793/77bd1Op1MJkN7IhNFEe5VV1fPnTv39OnTAGJiYt56662IiAgQdUyi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BDBo0CAQEXm8gQMHAjh//jyI2rHCwkIAgwYNAhG5DBw4EMD58+fxwAnnz58HEBISAiJyCQwMBFBeXg6iduz8+fMAQkJCQETk8QIDAwGUl5eDqB07f/48gJCQEBCRS2BgIIDy8nI8cMLFixcB9O7dG0Tk4ufnB+DSpUsgascuXrwIoHfv3iAi8nh+fn4ALl26BKJ27OLFiwB69+4NInLx8/MDcOnSJTxwwo0bNwBoNBoQkYtarQZgtVpB1I7duHEDgEajARGRx1Or1QCsViuI2rEbN24A0Gg0ICIXtVoNwGq14oETrFYrAJVKBSJy6dKlCwCr1QqidsxqtQJQqVQgIvJ4Xbp0AWC1WkHUjlmtVgAqlQpE5NKlSxcAVqsVD5xgt9sBeHl5gYhcvLy8ANjtdhC1Y3a7HYCXlxeIiDyel5cXALvdDqJ2zG63A/Dy8gIRuXh5eQGw2+144AQQEREREREREXkMAUREREREREREHkMAEREREREREZHHEEBERERERERE5DEEEBERERERERF5DAFERERERERERB5DABERERERERGRxxBAREREREREROQxBBAREREREREReQwBREREREREREQeQwARERERERERkccQQERERERERETkMQQQEREREREREXkMAUREREREREREHkMAEREREREREZHHEEBERERERERE5DEEEBERERERERF5DAFERERERERERB5DABERERERERGRxxBAREREREREROQxBBAREREREREReQwBREREREREREQeQwARERERERERkccQQERERERERETkMQQQEREREREREXkMAUREREREREREHkMAEREREREREZHHEEBERERERERE5DEEEBERERERERF5DAFERERERERERB5DABERERERERGRxxBAREREREREROQxBBAREREREREReQwBREREREREREQeQwARERERERERkccQQEREROQB3nvvvQBJRkYG7tb06dMDJB999BGIiIioo1HgkWCz2RoaGuDSuXPnTp06gYiIPJvFYrHb7XBRq9VyuRw31djYWF9fD5euXbuCbpvNZmtoaEArXl5eXbp0kclkeCCampqsVitcVCqVl5cXWiovL7906RJcLl68iLtVVlZ26dIluFy6dAlERJ7EZrM1NDTApXPnzp06dQJRByTgkfDcc889JklOTgYREXm8J554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)


## 例子

假如说我们有一个结构体：

```go
type User struct {
  ID uint64
  Email string
  FirstName string
  Age uint8
}
```
单个插入
```go
NewInserter[User]().Values(&User{ID: 1, Email: "xxx@xx"})
// INSERT INTO `user`(id, email, first_name, age) VALUES(?,?,?,?)
// []{1, "xxx@xx", "", 0}
```
批量插入
```go
NewInserter[User]().Values(&User{ID: 1, Email: "xxx@xx"}, 
&User{ID: 2, Email:"bb@aa", Age: 18})
// INSERT INTO `user`(id, email, first_name, age) VALUES(?,?,?,?),(?,?,?,?)
// []{1, "xxx@xx", "", 0, 2, "bb@aa", "", 18}
```
指定列
```go
NewInserter[User]().Values(&User{ID: 1, Email: "xxx@xx"}).Columns("id","email")
// INSERT INTO `user`(id, email) VALUES(?,?)
// []{1, "xxx@xx"}
```
不插入主键，只是在指定列的时候把主键排除掉
```go
NewInserter[User]().Values(&User{Email: "xxx@xx", FirstName:"Deng"}).Columns("email", "first_name", "age")
// INSERT INTO `user`(email, first_name, age) VALUES(?,?,?)
// []{"xxx@xx", "Deng", 0}
```
MySQL 冲突，更新列：
```go
NewInserter[User]().Values(&User{ID: 1,Email: "xxx@xx", FirstName:"Deng"}).
Upsert().Update("first_name")
// INSERT INTO `user`(email, first_name, age) VALUES(?,?,?,?) 
// ON DUPLICATE KEY UPDATE `first_name`=VALUES(`first_name`)
// []{1, "xxx@xx", "Deng", 0}
```
SQILite 冲突，DoNothing
```go
NewInserter[User]().Values(&User{ID: 1,Email: "xxx@xx", FirstName:"Deng"}).
Upsert().ConflictColumn("email")
// INSERT INTO `user`(email, first_name, age) VALUES(?,?,?,?) 
// ON CONFLICT(email) DO NOTHING
// []{1, "xxx@xx", "Deng", 0}
```
# 测试

## 单元测试

影响最终结果的主要有两个因素：

* 输入的形态：nil, 非指针，指针，多级指针
* 结构体的定义：
    * 普通结构体
    * 组合
    * 指针式组合
    * 多重组合
    * 深层组合
    * 组合字段冲突（即含有同名字段）
* 指定列：
    * 直接指定列
    * 指定列的表达式：例如 now()
* Upsert 语句：
    * 指定冲突列
    * 不同方言：MySQL, SQLite 和 PostgreSQL
    * 冲突时候更新列的值：
        * 更新为一个具体值
        * 更新为插入的值
将这些因素进行交叉组合，得到所有的测试用例。注意，这些因素从设计的角度来看是各自独立的，所以可以考虑针对不同的因素分别设计测试用例，而不需要进行笛卡尔积。

