# DMReader插件文档

## 快速介绍

DMReader插件实现了从达梦数据库读取数据。在底层实现上，DMReader通过gitee.com/chunanyong/dm以及database/sql连接远程达梦数据库，并执行相应的sql语句将数据从达梦库中查询出来。

## 实现原理

DMReader通过gitee.com/chunanyong/dm连接远程达梦数据库，并根据用户配置的信息生成查询SQL语句，然后发送到远程达梦数据库，并将该SQL执行返回结果使用go-etl自定义的数据类型拼装为抽象的数据集，并传递给下游Writer处理。

DMReader通过使用dbmsreader中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中达梦数据库采取了storage/database/dm实现的Dialect。

## 功能说明

### 配置样例

配置一个从达梦数据库同步抽取数据到本地的作业:

```json
{
  "job":{
    "content":[
      {
        "reader":{
          "name": "dmreader",
          "parameter": {
            "username": "",
            "password": "",
            "column": [],
            "connection": {
              "url": "",
              "table": {
                "db": "",
                "name": ""
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

- 描述 主要用于配置对端连接信息。基本配置格式：ip:port，ip:port代表达梦数据库的IP地址和端口。
- 必选：是
- 默认值: 无

#### username

- 描述 主要用于配置达梦数据库的用户
- 必选：是
- 默认值: 无

#### password

- 描述 主要用于配置达梦数据库的密码
- 必选：是
- 默认值: 无

#### table

描述达梦数据库表信息

##### db

- 描述 主要用于配置达梦数据库表的数据库名
- 必选：是
- 默认值: 无

##### name

- 描述 主要用于配置达梦数据库表的表名
- 必选：是
- 默认值: 无

#### column

- 描述：所配置的表中需要同步的列名集合，使用JSON的数组描述字段信息。用户使用*代表默认使用所有列配置，例如["\*"]。

  支持列裁剪，即列可以挑选部分列进行导出。

  支持列换序，即列可以不按照表schema信息进行导出。

- 必选：是

- 默认值: 无

#### where

- 描述 主要用于配置select的where条件
- 必选：否
- 默认值: 无

### 类型转换

目前DMReader支持大部分达梦数据库类型，但也存在部分个别类型没有支持的情况，请注意检查你的类型。
| go-etl 类型 | 达梦数据库类型 |
| --- | --- |
| bool  |bit |
| bigInt | BIGINT, INT, INTEGER, SMALLINT,TINYINT,BYTE |
| decimal | DOUBLE，DOUBLE PRECISION, REAL, DECIMAL，NUMERIC,NUMBER  |
| string | VARCHAR, CHAR, TEXT, LONG,CLOB，LONGVARCHAR   |
| time | DATE, TIME, DATETIME,TIMESTAMP |
| bytes | BLOB, VARBINARY, IMAGE, LONGVARBINARY  |
## 性能报告

待测试

## 约束限制

### 数据库编码问题
目前仅支持utf8字符集