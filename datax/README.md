# go-etl Data Synchronization Developer Guide

## 1 Introduction to the Synchronization Framework

go-etl is primarily an offline data synchronization framework, structured as follows:

```
readerPlugin(reader) —> Framework(Exchanger+Transformer) -> writerPlugin(riter)  
```

It is built using a Framework + plugin architecture. Data source reading and writing are abstracted into Reader/Writer plugins and integrated into the overall synchronization framework.

+ Reader: The Reader is the data acquisition module, responsible for collecting data from the data source and sending it to the Framework.
+ Writer: The Writer is the data writing module, responsible for continuously fetching data from the Framework and writing it to the destination.
+ Framework: The Framework connects the reader and writer, serving as a data transmission channel between them, and handles core technical issues such as buffering, flow control, concurrency, and data transformation.

## 2 Introduction to the Core Module (core)

A single data synchronization job completed by go-etl is called a Job. When go-etl receives a Job, it starts a process to complete the entire job synchronization process.

The go-etl Job module is the central management node for a single job, responsible for data cleanup, sub-task splitting (converting a single job calculation into multiple sub-Tasks), TaskGroup management, and other functions.

### 2.1 Scheduling Process

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

As shown above, after the go-etl Job starts, it splits the Job into multiple smaller Tasks (sub-tasks) based on different source-side splitting strategies to facilitate concurrent execution. A Task is the smallest unit of a go-etl job, and each Task is responsible for synchronizing a portion of the data. After splitting into multiple Tasks, the go-etl Job calls the Scheduler module. Based on the configured concurrent data volume, the split Tasks are recombined into TaskGroups. The number of Tasks and TaskGroups can be different (N:M). Each TaskGroup is responsible for running all allocated Tasks concurrently with a certain concurrency. The default concurrency for a single TaskGroup is 4. Each Task is started by a TaskGroup. Once a Task starts, it fixes the thread for Reader—>Channel—>Writer to complete the task synchronization work.

When a go-etl job is running, the Job monitors and waits for multiple TaskGroup modules to complete their tasks. Once all TaskGroup tasks are completed, the Job exits successfully. Otherwise, it exits abnormally with a non-zero exit value.

For example, a user submits a go-etl job with 20 concurrencies configured, aiming to synchronize data from 100 MySQL sharded tables to ODPS. The scheduling decision-making process of go-etl is as follows: The go-etl Job splits into 100 Tasks based on the sharding of tables. Based on 20 concurrencies, go-etl calculates that a total of 4 TaskGroups are needed. The 4 TaskGroups evenly distribute the 100 split Tasks, and each TaskGroup is responsible for running 25 Tasks with 5 concurrencies.

+ Job: A Job is go-etl's description of a synchronization job from a source to a destination. It is the smallest business unit for go-etl data synchronization. For example, synchronizing from a MySQL table to a specific partition of an ODPS table.
+ Task: A Task is the smallest execution unit obtained by maximizing the split of a Job. For example, reading a MySQL sharded table with 1024 sharded tables can be split into 1024 read Tasks and executed with several concurrencies.
+ TaskGroup: Describes a set of Task collections. A collection of Tasks executed under the same TaskGroupContainer is called a TaskGroup.
+ JobContainer: The Job executor, responsible for global Job splitting, scheduling, pre-statements, post-statements, and other work units. Similar to JobTracker in Yarn.
+ TaskGroupContainer: The TaskGroup executor, responsible for executing a set of Tasks. Similar to TaskTracker in Yarn.

## 3 Programming Interface

### 3.1 Reader Plugin Interface

The Reader needs to implement the following interfaces:

#### 3.1.1 Job

The Job combines *plugin.BaseJob and implements the following methods:

```golang
    Init(ctx context.Context) (err error)
    Destroy(ctx context.Context) (err error)
    Split(ctx context.Context, number int) ([]*config.JSON, error)
    Prepare(ctx context.Context) error  // Default empty method
    Post(ctx context.Context) error     // Default empty method
```

