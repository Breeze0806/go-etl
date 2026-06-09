# parquetReader插件文档

## 快速介绍

parquetReader插件实现了从parquet文件读取数据。在底层实现上，parquetReader通过`github.com/xitongsys/parquet-go`库读取文件。

## 实现原理

parquetReader通过`github.com/xitongsys/parquet-go`库读取文件，并将每一行结果使用go-etl自定义的数据类型拼装为抽象的数据集，并传递给下游Writer处理。

parquetReader通过使用file.Task中定义的读取流程调用go-etl自定义的storage/stream/file的file.InStreamer来实现具体的读取。

## 功能说明

### 配置样例

配置一个从parquet文件同步抽取数据到本地的作业:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "parquetreader",
                    "parameter": {
                        "path":["a.parquet"],
                        "column":[
                            {
                                "name":"col1",
                                "type":"string"
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

#### path

- 描述 主要用于配置parquet文件的绝对路径数组
- 必选：是
- 默认值: 无

#### column

- 描述 主要用于配置parquet文件的列信息数组，如不配置对应信息，则认为对应为string类型
- 必选：是
- 默认值: 无

##### name

- 描述 主要用于配置parquet文件的列名
- 必选：是
- 默认值: 无

##### type

- 描述 主要用于配置parquet文件的列类型，主要有boolen,bigInt,decimal,string,time等类型
- 必选：是
- 默认值: 无

### 类型转换

目前parquetReader支持的parquet数据类型需要在column配置中配置，请注意检查你的类型。

下面列出parquetReader针对parquet类型转换列表:

| go-etl的类型 | parquet数据类型 |
| ------------ | ----------- |
| bigInt       | INT32, INT64 |
| decimal      | FLOAT, DOUBLE |
| string       | BYTE_ARRAY, FIXED_LEN_BYTE_ARRAY |
| time         | INT64 |
| bool         | BOOLEAN |

## 性能报告

待测试

## 约束限制

## FAQ