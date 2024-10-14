# SQLServerWriter插件文档

## 快速介绍

SQLServerWriter插件实现了向sql server数据库写入数据。在底层实现上，SQLServerWriter通过github.com/microsoft/go-mssqldb以及database/sql连接远程sql server数据库，并执行相应的sql语句将数据写入sql server数据库。

## 实现原理

SQLServerWriter通过github.com/microsoft/go-mssqldb连接远程sql server数据库，并根据用户配置的信息和来自Reader的go-etl自定义的数据类型生成写入SQL语句，然后发送到远程sql server数据库执行。

SQLServerWriter通过使用dbmswriter中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中sqlserver采取了storage/database/sqlserver实现的Dialect。

根据你配置的 `writeMode` 生成

- `insert into...`(当主键/唯一性索引冲突时会写不进去冲突的行)

**或者**

- bulk copy 即`BULK INSERT ...` 与 insert into 行为一致，速度比insert into方式迅速。比起`insert into...`,我们更推荐这种写入模式。


## 功能说明

### 配置样例

配置一个向sql server数据库同步写入数据的作业:

```json
{
    "job":{
        "content":[
            {
               "writer":{
                    "name": "sqlserverwriter",
                    "parameter": {
                        "username": "sa",
                        "password": "Breeze_0806",
                        "writeMode": "insert",
                        "column": ["*"],
                        "preSql": ["create table a like b"],
                        "postSql": ["drop table a"],
                        "connection":  {
                                "url": "sqlserver://192.168.15.130:1433?database=test&encrypt=disable",
                                "table": {
                                    "db":"test",
                                    "schema":"dest",
                                    "name":"mytable"
                                }
                         },
                         "bulkOption":{
                            "KeepNulls":true
                         },
                        "batchTimeout": "1s",
                        "batchSize":1000
                    }
               }
            }
        ]
    }
}
```

### 参数说明

#### url

- 描述 主要用于配置对端连接信息。基本配置格式：sqlserver://ip:port?database=db&encrypt=disable"，ip:port代表mysql数据库的IP地址和端口，db表示要默认连接的数据库，详细见[go-mssqldb](https://github.com/microsoft/go-mssqldb)的连接配置信息.
- 必选：是
- 默认值: 无

#### username

- 描述 主要用于配置sql server数据库的用户
- 必选：是
- 默认值: 无

#### password

- 描述 主要用于配置sql server数据库的密码
- 必选：是
- 默认值: 无

#### table

描述sql server表信息

##### db

- 描述 主要用于配置sql server表的数据库名
- 必选：是
- 默认值: 无

##### schema

- 描述 主要用于配置sql server表的模式名
- 必选：是
- 默认值: 无

##### name

- 描述 主要用于配置sql server表的表名
- 必选：是
- 默认值: 无

#### column

- 描述：所配置的表中需要同步的列名集合，使用JSON的数组描述字段信息。用户使用*代表默认使用所有列配置，例如["\*"]。

  支持列裁剪，即列可以挑选部分列进行导出。

  支持列换序，即列可以不按照表schema信息进行导出。

  支持常量配置，用户需要按照sql server语法格式: ["id",  "true", "power(2,3)"] id为普通列名，'hello'::varchar为字符串常量，true为布尔值，2.5为浮点数, power(2,3)为函数。

- 必选：是

- 默认值: 无

#### writeMode

- 描述：写入模式，insert代表insert into方式写入数据， copyIn代表批量复制插入。
- 必选：否
- 默认值: insert

#### bulkOption

- 描述：主要用于copyIn的批量写入配置，作用于`BULK INSERT`的配置
- 必选：否
- 默认值: 无

##### CheckConstraints

+ 描述：表示`CHECK_CONSTRAINTS`。指定在批量导入操作期间，必须检查目标表或视图上的所有约束。如果不使用`CHECK_CONSTRAINTS`选项，则会忽略任何CHECK和FOREIGN KEY约束，并且在操作后，表上的约束将被标记为不受信任。
+ 必选：否
+ 默认值：无

