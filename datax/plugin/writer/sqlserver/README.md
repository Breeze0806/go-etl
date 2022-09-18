# SQLServerWriter插件文档

## 快速介绍

SQLServerWriter插件实现了向sql server数据库写入数据。在底层实现上，SQLServerWriter通过github.com/denisenkom/go-mssqldb以及database/sql连接远程sql server数据库，并执行相应的sql语句将数据写入sql server数据库。

## 实现原理

SQLServerWriter通过github.com/denisenkom/go-mssqldb连接远程sql server数据库，并根据用户配置的信息和来自Reader的go-etl自定义的数据类型生成写入SQL语句，然后发送到远程sql server数据库执行。

SQLServerWriter通过使用rdbmwriter中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中sqlserver采取了storage/database/sqlserver实现的Dialect。

根据你配置的 `writeMode` 生成

- `insert into...`(当主键/唯一性索引冲突时会写不进去冲突的行)

**或者**

- bulk copy 即`inster bulk ...` 与 insert into 行为一致，速度比insert into方式迅速，但是目前不知为何无法插入含有空值的记录

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
                        "preSql": [],
                        "connection":  {
                                "url": "sqlserver://192.168.15.130:1433?database=test&encrypt=disable",
                                "table": {
                                    "db":"test",
                                    "schema":"dest",
                                    "name":"mytable"
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

- 描述 主要用于配置对端连接信息。基本配置格式：sqlserver://ip:port?database=db&encrypt=disable"，ip:port代表mysql数据库的IP地址和端口，db表示要默认连接的数据库，详细见[go-mssqldb](https://github.com/denisenkom/go-mssqldb)的连接配置信息.
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

- 描述：写入模式，insert代表insert into方式写入数据,copyIn代表批量复制插入。
- 必选：否
- 默认值: insert

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

