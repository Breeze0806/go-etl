# OracleReader插件文档

## 快速介绍

OracleReader插件实现了从Oracle数据库读取数据。在底层实现上，OracleReader通过github.com/godror/godror以及database/sql连接远程Oracle数据库，并执行相应的sql语句将数据从Oracle库中查询出来,这里和其他数据库不同的是由于oracle未公开交互协议，oracle的golang驱动基于[ODPI-C](https://oracle.github.io/odpi/doc/installation.html)实现的,需要利用[Oracle Instant Client]( https://www.oracle.com/database/technologies/instant-client/downloads.html)进行连接,例如，连接oracle 11g需要12.x版本。

## 实现原理

OracleReader通过github.com/godror/godror使用的Oracle Instant Client连接远程oracle数据库，并根据用户配置的信息生成查询SQL语句，然后发送到远程oracle数据库，并将该SQL执行返回结果使用go-etl自定义的数据类型拼装为抽象的数据集，并传递给下游Writer处理。

OracleReader通过使用dbmsreader中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中Oracle采取了storage/database/oracle实现的Dialect。

## 功能说明

### 配置样例

配置一个从Oracle数据库同步抽取数据到本地的作业:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "oraclereader",
                    "parameter": {
                        "connection":  {
                            "url": "connectString=\"192.168.15.130:1521/xe\" heterogeneousPool=false standaloneConnection=true",
                            "table": {
                                "schema":"TEST",
                                "name":"SRC"
                            }
                        },
                        "username": "system",
                        "password": "oracle",
                        "column": ["*"],
                        "split" : {
                            "key":"id"
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

#### column

- 描述：所配置的表中需要同步的列名集合，使用JSON的数组描述字段信息。用户使用*代表默认使用所有列配置，例如["\*"]。

  支持列裁剪，即列可以挑选部分列进行导出。

  支持列换序，即列可以不按照表schema信息进行导出。

  支持常量配置，用户需要按照ORACLE SQL语法格式: ["id", "`table`", "1", "'bazhen.csy'", "null", "left(a,10)", "2.3" , "true"] id为普通列名，`table`为包含保留在的列名，1为整形数字常量，'bazhen.csy'为字符串常量，null为空指针，left(a,10)为表达式，2.3为浮点数，true为布尔值。

- 必选：是

- 默认值: 无

#### split

##### key

- 描述 主要用于配置oracle表的切分键，切分键必须为bigInt/string/time类型，假设数据按切分键分布是均匀的
- 必选：否
- 默认值: 无

##### timeAccuracy

- 描述 主要用于配置oracle表的时间切分键，主要用于描述时间最小单位，day（日）,min（分钟）,s（秒）,ms（毫秒）,us（微秒）,ns（纳秒）
- 必选：否
- 默认值: 无

##### range

###### type
- 描述 主要用于配置db2表的切分键默认值类型，值为bigInt/string/time，这里不会检查表切分键中的类型，但也请务必确保类型正确。
- 必选：否
- 默认值: 无

###### left
- 描述 主要用于配置db2表的切分键默认最大值
- 必选：否
- 默认值: 无

###### right
- 描述 主要用于配置db2表的切分键默认最小值
- 必选：否
- 默认值: 无

#### where

- 描述 主要用于配置select的where条件
- 必选：否
- 默认值: 无

### 类型转换

目前  OracleReader支持大部分  Oracle类型，但也存在部分个别类型没有支持的情况，请注意检查你的类型。
下面列出OracleReader针对  Oracle类型转换列表:

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