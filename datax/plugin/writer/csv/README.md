# CsvWriter Plugin Documentation

## Quick Introduction

The CsvWriter plugin enables data writing to CSV files. Internally, it utilizes the standard libraries `os` and `encoding/csv` for file writing. Additionally, it's important to ensure that the number of files matches the number of splits defined by the reader, as any mismatch can prevent the task from starting.

## Implementation Principles

The CsvWriter converts each record received from the reader into a string using the standard libraries `os` and `encoding/csv`, and then writes it to the file.

CsvWriter leverages the write process defined in `file.Task` to invoke `file.OutStreamer` from go-etl's custom `storage/stream/file` for specific reading operations.

## Functionality Description

### Configuration Example

Configuring a job to synchronously write data to a CSV file:

```json
{
    "job":{
        "content":[
            {
                "writer":{
                    "name": "cvswriter",
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
                        "delimiter":",",
                        "batchSize":1000,
                        "batchTimeout":"1s"
                    }
                }
            }
        ]
    }
}
```

### Parameter Description

#### path

- Description: Specifies the absolute path(s) of the CSV file(s). Multiple files can be configured.
- Required: Yes
- Default: None

#### column

- Description: Configures the column information array for the CSV file. If not specified, the corresponding data is assumed to be of type string.
- Required: Yes
- Default: None

##### index

- Description: Configures the column number in the CSV file, starting from 1.
- Required: Yes
- Default: None

##### type

- Description: Configures the data type of the CSV column, including options like boolean, bigInt, decimal, string, time, etc.
- Required: Yes
- Default: None

##### format

- Description: Configures the format of the CSV column, primarily used for the time type. It follows the Java Joda time format, such as "yyyy-MM-dd".
- Required: Yes, for time type
- Default: None

#### encoding

- Description: Configures the encoding type of the CSV file, currently supporting only UTF-8 and GBK.
- Required: No
- Default: None

#### delimiter

- Description: Configures the separator for the CSV file, currently supporting only visible symbols like spaces, commas, semicolons, etc.
- Required: No
- Default: None

#### nullFormat

- Description: Standard strings cannot define null (empty pointers) in text files. DataX provides nullFormat to define which strings can represent null. For example, if the user configures nullFormat as "\N", DataX treats "\N" in the source data as a null field.
- Required: No
- Default: Empty string

#### hasHeader

- Description: Determines whether to write the header to the CSV file. If a header exists, it writes the header; otherwise, it writes the column names.
- Required: No
- Default: false

#### header

- Description: Specifies the header array to write to the CSV file. This is only valid when hasHeader is true.
- Required: No
- Default: None

#### compress

- Description: Configures the compression method for the CSV file, currently supporting gz (gzip compression) and zip (zip compression).
- Required: No
- Default: No compression

#### batchTimeout

- Description: Configures the timeout interval for each batch write operation. Format: number + unit, where the unit can be s (seconds), ms (milliseconds), or us (microseconds). If the specified time interval elapses, the data is written directly. This parameter, along with batchSize, helps regulate write performance.
- Required: No
- Default: 1s

#### batchSize

- Description: Configures the size of each batch write operation. If the specified size is exceeded, the data is written directly. This parameter, along with batchTimeout, helps regulate write performance.
- Required: No
- Default: 1000

### Type Conversion

Currently, the supported CSV data types in CsvWriter need to be configured in the column settings. Please check your data types accordingly.

Below is a list of type conversions supported by CsvWriter for CSV data:

| go-etl Type | CSV Data Type |
| --- | --- |
| bigInt | bigInt |
| decimal | decimal |
| string | string |
| time | time |
| bool | bool |

## Performance Report

Pending testing.

## Constraints and Limitations

### Database Encoding Issues
Currently, only the UTF-8 character set is supported.

## FAQ
