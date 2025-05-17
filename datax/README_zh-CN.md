[TOC]



# go-etl数据同步开发者指南

## 1 同步框架简介

go-etl主要离线数据同步框架，框架如下

```
readerPlugin(reader)—> Framework(Exchanger+Transformer) ->writerPlugin(riter)  
```

采用Framework + plugin架构构建。将数据源读取和写入抽象成为Reader/Writer插件，纳入到整个同步框架中。

+ Reader：Reader为数据采集模块，负责采集数据源的数据，将数据发送给Framework。 

+ Writer：Writer为数据写入模块，负责不断向Framework取数据，并将数据写入到目的端。

+ Framework：Framework用于连接reader和writer，作为两者的数据传输通道，并处理缓冲，流控，并发，数据转换等核心技术问题

## 2 核心模块(core)介绍

go-etl完成单个数据同步的作业，我们称之为Job，go-etl接受到一个Job之后，将启动一个进程来完成整个作业同步过程。

go-etl Job模块是单个作业的中枢管理节点，承担了数据清理、子任务切分(将单一作业计算转化为多个子Task)、TaskGroup管理等功能。

### 2.1 调度流程
```
    JOB--split--+-- task1--+           +--taskGroup1--+   
                |-- task2--|           |--taskGroup2--|        
                |-- task3--|-schedule--|--taskGroup3--|             
                |  ......  |           |  ......      |           
                |-- taskN--|           |--taskGroupM--|      
                                            |
                                     +------+---------------------------------------+ 
                                     |   Reader1->Exchanger1(Transformer)->Writer1  |
                                     |   Reader2->Exchanger2(Transformer)->Writer2  |
                                     |   Reader3->Exchanger3(Transformer)->Writer3  |
                                     |     ...            ......            ...     |
                                     |   ReaderN->ExchangerX(Transformer)->WriterX  |
                                     +----------------------------------------------+     
```

如上所示，go-etl Job启动后，会根据不同的源端切分策略，将Job切分成多个小的Task(子任务)，以便于并发执行。Task便是go-etl作业的最小单元，每一个Task都会负 责一部分数据的同步工作。切分多个Task之后，go-etl Job会调用Scheduler模块，根据配置的并发数据量，将拆分成的Task重新组合，组装成TaskGroup(任务组),Task数和taskGroup数可以不同（N:M）。每一个TaskGroup负责以一定的并发运行完毕分配好的所有Task，默认单个任务组的并发数量为4。每一个Task都由TaskGroup负责启动，Task启动后，会固定启动Reader—>Channel—>Writer的线程来完成任务同步工作。

go-etl作业运行起来之后，Job监控并等待多个TaskGroup模块任务完成，等待所有TaskGroup任务完成后Job成功退出。否则，异常退出，进程退出值非0。

举例来说，用户提交了一个go-etl作业，并且配置了20个并发，目的是将一个100张分表的mysql数据同步到odps里面。go-etl的调度决策思路是：go-etlJob根据分库分表切分成了100个Task。 根据20个并发，go-etl计算共需要分配4个TaskGroup。4个TaskGroup平分切分好的100个Task，每一个TaskGroup负责以5个并发共计运行25个Task。

+ Job:Job是go-etl用以描述从一个源头到一个目的端的同步作业，是go-etl数据同步的最小业务单元。比如：从一张mysql的表同步到odps的一个表的特定分区。   
+ Task:Task是为最大化而把Job拆分得到的最小执行单元。比如：读一张有1024个分表的mysql分库分表的Job，拆分成1024个读Task，用若干个并发执行。        
+ TaskGroup: 描述的是一组Task集合。在同一个TaskGroupContainer执行下的Task集合称之为TaskGroup
+ JobContainer: Job执行器，负责Job全局拆分、调度、前置语句和后置语句等工作的工作单元。类似Yarn中的JobTracker
+ TaskGroupContainer: TaskGroup执行器，负责执行一组Task的工作单元，类似Yarn中的TaskTracker。

## 3 编程接口

### 3.1 Reader插件接口

Reader需要实现以下接口:

#### 3.1.1 Job

Job组合*plugin.BaseJob，实现方法

```golang
    Init(ctx context.Context) (err error)
    Destroy(ctx context.Context) (err error)
    Split(ctx context.Context, number int) ([]*config.JSON, error)
    Prepare(ctx context.Context) error  //默认为空方法
    Post(ctx context.Context) error     //默认为空方法
```

