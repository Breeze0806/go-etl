# parquetWriter Plugin Documentation

## Quick Introduction

The parquetWriter plugin enables data writing to Apache Parquet files. Internally, it utilizes the `github.com/xitongsys/parquet-go` library for file writing. Additionally, it's important to ensure that the number of files matches the number of splits defined by the reader, as any mismatch can prevent the task from starting.

## Implementation Principles

The parquetWriter converts each record received from the reader into parquet-compatible data structures using the `github.com/xitongsys/parquet-go` library, and then writes it to the file.

parquetWriter leverages the write process defined in `file.Task` to invoke `file.OutStreamer` from go-etl's custom `storage/stream/file` for specific writing operations.

## Functionality Description

### Configuration Example

Configuring a job to synchronously write data to a parquet file:

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

### Parameter Description

#### path

- Description: Specifies the absolute path(s) of the parquet file(s).
- Required: Yes
- Default: None

#### column

- Description: Configures the column information array for the parquet file. If not specified, the corresponding data is assumed to be of type string.
- Required: Yes
- Default: None

##### name

- Description: Configures the column name in the parquet file.
- Required: Yes
- Default: None

##### type

- Description: Configures the data type of the parquet column, including options like boolean, bigInt, decimal, string, time, etc.
- Required: Yes
- Default: None

### Type Conversion

Currently, the supported parquet data types in parquetWriter need to be configured in the column settings. Please check your data types accordingly.

Below is a list of type conversions supported by parquetWriter for parquet data:

| go-etl Type | parquet Data Type |
| --- | --- |
| bigInt | INT32, INT64 |
| decimal | FLOAT, DOUBLE |
| string | BYTE_ARRAY, FIXED_LEN_BYTE_ARRAY |
| time | INT64 |
| bool | BOOLEAN |

## Performance Report

Pending testing.

## Constraints and Limitations

## FAQ