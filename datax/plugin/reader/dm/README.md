# DMReader Plugin Documentation

## Quick Introduction

The DMReader plugin enables data extraction from Dameng (DM) databases. Internally, DMReader connects to a remote DM database using the official DM Go driver `gitee.com/chunanyong/dm` and executes corresponding queries to retrieve data from the DM server.

## Implementation Principles

DMReader connects to a remote DM database using the official DM Go driver. Based on user-provided configuration information, it generates queries and sends them to the remote DM server. The returned results from these queries are assembled into an abstract dataset using go-etl's custom data types and passed to downstream Writers for processing.

The plugin splits data reading tasks based on a split key. It first determines the minimum and maximum values of the split key in the table, then divides the range into multiple segments based on the number of channels configured, and assigns each segment to a separate task for parallel processing.

## Functionality Description

### Configuration Example

Configuring a job to synchronize data from a DM database to another system:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "dmreader",
                    "parameter": {
                        "connection":  {
                            "url": "dm://username:password@ip:port/database",
                            "table": {
                                "db":"dbname",
                                "name":"table_name"
                            }
                        },
                        "username": "username",
                        "password": "password",
                        "column": ["*"],
                        "split" : {
                            "key":"id"
                        },
                        "where": "",
                        "querySql":["select a,b from table_a join table_b on table_a.id = table_b.id"]
                    }
                }
            }
        ]
    }
}
```

### Parameter Details

#### url

* Description: Specifies the connection information for the remote DM database. The basic format is: `dm://username:password@ip:port/database`. This configuration separates the username and password for easier encryption.
* Required: Yes
* Default: None

#### username

* Description: The username for the DM database.
* Required: Yes
* Default: None

#### password

* Description: The password for the DM database.
* Required: Yes
* Default: None

#### table

Describes the DM table information.

##### db

* Description: Primarily used to configure the database name for the DM table.
* Required: Yes
* Default: None

##### name

* Description: Primarily used to configure the table name for the DM table.
* Required: Yes
* Default: None

#### split

##### key

* Description: Primarily used to configure the splitting key for the DM table. The splitting key must be of type bigInt/string/time, assuming that the data is evenly distributed based on the splitting key.
* Required: No
* Default: None

##### timeAccuracy

* Description: Primarily used to configure the time splitting key for the DM table. It is mainly used to describe the smallest unit of time, such as day, minute, second, millisecond, microsecond, nanosecond.
* Required: No
* Default: None

##### range

###### type

* Description: Primarily used to configure the default type for the splitting key of the DM table. The value can be bigInt/string/time. The system does not check the type in the table splitting key, but it is important to ensure the correct type.
* Required: No
* Default: None

###### left

* Description: Primarily used to configure the default minimum value for the splitting key of the DM table.
* Required: No
* Default: None

###### right

* Description: Primarily used to configure the default maximum value for the splitting key of the DM table.
* Required: No
* Default: None

#### column

* Description: An array of column names to be synchronized from the configured table. Users can use the asterisk (*) to indicate that all columns should be used by default, for example, ["*"]. Column pruning is supported, meaning that only a subset of columns can be exported. Column reordering is also supported, meaning that columns do not need to be exported in the same order as the table schema. Constant configuration is supported, where users need to follow the DM SQL syntax format: ["id", "`table`", "1", "'bazhen.csy'", "null", "left(a,10)", "2.3", "true"]. In this example, "id" is a regular column name, "`table`" is a column name that contains reserved words, "1" is an integer constant, "'bazhen.csy'" is a string constant, "null" is a null pointer, "left(a,10)" is an expression, "2.3" is a floating-point number, and "true" is a boolean value.
* Required: Yes
* Default: None

#### where

* Description: Primarily used to configure the WHERE condition for the SELECT statement.
* Required: No
* Default: None

#### querySql

* Description: Allows the user to define a custom SQL query for data filtering. When this option is configured, the system ignores the table, column, and where settings, and directly uses the content of this configuration for data filtering. This is useful for scenarios that require joining multiple tables or complex filtering conditions.
* Required: No
* Default: None

### Type Conversions

DMReader supports most DM data types, but there may be some unsupported types. Please check your data types carefully.

The following table lists the type conversion mappings supported by DMReader:

| go-etl Type | DM Data Type     |
| ----------- | ---------------- |
| bool        | BOOLEAN          |
| bigInt      | BIGINT, INTEGER  |
| decimal     | DECIMAL, NUMERIC |
| string      | VARCHAR, CHAR    |
| time        | DATE, TIME, TIMESTAMP |
| bytes       | BLOB, CLOB       |

## Performance Report

Pending testing.

## Constraints and Limitations

### Database Encoding Issues

Currently, only the UTF-8 character set is supported.

## FAQ

1. How to configure the DM database connection?

   - Make sure the DM database service is running and accessible over the network.
   - Ensure the connection URL follows the format: `dm://username:password@ip:port/database`.