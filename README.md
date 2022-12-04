# go-etl
[![Go Report Card][report-img]][report][![GoDoc][doc-img]][doc][![LICENSE][license-img]][license][![Build Status][ci-img]][ci][![Coverage Status][cov-img]][cov]

go-etl是一个集数据源抽取，转化，加载，同步校验的工具集，提供强大的数据同步，数据转化甚至数据校验的功能。

go-etl将提供的etl能力如下：

1. 主流数据库的数据抽取以及数据加载的能力，在storage包中实现
2. 类二维表的数据流的数据抽取以及数据加载的能力，在stream包中实现
3. 类似datax的数据同步能力，在datax包中实现
4. 以mysql sql语法为基础的数据筛选、转化能力，在transform包中实现
5. 数据库间的数据校验能力，在libra包中实现 (目前计划中)

鉴于本人实在精力有限，欢迎大家来提交issue或者***加QQ群185188648***来讨论go-etl，让我们一起进步!

## 数据同步工具

本数据数据同步工具以下数据源的同步能力

| 类型         | 数据源             | Reader（读） | Writer(写) | 文档                                                         |
| ------------ | ------------------ | ------------ | ---------- | ------------------------------------------------------------ |
| 关系型数据库 | MySQL/Mariadb/Tidb | √            | √          | [读](datax/plugin/reader/mysql/README.md)、[写](datax/plugin/writer/mysql/README.md) |
|              | Postgres/Greenplum | √            | √          | [读](datax/plugin/reader/postgres/README.md)、[写](datax/plugin/writer/postgres/README.md) |
|              | DB2 LUW            | √            | √          | [读](datax/plugin/reader/db2/README.md)、[写](datax/plugin/writer/db2/README.md) |
|              | SQL Server            | √            | √          | [读](datax/plugin/reader/sqlserver/README.md)、[写](datax/plugin/writer/sqlserver/README.md) |
|              | Oracle            | √            | √          | [读](datax/plugin/reader/oracle/README.md)、[写](datax/plugin/writer/oracle/README.md) |
| 无结构流     | CSV                | √            | √          | [读](datax/plugin/reader/csv/README.md)、[写](datax/plugin/writer/csv/README.md) |
|              | XLSX（excel）      | √            | √          | [读](datax/plugin/reader/xlsx/README.md)、[写](datax/plugin/writer/xlsx/README.md) |

### 安装和使用

在使用下列命令前请确保你已经安装go的编译环境并且设置好了GOPATH

#### linux

```bash
make dependencies
make release
```

#### windows

需要mingw-w64 with gcc 7.2.0以上的环境进行编译

```bash
release.bat
```


### 使用方式

```bash
Usage of datax.exe:
  -c string
        config (default "config.json")
  -w string
         (default "wizard.csv")
```

#### 批量生成配置集

```bash
datax -c testData/xlsx.json -w testData/wizard.csv 
```
-c 指定数据源配置文件 -w 指定源目的配置向导文件。

##### 数据源配置文件

数据源配置文件是json文件，使用数据源相互组合，如从mysql同步到postgres中
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

##### 源目的配置向导文件
源目的配置向导文件是csv文件，每行配置可以配置如下:
```csv
path[table],path[table]
```
每一列可以是路径或者是表名，注意所有的表要配置库名或者模式名，需要在数据源配置文件配置。

#### 数据同步

调用datax十分简单，只要直接调用它即可

```bash
data -c config.json
```
-c 指定数据源配置文件

当返回值是0，并且显示run success,表示执行成功

当返回值是1，并且显示run fail,并告知执行失败的原因

#### 数据库全局配置

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
##### 连接池pool
+ maxOpenConns: 最大连接打开数
+ maxIdleConns: 最大空闲连接打开数
+ connMaxIdleTime： 最大空闲时间
+ connMaxLifetime： 最大生存时间

##### 重试retry
ignoreOneByOneError 是否忽略一个个重试错误
+ 重试类型type和重试策略
1. 类型有`ntimes`,指n次数重复重试策略,`"strategy":{"n":3,"wait":"1s"}`,n代表重试次数，wait代表等待时间
2. 类型有`forever`,指永久重复重试策略,`"strategy":{"wait":"1s"}`,wait代表等待时间
3. 类型有`exponential`,指幂等重复重试策略,`"strategy":{"init":"100ms","max":"4s"}`,init代表开始时间，max代表最大时间

### 使用示例

#### 使用mysql同步

- 使用cmd/datax/examples/mysql/init.sql初始化数据库**用于测试**
- 开启同步mysql命令

```bash
datax -c examples/mysql/config.json
```

#### 使用postgres同步

- 使用cmd/datax/examples/postgres/init.sql初始化数据库**用于测试**
- 开启同步postgres命令

```bash
datax -c examples/postgres/config.json
```

#### 使用db2同步

