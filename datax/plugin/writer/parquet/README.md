# parquetWriter Plugin Documentation

The parquet writer plugin allows you to write data to parquet files using the `github.com/xitongsys/parquet-go` library.

## Configuration Parameters

The following configuration parameters are supported:

- `path`: The path to the output parquet file.
- `rowGroupSize`: Row group size in bytes (default: 134217728 - 128MB).
- `pageSize`: Page size in bytes (default: 8192 - 8KB).
- `batchSize`: Number of records per batch (default: 1000).
- `batchTimeout`: Batch write timeout (default: "1s").

## Implementation Details

This plugin uses the `github.com/xitongsys/parquet-go` library to write parquet files. It converts go-etl's element.Record objects to parquet-compatible data structures and writes them to parquet files.

The implementation leverages the stream file framework in go-etl, making it consistent with other file-based writers like CSV.