- `Init`: Initializes the Job object. At this point, the configuration related to this plugin can be obtained through `PluginJobConf()`. The Reader plugin obtains the `reader` part of the configuration.
- `Prepare`: Global preparation work.
- `Split`: Splits the `Task`. The parameter `number` suggests the number of splits, which is generally the concurrency configured during runtime. The return value is a list of `Task` configurations.
- `Post`: Global post-processing work.
- `Destroy`: Destroys the Job object itself.

#### 3.1.2 Task

The Task combines *plugin.BaseTask and implements the following methods:

```golang
    Init(ctx context.Context) (err error)
    Destroy(ctx context.Context) (err error)
    StartRead(ctx context.Context, sender plugin.RecordSender) error 
    Prepare(ctx context.Context) error  // Default empty method
    Post(ctx context.Context) error     // Default empty method
```

- `Init`: Initializes the Task object. At this point, the configuration related to this `Task` can be obtained through `PluginJobConf()`. The configuration here is one of the configuration lists returned by the `Split` method of the `Job`.
- `Prepare`: Local preparation work.
- `StartRead`: Reads data from the data source and writes it to `RecordSender`. `RecordSender` writes the data to the cache queue connecting Reader and Writer.
- `Post`: Local post-processing work.
- `Destroy`: Destroys the Task object itself.

#### 3.1.3 Reader

```golang
    Job() reader.Job
    Task() reader.Task
```

+ `Job`: Gets an instance of the aforementioned Job.
+ `Task`: Gets an instance of the aforementioned Task.

#### 3.1.4 Command Generation

```bash
cd tools/go-etl/plugin
# Adds a new Reader named Mysql. The -p command can be in any case and is used to specify the name of the Reader. If -d is added, it means the original template will be deleted.
go run main.go -t reader -p Mysql
```

This command automatically generates the following Mysql Reader template in datax/plugin/reader to assist in development:

```
    reader---mysql--+-----resources--+--plugin.json
                    |--job.go        |--plugin_job_template.json
                    |--reader.go
                    |--README.md
                    |--task.go
```

As shown below, don't forget to add the developer's name and description in plugin.json:

```json
{
    "name" : "mysqlreader",
    "developer":"Breeze0806",
    "description":"use github.com/go-sql-driver/mysql."
}
```

Additionally, this helps developers avoid compilation errors after using the plugin registration command.

#### 3.1.5 Database

If you want to help implement a data source for a relational database, following these guidelines will make the implementation of your data source more convenient.

##### 3.1.5.1 Database Storage

Refer to the [Database Storage Developer Guide](../storage/database/README.md). This will not only assist you in implementing the Reader plugin interface more quickly but also aid in the implementation of the Writer plugin interface.

##### 3.1.5.2 Database Reader

The dbms reader abstracts the DBWrapper structure of database storage into a Querier as follows and then utilizes the Querier to implement Job and Task functionalities.

```go
// Querier Interface
type Querier interface {
 // Obtains a specific table based on basic table information
 Table(*database.BaseTable) database.Table
 // Checks connectivity
 PingContext(ctx context.Context) error
 // Queries using a query statement
 QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
 // Obtains a specific table based on parameters
 FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error)
 // Retrieves records using parameters and a handler
 FetchRecord(ctx context.Context, param database.Parameter, handler database.FetchHandler) error
 // Retrieves records using parameters, a handler, and a transaction
 FetchRecordWithTx(ctx context.Context, param database.Parameter, handler database.FetchHandler) error
 // Closes resources
 Close() error
}
```

For implementing Job and Reader in the context of MySQL, the Task requires the use of the `dbms.StartRead` function to implement the `StartRead` method.

#### 3.1.6 Two-dimensional Table File Stream

##### 3.1.6.1 Two-dimensional Table File Stream Storage

Refer to the [Two-dimensional Table File Stream Storage Developer Guide](../storage/stream/file/README.md). This will assist you in implementing both the Reader and Writer plugin interfaces more quickly.

##### 3.1.6.2 File Reader

