# XlsxWriter Plugin Documentation

## Quick Introduction

The XlsxWriter plugin enables writing data to xlsx files. Under the hood, XlsxWriter utilizes the streaming write method of github.com/xuri/excelize/v2 for file writing. It's important to note that the maximum amount of data allowed per sheet is 1,048,576. Therefore, the number of sheets to export must be calculated appropriately to avoid errors and export failures. Additionally, the number of files must align with the number of splits in the reader, or else the task cannot commence.

## Implementation Principles

XlsxWriter takes each record passed from the reader and writes it to the file using the streaming write method of github.com/xuri/excelize/v2. This streaming write approach offers the advantages of fast write speeds and low memory usage.

XlsxWriter achieves specific reads by utilizing the writing process defined in file.Task, which calls the file.OutStreamer of go-etl's custom storage/stream/file.

## Functionality Description

### Configuration Example

Configuring a job to synchronously write data to an xlsx file:

```json
{
    "job":{
        "content":[
            {
                "writer":{
                    "name": "xlsxwriter",
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
                                "path":"Book1.xlsx",
                                "sheets":["Sheet1"]
                            }
                        ],
                        "batchTimeout": "1s",
                        "batchSize":1000
                    }
                }
            }
        ]
    }
}
```

### Parameter Description

#### column

- Description: Configures the column information array for the xlsx file. If corresponding information is not configured, it is assumed to be of type string.
- Required: Yes
- Default: None

##### index

- Description: Configures the column index for the xlsx file, starting from A.
- Required: Yes
- Default: None

##### type

- Description: Configures the column type for the xlsx file, including types such as boolean, bigInt, decimal, string, and time. Currently, only the string type can be used for time.
- Required: Yes
- Default: None

##### format

- Description: Configures the format for the column type in the xlsx file, primarily used for configuring the format of the time type using Java's joda time format, e.g., yyyy-MM-dd.
- Required: Yes
- Default: None

#### xlsxs

- Description: Configures information about the xlsx file(s), allowing for the configuration of multiple files.
- Required: Yes
- Default: None

##### path

- Description: Configures the absolute path of the xlsx file.
- Required: Yes
- Default: None

###### sheets

- Description: Configures an array of sheet names for the xlsx file.
- Required: Yes
- Default: None

#### nullFormat

- Description: Standard strings cannot define null (null pointers) in text files. DataX provides nullFormat to define which strings can represent null. For example, if the user configures: nullFormat="\N", then DataX treats "\N" as a null field if it appears in the source data.
- Required: No
- Default: Empty string

#### hasHeader

- Description: Determines whether to write the column headers to the csv file. When headers exist, they are written; otherwise, column names are written.
- Required: No
- Default: false

#### header

- Description: Writes an array of column headers to the csv file, effective only when hasHeader is true.
- Required: No
- Default: None

#### sheetRow

- Description: Specifies the maximum number of rows per sheet, with a maximum of 1,048,576.
- Required: No
- Default: 1048576

#### batchTimeout

- Description: Configures the timeout interval for each batch write operation. Format: number + unit, where the unit can be s for seconds, ms for milliseconds, or us for microseconds. If the specified time interval is exceeded, the data is written directly. This parameter, along with batchSize, helps regulate write performance.
- Required: No
- Default: 1s

#### batchSize

- Description: Configures the size of each batch write operation. If the specified size is exceeded, the data is written directly. This parameter, along with batchTimeout, helps regulate write performance.
- Required: No
- Default: 1000

### Type Conversion

Currently, the xlsx data types supported by XlsxWriter need to be configured in the column settings. Only text-formatted cells are supported in xlsx files, so please check your types accordingly.

Below is a list of type conversions supported by XlsxWriter for xlsx data types:

| go-etl Type | xlsx Data Type |
| --- | --- |
| bigInt | bigInt |
| decimal | decimal |
| string | string |
| time | time |
| bool | bool |

## Performance Report

Pending testing.

## Constraints and Limitations


## FAQ
