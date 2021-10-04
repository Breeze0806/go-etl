# PostgresReader插件文档

## 快速介绍

PostgresReader插件实现了从Postgres/Greenplum数据库读取数据。在底层实现上，PostgresReader通过github.com/lib/pq连接远程Mysql数据库，并执行相应的sql语句将数据从mysql库中查询出来。

## 实现原理

PostgresReader通过github.com/lib/pq连接远程Postgres/Greenplum数据库，并根据用户配置的信息生成查询SQL语句，然后发送到远程postgres/greenplum数据库，并将该SQL执行返回结果使用go-etl自定义的数据类型拼装为抽象的数据集，并传递给下游Writer处理。和直接使用github.com/lib/pq连接数据库不同的是，这里采用了github.com/Breeze0806/go/database/pqto以便能设置读写超时。

PostgresReader通过使用rdbmreader中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中Mysql采取了storage/database/postgres实现的Dialect。

## 功能说明

### 配置样例

配置一个从Postgres/Greenplum数据库同步抽取数据到本地的作业:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "postgresreader",
                    "parameter": {
                        "username": "postgres",
                        "password": "123456",
                        "column": ["*"],
                        "connection":  {
                                "url": "postgres://192.168.15.130:5432/postgres?sslmode=disable",
                                "table": {
                                    "db":"postgres",
                                    "schema":"source",
                                    "name":"type_table"
                                }
                            },
                        "where": ""
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

#### table

描述postgres表信息

##### db

- 描述 主要用于配置postgres表的实例名
- 必选：是
- 默认值: 无

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

#### where

- 描述 主要用于配置select的where条件
- 必选：否
- 默认值: 无

### 类型转换

目前PostgresReader支持大部分Postgres类型，但也存在部分个别类型没有支持的情况，请注意检查你的类型。

下面列出PostgresReader针对Postgres类型转换列表:

| go-etl的类型 | Postgres数据类型                                         |
| ------------ | -------------------------------------------------------- |
| bool         | boolen                                                   |
| bigInt       | bigint, bigserial, integer, smallint, serial,smallserial |
| decimal      | double precision, decimal, numeric, real                 |
| string       | varchar, text                                            |
| time         | date, time, timestamp                                    |
| bytes        | char                                                     |

## 性能报告

待测试