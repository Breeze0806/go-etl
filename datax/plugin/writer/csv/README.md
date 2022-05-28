# CsvWriter插件文档

## 快速介绍

CsvWriter插件实现了向csv文件写入数据。在底层实现上，CsvWriter通过标准库os以及encoding/csv写入文件。此外，对于文件数目的大小要和reader的切分数一致，否则会导致任务无法开始。

## 实现原理

CsvWriter将reader传来的每一个记录，通过标准库os以及encoding/csv转换成字符串写入文件。

CsvWriter通过使用file.Task中定义的写入流程调用go-etl自定义的storage/stream/file的file.OutStreamer来实现具体的读取。

## 功能说明

### 配置样例

配置一个向csv文件同步写入数据的作业:

```json
{
    "job":{
        "content":[
            {
                "writer":{
                    "name": "cvswriter",
                    "parameter": {
                        "path":["a.txt","b.txt"],
                        "column":[
                            {
                                "index":"1",
                                "type":"time",
                                "format":"yyyy-MM-dd"
                            }
                        ],
                        "encoding":"utf-8",
                        "delimiter":",",
                        "batchSize":1000,
                        "batchTimeout":"1s"
                    }
                }
            }
        ]
    }
}
```

### 参数说明

#### path

- 描述 主要用于配置csv文件的绝对路径，可以配置多个文件
- 必选：是
- 默认值: 无

#### column

- 描述 主要用于配置csv文件的列信息数组，如不配置对应信息，则认为对应为string类型
- 必选：是
- 默认值: 无

##### index

- 描述 主要用于配置csv文件的列编号，从1开始
- 必选：是
- 默认值: 无

##### type

- 描述 主要用于配置csv文件的列类型，主要有boolen,bigInt,decimal,string,time等类型
- 必选：是
- 默认值: 无

##### format

- 描述 主要用于配置csv文件的列类型，主要用于配置time类型的格式，使用的是java的joda time格式，如yyyy-MM-dd
- 必选：是
- 默认值: 无

#### encoding

- 描述 主要用于配置csv文件的编码类型，目前仅支持utf-8和gbk
- 必选：否
- 默认值: 无

#### delimiter

- 描述 主要用于配置csv文件的分隔符，目前仅支持空格和可见的符号，如逗号，分号等
- 必选：否
- 默认值: 无

#### nullFormat

- 描述：文本文件中无法使用标准字符串定义null(空指针)，DataX提供nullFormat定义哪些字符串可以表示为null。例如如果用户配置: nullFormat="\N"，那么如果源头数据是"\N"，DataX视作null字段。
- 必选：否
- 默认值：空字符串

#### hasHeader

- 描述：是否写入csv文件的列头，当存在header，写入header，而不存在时，写入列名。
- 必选：否
- 默认值：false

#### header

- 描述：写入csv文件的列头数组，仅hasHeader有效。
- 必选：否
- 默认值：无

#### batchTimeout

- 描述 主要用于配置每次批量写入超时时间间隔，格式：数字+单位， 单位：s代表秒，ms代表毫秒，us代表微妙。如果超过该时间间隔就直接写入，和batchSize一起调节写入性能。
- 必选：否
- 默认值: 1s

#### batchSize

- 描述 主要用于配置每次批量写入大小，如果超过该大小就直接写入，和batchTimeout一起调节写入性能。
- 必选：否
- 默认值: 1000

### 类型转换

目前CsvWriter支持的csv数据类型需要在column配置中配置，请注意检查你的类型。

下面列出CsvWriter针对csv类型转换列表:

| go-etl的类型 | csv数据类型 |
| ------------ | ----------- |
| bigInt       | bigInt      |
| decimal      | decimal     |
| string       | string      |
| time         | time        |
| bool         | bool        |

## 性能报告

待测试

## 约束限制

### 数据库编码问题
目前仅支持utf8字符集

## FAQ
