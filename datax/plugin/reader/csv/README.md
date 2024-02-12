# CsvReader Plugin Documentation

## Quick Introduction

The CsvReader plugin enables data extraction from CSV files. Under the hood, it utilizes the standard libraries `os` and `encoding/csv` for file reading.

## Implementation Principles

CsvReader leverages the `os` and `encoding/csv` standard libraries to read files. Each row is assembled into an abstract dataset using go-etl's custom data types and passed downstream for further processing by a Writer.

The specific reading process is implemented by invoking go-etl's custom `file.InStreamer` from the reading flow defined in `file.Task`.

## Functionality Description

### Configuration Example

Configuring a job to synchronously extract data from a CSV file to a local destination:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "csvreader",
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
                        "delimiter":","
                    }
                }
            }
        ]
    }
}
```

### Parameter Explanation

#### path

- Description: Specifies the absolute path(s) of the CSV file(s). Multiple files can be configured.
- Required: Yes
- Default: None

#### column

- Description: Configures the column information array for the CSV file. If not specified, the corresponding columns are assumed to be of type string.
- Required: Yes
- Default: None

##### index

- Description: Specifies the column number in the CSV file, starting from 1.
- Required: Yes
- Default: None

##### type

- Description: Configures the data type of the CSV column, including options like boolean, bigInt, decimal, string, time, etc.
- Required: Yes
- Default: None

##### format

- Description: Specifies the format for the column type, particularly useful for the time type. It uses the Java Joda time format, e.g., yyyy-MM-dd.
- Required: Yes, for time type
- Default: None

#### encoding

- Description: Configures the encoding type of the CSV file, currently supporting utf-8 and gbk.
- Required: No
- Default: utf-8

#### delimiter

- Description: Specifies the delimiter used in the CSV file. It supports not only visible symbols like commas or semicolons but also invisible characters such as 0x10 (configured as "\u0010").
- Required: No
- Default: , (comma)

#### nullFormat

- Description: CSV files cannot represent null (empty pointers) using standard strings. The nullFormat parameter defines which strings can be interpreted as null. For example, if nullFormat is set to "\N", then DataX will treat the source data "\N" as a null field.
- Required: No
- Default: Empty string

#### startRow

- Description: Specifies the row number from which to start reading in the CSV file, starting from 1.
- Required: No
- Default: 1

#### comment

- Description: Provides a comment for the CSV file.
- Required: No
- Default: None

#### compress

- Description: Specifies the compression method used for the CSV file, currently supporting gz (gzip compression) and zip (zip compression).
- Required: No
- Default: No compression

### Type Conversion

The CsvReader currently supports CSV data types that need to be configured in the "column" setting. Please ensure you check your data types.

Below is a list of type conversions supported by CsvReader for CSV data:

| go-etl Type | CSV Data Type |
| --- | --- |
| bigInt | bigInt |
| decimal | decimal |
| string | string |
| time | time |
| bool | bool |

## Performance Report

Pending testing.

## Limitations and Constraints


## Frequently Asked Questions (FAQ)

(Note: The FAQ section would typically include common questions and answers related to the plugin's usage, troubleshooting, or best practices. However, as no specific questions were provided, this section remains empty. It can be populated as questions arise from users or developers.)