# parquetWriter 插件文档

parquet writer 插件允许您将数据写入 parquet 文件，使用 `github.com/xitongsys/parquet-go` 库实现。

## 配置参数

支持以下配置参数：

- `path`：输出 parquet 文件的路径。
- `rowGroupSize`：行组大小（字节）（默认值：134217728 - 128MB）。
- `pageSize`：页面大小（字节）（默认值：8192 - 8KB）。
- `batchSize`：单个批次写入的记录数（默认值：1000）。
- `batchTimeout`：批处理写入的超时时间（默认值："1s"）。

### 配置示例

```json
{
  "name": "parquetwriter",
  "parameter": {
    "path": "/path/to/output.parquet",
    "rowGroupSize": 134217728,
    "pageSize": 8192,
    "batchSize": 1000,
    "batchTimeout": "1s"
  }
}
```

## 实现细节

该插件使用 `github.com/xitongsys/parquet-go` 库来写入 parquet 文件。它将 go-etl 的 element.Record 对象转换为与 parquet 兼容的数据结构，并将其写入 parquet 文件。

该实现利用了 go-etl 中的流文件框架，使其与其他基于文件的写入器（如 CSV）保持一致。