For Tasks and Readers like CSV, independent implementation of Job is required, specifically implementing the `Split` method for splitting and the `Init` method for initialization.

### 3.2 Writer Plugin Interface

Writers need to implement the following interfaces:

#### 3.2.1 Job

The Job combines `*plugin.BaseJob` and implements the following methods:

```golang
 Init(ctx context.Context) (err error)
 Destroy(ctx context.Context) (err error)
 Split(ctx context.Context, number int) ([]*config.JSON, error)
 Prepare(ctx context.Context) error // Default empty method
 Post(ctx context.Context) error    // Default empty method
```

- `Init`: Initializes the Job object. At this point, plugin-related configurations can be obtained through `PluginJobConf()`. The writer section is obtained for the write plugin.
- `Prepare`: Performs global preparation work.
- `Split`: Splits the Task. The parameter `number` suggests the number of splits, generally the configured concurrency level during runtime. The return value is a list of Task configurations.
- `Post`: Performs global post-processing work.
- `Destroy`: Performs destruction work for the Job object itself.

#### 3.2.2 Task

The Task combines `*plugin.BaseTask` and implements the following methods:

```golang
 Init(ctx context.Context) (err error)
 Destroy(ctx context.Context) (err error)
 StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error
 Prepare(ctx context.Context) error     // Default empty method
 Post(ctx context.Context) error        // Default empty method
 SupportFailOver() bool                 // Default empty method
```

- `Init`: Initializes the Task object. At this point, the configuration related to this Task can be obtained through `PluginJobConf()`. The configuration here is one of the configuration lists returned by the Job's `split` method.
- `Prepare`: Performs local preparation work.
- `StartWrite`: Reads data from the `RecordReceiver` and writes it to the target data source. The data in the `RecordReceiver` comes from the cache queue between the Reader and Writer.
- `Post`: Performs local post-processing work for the Task.
- `Destroy`: Performs destruction work for the Task itself.
- `SupportFailOver`: Indicates whether the Task supports failover.

#### 3.2.3 Writer

```golang
 Job() writer.Job
 Task() writer.Task
```

- `Job`: Obtains an instance of the aforementioned Job.
- `Task`: Obtains an instance of the aforementioned Task.

#### 3.2.4 Command Generation

```bash
cd tools/go-etl/plugin
# Adds a new writer named Mysql. The -p command can be in any case and is used to specify the name of the writer. If -d is added, it will delete the original template.
go run main.go -t writer -p Mysql
```

This command automatically generates the following template for a mysql writer in `datax/plugin/writer` to assist in development:

```
    writer--mysql---+-----resources--+--plugin.json
                    |--job.go        |--plugin_job_template.json
                    |--README.md
                    |--task.go
                    |--writer.go
```

Additionally, don't forget to add the developer's name and description to `plugin.json` as follows:

```json
{
    "name" : "mysqlwriter",
    "developer":"Breeze0806",
    "description":"Uses github.com/go-sql-driver/mysql. The database/sql DB executes select SQL and retrieves data from the ResultSet. Warning: The more you know about the database, the fewer problems you will encounter."
}
```

This helps developers avoid compilation errors after using the plugin registration command.

#### 3.2.5 Database

If you want to help implement a data source for a database, following these guidelines will make the implementation of your data source more convenient. However, it is essential that the driver library you use implements the `database/sql` interface of the Golang standard library.

##### 3.2.5.1 Database Storage

Refer to the [Database Storage Developer Guide](../storage/database/README.md). This will assist you in implementing both the Reader and Writer plugin interfaces more quickly.

##### 3.2.5.2 Database Writer

The dbms writer abstracts the DBWrapper structure of database storage into an Execer as follows and then utilizes the Execer to implement Job and Task functionalities.