##### FireTriggers

+ 描述：表示`FIRE_TRIGGERS`。指定在批量导入操作期间，目标表上定义的任何插入触发器都会执行。如果为目标表的INSERT操作定义了触发器，则它们会在每个完成的批次上触发。如果未指定`FIRE_TRIGGERS`，则不会执行任何插入触发器。
+ 必选：否
+ 默认值：无

##### KeepNulls

+ 描述：表示`KEEPNULLS`。指定在批量导入操作期间，空列应保留空值，而不是为插入的列插入任何默认值。
+ 必选：否
+ 默认值：无

##### KilobytesPerBatch

+ 描述：表示`KILOBYTES_PER_BATCH`。指定每批数据的大致千字节（KB）数为*kilobytes_per_batch*。默认情况下，`KILOBYTES_PER_BATCH`是未知的。
+ 必选：否
+ 默认值：无

##### RowsPerBatch

+ 描述：表示`ROWS_PER_BATCH`。指示数据文件中数据的大致行数。默认情况下，数据文件中的所有数据都作为单个事务发送到服务器，并且查询优化器不知道批次中的行数。如果指定了`ROWS_PER_BATCH`（值大于0），则服务器将使用此值来优化批量导入操作。为`ROWS_PER_BATCH`指定的值应大致与实际行数相同。
+ 必选：否
+ 默认值：无

##### Order

+ 描述：表示`ORDER`。指定数据文件中数据的排序方式。如果导入的数据根据表上的聚集索引（如果有）进行排序，则可以提高批量导入性能。如果数据文件以不同的顺序排序，即不是聚集索引键的顺序，或者表上没有聚集索引，则忽略`ORDER`子句。提供的列名称必须是目标表中的有效列名称。默认情况下，批量插入操作假定数据文件是无序的。为了优化批量导入，SQL Server还会验证导入的数据是否已排序。
+ 必选：否
+ 默认值：无

##### Tablock

+ 描述：表示`TABLOCK`。指定在批量导入操作期间获取表级锁。如果表没有索引并且指定了TABLOCK，则多个客户端可以同时加载表。默认情况下，锁定行为由表的**在批量加载时表锁定**选项确定。在批量导入操作期间持有锁可以减少表上的锁争用，在某些情况下可以显著提高性能。
+ 必选：否
+ 默认值：无

#### batchTimeout

- 描述 主要用于配置每次批量写入超时时间间隔，格式：数字+单位， 单位：s代表秒，ms代表毫秒，us代表微妙。如果超过该时间间隔就直接写入，和batchSize一起调节写入性能。
- 必选：否
- 默认值: 1s

#### batchSize

- 描述 主要用于配置每次批量写入大小，如果超过该大小就直接写入，和batchTimeout一起调节写入性能。
- 必选：否
- 默认值: 1000

#### preSql

- 描述 主要用于在写入数据前的sql语句组,不要使用select语句，否则会报错。
- 必选：否
- 默认值: 无

#### postSql

- 描述 主要用于在写入数据后的sql语句组,不要使用select语句，否则会报错。
- 必选：否
- 默认值: 无

### 类型转换

目前SQLServerReader支持大部分SQLServer类型，但也存在部分个别类型没有支持的情况，请注意检查你的类型。

下面列出SQLServerReader针对sql server类型转换列表:

| go-etl的类型 | sql server数据类型                                          |
| ------------ | ----------------------------------------------------------- |
| bool         | bit                                                         |
| bigInt       | bigint, int, smallint, tinyint                              |
| decimal      | numeric, real,float                                         |
| string       | char, varchar, text, nchar, nvarchar, ntext                 |
| time         | date, time, datetimeoffset,datetime2,smalldatetime,datetime |
| bytes        | binary，varbinary，varbinary(max)                           |

## 性能报告

待测试

## 约束限制

## FAQ

