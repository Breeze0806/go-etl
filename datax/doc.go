// Package datax 主要离线数据同步框架，框架如下
//        Database->readerPlugin（reader）->Framework(Exchanger+Transformer) ->writerPlugin（writer）->Database
//
// 采用Framework + plugin架构构建。将数据源读取和写入抽象成为Reader/Writer插件，纳入到整个同步框架中。
// Reader：Reader为数据采集模块，负责采集数据源的数据，将数据发送给Framework。
// Writer： Writer为数据写入模块，负责不断向Framework取数据，并将数据写入到目的端。
// Framework：Framework用于连接reader和writer，作为两者的数据传输通道，并处理缓冲，流控，并发，数据转换等核心技术问题
//	JOB--split--+-- task1--+           +--taskGroup1-------+______________________________________________
//	            |-- task2--|           |--taskGroup2       |   Reader1->Exchanger1(Transformer)->Writer1  |
//	            |-- task3--|-schedule--|--taskGroup3       |   Reader2->Exchanger2(Transformer)->Writer2  |
//	               ......                 ......           |   Reader3->Exchanger3(Transformer)->Writer3  |
//	            |-- taskN--|           |--taskGroupM       |          ......                              |
//	                                                       |   ReaderN->ExchangerN(Transformer)->WriterN  |
//	                                                       |______________________________________________|
//
//
// 核心模块介绍:
//
// DataX完成单个数据同步的作业，我们称之为Job，DataX接受到一个Job之后，将启动一个进程来完成整个作业同步过程。DataX Job模块是单个作业的中枢管理节点，承担了数据清理、子任务切分(将单一作业计算转化为多个子Task)、TaskGroup管理等功能。
// DataXJob启动后，会根据不同的源端切分策略，将Job切分成多个小的Task(子任务)，以便于并发执行。Task便是DataX作业的最小单元，每一个Task都会负责一部分数据的同步工作。
// 切分多个Task之后，DataX Job会调用Scheduler模块，根据配置的并发数据量，将拆分成的Task重新组合，组装成TaskGroup(任务组)。每一个TaskGroup负责以一定的并发运行完毕分配好的所有Task，默认单个任务组的并发数量为5。
// 每一个Task都由TaskGroup负责启动，Task启动后，会固定启动Reader—>Channel—>Writer的线程来完成任务同步工作。
// DataX作业运行起来之后， Job监控并等待多个TaskGroup模块任务完成，等待所有TaskGroup任务完成后Job成功退出。否则，异常退出，进程退出值非0
//
// DataX调度流程：
// 举例来说，用户提交了一个DataX作业，并且配置了20个并发，目的是将一个100张分表的mysql数据同步到odps里面。 DataX的调度决策思路是：
//
// DataXJob根据分库分表切分成了100个Task。
// 根据20个并发，DataX计算共需要分配4个TaskGroup。
// 4个TaskGroup平分切分好的100个Task，每一个TaskGroup负责以5个并发共计运行25个Task。
//
// Job: Job是DataX用以描述从一个源头到一个目的端的同步作业，是DataX数据同步的最小业务单元。比如：从一张mysql的表同步到odps的一个表的特定分区。
// Task: Task是为最大化而把Job拆分得到的最小执行单元。比如：读一张有1024个分表的mysql分库分表的Job，拆分成1024个读Task，用若干个并发执行。
// TaskGroup:  描述的是一组Task集合。在同一个TaskGroupContainer执行下的Task集合称之为TaskGroup
// JobContainer:  Job执行器，负责Job全局拆分、调度、前置语句和后置语句等工作的工作单元。类似Yarn中的JobTracker
// TaskGroupContainer: TaskGroup执行器，负责执行一组Task的工作单元，类似Yarn中的TaskTracker。
//
// Reader需要实现以下接口:
//
// 1. Reader的Job组合*plugin.BaseJob，实现方法Init(ctx context.Context) (err error)
// Destroy(ctx context.Context) (err error)，Split(ctx context.Context, number int) ([]*config.JSON, error)
// Prepare(ctx context.Context) error以及Post(ctx context.Context) error
//
// 2. Reader的Task组合*plugin.BaseTash,实现方法Init(ctx context.Context) (err error)
// Destroy(ctx context.Context) (err error)，StartRead(ctx context.Context, sender plugin.RecordSender) error
// Prepare(ctx context.Context) error以及Post(ctx context.Context) error
//
// 3. Reader的本身实现Job() reader.Job以及Task() reader.Task
//
// Writer 需要实现以下接口:
//
// 1. Writer的Job组合*plugin.BaseJob，实现方法Init(ctx context.Context) (err error)，
// Destroy(ctx context.Context) (err error)，Split(ctx context.Context, number int) ([]*config.JSON, error)，
// Prepare(ctx context.Context) error以及Post(ctx context.Context) error
//
// 2. Writer的Task组合*plugin.BaseTash,实现方法Init(ctx context.Context) (err error)
// Destroy(ctx context.Context) (err error)，StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error
// ，Prepare(ctx context.Context) error，Post(ctx context.Context) error以及SupportFailOver() bool
//
// 3. Writer的本身实现Job() writer.Job以及Task() writer.Task。
//
// 实现上述接口后，需要将上述的reader和writer通过loader.RegisterReader以及loader.RegisterWriter注册到
// 这些插件通过配置文件以下json配置描述：
//
//	{
//	     "name" : "mysqlreader",
//	     "developer":"Breeze0806",
//	     "description":"use github.com/go-sql-driver/mysql. database/sql DB execute select sql, retrieve data from the ResultSet. warn: The more you know about the database, the less problems you encounter."
//	}
//
// 上述接口配置目录按照以下
//
//	plugin+--- reader--mysql---+-----resources--+--plugin.json
//	      |                    |--job.go        |--plugin_job_template.json
//	      |                    |--reader.go
//	      |                    |--README.md
//	      |                    |--task.go
//	      |
//	      +--- reader--mysql---+-----resources--+--plugin.json
//	                           |--job.go        |--plugin_job_template.json
//	                           |--reader.go
//	                           |--README.md
//	                           |--task.go
package datax