- `Init`: Job对象初始化工作，此时可以通过`PluginJobConf()`获取与本插件相关的配置。读插件获得配置中`reader`部分。
- `Prepare`: 全局准备工作。
- `Split`: 拆分`Task`。参数number框架建议的拆分数，一般是运行时所配置的并发度。值返回的是`Task`的配置列表。
- `Post`: 全局的后置工作。
- `Destroy`: Job对象自身的销毁工作。

#### 3.1.2 Task

Task组合*plugin.BaseTask,实现方法

```golang
    Init(ctx context.Context) (err error)
    Destroy(ctx context.Context) (err error)
    StartRead(ctx context.Context,sender plugin.RecordSender) error 
    Prepare(ctx context.Context) error  //默认为空方法
    Post(ctx context.Context) error     //默认为空方法
```

- `Init`：Task对象的初始化。此时可以通过`PluginJobConf()`获取与本`Task`相关的配置。这里的配置是`Job`的`Split`方法返回的配置列表中的其中一个。
- `Prepare`：局部的准备工作。
- `StartRead`: 从数据源读数据，写入到`RecordSender`中。`RecordSender`会把数据写入连接Reader和Writer的缓存队列。
- `Post`: 局部的后置工作。
- `Destroy`: Task象自身的销毁工作。

#### 3.1.3 Reader

```golang
    Job() reader.Job
    Task() reader.Task
```

+ `Job`: 获取上述的Job的实例

+ `Task`: 获取上述的Task的实例

#### 3.1.4 命令生成 

```bash
cd tools/datax/plugin
#新增一个名为Mysql的reader -p命令可以时任意大小写，用于指定reader的名字，如果新增-d 代表会删除原来的模板
go run main.go -t reader -p Mysql
```

这个命令会在datax/plugin/reader中自动生成一个如下mysql的reader模板来帮助开发

```
    reader---mysql--+-----resources--+--plugin.json
                    |--job.go        |--plugin_job_template.json
                    |--reader.go
                    |--README.md
                    |--task.go
```

如下，不要忘了在plugin.json加入开发者名字和描述

```json
{
    "name" : "mysqlreader",
    "developer":"Breeze0806",
    "description":"use github.com/go-sql-driver/mysql."
}
```

另外，以帮助开发者避免在使用插件注册命令后编译时报错。

#### 3.1.5 数据库

如果你想帮忙实现关系型数据库的数据源，根据以下方式去实现你的数据源将更加方便

##### 3.1.5.1 数据库存储

查看[数据库存储开发者指南](../storage/database/README_zh-CN.md),不仅能帮助你更快地实现Reader插件接口，而且能帮助你更快地实现Writer插件接口

##### 3.1.5.2 数据库读取器

dbms reader通过抽象数据库存储的DBWrapper结构体成如下Querier，然后利用Querier完成Job和Task的实现

```go
//Querier 查询器
type Querier interface {
	//通过基础表信息获取具体表
	Table(*database.BaseTable) database.Table
	//检测连通性
	PingContext(ctx context.Context) error
	//通过query查询语句进行查询
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	//通过参数param获取具体表
	FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error)
	//通过参数param，处理句柄handler获取记录
	FetchRecord(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error)
	//通过参数param，处理句柄handler使用事务获取记录
	FetchRecordWithTx(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error)
	//关闭资源
	Close() error
}
```

像mysql实现Job和Reader,对于Task需要使用dbms.StartRead函数实现StartRead方法

#### 3.1.6 二维表文件流

##### 3.1.6.1 二维表文件流存储

查看[二维表文件流存储开发者指南](../storage/stream/file/README_zh-CN.md),不仅能帮助你更快地实现Reader插件接口，而且能帮助你更快地实现Writer插件接口

##### 3.1.6.2 文件读取器

像cvs那样Task和Reader,这里需要独立实现Job，实现切分方法Split和初始化方法Init

### 3.2 Writer插件接口

Writer 需要实现以下接口:

#### 3.2.1 Job

Job组合*plugin.BaseJob,实现方法:

```golang
    Init(ctx context.Context) (err error)
    Destroy(ctx context.Context) (err error)
    Split(ctx context.Context, number int) ([]*config.JSON, error) 
    Prepare(ctx context.Context) error //默认为空方法
    Post(ctx context.Context) error    //默认为空方法
```

- `Init`: Job对象初始化工作，此时可以通过`PluginJobConf()`获取与本插件相关的配置。写插件获得`writer`部分。
- `Prepare`: 全局准备工作。
- `Split`: 拆分`Task`。参数`number`框架建议的拆分数，一般是运行时所配置的并发度。值返回的是`Task`的配置列表。
- `Post`: 全局的后置工作。
- `Destroy`: Job对象自身的销毁工作。