```go
// Execer Interface
type Execer interface {
 // Obtains a specific table based on basic table information
 Table(*database.BaseTable) database.Table
 // Checks connectivity
 PingContext(ctx context.Context) error
 // Queries using a query statement
 QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
 // Executes a query statement
 ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
 // Obtains a specific table based on parameters
 FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error)
 // Performs batch execution
 BatchExec(ctx context.Context, opts *database.ParameterOptions) error
 // Performs batch execution using prepare/exec
 BatchExecStmt(ctx context.Context, opts *database.ParameterOptions) error
 // Performs batch execution within a transaction
 BatchExecWithTx(ctx context.Context, opts *database.ParameterOptions) error
 // Performs batch execution using prepare/exec within a transaction
 BatchExecStmtWithTx(ctx context.Context, opts *database.ParameterOptions) error
 // Closes the connection
 Close() error
}
```

For implementing Job and Writer in the context of MySQL, the Task requires the use of the `dbms.StartWrite` function to implement the `StartWrite` method.

#### 3.2.6 Two-dimensional Table File Stream

##### 3.2.6.1 Two-dimensional Table File Stream Storage

Refer to the [Two-dimensional Table File Stream Storage Developer Guide](../storage/stream/file/README.md). This will assist you in implementing both the Reader and Writer plugin interfaces more quickly.

##### 3.2.6.2 File Writer

For Tasks and Writers like CSV, independent implementation of Job is required, specifically implementing the `Split` method for splitting and the `Init` method for initialization.

## 4 Plugin Configuration File

`go-etl` uses `json` as the format for its configuration files. A typical `go-etl` task configuration looks like this:

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

In the task configuration, the `value` part of `job.content.reader.parameter` is passed to `Reader.Job`, and the `value` part of `job.content.writer.parameter` is passed to `Writer.Job`. Both `Reader.Job` and `Writer.Job` can access these values through `super.getPluginJobConf()`.

### 4.1 Designing Configuration Parameters

> Designing the configuration file is the first step in plugin development!

The `parameter` sections under `reader` and `writer` in the task configuration are the configuration parameters for the plugins. These configuration parameters should follow these principles:

- Camel Case Naming: All configuration items should use camel case naming, with the first letter lowercase and the first letter of each word capitalized.

- Orthogonality: Configuration items must be orthogonal, with no overlapping functionality and no hidden rules.

- Rich Types: Reasonably use JSON types to reduce unnecessary processing logic and potential for errors.

  - Use the correct data type. For example, use `true`/`false` for bool type values, not `"yes"`/`"true"`/`0`, etc.
  - Reasonably use collection types. For example, use arrays instead of delimited strings.

- Similar and Universal: Follow the conventions of the same type of plugins. For example, the `connection` parameter for databases typically has the following structure:

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

### 4.2 Using the `config.JSON` Struct

```json
{
    "a":{
        "b":[{
            "c":"x"
        }]
    }
}
```

To access the string "x" in `GetConfig`, the path would be a, a.b, a.b.0, a.b.0.c.

Note that because plugins only see a portion of the overall configuration, when using the `json.Config` object, it is important to be aware of the current root path.

For more operations with `json.Config`, please refer to the documentation of the `config` package.

## 5 Plugin Packaging and Release

### 5.1 Adding a License

After developing a feature and before submitting it, please run the following command to automatically add a license and format the code using `gofmt -s -w`:

```bash
go run tools/license/main.go
```

### 5.2 Plugin Registration

Before compiling with Golang, plugins need to be registered within the codebase.

Due to Golang's static compilation, the go-etl framework cannot dynamically load plugins at runtime. Therefore, plugins developed by developers, specifically readers and writers, need to be registered via generated code. The following command facilitates this:

```bash
go generate ./...
```

The main principle is to generate `plugin.go` files from the `plugin.json` resources found in the corresponding go-etl/plugin directory for readers and writers. Additionally, a `plugin.go` file is generated in the go-etl directory to import these plugins. This process is implemented in `tools/go-etl/build`. Optionally, the `-i` command can be used to ignore compiling certain data sources, such as DB2, which uses ODBC for database access and requires additional Linux dependencies.

## 6. Plugin Data Transfer

Similar to the typical "producer-consumer" pattern, data transfer between the `Reader` and `Writer` plugins occurs through `channels`. These channels can be in-memory or persistent, and plugins do not need to concern themselves with the implementation details. Plugins write data to the channel using `RecordSender` and read data from the channel using `RecordReceiver`.

