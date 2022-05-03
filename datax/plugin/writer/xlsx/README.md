# XlsxWriter插件文档

## 快速介绍

XlsxWriter插件实现了向xlsx文件写入数据。在底层实现上，XlsxWriter通过github.com/xuri/excelize/v2的流式写入方式写入文件。同时，需要注意的是一个sheet允许保存的数据不应超过1048576，需要计算好导出的sheet数，否则会报错，导致导出失败。此外，对于文件数目的大小要和reader的切分数一致，否则会导致任务无法开始。

## 实现原理

XlsxWriter将reader传来的每一个记录，通过github.com/xuri/excelize/v2的流式写入方式写入文件，这种流式写入方式具有写入速度快，占用内存少的优点。

XlsxWriter通过使用file.Task中定义的写入流程调用go-etl自定义的storage/stream/file的file.OutStreamer来实现具体的读取。

## 功能说明

### 配置样例

配置一个向xlsx文件同步写入数据的作业:

```json
{
    "job":{
        "content":[
            {
                "writer":{
                    "name": "xlsxwriter",
                    "parameter": {
                        "column" :[
                            {
                                "index":"A",
                                "type":"time",
                                "format":"yyyy-MM-dd"
                            }
                        ],
                        "xlsxs":[
                            {
                                "path":"Book1.xlsx",
                                "sheets":["Sheet1"]
                            }
                        ],
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

#### column

- 描述 主要用于配置xlsx文件的列信息数组，如不配置对应信息，则认为对应为string类型
- 必选：是
- 默认值: 无

##### index

- 描述 主要用于配置xlsx文件的列编号，从A开始
- 必选：是
- 默认值: 无

##### type

- 描述 主要用于配置xlsx文件的列类型，主要有boolen,bigInt,decimal,string,time等类型，目前对time仅能使用string类型读取
- 必选：是
- 默认值: 无

##### format

- 描述 主要用于配置xlsx文件的列类型，主要用于配置time类型的格式，使用的是java的joda time格式，如yyyy-MM-dd
- 必选：是
- 默认值: 无

#### xlsxs

- 描述 主要用于配置xlsx文件的信息，可以配置多个文件
- 必选：是
- 默认值: 无

##### path

- 描述 主要用于配置xlsx文件的绝对路径
- 必选：是
- 默认值: 无

###### sheets

- 描述 主要用于配置xlsx文件的sheet名数组
- 必选：是
- 默认值: 无

#### nullFormat

- 描述：文本文件中无法使用标准字符串定义null(空指针)，DataX提供nullFormat定义哪些字符串可以表示为null。例如如果用户配置: nullFormat="\N"，那么如果源头数据是"\N"，DataX视作null字段。
- 必选：否
- 默认值：空字符串

#### batchTimeout

- 描述 主要用于配置每次批量写入超时时间间隔，格式：数字+单位， 单位：s代表秒，ms代表毫秒，us代表微妙。如果超过该时间间隔就直接写入，和batchSize一起调节写入性能。
- 必选：否
- 默认值: 1s

#### batchSize

- 描述 主要用于配置每次批量写入大小，如果超过该大小就直接写入，和batchTimeout一起调节写入性能。
- 必选：否
- 默认值: 1000

### 类型转换

目前XlsxWriter支持的xlsx数据类型需要在column配置中配置，目前xlsx仅支持文本格式的单元格，请注意检查你的类型。

下面列出XlsxWriter针对xlsx类型转换列表:

| go-etl的类型 | xlsx数据类型 |
| ------------ | ----------- |
| bigInt       | bigInt      |
| decimal      | decimal     |
| string       | string      |
| time         | time        |
| bool         | bool        |

## 性能报告

待测试

## 约束限制


## FAQ