#### 3.2.2 Task

Task组合*plugin.BaseTask,实现方法:

```golang
    Init(ctx context.Context) (err error)
    Destroy(ctx context.Context) (err error)
    StartWrite(ctx context.Context,receiver plugin.RecordReceiver) error
    Prepare(ctx context.Context) error     //默认为空方法
    Post(ctx context.Context) error        //默认为空方法
    SupportFailOver() bool                 //默认为空方法
```

- `Init`：Task对象的初始化。此时可以通过`PluginJobConf()`获取与本`Task`相关的配置。这里的配置是`Job`的`split`方法返回的配置列表中的其中一个。
- `Prepare`：局部的准备工作。
- `StartWrite`：从`RecordReceiver`中读取数据，写入目标数据源。`RecordReceiver`中的数据来自Reader和Writer之间的缓存队列。
- `Post`: 局Task部的后置工作。
- `Destroy`: Task自身的销毁工作。
- `SupportFailOver`: Task是否支持故障转移。

#### 3.2.3 Writer

```golang
    Job() writer.Job
    Task() writer.Task
```

+ `Job`: 获取上述的Job的实例

+ `Task`: 获取上述的Task的实例

#### 3.2.4 命令生成

```bash
cd tools/go-etl/plugin
#新增一个名为Mysql的writer -p命令可以时任意大小写，用于指定writer的名字，如果新增-d 代表会删除原来的模板
go run main.go -t writer -p Mysql
```

这个命令会在datax/plugin/writer中自动生成如下一个mysql的writer模板来帮助开发
```
    writer--mysql---+-----resources--+--plugin.json
                    |--job.go        |--plugin_job_template.json
                    |--README.md
                    |--task.go
                    |--writer.go
```

如下，不要忘了在plugin.json加入开发者名字和描述

```json
{
    "name" : "mysqlwriter",
    "developer":"Breeze0806",
    "description":"use github.com/go-sql-driver/mysql. database/sql DB execute select sql, retrieve data from the ResultSet. warn: The more you know about the database, the less problems you encounter."
}
```

另外，这个可以帮助开发者避免在使用插件注册命令后编译时报错。

#### 3.2.5 数据库

如果你想帮忙实现数据库的数据源，根据以下方式去实现你的数据源将更加方便，当然前提你所使用的驱动库必须实现golang标准库的database/sql的接口。

##### 3.2.5.1 数据库存储

查看[数据库存储开发者指南](../storage/database/README_zh-CN.md),不仅能帮助你更快地实现Reader插件接口，而且能帮助你更快地实现Writer插件接口

##### 3.2.5.2 数据库写入器

dbms writer通过抽象数据库存储的DBWrapper结构体成如下Execer，然后利用Execer完成Job和Task的实现

```go
//Execer 执行器
type Execer interface {
	//通过基础表信息获取具体表
	Table(*database.BaseTable) database.Table
	//检测连通性
	PingContext(ctx context.Context) error
	//通过query查询语句进行查询
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	//通过query查询语句进行查询
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	//通过参数param获取具体表
	FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error)
	//批量执行
	BatchExec(ctx context.Context, opts *database.ParameterOptions) (err error)
	//prepare/exec批量执行
	BatchExecStmt(ctx context.Context, opts *database.ParameterOptions) (err error)
	//事务批量执行
	BatchExecWithTx(ctx context.Context, opts *database.ParameterOptions) (err error)
	//事务prepare/exec批量执行
	BatchExecStmtWithTx(ctx context.Context, opts *database.ParameterOptions) (err error)
	//关闭
	Close() error
}
```

像mysql实现Job和Writer,对于Task需要使用dbms.StartWrite函数实现StartWrite方法

#### 3.2.6 二维表文件流

##### 3.2.6.1 二维表文件流存储

查看[二维表文件流存储开发者指南](../storage/stream/file/README_zh-CN.md),不仅能帮助你更快地实现Reader插件接口，而且能帮助你更快地实现Writer插件接口

##### 3.2.6.2 文件读取器

像cvs那样Task和Writer,这里需要独立实现Job，实现切分方法Split和初始化方法Init

## 4 插件配置文件

`go-etl`使用`json`作为配置文件的格式。一个典型的`go-etl`任务配置如下：

