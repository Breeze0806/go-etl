# PostgresWriter插件文档

## 快速介绍

PostgresWriter插件实现了向Postgres/Greenplum数据库写入数据。在底层实现上，PostgresWriter通过github.com/lib/pq以及database/sql连接远程Postgres/Greenplum数据库，并执行相应的sql语句将数据写入Postgres/Greenplum数据库。

## 实现原理

PostgresWriter通过github.com/lib/pq连接远程Postgres/Greenplum数据库，并根据用户配置的信息和来自Reader的go-etl自定义的数据类型生成写入SQL语句，然后发送到远程Postgres/Greenplum数据库执行。

Postgres/Greenplum通过使用dbmswriter中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中Postgres/Greenplum采取了storage/database/postgres实现的Dialect。

根据你配置的 `writeMode` 生成

- `insert into...`(当主键/唯一性索引冲突时会写不进去冲突的行)

**或者**

- `copy in ...` 与 insert into 行为一致，速度比insert into方式迅速。出于性能考虑，将数据缓冲到内存 中，当 内存累计到预定阈值时，才发起写入请求。

**或者**

- `on conflict ... do update set ...` 允许你在插入数据时处理冲突情况。当插入的数据违反唯一约束（如主键或唯一索引）时，可以选择更新现有记录而不是抛出错误

## 功能说明

### 配置样例

配置一个向Postgres/Greenplum数据库同步写入数据的作业:

```json
{
    "job":{
        "content":[
            {
               "writer":{
                    "name": "postgreswriter",
                    "parameter": {
                        "username": "postgres",
                        "password": "123456",
                        "writeMode": "insert",
                        "column": ["*"],
                        "session": [],
                        "preSql": [],
                        "connection":  {
                                "url": "postgres://192.168.15.130:5432/postgres?sslmode=disable",
                                "table": {
                                    "schema":"destination",
                                    "name":"type_table"
                                }
                         },
                        "preSql": ["create table a like b"],
                        "postSql": ["drop table a"],
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

- 描述 主要用于配置对端连接信息。基本配置格式：postgres://ip:port/db，ip:port代表mysql数据库的IP地址和端口，db表示要默认连接的数据库，和[pq](https://pkg.go.dev/github.com/lib/pq)的连接配置信息基本相同，只是将用户名和密码从连接配置信息提出，方便之后对这些信息加密。与[pq](https://pkg.go.dev/github.com/lib/pq)不同的是，可以使用readTimeout/writeTimeout配置读/写超时,格式与batchTimeout相同。
- 必选：是
- 默认值: 无

#### username

- 描述 主要用于配置postgres数据库的用户
- 必选：是
- 默认值: 无

#### password

- 描述 主要用于配置postgres数据库的密码
- 必选：是
- 默认值: 无

#### name

描述postgres表信息

##### schema

- 描述 主要用于配置postgres表的模式名
- 必选：是
- 默认值: 无

##### table

- 描述 主要用于配置postgres表的表名
- 必选：是
- 默认值: 无

#### column

- 描述：所配置的表中需要同步的列名集合，使用JSON的数组描述字段信息。用户使用*代表默认使用所有列配置，例如["\*"]。

  支持列裁剪，即列可以挑选部分列进行导出。

  支持列换序，即列可以不按照表schema信息进行导出。

  支持常量配置，用户需要按照PostgreSQL语法格式: ["id", "'hello'::varchar", "true", "2.5::real", "power(2,3)"] id为普通列名，'hello'::varchar为字符串常量，true为布尔值，2.5为浮点数, power(2,3)为函数。

- 必选：是

- 默认值: 无

#### writeMode

- 描述：写入模式，`insert`代表`insert into`方式写入数据，`copyIn`代表`copy in`方式写入数据, `upsert`代表`insert into ... on conflict ... do update set ...` 方式写入数据。
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

- 描述 主要用于在写入数据前的sql语句组,不要使用select语句，否则会报错。
- 必选：否
- 默认值: 无

#### postSql

- 描述 主要用于在写入数据后的sql语句组,不要使用select语句，否则会报错。
- 必选：否
- 默认值: 无

#### upsertSql

- 描述 主要用于配置`upsert`的`on conflict ... do update set ...` 语句
- 必选：否
- 默认值: 无

### 类型转换

目前PostgresWriter支持大部分Postgres类型，但也存在部分个别类型没有支持的情况，请注意检查你的类型。

下面列出PostgresWriter针对Postgres类型转换列表:

| go-etl的类型 | Postgres数据类型                                         |
| ------------ | -------------------------------------------------------- |
| bool         | boolen                                                   |
| bigInt       | bigint, bigserial, integer, smallint, serial,smallserial |
| decimal      | double precision, decimal, numeric, real                 |
| string       | varchar, text, uuid                                     |
| time         | date, time, timestamp                                    |
| bytes        | char                                                     |

## 性能报告

待测试

## 约束限制

### 数据库编码问题
目前仅支持utf8字符集

## FAQ
1. upsert 模式支持 postgres 9.6+ 版本，支持哪些 PostgreSQL 和 Greenplum 版本？
   - PostgreSQL: upsert 功能（`INSERT ... ON CONFLICT ... DO UPDATE SET`）在 PostgreSQL 9.5+ 中引入，但 go-etl 特别要求 PostgreSQL 9.6+ 以获得稳定的 upsert 操作，因为 9.6 版本对此功能进行了改进和完善。
   - Greenplum: 由于 Greenplum 基于 PostgreSQL 构建，upsert 功能在基于 PostgreSQL 12.12 的 Greenplum 7.x 及更高版本中受支持。