# SQLServerReader Plugin Documentation

## Quick Introduction

The SQLServerReader plugin enables data extraction from SQL Server databases. Under the hood, SQLServerReader connects to a remote SQL Server database via `github.com/denisenkom/go-mssqldb` and executes SQL queries to retrieve data from the SQL Server.

## Implementation Details

SQLServerReader connects to the remote SQL Server database using `github.com/denisenkom/go-mssqldb` and generates SQL queries based on user-provided information. These queries are then sent to the remote SQL Server, and the returned results are assembled into an abstract dataset using go-etl's custom data types before being passed to downstream Writer processing. This differs from directly using `github.com/denisenkom/go-mssqldb`.

SQLServerReader implements specific queries by invoking the query process defined in `dbmsreader` using go-etl's custom `storage/database` DBWrapper. DBWrapper encapsulates many `database/sql` interfaces and abstracts the database dialect. For SQL Server, it uses the dialect implemented in `storage/database/sqlserver`.

## Functionality Overview

### Configuration Example

Configuring a job to synchronize data from a SQL Server database to a local destination:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "sqlserverreader",
                    "parameter": {
                        "username": "sa",
                        "password": "Breeze_0806",
                        "column": ["*"],
                        "connection":  {
                                "url": "sqlserver://192.168.15.130:1433?database=test&encrypt=disable",
                                "table": {
                                    "db":"test",
                                    "schema":"SOURCE",
                                    "name":"mytable"
                                }
                            },
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

### Parameter Explanation

#### url

- Description: Specifies the connection information for the remote SQL Server. The basic format is `sqlserver://ip:port?database=db&encrypt=disable`, where `ip:port` represents the IP address and port of the SQL Server, and `db` is the default database to connect to. See [go-mssqldb](https://github.com/denisenkom/go-mssqldb) for more connection configuration details.
- Required: Yes
- Default: None

#### username

- Description: Specifies the SQL Server database user.
- Required: Yes
- Default: None

#### password

- Description: Specifies the password for the SQL Server database user.
- Required: Yes
- Default: None

#### table

Describes the SQL Server table information.

##### db

- Description: Specifies the database name of the SQL Server table.
- Required: Yes
- Default: None

##### schema

- Description: Specifies the schema name of the SQL Server table.
- Required: Yes
- Default: None

##### name

- Description: Specifies the table name of the SQL Server table.
- Required: Yes
- Default: None

#### column

- Description: Specifies the set of column names to synchronize from the configured table. Use a JSON array to describe the column information. Users can use `*` to select all columns by default, e.g., `["*"]`. Supports column pruning (selecting only specific columns for export) and column reordering (exporting columns in a different order than the table schema). Also supports constant configuration, where users need to follow SQL Server syntax, e.g., `["id", "true", "power(2,3)"]`, where `id` is a regular column name, `'hello'::varchar` is a string constant, `true` is a boolean value, `2.5` is a floating-point number, and `power(2,3)` is a function.
- Required: Yes
- Default: None

#### split

##### key

- Description: Specifies the split key for the SQL Server table. The split key must be of type bigInt/string/time, assuming the data is evenly distributed based on the split key.
- Required: No
- Default: None

##### timeAccuracy

- Description: Specifies the time precision for the SQL Server table's time split key. Used to describe the smallest unit of time, such as day, minute, second, millisecond, microsecond, or nanosecond.
- Required: No
- Default: None

##### range

###### type

- Description: Specifies the default data type for the SQL Server table's split key. Values can be bigInt/string/time. This will check the type of the table's split key, so it's important to ensure the type is correct.
- Required: No
- Default: None

###### left

- Description: Specifies the default minimum value for the SQL Server table's split key.
- Required: No
- Default: None

###### right

- Description: Specifies the default maximum value for the SQL Server table's split key.
- Required: No
- Default: None

#### where

- Description: Specifies the WHERE condition for the SELECT statement.
- Required: No
- Default: None

#### querySql

- Description: In some scenarios, the `where` configuration may not be sufficient to describe the filtering conditions. Users can use this configuration to define custom SQL queries. When this option is configured, the DataX system will ignore the `table`, `column`, and other configurations and directly use the content of this configuration to filter the data. For example, it can be used for data synchronization after performing a join operation on multiple tables, such as `select a,b from table_a join table_b on table_a.id = table_b.id`. When `querySql` is configured, SQLServerReader ignores the configuration of `table`, `column`, and `where` options, and `querySql` takes priority over these options.
- Required: No
- Default: None

#### trimChar

- Description: Specifies whether to remove leading and trailing spaces for SQL Server's char and nchar types.
- Required: No
- Default: false

### Type Conversion

Currently, SQLServerReader supports most SQL Server data types, but there may be some unsupported types. Please check your data types accordingly.

Below is a conversion table for SQLServerReader with respect to SQL Server data types:

| go-etl Type | SQL Server Data Type                                          |
| ----------- | ----------------------------------------------------------- |
| bool        | bit                                                         |
| bigInt      | bigint, int, smallint, tinyint                              |
| decimal     | numeric, real, float                                        |
| string      | char, varchar, text, nchar, nvarchar, ntext                 |
| time        | date, time, datetimeoffset, datetime2, smalldatetime, datetime |
| bytes       | binary, varbinary, varbinary(max)                           |

## Performance Report

Pending testing.

## Constraints and Limitations


## Frequently Asked Questions (FAQ)