```json
{
    "core" : {
        "container": {
            "job":{
                "id": 1,
                "sleepInterval":100
            },
            "taskGroup":{
                "id": 1,
                "failover":{
                    "retryIntervalInMsec":0
                }
            }
        },
        "transport":{
            "channel":{
                "speed":{
                    "byte": 100,
                    "record":100
                }
            }
        }
    },
    "job":{
        "content":[
            {
                "reader":{
                    "name": "csvreader",
                    "parameter": {
                        "path":["d:\\a.txt"],
                        "column":[
                            {
                                "index":"1",
                                "type":"time",
                                "format":"yyyy-MM-dd"
                            }
                        ],
                        "encoding":"utf-8",
                        "delimiter":","
                    }
                },
                "writer":{
                    "name": "postgreswriter",
                    "parameter": {
                        "username": "postgres",
                        "password": "123456",
                        "writeMode": "copyIn",
                        "column": ["*"],
                        "preSql": [],
                        "connection":  {
                                "url": "postgres://192.168.15.130:5432/postgres?sslmode=disable&connect_timeout=2",
                                "table": {
                                    "db":"postgres",
                                    "schema":"public",
                                    "name":"cvs"
                                }
                         },
                        "batchTimeout": "1s",
                        "batchSize":1000
                    }
                },
               "transformer":[]
            }
        ],
        "setting":{
            "speed":{
                "byte":3000,
                "record":400,
                "channel":4
            }
        }
    }
}
```

任务的**配置中`job.content.reader.parameter`的value部分会传给`Reader.Job`；`job.content.writer.parameter`的value部分会传给`Writer.Job`** ，`Reader.Job`和`Writer.Job`可以通过`super.getPluginJobConf()`来获取。

### 4.1 如何设计配置参数

> 配置文件的设计是插件开发的第一步！

任务配置中`reader`和`writer`下`parameter`部分是插件的配置参数，插件的配置参数应当遵循以下原则：

- 驼峰命名：所有配置项采用驼峰命名法，首字母小写，单词首字母大写。

- 正交原则：配置项必须正交，功能没有重复，没有潜规则。

- 富类型：合理使用json的类型，减少无谓的处理逻辑，减少出错的可能。

  - 使用正确的数据类型。比如，bool类型的值使用`true`/`false`，而非`"yes"`/`"true"`/`0`等。
  - 合理使用集合类型，比如，用数组替代有分隔符的字符串。

- 类似通用：遵守同一类型的插件的习惯，比如数据库的`connection`参数都是如下结构：

  ```json
  {
      "connection":  {
      	"url": "tcp(192.168.0.1:3306)/mysql?parseTime=false",
          "table": {
              "db":"source",
              "name":"type_table"
          }
      }
  }
  ```

### 4.2 如何使用`config.JSON`结构体

```josn
{
    "a":{
        "b":[{
            c:"x"
        }]
    }
}
```

GetConfig中要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c

注意，因为插件看到的配置只是整个配置的一部分。使用`json.Config`对象时，需要注意当前的根路径是什么。

更多`json.Config`的操作请参考`config`包的文档。

## 5 插件打包发布

### 5.1 新增许可证（license）

当你开发完一个功能后在提交前，请运行如下命令用于自动加入许可证并使用gofmt -s -w格式化代码

```bash
go run tools/license/main.go
```

### 5.2 插件注册

在使用golang编译前，需要将插件注册到代码中去。

golang静态编译的方式决定了go-etl框架不能用运行时动态加载插件的方式去获取插件，为此这里只能使用注册代码的方式，以下命令会生成将由开发者开发的reader和writer插件注册到程序中的代码。

```bash
go generate ./...
```
主要的原理如下会将对应go-etl/plugin插件中的reader和writer的resources的plugin.json生成plugin.go，同时在go-etl目录下生成plugin.go用于导入这些插件， 具体在tools/go-etl/build实现,另外通过-i命令可以忽略编译数据源来源，可以忽略db2， 由于db2会使用odbc去访问数据库，并且需要在linux中被依赖，如果不需要用这个直接忽略。

## 6. 插件数据传输

跟一般的`生产者-消费者`模式一样，`Reader`插件和`Writer`插件之间也是通过`channel`来实现数据的传输的。`channel`可以是内存的，也可能是持久化的，插件不必关心。插件通过`RecordSender`往`channel`写入数据，通过`RecordReceiver`从`channel`读取数据。

`channel`中的一条数据为一个`Record`的对象，`Record`中可以放多个`Column`对象，这可以简单理解为数据库中的记录和列，`Record`原型具体见[文档](../element/README_zh-CN.md)的《记录》一章。

