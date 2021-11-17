# MysqlWriter插件文档

## 快速介绍

MysqlWriter插件实现了向Postgres/Greenplum数据库写入数据。在底层实现上，MysqlWriter通过github.com/go-sql-driver/mysql以及database/sql连接远程Mysql数据库，并执行相应的sql语句将数据写入mysql库。

## 实现原理

MysqlWriter通过github.com/go-sql-driver/mysql连接远程Mysql数据库，并根据用户配置的信息和来自Reader的go-etl自定义的数据类型生成写入SQL语句，然后发送到远程Mysql数据库执行。

MysqlWriter通过使用rdbmwriter中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中Mysql采取了storage/database/mysql实现的Dialect。

根据你配置的 `writeMode` 生成

- `insert into...`(当主键/唯一性索引冲突时会写不进去冲突的行)

**或者**

- `replace into...`(没有遇到主键/唯一性索引冲突时，与 insert into 行为一致，冲突时会用新行替换原有行所有字段) 的语句写入数据到 Mysql。出于性能考虑，将数据缓冲到内存 中，当 内存累计到预定阈值时，才发起写入请求。

## 功能说明

### 配置样例

配置一个从内存写入Mysql数据库数据的作业:

```json
{
    "job":{
        "content":[
            {
               "writer":{
                    "name": "mysqlwriter",
                    "parameter": {
                        "username": "root",
                        "password": "123456",
                        "writeMode": "insert",
                        "column": ["*"],
                        "preSql": [],
                        "connection":  {
                                "url": "tcp(192.168.0.1:3306)/mysql?parseTime=false",
                                "table": {
                                    "db":"destination",
                                    "name":"type_table"
                                }
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

- 描述 主要用于配置对端连接信息。基本配置格式：tcp(ip:port)/db，ip:port代表mysql数据库的IP地址和端口，db表示要默认连接的数据库，和[mysql](https://github.com/go-sql-driver/mysql)的连接配置信息基本相同，只是将用户名和密码从连接配置信息提出，方便之后对这些信息加密。
- 必选：是
- 默认值: 无

#### username

- 描述 主要用于配置mysql数据库的用户
- 必选：是
- 默认值: 无

#### password

- 描述 主要用于配置mysql数据库的密码
- 必选：是
- 默认值: 无

#### table

描述mysql表信息

##### db

- 描述 主要用于配置mysql表的数据库名
- 必选：是
- 默认值: 无

##### table

- 描述 主要用于配置mysql表的表名
- 必选：是
- 默认值: 无

#### writeMode

- 描述：写入模式，insert代表insert into方式写入数据，replace代表replace into方式写入数据。
- 必选：否
- 默认值: insert

#### column

- 描述：所配置的表中需要同步的列名集合，使用JSON的数组描述字段信息。用户使用*代表默认使用所有列配置，例如["\*"]。

  支持列裁剪，即列可以挑选部分列进行插入。

  支持列换序，即列可以不按照表schema信息进行插入。

- 必选：是

- 默认值: 无

#### batchTimeout

- 描述 主要用于配置每次批量写入超时时间间隔，格式：数字+单位， 单位：s代表秒，ms代表毫秒，us代表微妙。如果超过该时间间隔就直接写入，和batchSize一起调节写入性能。
- 必选：否
- 默认值: 1s

#### batchSize

- 描述 主要用于配置每次批量写入大小，如果超过该大小就直接写入，和batchTimeout一起调节写入性能。
- 必选：否
- 默认值: 1000

#### preSql

- 描述 主要用于在写入数据前的sql语句组，目前还没支持
- 必选：否
- 默认值: 无

### 类型转换

目前MysqlWriter支持大部分Mysql类型，但也存在部分个别类型没有支持的情况，请注意检查你的类型。

下面列出MysqlWriter针对Mysql类型转换列表:

| go-etl的类型 | mysql数据类型                                       |
| ------------ | --------------------------------------------------- |
| bigInt       | int, tinyint, smallint, mediumint, int, bigint,year |
| decimal      | float, double, decimal                              |
| string       | varchar, char, tinytext, text, mediumtext, longtext |
| time         | date, datetime, timestamp, time                     |
| bytes        | tinyblob, mediumblob, blob, longblob, varbinary,bit |

## 性能报告

待测试