A single item in the channel is a `Record` object, which can hold multiple `Column` objects. This can be simply understood as a record and its columns in a database. For more details on the `Record` prototype, refer to the "Records" chapter in the [documentation](../element/README.md).

Since `Record` is an interface, the `Reader` plugin first calls `RecordSender.createRecord()` to create a `Record` instance and then adds `Column` objects to it.

The `Writer` plugin calls `RecordReceiver.getFromReader()` to retrieve a `Record` and then iterates over the `Column` objects to write them to the target storage. While the `Reader` is still active and transmission is ongoing, if there is no data currently available, `RecordReceiver.getFromReader()` will block until data becomes available. Once transmission has ended, it will return `ErrTerminate`, allowing the `Writer` plugin to determine when to end its `startWrite` method.

### 6.1 Data Type Conversion

To standardize data type conversion operations between the source and destination, and to ensure data fidelity, go-etl supports six internal data types. For more details, refer to the "Data Type Conversion" chapter in the [documentation](../element/README.md).

## 7. Plugin Documentation

Include the following chapters in the plugin's README.md documentation:

1. **Quick Introduction**: Describes the plugin's use cases and features.
2. **Implementation Principles**: Explains the underlying principles of the plugin's implementation.
3. **Configuration Instructions**:
	* Provides a sample JSON configuration file for a typical synchronization task.
	* Describes the meaning, requirement, default value, range, and other constraints of each parameter.
4. **Type Conversion**:
	* Explains how the plugin converts between the actual storage type and go-etl's internal type.
	* Mentions any special handling.
5. **Performance Report**:
	* Details the hardware and software environment, system version, Java version, CPU, memory, etc.
	* Describes the data characteristics, such as record size.
	* Lists the test parameter sets, system parameters (e.g., concurrency), and plugin parameters (e.g., batchSize).
	* Provides synchronization speeds (Rec/s, MB/s), machine load (load, cpu), and the impact on the data source (load, cpu, mem) for different parameters.
6. **Constraints and Limitations**: Mentions any additional usage restrictions.
7. **FAQ**: Addresses commonly asked questions by users.

## 8. Compiling from Source Code

### 8.1 Linux

#### 8.1.1 Compilation Dependencies

1. golang 1.20 or higher.

#### 8.1.2 Building

```bash
make dependencies
make release
```

#### 8.1.3 Excluding DB2 Dependencies

Before compiling, set the `IGNORE_PACKAGES` environment variable to `db2`:

```bash
export IGNORE_PACKAGES=db2
make dependencies
make release
```

### 8.2 Windows

#### 8.2.1 Compilation Dependencies

1. A MinGW-w64 environment with GCC 7.2.0 or higher is required for compilation.
2. golang 1.20 or higher.
3. The minimum supported compilation environment is Windows 7.

#### 8.2.2 Building

```bash
release.bat
```

#### 8.2.3 Excluding DB2 Dependencies

Before compiling, set the `IGNORE_PACKAGES` environment variable to `db2`:

```bash
set IGNORE_PACKAGES=db2
release.bat
```

### 8.3 Compilation Output

```
+---datax---|---plugin---+---reader--mysql---|--README.md
|                        | .......
|                        |
|                        |---writer--mysql---|--README.md
|                        | .......
|
+---bin----datax
+---examples---+---csvpostgres----config.json
|               |---db2------------config.json
|               | .......
|
+---README_USER.md
```

* The `datax/plugin` directory contains documentation for each plugin.
* The `bin` directory contains the data synchronization program `datax`.
* The `examples` directory contains configuration files for various data synchronization scenarios.
* `README_USER.md` is the user manual.

## 9. Debugging HTTP Interfaces

```bash
datax -http :8443 -c examples/limit/config.json
```

### 9.1 Accessing Current Debug Data

Use a web browser to access `http://127.0.0.1:8443/debug/pprof` to retrieve debug information.

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