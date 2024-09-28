# Sqlite3Reader插件文档

## 快速介绍

Sqlite3Reader插件实现了从Sqlite3数据库读取数据。在底层实现上，Sqlite3Reader通过github.com/mattn/go-sqlite3连接远程Sqlite3数据库，并执行相应的sql语句将数据从数据库库中查询出来。

## 实现原理

Sqlite3Reader通过github.com/mattn/go-sqlite3连接远程Sqlite3数据库，并根据用户配置的信息生成查询SQL语句，然后发送到远程Sqlite3数据库，并将该SQL执行返回结果使用go-etl自定义的数据类型拼装为抽象的数据集，并传递给下游Writer处理。

Sqlite3Reader通过使用dbmsreader中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中sqlite3采取了storage/database/sqlite3实现的Dialect。

## 功能说明

### 配置样例

配置一个从Sqlite3数据库同步抽取数据到本地的作业:

```json
{
  "job": {
    "content": [
      {
        "reader": {
          "name": "sqlite3reader",
          "parameter": {
            "column": [
              "*"
            ],
            "connection": {
              "url": "E:\\Sqlite3\\test.db",
              "table": {
                "db": "main",
                "name": "type_table"
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

- 描述 主要用于配置sqlite3数据库文件路径。
- 必选：是
- 默认值: 无

#### table

描述sqlite3表信息

##### name

- 描述 主要用于配置sqlite3表的表名
- 必选：是
- 默认值: 无

#### column

- 描述：所配置的表中需要同步的列名集合，使用JSON的数组描述字段信息。用户使用*代表默认使用所有列配置，例如["\*"]。

  支持列裁剪，即列可以挑选部分列进行导出。

  支持列换序，即列可以不按照表schema信息进行导出。

  支持常量配置，用户需要按照sqlite3语法格式。

- 必选：是

- 默认值: 无

#### split

##### key

- 描述 主要用于配置sqlite3表的切分键，切分键必须为bigInt/string/time类型，假设数据按切分键分布是均匀的
- 必选：否
- 默认值: 无

##### timeAccuracy

- 描述 主要用于配置sqlite3ql表的时间切分键，主要用于描述时间最小单位，day（日）,min（分钟）,s（秒）,ms（毫秒）,us（微秒）,ns（纳秒）
- 必选：否
- 默认值: 无

##### range

###### type
- 描述 主要用于配置sqlite3表的切分键默认值类型，值为bigInt/string/time，这里会检查表切分键中的类型，请务必确保类型正确。
- 必选：否
- 默认值: 无

###### left
- 描述 主要用于配置sqlite3表的切分键默认最大值
- 必选：否
- 默认值: 无

###### right
- 描述 主要用于配置sqlite3表的切分键默认最小值
- 必选：否
- 默认值: 无

#### where

- 描述 主要用于配置select的where条件
- 必选：否
- 默认值: 无

#### querySql

- 描述：在有些业务场景下，where这一配置项不足以描述所筛选的条件，用户可以通过该配置型来自定义筛选SQL。当用户配置了这一项之后，DataX系统就会忽略table，column这些配置型，直接使用这个配置项的内容对数据进行筛选，例如需要进行多表join后同步数据，使用select a,b from table_a join table_b on table_a.id = table_b.id
当用户配置querySql时，Sqlite3Reader直接忽略table、column、where条件的配置，querySql优先级大于table、column、where选项。
- 必选：否
- 默认值：无

#### trimChar

- 描述：对于sqlite3的char类型是否去掉其前后的空格
- 必选：否
- 默认值：false

### 类型转换

目前Sqlite3Reader支持大部分sqlite3类型，但也存在部分个别类型没有支持的情况，请注意检查你的类型。

下面列出Sqlite3Reader针对sqlite3类型转换列表:

| go-etl的类型 | sqlite3数据类型        |
| ------------ |--------------------|
| string       | INTEGER、TEXT、REAL、BLOB |

## 性能报告

待测试

## 约束限制

### 数据库编码问题
目前仅支持utf8字符集

## FAQ
