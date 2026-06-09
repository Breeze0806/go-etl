# parquetWriter 插件文档

## 快速介绍

parquetWriter插件实现了将数据写入到parquet文件。在底层实现上，parquetWriter通过`github.com/xitongsys/parquet-go`库写入文件。此外，需要注意的是，文件数量必须与reader切分的split数量一致，否则无法正确执行任务。

## 实现原理

parquetWriter通过将读取到的每条记录转换为parquet兼容的数据结构，使用`github.com/xitongsys/parquet-go`库写入文件。

parquetWriter通过使用file.Task中定义的写入流程调用go-etl自定义的storage/stream/file的file.OutStreamer来实现具体的写入。

## 功能说明

### 配置样例

配置一个写入parquet文件同步作业:

```json
{
    "job":{
        "content":[
            {
                "writer":{
                    "name": "parquetwriter",
                    "parameter": {
                        "path": ["output.parquet"],
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

- 描述：指定parquet文件的绝对路径数组
- 必选：是
- 默认值：无

#### column

- 描述：配置parquet文件的列信息数组，如不配置则对应数据为string类型
- 必选：是
- 默认值：无

##### name

- 描述：配置parquet文件的列名
- 必选：是
- 默认值：无

##### type

- 描述：配置parquet文件的列类型，主要有boolean,bigInt,decimal,string,time等类型
- 必选：是
- 默认值：无

### 类型转换

目前parquetWriter支持的parquet数据类型需要在column配置中配置，请注意检查你的类型。

下面列出parquetWriter针对parquet类型转换列表:

| go-etl的类型 | parquet数据类型 |
| --- | --- |
| bigInt | INT32, INT64 |
| decimal | FLOAT, DOUBLE |
| string | BYTE_ARRAY, FIXED_LEN_BYTE_ARRAY |
| time | INT64 |
| bool | BOOLEAN |

## 性能报告

待测试

## 约束限制

## FAQ