- 注意使用前请下载相应的db2的odbc库，如linux的make dependencies和release.bat
- 注意在linux下如Makefile所示export LD_LIBRARY_PATH=${DB2HOME}/lib
- 注意在windows下如release.bat所示set path=%path%;%GOPATH%\src\github.com\ibmdb\go_ibm_db\clidriver\bin
- 使用cmd/datax/examples/db2/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/db2/config.json
```

#### 使用oracle同步

- 注意使用前请下载相应的[Oracle Instant Client]( https://www.oracle.com/database/technologies/instant-client/downloads.html)，例如，连接oracle 11g最好下载12.x版本。
- 注意在linux下如export LD_LIBRARY_PATH=/opt/oracle/instantclient_21_1:$LD_LIBRARY_PATH，另需要安装libaio
- 注意在windows下如set path=%path%;%GOPATH%\oracle\instantclient_21_1，
Oracle Instant Client 19不再支持windows7
- 使用cmd/datax/examples/oracle/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/oracle/config.json
```

#### 使用sql server同步

- 使用cmd/datax/examples/sqlserver/init.sql初始化数据库**用于测试**
- 开启同步sql server命令

```bash
datax -c examples/sqlserver/config.json
```

#### 使用csv同步到postgres

- 使用cmd/datax/examples/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/csvpostgres/config.json
```

#### 使用xlsx同步到postgres

- 使用cmd/examples/datax/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/xlsxpostgres/config.json
```

#### 使用postgres同步csv

- 使用cmd/datax/examples/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/postgrescsv/config.json
```

#### 使用postgres同步xlsx

- 使用cmd/datax/examples/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/postgresxlsx/config.json
```

#### 其他数据源

你可以选择：
1. 提交issue让别人帮你开发来支持新的数据源
2. 参考开发宝典自己开发来支持新的数据源

### 开发宝典

可以参考[Datax开发者文档](datax/README.md)来帮助开发

## 模块简介
### datax

本包将提供类似于阿里巴巴[DataX](https://github.com/alibaba/DataX)的接口去实现go的etl框架，目前主要实现了job框架内的数据同步能力，监控等功能还未实现.

#### plan

- [x] 实现db2数据库reader/writer插件
- [x] 实现sql server数据库reader/writer插件
- [X] 实现oracle数据库reader/writer插件
- [x] 实现cvs文件reader/writer插件
- [x] 实现xlsx文件reader/writer插件
- [ ] 实现监控模块
- [ ] 实现流控模块
- [x] 实现关系型数据库入库重试保证数据不丢失功能
- [ ] 实现关系型数据库入库断点续传

### storage

#### database

目前已经实现了数据库的基础集成，已有mysql和postgresql的实现，如何实现可以查看godoc文档，利用它能非常方便地实现datax数据库间的同步，欢迎大家来提交新的数据同步方式，可以在下面选择新的数据库来同步

##### plan

- [x] 实现db2数据库的dialect 
- [x] 实现sql server数据库的dialect
- [X] 实现oracle数据库的dialect

#### stream

主要用于字节流的解析，如文件，消息队列，elasticsearch等，字节流格式可以是cvs，json, xml等

#### file
##### plan

- [x] 实现文件流的数据传输框架
- [x] 单元测试文件流的数据传输框架
- [x] 实现cvs文件字节流的数据传输框架
- [x] 单元测试cvs文件字节流的数据传输框架
- [x] 实现xlsx文件字节流的数据传输框架
- [x] 单元测试xlsx文件字节流的数据传输框架

#### mq

##### plan

暂无时间安排计划，欢迎来实现

#### elasticsearch

##### plan

暂无时间安排计划，欢迎来实现

### transform

主要用于类sql数据转化

#### plan

- [ ] 引入tidb数据库的mysql解析能力
- [ ] 引入tidb数据库的mysql函数计算能力
- [ ] 运用mysql解析能力和mysql函数计算能力实现数据转化能力

### tools

工具集用于编译，新增许可证等

#### datax

##### build

发布命令，用于将由开发者开发的reader和writer插件注册到程序中的代码

##### plugin

数据源插件模板新增工具，用于新增一个reader或writer模板，配合发布命令使用，减少开发者负担

#### license

用于自动新增go代码文件中许可证

```bash
go run tools/license/main.go
```

### libra

主要用于数据库间数据校验

[report-img]:https://goreportcard.com/badge/github.com/Breeze0806/go-etl
[report]:https://goreportcard.com/report/github.com/Breeze0806/go-etl
[doc-img]:https://godoc.org/github.com/Breeze0806/go-etl?status.svg
[doc]:https://godoc.org/github.com/Breeze0806/go-etl
[license-img]: https://img.shields.io/badge/License-Apache%202.0-blue.svg
[license]: https://github.com/Breeze0806/go-etl/blob/main/LICENSE
[ci-img]: https://app.travis-ci.com/Breeze0806/go-etl.svg?branch=main
[ci]: https://app.travis-ci.com/Breeze0806/go-etl
[cov-img]: https://codecov.io/gh/Breeze0806/go-etl/branch/main/graph/badge.svg?token=UGb27Nysga
[cov]: https://codecov.io/gh/Breeze0806/go-etl