# DMWriter插件文档

## 快速介绍

DMWriter插件实现了向达梦数据库写入数据。在底层实现上，DMWriter通过gitee.com/chunanyong/dm以及database/sql连接远程达梦数据库，并执行相应的sql语句将数据写入达梦库。

## 实现原理

DMWriter通过gitee.com/chunanyong/dm连接远程达梦数据库，并根据用户配置的信息和来自Reader的go-etl自定义的数据类型生成写入SQL语句，然后发送到远程达梦数据库执行。

DMWriter通过使用dbmswriter中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中达梦数据库采取了storage/database/dm实现的Dialect。

根据你配置的 `writeMode` 生成

- `insert into...`(当主键/唯一性索引冲突时会写不进去冲突的行)

## 功能说明

### 配置样例

配置一个从内存写入达梦数据库数据的作业:

```json
{
  "job": {
    "content": [
      {
        "writer":{
          "name": "dmwriter",
          "parameter": {
            "username": "",
            "password": "",
            "column": ["*"],
            "preSql": [],
            "writeMode": "insert",
            "connection":  {
              "url": "",
              "table": {
                "db":"",
                "name":""
              }
            },
            "batchTimeout": "1s",
            "batchSize":1000
          }
        },
        "transformer": []
      }
    ]
  }
}
```

### 参数说明

#### url

- 描述 主要用于配置对端连接信息。基本配置格式：ip:port，ip:port代表达梦数据库的IP地址和端口号。
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

目前DMWriter支持大部分达梦数据库类型，但也存在部分个别类型没有支持的情况，请注意检查你的类型。
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

## FAQ