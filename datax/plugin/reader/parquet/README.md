# parquetReader Plugin Documentation

## Quick Introduction

The parquetReader plugin enables data extraction from Apache Parquet files. Under the hood, it utilizes the `github.com/xitongsys/parquet-go` library for file reading.

## Implementation Principles

parquetReader leverages the `github.com/xitongsys/parquet-go` library to read files. Each row is assembled into an abstract dataset using go-etl's custom data types and passed downstream for further processing by a Writer.

The specific reading process is implemented by invoking go-etl's custom `file.InStreamer` from the reading flow defined in `file.Task`.

## Functionality Description

### Configuration Example

Configuring a job to synchronously extract data from a parquet file to a local destination:

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

### Parameter Explanation

#### path

- Description: Specifies the absolute path(s) of the parquet file(s).
- Required: Yes
- Default: None

#### column

- Description: Configures the column information array for the parquet file. If not specified, the corresponding columns are assumed to be of type string.
- Required: Yes
- Default: None

##### name

- Description: Specifies the column name in the parquet file.
- Required: Yes
- Default: None

##### type

- Description: Configures the data type of the parquet column, including options like boolean, bigInt, decimal, string, time, etc.
- Required: Yes
- Default: None

### Type Conversion

The parquetReader currently supports parquet data types that need to be configured in the "column" setting. Please ensure you check your data types.

Below is a list of type conversions supported by parquetReader for parquet data:

| go-etl Type | parquet Data Type |
| --- | --- |
| bigInt | INT32, INT64 |
| decimal | FLOAT, DOUBLE |
| string | BYTE_ARRAY, FIXED_LEN_BYTE_ARRAY |
| time | INT64 |
| bool | BOOLEAN |

## Performance Report

Pending testing.

## Limitations and Constraints

## Frequently Asked Questions (FAQ)