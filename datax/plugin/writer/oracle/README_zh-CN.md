# OracleWriter插件文档

## 快速介绍

OracleReader插件实现了向Oracle数据库写入数据。在底层实现上，OracleReader通过github.com/godror/godror以及database/sql连接远程Oracle数据库，Oracle,这里和其他数据库不同的是由于oracle未公开交互协议，oracle的golang驱动基于[ODPI-C](https://oracle.github.io/odpi/doc/installation.html)实现的,需要利用[Oracle Instant Client]( https://www.oracle.com/database/technologies/instant-client/downloads.html)进行连接,例如，连接oracle 11g需要12.x版本。

## 实现原理

OracleReader通过github.com/godror/godror使用的Oracle Instant Client连接远程oracle数据库，并根据用户配置的信息和来自Reader的go-etl自定义的数据类型生成写入SQL语句，然后发送到远程Oracle数据库执行。

OracleReader通过使用dbmswriter中中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中Oracle采取了storage/database/oracle实现的Dialect。

根据你配置的 `writeMode` 生成

- `insert into...`(当主键/唯一性索引冲突时会写不进去冲突的行)

注意这里的insert写入方式已经不是通常的storage/database的insert实现方式，而是oracle特有的方式,例如在这里的实现中，query为`insert into a(x,y,x) values(:1,:2,:3)`，而x,y,z的args输入是该列数值组成的三个数组。
## 功能说明

### 配置样例

配置一个从内存写入Oracle数据库数据的作业:

```json
{
    "job":{
        "content":[
            {
                 "writer":{
                    "name": "oraclewriter",
                    "parameter": {
                        "connection":  {
                            "url": "connectString=\"192.168.15.130:1521/xe\" heterogeneousPool=false standaloneConnection=true",
                            "table": {
                                "schema":"TEST",
                                "name":"DEST"
                            }
                        },
                        "username": "system",
                        "password": "oracle",
                        "writeMode": "insert",
                        "column": ["*"],
                        "preSql": ["create table a like b"],
                        "postSql": ["drop table a"],
                        "batchTimeout": "1s",
                        "batchSize":1000
                    }
                },
            }
        ]
    }
}
```

### 参数说明

#### url

- 描述 主要用于配置对端连接信息。oracle连接数据库的基本配置格式：`connectString="192.168.15.130:1521/xe" heterogeneousPool=false standaloneConnection=true`，connectString代表连接oracle数据库的信息，如果使用servername连接请使用`ip:port/servername`或者`(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=ip)(PORT=port))(CONNECT_DATA=(SERVICE_NAME=servername)))`，如果使用sid连接，那么请使用`(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=ip)(PORT=port))(CONNECT_DATA=(SID=sid)))`,和[Godror User Guide](https://godror.github.io/godror/doc/contents.html)的连接配置信息基本相同，只是将用户名和密码从连接配置信息提出，方便之后对这些信息加密。
- 必选：是
- 默认值: 无

#### username

- 描述 主要用于配置oracle数据库的用户
- 必选：是
- 默认值: 无

#### password

- 描述 主要用于配置oracle数据库的密码
- 必选：是
- 默认值: 无

#### table

描述oracle表信息

##### schema

- 描述 主要用于配置oracle表的模式名
- 必选：是
- 默认值: 无

##### name

- 描述 主要用于配置oracle表的表名
- 必选：是
- 默认值: 无

#### writeMode

- 描述：写入模式，insert代表insert into方式写入数据。
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

- 描述 主要用于在写入数据前的sql语句组,不要使用select语句，否则会报错。
- 必选：否
- 默认值: 无

#### postSql

- 描述 主要用于在写入数据后的sql语句组,不要使用select语句，否则会报错。
- 必选：否
- 默认值: 无

### 类型转换

目前  OracleWriter支持大部分  Oracle类型，但也存在部分个别类型没有支持的情况，请注意检查你的类型。
下面列出OracleWriter针对  Oracle类型转换列表:

| go-etl的类型 | Oracle数据类型               |
| ------------ | ------------------------- |
| bool         | BOOLEAN                   |
| bigInt       | NUMBER,INTEGER,SMALLINT |
| decimal      | BINARY_FLOAT, FLOAT, BINARY_DOUBLE,REAL, DECIMAL,NUMBRIC     |
| string       | VARCHAR,CHAR,NCHAR,VARCHAR2,NVARCHAR2,CLOB,NCLOB              |
| time         | DATE,TIMESTAMP       |
| bytes        | BLOB,RAW,LONG RAW,LONG                      |

## 性能报告

待测试

## 约束限制

### 数据库编码问题
目前仅支持utf8字符集

## FAQ

1.如何配置oracle的Oracle Instant Client

例子如下：

- 注意在linux下如export LD_LIBRARY_PATH=/opt/oracle/instantclient_21_1:$LD_LIBRARY_PATH，另需要安装libaio

- 注意在windows下如set path=%path%;%GOPATH%\oracle\instantclient_21_1，
Oracle Instant Client 19不再支持windows7

1.如何配置oracle的Oracle Instant Client

例子如下：

- 注意在linux下如export LD_LIBRARY_PATH=/opt/oracle/instantclient_21_1:$LD_LIBRARY_PATH，另需要安装libaio

- 注意在windows下如set path=%path%;%GOPATH%\oracle\instantclient_21_1，
Oracle Instant Client 19不再支持windows7，另外，需要安装[Oracle Instant Client以及对应的Visual Studio redistributable](https://odpi-c.readthedocs.io/en/latest/user_guide/installation.html#windows)


2.如何消除`godor WARNING: discrepancy between SESSIONTIMEZONE and SYSTIMESTAMP`

您可以与您的数据库管理员（DBA）沟通，以将数据库的时区（DBTIMEZONE）与底层操作系统的时区同步，或者使用以下 SQL 语句：

```sql
ALTER SESSION SET TIME_ZONE='Europe/Berlin'
```

或者在[./connection.md]（连接字符串）中设置一个选定的时区：

```ini
timezone="Europe/Berlin"
```

（它使用 time.LoadLocation 进行解析，因此可以使用这样的名称，或者使用 local，或者使用数字形式如 +0500 表示固定时区）。

警告：使用 ALTER SESSION 更改的时区可能不会被每次读取，因此请始终将 ALTER SESSION 一致地设置为相同的时区，或者使用以下连接参数：

```ini
perSessionTimezone=1
```

以强制每次会话都检查时区（而不是每个数据库只缓存一次）。

3.写入时间格式为何这么慢？

目前OracleWriter写入时间格式是通过函数to_date或者to_timestamp转换的，而不是通过绑定变量的方式，因此写入时间格式会比较慢。
