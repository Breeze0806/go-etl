# CsvReader插件文档

## 快速介绍

CsvReader插件实现了从csv文件读取数据。在底层实现上，CsvReader通过标准库os以及encoding/csv读取文件。

## 实现原理

CsvReader通过标准库os以及encoding/csv读取文件，并将每一行结果使用go-etl自定义的数据类型拼装为抽象的数据集，并传递给下游Writer处理。

CsvReader通过使用file.Task中定义的读取流程调用go-etl自定义的storage/stream/file的file.InStreamer来实现具体的读取。

## 功能说明

### 配置样例

配置一个从csv文件同步抽取数据到本地的作业:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "cvsreader",
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
                        "delimiter":","
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

- 描述 主要用于配置csv文件的编码类型，目前仅支持utf-8
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

### 类型转换

目前CsvReader支持的csv数据类型需要在column配置中配置，请注意检查你的类型。

下面列出CsvReader针对csv类型转换列表:

| go-etl的类型 | csv数据类型 |
| ------------ | ----------- |
| bigInt       | bigInt      |
| decimal      | decimal     |
| string       | string      |
| time         | time        |
| bool         | bool        |

## 性能报告

待测试