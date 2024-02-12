# XlsxReader Plugin Documentation

## Quick Introduction

The XlsxReader plugin enables data reading from XLSX files. Under the hood, XlsxReader utilizes the streaming read method from github.com/xuri/excelize/v2 to process files.

## Implementation Principle

XlsxReader employs the streaming read method from github.com/xuri/excelize/v2 to read files, assembling each row of data into an abstract dataset using go-etl's custom data types. This dataset is then passed downstream to the Writer for further processing. This streaming approach offers fast reading speeds and low memory usage.

XlsxReader implements the specific reading process by utilizing the reading workflow defined in file.Task and invoking go-etl's custom file.InStreamer from storage/stream/file.

## Functionality Description

### Configuration Example

Configuring a job to synchronize data extraction from an XLSX file to a local destination:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "xlsxreader",
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
                                "path":"",
                                "sheets":["",""]   
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

#### column

- Description: Primarily used to configure the column information array for the XLSX file. If corresponding information is not configured, the columns are assumed to be of the string type.
- Required: Yes
- Default: None

##### index

- Description: Primarily used to configure the column index for the XLSX file, starting from A.
- Required: Yes
- Default: None

##### type

- Description: Primarily used to configure the column type for the XLSX file, including types such as boolean, bigInt, decimal, string, and time. Currently, only the string type can be used for reading time.
- Required: Yes
- Default: None

##### format

- Description: Primarily used to configure the column format for the XLSX file, specifically for configuring the format of the time type. It uses the Java Joda Time format, such as "yyyy-MM-dd".
- Required: Yes, if the type is time.
- Default: None

#### xlsxs

- Description: Primarily used to configure information about the XLSX file(s), allowing for the configuration of multiple files.
- Required: Yes
- Default: None

##### path

- Description: Primarily used to configure the absolute path of the XLSX file.
- Required: Yes
- Default: None

##### sheets

- Description: Primarily used to configure an array of sheet names within the XLSX file.
- Required: Yes
- Default: None

#### nullFormat

- Description: XLSX files cannot define null (empty pointers) using standard strings. DataX provides the nullFormat parameter to define which strings can represent null. For example, if the user configures nullFormat="\N", then DataX treats "\N" in the source data as a null field.
- Required: No
- Default: Empty string

#### startRow

- Description: Specifies the row number from which to start reading in the XLSX file, starting from 1.
- Required: No
- Default: 1

### Type Conversion

Currently, the XLSX data types supported by XlsxReader need to be configured in the "column" setting. It should be noted that XLSX currently only supports text-formatted cells, so please check your data types accordingly.

Below is a list of type conversions supported by XlsxReader for XLSX data:

| go-etl Type | XLSX Data Type |
| --- | --- |
| bigInt | bigInt |
| decimal | decimal |
| string | string |
| time | time |
| bool | bool |

## Performance Report

Pending testing.

## Constraints and Limitations

- Currently, only text-formatted cells are supported in XLSX files.
- The time type in XLSX can only be read as a string type due to limitations in the underlying library.
- Memory usage and reading speeds may vary depending on the size and complexity of the XLSX file.

## FAQ
