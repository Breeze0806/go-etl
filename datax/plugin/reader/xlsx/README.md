# XlsxReader插件文档

## 快速介绍

XlsxReader插件实现了从csv文件读取数据。在底层实现上，XlsxReader通过github.com/xuri/excelize/v2的流式读取方式读取文件。

## 实现原理

XlsxReader通过github.com/xuri/excelize/v2的流式读取方式读取文件，并将每一行结果使用go-etl自定义的数据类型拼装为抽象的数据集，并传递给下游Writer处理。这种流式读取方式具有读取速度快，占用内存少的优点。

XlsxReader通过使用file.Task中定义的读取流程调用go-etl自定义的storage/stream/file的file.InStreamer来实现具体的读取。

## 功能说明

### 配置样例

配置一个从xlsx文件同步抽取数据到本地的作业:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "xlsxreader",
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
                                "path":"",
                                "sheets":["",""]   
                            }
                        ]
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

### 类型转换

目前XlsxReader支持的cvs数据类型需要在column配置中配置，目前xlsx仅支持文本格式的单元格，请注意检查你的类型。

下面列出XlsxReader针对 xlsx类型转换列表:

| go-etl的类型 | xlsx数据类型 |
| ------------ | ------------ |
| bigInt       | bigInt       |
| decimal      | decimal      |
| string       | string       |
| time         | time         |
| bool         | bool         |

## 性能报告

待测试

## 约束限制

## FAQ
