# MysqlReader插件文档

## 快速介绍

MysqlReader插件实现了从mysql数据库读取数据。在底层实现上，MysqlReader通过github.com/go-sql-driver/mysql以及database/sql连接远程Mysql数据库，并执行相应的sql语句将数据从mysql库中查询出来。

## 实现原理

MysqlReader通过github.com/go-sql-driver/mysql连接远程Mysql数据库，并根据用户配置的信息生成查询SQL语句，然后发送到远程Mysql数据库，并将该SQL执行返回结果使用go-etl自定义的数据类型拼装为抽象的数据集，并传递给下游Writer处理。

MysqlReader通过使用dbmsreader中定义的查询流程调用go-etl自定义的storage/database的DBWrapper来实现具体的查询。DBWrapper封装了database/sql的众多接口，并且抽象出了数据库方言Dialect。其中Mysql采取了storage/database/mysql实现的Dialect。

## 功能说明

### 配置样例

配置一个从Mysql数据库同步抽取数据到本地的作业:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "mysqlreader",
                    "parameter": {
                        "username": "root",
                        "password": "123456",
                        "column": ["*"],
                        "connection":  {
                                "url": "tcp(192.168.0.1:3306)/mysql?parseTime=false",
                                "table": {
                                    "db":"source",
                                    "name":"type_table"
                                }
                            },
                        "split" : {
                            "key":"id"
                        },                       
                        "where": "",
                        "querySql":["select a,b from table_a join table_b on table_a.id = table_b.id"]
                    }
                }
            }
        ]
    }
}
```

### 参数说明

#### url

- 描述 主要用于配置对端连接信息。基本配置格式：tcp(ip:port)/db，ip:port代表mysql数据库的IP地址和端口，db表示要默认连接的数据库，和[mysql](https://github.com/go-sql-driver/mysql)的连接配置信息基本相同，只是将用户名和密码从连接配置信息提出，方便之后对这些信息加密。
- 必选：是
- 默认值: 无

#### username

- 描述 主要用于配置mysql数据库的用户
- 必选：是
- 默认值: 无

#### password

- 描述 主要用于配置mysql数据库的密码
- 必选：是
- 默认值: 无

#### table

描述mysql表信息

##### db

- 描述 主要用于配置mysql表的数据库名
- 必选：是
- 默认值: 无

##### name

- 描述 主要用于配置mysql表的表名
- 必选：是
- 默认值: 无

#### column

- 描述：所配置的表中需要同步的列名集合，使用JSON的数组描述字段信息。用户使用*代表默认使用所有列配置，例如["\*"]。

  支持列裁剪，即列可以挑选部分列进行导出。

  支持列换序，即列可以不按照表schema信息进行导出。

  支持常量配置，用户需要按照Mysql SQL语法格式: ["id", "`table`", "1", "'bazhen.csy'", "null", "to_char(a + 1)", "2.3" , "true"] id为普通列名，`table`为包含保留在的列名，1为整形数字常量，'bazhen.csy'为字符串常量，null为空指针，to_char(a + 1)为表达式，2.3为浮点数，true为布尔值。

- 必选：是

- 默认值: 无

#### split

##### key

- 描述 主要用于配置mysql表的切分键，切分键必须为bigInt/string/time类型，假设数据按切分键分布是均匀的
- 必选：否
- 默认值: 无

##### timeAccuracy

- 描述 主要用于配置mysql表的时间切分键，主要用于描述时间最小单位，day（日）,min（分钟）,s（秒）,ms（毫秒）,us（微秒）,ns（纳秒）
- 必选：否
- 默认值: 无

##### range

###### type
- 描述 主要用于配置mysql表的切分键默认值类型，值为bigInt/string/time，这里会检查表切分键中的类型，请务必确保类型正确。
- 必选：否
- 默认值: 无

###### left
- 描述 主要用于配置mysql表的切分键默认最大值
- 必选：否
- 默认值: 无

###### right
- 描述 主要用于配置mysql表的切分键默认最小值
- 必选：否
- 默认值: 无

#### where

- 描述 主要用于配置select的where条件
- 必选：否
- 默认值: 无

#### querySql

- 描述：在有些业务场景下，where这一配置项不足以描述所筛选的条件，用户可以通过该配置型来自定义筛选SQL。当用户配置了这一项之后，DataX系统就会忽略table，column这些配置型，直接使用这个配置项的内容对数据进行筛选，例如需要进行多表join后同步数据，使用select a,b from table_a join table_b on table_a.id = table_b.id
当用户配置querySql时，MysqlReader直接忽略table、column、where条件的配置，querySql优先级大于table、column、where选项。
- 必选：否
- 默认值：无

#### trimChar

- 描述：对于mysql的char类型是否去掉其前后的空格
- 必选：否
- 默认值：false

### 类型转换

目前MysqlReader支持大部分Mysql类型，但也存在部分个别类型没有支持的情况，请注意检查你的类型。

下面列出MysqlReader针对Mysql类型转换列表:

| go-etl的类型 | mysql数据类型                                       |
| ------------ | --------------------------------------------------- |
| bigInt       | int, tinyint, smallint, mediumint, bigint,year,unsigned int, unsigned bigint, unsigned smallint, unsigned tinyint     |
| decimal      | float, double, decimal                              |
| string       | varchar, char, tinytext, text, mediumtext, longtext |
| time         | date, datetime, timestamp, time                     |
| bytes        | tinyblob, mediumblob, blob, longblob, varbinary,bit |

## 性能报告

待测试

## 约束限制

### 数据库编码问题
目前仅支持utf8字符集

## FAQ
