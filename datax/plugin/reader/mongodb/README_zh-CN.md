# MongodbReader 插件文档

## 快速介绍

MongodbReader 插件能够从 MongoDB 集合中读取数据。内部实现上，MongodbReader 通过 MongoDB 官方 Go 驱动 `go.mongodb.org/mongo-driver` 连接远程 MongoDB 数据库，并执行相应的查询语句从 MongoDB 服务器获取数据。

## 实现原理

MongodbReader 通过 MongoDB 官方 Go 驱动连接远程 MongoDB 数据库。根据用户配置的信息，生成查询语句并发送到远程 MongoDB 服务器。查询返回的结果会被组装成抽象数据集，使用 go-etl 的自定义数据类型传递给下游 Writer 处理。

插件基于 ObjectId 范围进行数据读取任务的拆分。首先确定集合中的最小和最大 ObjectId，然后根据配置的通道数将 ObjectId 范围划分为多个段，每个段分配给一个独立的任务进行并行处理。

## 功能说明

### 配置示例

配置一个从 MongoDB 集合同步数据到另一个 MongoDB 集合的作业：

```json
{
    "job": {
        "content": [
            {
                "reader": {
                    "name": "mongodbreader",
                    "parameter": {
                        "connection": {
                            "address": "localhost:27017",
                            "username": "root",
                            "password": "123456",
                            "table": {
                                "db": "test_database",
                                "collection": "users"
                            }
                        },
                        "column": [
                            {
                                "name": "_id",
                                "type": "string"
                            },
                            {
                                "name": "name",
                                "type": "string"
                            },
                            {
                                "name": "email",
                                "type": "string"
                            },
                            {
                                "name": "age",
                                "type": "int"
                            },
                            {
                                "name": "created_at",
                                "type": "date"
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

#### connection

##### address

- 描述：用于配置 MongoDB 服务器地址，格式为 `host:port`。
- 必选：是
- 默认值：无

##### username

- 描述：用于配置 MongoDB 认证用户名。
- 必选：否
- 默认值：无

##### password

- 描述：用于配置 MongoDB 认证密码。
- 必选：否
- 默认值：无

##### table

###### db

- 描述：用于配置需要读取的集合所在的数据库名。
- 必选：是
- 默认值：无

###### collection

- 描述：用于配置需要读取的集合名。
- 必选：是
- 默认值：无

#### column

- 描述：描述需要从配置的集合中同步的字段数组。用户可以使用 JSON 数组格式来描述字段信息。

  每个字段配置包括：
  - `name`：MongoDB 文档中的字段名
  - `type`：目标数据类型（string, int, date 等）
  - `spliter`：对于数组字段，转换为字符串时使用的分隔符

  示例：
  ```json
  {
    "column": [
        {
            "name": "_id",
            "type": "string"
        },
        {
            "name": "tags",
            "type": "string",
            "spliter": ","
        }
    ]
  }
  ```

  支持的数据类型：
  - `string`：字符串类型
  - `int`：整数类型
  - `date`：日期/时间类型
  - `Array`：数组类型（使用分隔符转换为字符串）

- 必选：是
- 默认值：无

#### split_key

- 描述：用于配置数据分片的字段。目前仅支持 "_id" 作为分片键。
- 必选：否
- 默认值："_id"

### 类型转换

目前 MongodbReader 支持大多数 MongoDB 类型，并适当转换为 go-etl 内部类型：

| go-etl 类型 | MongoDB 数据类型           | 说明                                       |
| ----------- | -------------------------- | ------------------------------------------ |
| bigInt      | int, int32, int64          |                                            |
| decimal     | float32, float64           |                                            |
| string      | string                     | 也用于 ObjectId 和其他类型                 |
| time        | primitive.DateTime         |                                            |
| bool        | bool                       |                                            |
| bytes       | []byte                     |                                            |

特殊处理：
1. `primitive.ObjectID`：默认转换为十六进制字符串
2. `[]interface{}`（数组）：可使用自定义分隔符转换为字符串
3. `bson.M`（嵌入文档）：序列化为 JSON 字符串
4. 其他类型：转换为字符串表示

## 性能报告

待测试。

## 约束限制

### 数据库编码问题

目前仅支持 utf8 字符集。

### 分片键限制

目前仅支持 "_id" 字段作为并行处理的分片键。