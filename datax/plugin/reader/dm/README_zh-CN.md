# DMReader 插件文档

## 快速介绍

DMReader插件实现了从达梦(DM)数据库读取数据。在底层实现上，DMReader通过[gitee.com/chunanyong/dm](gitee.com/chunanyong/dm)的官方Go驱动连接远程DM数据库，并执行相应的SQL语句将数据从DM库中SELECT出来。

## 实现原理

DMReader通过官方的DM Go驱动连接到远程DM数据库，根据用户配置的信息生成查询语句并发送到远程DM数据库，然后将数据库返回的结果集使用go-etl自定义的数据类型拼装成抽象的数据集，并传递给下游Writer处理。

插件通过切分键进行数据读取任务的切分。首先确定切分键的最小值和最大值，然后根据配置的通道数量将范围划分为多个段，每个段分配给一个独立的任务并行处理。

## 功能说明

### 配置示例

配置一个从DM数据库同步抽取数据到其他系统的作业:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "dmreader",
                    "parameter": {
                        "connection":  {
                            "url": "ip:port",
                            "table": {
                                "db":"dbname",
                                "name":"table_name"
                            }
                        },
                        "username": "username",
                        "password": "password",
                        "column": ["*"],
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

* 描述：用于配置远程DM数据库的连接信息。基本格式为：`dm://username:password@ip:port/database`。该配置将用户名和密码分离出来，便于加密处理。
* 必选：是
* 默认值：无

#### username

* 描述：用于配置DM数据库的用户名。
* 必选：是
* 默认值：无

#### password

* 描述：用于配置DM数据库的密码。
* 必选：是
* 默认值：无

#### table

描述DM表的信息。

##### db

* 描述：主要用于配置DM表的数据库名。
* 必选：是
* 默认值：无

##### name

* 描述：主要用于配置DM表的表名。
* 必选：是
* 默认值：无

#### split

##### key

* 描述：主要用于配置DM表的切分键。切分键必须是bigInt/string/time类型，且数据需要根据切分键均匀分布。
* 必选：否
* 默认值：无

##### timeAccuracy

* 描述：主要用于配置DM表的时间切分键。主要描述时间的最小单位，如day、min、s、ms、us、ns。
* 必选：否
* 默认值：无

##### range

###### type

* 描述：主要用于配置DM表切分键的默认值类型。值可以是bigInt/string/time。系统不会检查表切分键的类型，但需要确保正确的类型。
* 必选：否
* 默认值：无

###### left

* 描述：主要用于配置DM表切分键的默认最小值。
* 必选：否
* 默认值：无

###### right

* 描述：主要用于配置DM表切分键的默认最大值。
* 必选：否
* 默认值：无

#### column

* 描述：所配置的表中需要同步的列名集合。用户可以使用*来表示默认同步所有列，例如["*"]。支持列裁剪，即只导出部分列。支持列重新排序，即列不需要按照表schema的顺序导出。支持常量配置，用户需要按照DM SQL语法格式：["id", "`table`", "1", "'bazhen.csy'", "null", "left(a,10)", "2.3", "true"]。其中"id"为常规列名，"`table`"为包含保留字的列名，"1"为整型常量，"'bazhen.csy'"为字符串常量，"null"为空指针，"left(a,10)"为表达式，"2.3"为浮点数，"true"为布尔值。
* 必选：是
* 默认值：无

#### where

* 描述：主要用于配置SELECT语句的WHERE条件。
* 必选：否
* 默认值：无

#### querySql

* 描述：允许用户自定义SQL查询语句来过滤数据。当配置该选项时，系统会忽略table、column和where配置，直接使用该配置项的内容进行数据过滤。这在需要多表关联查询或复杂过滤条件时非常有用。
* 必选：否
* 默认值：无

### 类型转换

DMReader支持大部分DM数据类型，但可能存在少量不支持的类型，请仔细检查您的数据类型。

下表列出了DMReader支持的类型转换映射关系：

| go-etl 类型 | DM 数据类型      |
| ----------- | ---------------- |
| bool        | BOOLEAN          |
| bigInt      | BIGINT, INTEGER  |
| decimal     | DECIMAL, NUMERIC |
| string      | VARCHAR, CHAR    |
| time        | DATE, TIME, TIMESTAMP |
| bytes       | BLOB, CLOB       |

## 性能报告

待测试。

## 约束限制

### 数据库编码问题

目前仅支持UTF-8字符集。

## FAQ

1. 如何配置DM数据库连接？

   - 确保DM数据库服务正在运行且可以通过网络访问。
   - 确保连接URL遵循格式：`dm://username:password@ip:port/database`。