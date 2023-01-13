# DB2Writer插件文档

## 快速介绍

DB2Writer插件实现了向DB2 LUW 数据库写入数据。在底层实现上，DB2Writer通过github.com/ibmdb/go_ibm_db以及database/sql连接远程DB2 LUW 数据库，并执行相应的sql语句将数据写入db2库。这里和其他数据库不同的是由于db2未公开交互协议，db2的golang驱动利用db2的odbc库来连接数据库。

## 实现原理
DB2Writer通过github.com/ibmdb/go_ibm_db利用db2的odbc库连接远程DB2 LUW数据库，并根据用户配置的信息和来自Reader的go-etl自定义的数据类型生成写入SQL语句，然后发送到远程DB2数据库执行。

DB2Writer通过使用rdbmwriter中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中DB2采取了storage/database/db2实现的Dialect。

根据你配置的 `writeMode` 生成

- `insert into...`(当主键/唯一性索引冲突时会写不进去冲突的行)


## 功能说明

### 配置样例

配置一个从内存写入DB2数据库数据的作业:

```json
{
    "job":{
        "content":[
            {
               "writer":{
                    "name": "db2writer",
                    "parameter": {
                        "connection":  {
                            "url": "HOSTNAME=127.0.0.1;PORT=50000;DATABASE=db",
                            "table": {
                                "schema":"SOURCE",
                                "name":"TEST"
                            }
                        },
                        "username": "root",
                        "password": "12345678",
                        "writeMode": "insert",
                        "column": ["*"],
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

- 描述 主要用于配置对端连接信息。基本配置格式：HOSTNAME=ip;PORT=port;DATABASE=db，ip代表db2数据库的IP地址和port端口，db表示要默认连接的数据库，和[ibm db2](https://github.com/ibmdb/go_ibm_db)的连接配置信息基本相同，只是将用户名和密码从连接配置信息提出，方便之后对这些信息加密。
- 必选：是
- 默认值: 无

#### username

- 描述 主要用于配置db2数据库的用户
- 必选：是
- 默认值: 无

#### password

- 描述 主要用于配置db2数据库的密码
- 必选：是
- 默认值: 无

#### table

描述db2表信息

##### schema

- 描述 主要用于配置db2表的模式名
- 必选：是
- 默认值: 无

##### name

- 描述 主要用于配置db2表的表名
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

目前  DB2Reader支持大部分  DB2类型，但也存在部分个别类型没有支持的情况，请注意检查你的类型。
下面列出DB2Reader针对  DB2类型转换列表:

| go-etl的类型 | DB2数据类型               |
| ------------ | ------------------------- |
| bool         | BOOLEAN                   |
| bigInt       | BIGINT, INTEGER, SMALLINT |
| decimal      | DOUBLE, REAL, DECIMAL     |
| string       | VARCHAR,CHAR              |
| time         | DATE,TIME,TIMESTAMP       |
| bytes        | BLOB                      |

## 性能报告

待测试

## 约束限制


### 数据库编码问题
目前仅支持utf8字符集

## FAQ

1.如何配置db2的odbc库

- 注意在linux下如Makefile所示export LD_LIBRARY_PATH=${DB2HOME}/lib
- 注意在windows下如release.bat所示set path=%path%;%GOPATH%\src\github.com\ibmdb\go_ibm_db\clidriver\bin