因为`Record`是一个接口，`Reader`插件首先调用`RecordSender.createRecord()`创建一个`Record`实例，然后把`Column`一个个添加到`Record`中。

`Writer`插件调用`RecordReceiver.getFromReader()`方法获取`Record`，然后把`Column`遍历出来，写入目标存储中。当`Reader`尚未退出，传输还在进行时，如果暂时没有数据`RecordReceiver.getFromReader()`方法会阻塞直到有数据。如果传输已经结束，会返回`ErrTerminate`，`Writer`插件可以据此判断是否结束`startWrite`方法。

### 6.1 数据类型转化

为了规范源端和目的端类型转换操作，保证数据不失真，go-etl支持六种内部数据类型,具体见[文档](../element/README_zh-CN.md)的《数据类型转化》一章。

## 7. 插件文档

在插件文档README.md文档中加入以下几章内容

1. **快速介绍**：介绍插件的使用场景，特点等。
2. **实现原理**：介绍插件实现的底层原理，比如`mysqlwriter`通过`insert into`和`replace into`来实现插入，`tair`插件通过tair客户端实现写入。
3. **配置说明**
   - 给出典型场景下的同步任务的json配置文件。
   - 介绍每个参数的含义、是否必选、默认值、取值范围和其他约束。
4. **类型转换**
   - 插件是如何在实际的存储类型和`go-etl`的内部类型之间进行转换的。
   - 以及是否存在特殊处理。
5. **性能报告**
   - 软硬件环境，系统版本，java版本，CPU、内存等。
   - 数据特征，记录大小等。
   - 测试参数集（多组），系统参数（比如并发数），插件参数（比如batchSize）
   - 不同参数下同步速度（Rec/s, MB/s），机器负载（load, cpu）等，对数据源压力（load, cpu, mem等）。
6. **约束限制**：是否存在其他的使用限制条件。
7. **FAQ**：用户经常会遇到的问题。

## 8. 从源码进行编译

### 8.1 linux

#### 8.1.1 编译依赖

1. golang 1.20以及以上版本

#### 8.1.2 构建
```bash
make dependencies
make release
```

#### 8.1.3 去掉db2依赖

在编译前需要export IGNORE_PACKAGES=db2 

```bash
export IGNORE_PACKAGES=db2
make dependencies
make release
```

### 8.2 windows

####  8.2.1 编译依赖
1. 需要mingw-w64 with gcc 7.2.0以上的环境进行编译
2. golang 1.20以及以上
3. 最小编译环境为win7 

####  8.2.2 构建
```bash
release.bat
```

#### 8.1.3 去掉db2依赖

在编译前需要set IGNORE_PACKAGES=db2

```bash
set IGNORE_PACKAGES=db2
release.bat
```


### 8.3 编译产物

```
    +---go-etl---|---plugin---+---reader--mysql---|--README_zh-CN.md
    |                        | .......
    |                        |
    |                        |---writer--mysql---|--README_zh-CN.md
    |                        | .......
    |
    +---bin----go-etl
    +---exampales---+---csvpostgres----config.json
    |               |---db2------------config.json
    |               | .......
    |
    +---README_USER_zh-CN.md

```
+ go-etl/plugin下是各插件的文档
+ bin下的是数据同步程序dgo-etl
+ exampales下是各场景的数据同步的配置文档
+ README_USER_zh-CN.md是用户使用手册

## 9. 调试http接口

```bash
./go-etl -http :6080 -c examples/limit/config.json
```

### 9.1 获取当前调试数据
使用浏览器访问http://127.0.0.1:6080/debug/pprof获取调试信息
```
/debug/pprof/

Types of profiles available:
Count	Profile
19	allocs
0	block
0	cmdline
18	goroutine
19	heap
0	mutex
0	profile
10	threadcreate
0	trace
full goroutine stack dump
Profile Descriptions:

allocs: A sampling of all past memory allocations
block: Stack traces that led to blocking on synchronization primitives
cmdline: The command line invocation of the current program
goroutine: Stack traces of all current goroutines
heap: A sampling of memory allocations of live objects. You can specify the gc GET parameter to run GC before taking the heap sample.
mutex: Stack traces of holders of contended mutexes
profile: CPU profile. You can specify the duration in the seconds GET parameter. After you get the profile file, use the go tool pprof command to investigate the profile.
threadcreate: Stack traces that led to the creation of new OS threads
trace: A trace of execution of the current program. You can specify the duration in the seconds GET parameter. After you get the trace file, use the go tool trace command to investigate the trace.
```