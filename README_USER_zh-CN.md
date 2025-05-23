# go-etl数据同步用户手册

go-etl是一个数据同步工具，目前支持MySQL,postgres,oracle,SQL SERVER,Sqlite3,DB2等主流关系型数据库以及csv，xlsx文件之间的数据同步。

## 1 如何获取

参考[项目文档](README_zh-CN.md)，从获取二进制程序开始，从源代码开始或者从编译docker镜像开始。

## 2 如何开始

获取go-etl二进制程序或docker镜像，在linux下如Makefile所示`export LD_LIBRARY_PATH=/home/ibmdb/clidriver/lib`，这个库从[ibm db2](https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli)下载，否则无法运行。

另外oracle需要下载[oracle](https://www.oracle.com/database/technologies/instant-client/downloads.html)下载到对应64位版本odbc依赖，也可以在**QQ群185188648**群共享中中下载到。

### 2.1 单任务数据同步
调用go-etl十分简单，只要直接调用它即可

```bash
./go-etl -c config.json
```
`-c` 指定数据源配置文件

当返回值是`0`，并且在最后显示`run success`,表示执行成功

当返回值是`1`，并且在最后显示`run fail`,并告知执行失败的原因

运行时会输出

```bash
datax_channel_total_byte(job_id=1,task_group_id=0,task_id=0) 18.97 MiB                                       [   =]                                        3.61 MiB/s-5s
datax_channel_total_record(job_id=1,task_group_id=0,task_id=0) 1000000                                       [   =]                                        192145.7/s-5s
datax_channel_byte(job_id=1,task_group_id=0,task_id=0) 0.00 b                                                 [   =]                                                  0s
datax_channel_record(job_id=1,task_group_id=0,task_id=0) 0                                                   [   =]                                                   0s
datax_channel_total_byte(job_id=1,task_group_id=0,task_id=1) 18.97 MiB                                       [   =]                                        2.72 MiB/s-6s
datax_channel_total_record(job_id=1,task_group_id=0,task_id=1) 1000000                                       [   =]                                        146067.9/s-6s
datax_channel_byte(job_id=1,task_group_id=0,task_id=1) 0.00 b                                                 [   =]                                                  0s
datax_channel_record(job_id=1,task_group_id=0,task_id=1) 0                                                   [   =]                                                   0s
```

上述开头的部分的是参数名称

- `datax_channel_total_byte`总共数据同步的字节数
- `datax_channel_total_record`总共数据同步的记录数
- `datax_channel_byte`,在通道里数据同步的字节数
- `datax_channel_record`在通道里数据同步的记录数

括号中的`job_id=1,task_group_id=0,task_id=0`能够标识那个任务

+ `job_id` 工作号
+ `task_group_id ` 任务组号
+ `task_id` 任务号

#### 2.1.1 数据源配置文件

数据源配置文件是json文件，使用数据源相互组合，如从mysql同步到postgres中

```json
{
    "core" : {
        "container": {
            "job":{
                "id": 1,
                "sleepInterval":100
            }
        }
    },
    "job":{
        "content":[
            {
                "reader":{
                    "name": "mysqlreader",
                    "parameter": {
                        "username": "test:",
                        "password": "test:",
                        "column": ["*"],
                        "connection":  {
                                "url": "tcp(192.168.15.130:3306)/source?parseTime=false",
                                "table": {
                                    "db":"source",
                                    "name":"type_table"
                                }
                            },
                        "where": ""
                    }
                },
                "writer":{
                    "name": "postgreswriter",
                    "parameter": {
                        "username": "postgres",
                        "password": "123456",
                        "writeMode": "insert",
                        "column": ["*"],
                        "preSql": [],
                        "connection":  {
                                "url": "postgres://192.168.15.130:5432/postgres?sslmode=disable&connect_timeout=2",
                                "table": {
                                    "schema":"destination",
                                    "name":"type_table"
                                }
                         },
                        "batchTimeout": "1s",
                        "batchSize":1000
                    }
                },
               "transformer":[]
            }
        ]
    }
}
```

`reader`和`writer`的配置如下：

| 类型         | 数据源             | Reader（读） | Writer(写) | 文档                                                         |
| ------------ | ------------------ | ------------ | ---------- | ------------------------------------------------------------ |
| 关系型数据库 | MySQL/Mariadb/Tidb | √            | √          | [读](datax/plugin/reader/mysql/README_zh-CN.md)、[写](datax/plugin/writer/mysql/README_zh-CN.md) |
|              | Postgres/Greenplum | √            | √          | [读](datax/plugin/reader/postgres/README_zh-CN.md)、[写](datax/plugin/writer/postgres/README_zh-CN.md) |
|              | DB2 LUW            | √            | √          | [读](datax/plugin/reader/db2/README_zh-CN.md)、[写](datax/plugin/writer/db2/README_zh-CN.md) |
|              | SQL Server         | √            | √          | [读](datax/plugin/reader/sqlserver/README_zh-CN.md)、[写](datax/plugin/writer/sqlserver/README_zh-CN.md) |
|              | Oracle             | √            | √          | [读](datax/plugin/reader/oracle/README_zh-CN.md)、[写](datax/plugin/writer/oracle/README_zh-CN.md) |
|              | Sqlite3            | √            | √          | [读](datax/plugin/reader/sqlite3/README.md)、[写](datax/plugin/writer/sqlite3/README.md) |
| 无结构流     | CSV                | √            | √          | [读](datax/plugin/reader/csv/README_zh-CN.md)、[写](datax/plugin/writer/csv/README_zh-CN.md) |
|              | XLSX（excel）      | √            | √          | [读](datax/plugin/reader/xlsx/README_zh-CN.md)、[写](datax/plugin/writer/xlsx/README_zh-CN.md) |

#### 2.1.2 使用示例

注意在linux下如Makefile所示export LD_LIBRARY_PATH=${DB2HOME}/lib

##### 2.1.2.1 使用mysql同步

- 使用cmd/datax/examples/mysql/init.sql初始化数据库**用于测试**
- 开启同步mysql命令

```bash
./go-etl -c examples/mysql/config.json
```

##### 2.1.2.2 使用postgres同步

- 使用cmd/datax/examples/postgres/init.sql初始化数据库**用于测试**
- 开启同步postgres命令

```bash
./go-etl -c examples/postgres/config.json
```

##### 2.1.2.3 使用db2同步

- 注意使用前请下载相应的db2的odbc库，如linux的make dependencies和release.bat
- 注意在linux下如Makefile所示export LD_LIBRARY_PATH=${DB2HOME}/lib
- 注意在windows下如release.bat所示set path=%path%;%GOPATH%\src\github.com\ibmdb\go_ibm_db\clidriver\bin
- 使用cmd/datax/examples/db2/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
./go-etl -c examples/db2/config.json
```

##### 2.1.2.4 使用oracle同步

- 注意使用前请下载相应的[Oracle Instant Client]( https://www.oracle.com/database/technologies/instant-client/downloads.html)，例如，连接oracle 11g最好下载12.x版本。
- 注意在linux下如export LD_LIBRARY_PATH=/opt/oracle/instantclient_21_1:$LD_LIBRARY_PATH，另需要安装libaio
- 注意在windows下如set path=%path%;%GOPATH%\oracle\instantclient_21_1，
Oracle Instant Client 19不再支持windows7
- 使用cmd/datax/examples/oracle/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
./go-etl -c examples/oracle/config.json
```

##### 2.1.2.5 使用sql server同步

- 使用cmd/datax/examples/sqlserver/init.sql初始化数据库**用于测试**
- 开启同步sql server命令

```bash
./go-etl -c examples/sqlserver/config.json
```

##### 2.1.2.6 使用csv同步到postgres

- 使用cmd/datax/examples/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
./go-etl -c examples/csvpostgres/config.json
```

##### 2.1.2.7 使用xlsx同步到postgres

- 使用cmd/datax/examples/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
./go-etl -c examples/xlsxpostgres/config.json
```

##### 2.1.2.8 使用postgres同步csv

- 使用cmd/datax/examples/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
./go-etl -c examples/postgrescsv/config.json
```

##### 2.1.2.9 使用postgres同步xlsx

- 使用cmd/datax/examples/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
./go-etl -c examples/postgresxlsx/config.json
```

##### 2.1.2.10 与 sqlite3 同步

* 在使用前，请下载相应的[SQLite驱动](https://www.sqlite.org/download.html). 
* 注意：在 Windows 系统上，设置 `path=%path%;/opt/sqlite/sqlite3.dll`。
* 使用 `cmd/datax/examples/sqlite3/init.sql` **用于测试目的** 初始化数据库
* 在 `examples/sqlite3/config.json` 文件中，`url` 表示 sqlite3 数据库文件的路径。在 Windows 系统上，它可以是 `E:\sqlite3\test.db`，而在 Linux 系统上，它可以是 `/sqlite3/test.db`。
* 启动 sqlite3 同步命令：

##### 2.1.2.11 其他同步例子

除了上述例子外，在go-etl特性中所列出的数据源都可以交叉使用，还配置例如mysql到postgresql数据源，mysql到oracle,oracle到db2等等，

#### 2.1.3 数据库全局配置

```json
{
    "job":{
        "setting":{
            "pool":{
              "maxOpenConns":8,
              "maxIdleConns":8,
              "connMaxIdleTime":"40m",
              "connMaxLifetime":"40m"
            },
            "retry":{
              "type":"ntimes",
              "strategy":{
                "n":3,
                "wait":"1s"
              },
              "ignoreOneByOneError":true
            }
        }
    }
}
```
##### 2.1.3.1 连接池pool
+ maxOpenConns: 最大连接打开数
+ maxIdleConns: 最大空闲连接打开数
+ connMaxIdleTime： 最大空闲时间
+ connMaxLifetime： 最大生存时间

##### 2.1.3.2 重试retry
ignoreOneByOneError 是否忽略一个个重试错误
+ 重试类型type和重试策略
1. 类型有`ntimes`,指n次数重复重试策略,`"strategy":{"n":3,"wait":"1s"}`,n代表重试次数，wait代表等待时间
2. 类型有`forever`,指永久重复重试策略,`"strategy":{"wait":"1s"}`,wait代表等待时间
3. 类型有`exponential`,指幂等重复重试策略,`"strategy":{"init":"100ms","max":"4s"}`,init代表开始时间，max代表最大时间

##### 2.1.3.3 测试数据

```bash
./go-etl -c examples/global/config.json
```

#### 2.1.4 使用切分键

这里假设数据按切分键分布是均匀的，合理使用这样的切分键可以使同步更快，另外为了加快对最大值和最小值的查询，这里对于大表可以预设最大最小值

##### 2.1.4.1 测试方式
- 使用程序生成mysql数据产生split.csv
```bash
cd cmd/datax/examples/split
go run main.go
```
- 使用init.sql建表
- 同步至mysql数据库
```bash
cd ../..
./go-etl -c examples/split/csv.json
```
- 修改examples/split/config.json的split的key为id,dt,str
- mysql数据库切分同步整形，日期，字符串类型
```bash
./go-etl -c examples/split/config.json
```

#### 2.1.5 使用preSql和postSql

preSql和postSql分别是写入数据前和写入数据后的sql语句组

##### 2.1.5.1 测试方式
在本例子中，采用了全量导入的方式
1.写入数据前先建立了一个临时表
2.在写入数据后，将原表删除，将临时表重名为新表

```bash
./go-etl -c examples/prePostSql/config.json
```

#### 2.1.6 流控配置

之前speed的byte和record配置并不会生效，现在加入流控特性后，byte和record将会生效，byte会限制缓存消息字节数，而record会限制缓存消息条数，如果byte设置过小会导致缓存过小而导致同步数据失败。当byte为0或负数时，限制器将不会工作, 例如byte为10485760，即10Mb(10x1024x1024)。

```json
{
    "job":{
        "setting":{
            "speed":{
                "byte":10485760,
                "record":1024,
                "channel":4
            }
        }
    }
}    
```
##### 2.1.6.1 流控测试
- 使用程序生成src.csv,发起流控测试
```bash
cd cmd/datax/examples/limit
go run main.go
cd ../..
./go-etl -c examples/limit/config.json
```

#### 2.1.7 querySql配置

数据库读取器使用querySql去查询数据库


##### 2.1.7.1 querySql配置测试

```bash
./go-etl -c examples/querySql/config.json
```

### 2.2 多任务数据同步

#### 2.2.1 使用方式

##### 2.2.1.1 数据源配置文件

配置数据源配置文件，如从mysql同步到postgres中

```json
{
    "core" : {
        "container": {
            "job":{
                "id": 1,
                "sleepInterval":100
            }
        }
    },
    "job":{
        "content":[
            {
                "reader":{
                    "name": "mysqlreader",
                    "parameter": {
                        "username": "test:",
                        "password": "test:",
                        "column": ["*"],
                        "connection":  {
                                "url": "tcp(192.168.15.130:3306)/source?parseTime=false",
                                "table": {
                                    "db":"source",
                                    "name":"type_table"
                                }
                            },
                        "where": ""
                    }
                },
                "writer":{
                    "name": "postgreswriter",
                    "parameter": {
                        "username": "postgres",
                        "password": "123456",
                        "writeMode": "insert",
                        "column": ["*"],
                        "preSql": [],
                        "connection":  {
                                "url": "postgres://192.168.15.130:5432/postgres?sslmode=disable&connect_timeout=2",
                                "table": {
                                    "schema":"destination",
                                    "name":"type_table"
                                }
                         },
                        "batchTimeout": "1s",
                        "batchSize":1000
                    }
                },
               "transformer":[]
            }
        ]
    }
}
```

##### 2.2.1.2 源目的配置向导文件

源目的配置向导文件是csv文件，每行配置可以配置如下:

```csv
path[table],path[table]
```

每一列可以是路径或者是表名，注意所有的表要配置库名或者模式名，需要在数据源配置文件配置。

##### 2.2.1.3 批量生成数据配置集和执行脚本

```bash
./go-etl -c tools/testData/xlsx.json -w tools/testData/wizard.csv 
```
-c 指定数据源配置文件 -w 指定源目的配置向导文件。

执行结果会在数据源配置文件目录文件生成源目的配置向导文件行数的配置集，分别以指定数据源配置文件1.json,指定数据源配置文件2.json,...,指定数据源配置文件[n].json的配置集。

另外，在当前目录会生成执行脚本run.bat或者run.sh。

##### 2.2.1.4 批量生成数据配置集和执行脚本

###### windows

```bash
run.bat
```

linux

```bash
run.sh
```

#### 2.2.2 测试结果
可以运行cmd/datax/testData的测试数据
```bash
cd cmd/datax
./go-etl -c testData/xlsx.json -w testData/wizard.csv 
```
结果会在testData下生成wizard.csv行数的配置文件，分别以xlsx1.json,xlsx2.json,...,xlsx[n].json的配置集。

### 2.3 数据同步帮助手册

#### 2.3.1 帮助命令

```
./go-etl -h
```

帮助显示

```bash
Usage of go-etl:
  -c string
        config (default "config.json")
  -http string
        http
  -w string
        wizard
```

-http 新增监听端口，如:6080, 开启后访问127.0.0.1:6080/metrics获取实时的吞吐量

#### 2.3.2 查看版本

```bash
./go-etl version
```

显示`版本号`(git commit:  `git提交号`） complied by go version `go版本号`

```bash
v0.1.0 (git commit: c82eb302218f38cd3851df4b425256e93f85160d) complied by go version go1.16.5 windows/amd64
```

#### 2.3.3 开启监控端口
```bash
./go-etl -http :6080 -c examples/limit/config.json
```

##### 2.3.3.1 获取当前监控数据

使用浏览器访问http://127.0.0.1:6080/metrics获取类似[prometheus exporters](https://prometheus.io/docs/instrumenting/writing_exporters/)的监控数据

```bash
# HELP datax_channel_byte the number of bytes currently being synchronized in the channel
# TYPE datax_channel_byte gauge
datax_channel_byte{job_id="1",task_group_id="0",task_id="0"} 20480
datax_channel_byte{job_id="1",task_group_id="0",task_id="1"} 20500
# HELP datax_channel_record the number of records currently being synchronized in the channel
# TYPE datax_channel_record gauge
datax_channel_record{job_id="1",task_group_id="0",task_id="0"} 1024
datax_channel_record{job_id="1",task_group_id="0",task_id="1"} 1025
# HELP datax_channel_total_byte the total number of bytes synchronized
# TYPE datax_channel_total_byte counter
datax_channel_total_byte{job_id="1",task_group_id="0",task_id="0"} 2.75355e+06
datax_channel_total_byte{job_id="1",task_group_id="0",task_id="1"} 5.29381e+06
# HELP datax_channel_total_record the total number of records synchronized
# TYPE datax_channel_total_record counter
datax_channel_total_record{job_id="1",task_group_id="0",task_id="0"} 143233
datax_channel_total_record{job_id="1",task_group_id="0",task_id="1"} 270246
```

另外使用http://127.0.0.1:6080/metrics?t=json也能获取`json`格式的监控数据

```json
{
    "jobID": 1,
    "metrics": [
        {
            "taskGroupID": 0,
            "metrics": [
                {
                    "taskID": 0,
                    "channel": {
                        "totalByte": 7069190,
                        "totalRecord": 359015,
                        "byte": 20500,
                        "record": 1025
                    }
                },
                {
                    "taskID": 1,
                    "channel": {
                        "totalByte": 13245910,
                        "totalRecord": 667851,
                        "byte": 20460,
                        "record": 1023
                    }
                }
            ]
        }
    ]
}
```

- `totalByte` 即`datax_channel_total_byte`总共数据同步的字节数
- `totalRecord` 即`datax_channel_total_record`总共数据同步的记录数
- `byte` 即`datax_channel_byte`,在通道里数据同步的字节数
- `record` 即`datax_channel_record`在通道里数据同步的记录数