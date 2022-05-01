# DB2Reader插件文档

## 快速介绍

DB2Reader插件实现了从DB2 LUW数据库读取数据。在底层实现上，DB2Reader通过github.com/ibmdb/go_ibm_db以及database/sql连接远程DB2 LUW数据库，并执行相应的sql语句将数据从DB2库中查询出来,这里和其他数据库不同的是由于db2未公开交互协议，db2的golang驱动利用db2的odbc库来连接数据库。

## 实现原理

DB2Reader通过github.com/ibmdb/go_ibm_db使用db2的odbc库连接远程DB2 LUW数据库，并根据用户配置的信息生成查询SQL语句，然后发送到远程DB2 LUW数据库，并将该SQL执行返回结果使用go-etl自定义的数据类型拼装为抽象的数据集，并传递给下游Writer处理。

DB2Reader通过使用rdbmreader中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中DB2采取了storage/database/db2实现的Dialect。

## 功能说明

### 配置样例

配置一个从DB2数据库同步抽取数据到本地的作业:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "db2reader",
                    "parameter": {
                        "connection":  {
                            "url": "HOSTNAME=127.0.0.1;PORT=50000;DATABASE=db",
                            "table": {
                                "schema":"SOURCE",
                                "name":"TEST"
                            }
                        },
                        "username": "user",
                        "password": "passwd",
                        "column": ["*"],
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

#### column

- 描述：所配置的表中需要同步的列名集合，使用JSON的数组描述字段信息。用户使用*代表默认使用所有列配置，例如["\*"]。

  支持列裁剪，即列可以挑选部分列进行导出。

  支持列换序，即列可以不按照表schema信息进行导出。

  支持常量配置，用户需要按照DB2 SQL语法格式: ["id", "`table`", "1", "'bazhen.csy'", "null", "left(a,10)", "2.3" , "true"] id为普通列名，`table`为包含保留在的列名，1为整形数字常量，'bazhen.csy'为字符串常量，null为空指针，left(a,10)为表达式，2.3为浮点数，true为布尔值。

- 必选：是

- 默认值: 无

#### where

- 描述 主要用于配置select的where条件
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