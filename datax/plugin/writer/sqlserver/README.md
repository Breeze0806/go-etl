# SQLServerWriter Plugin Documentation

## Quick Introduction

The SQLServerWriter plugin enables writing data to SQL Server databases. Under the hood, SQLServerWriter connects to remote SQL Server databases using github.com/microsoft/go-mssqldb and database/sql, executing corresponding SQL statements to write data into the SQL Server database.

## Implementation Principles

SQLServerWriter connects to remote SQL Server databases using github.com/microsoft/go-mssqldb. It generates write SQL statements based on user-configured information and go-etl's custom data types from the Reader. These statements are then sent to the remote SQL Server database for execution.

SQLServerWriter implements specific queries by invoking go-etl's custom storage/database DBWrapper, which is defined in the dbmswriter query process. DBWrapper encapsulates many interfaces of database/sql and abstracts the database dialect, Dialect. For SQL Server, it adopts the Dialect implemented by storage/database/sqlserver.

Based on your configured `writeMode`, it generates:

- `insert into...` (If there is a conflict with the primary key/unique index, the conflicting row will not be inserted.)

**Or**

- bulk copy, i.e., `insert bulk ...` which behaves similarly to insert into but is much faster. However, currently, it cannot insert records containing null values for unknown reasons.

## Functionality Description

### Configuration Example

Configuring a job to synchronously write data to a SQL Server database:

```json
{
    "job":{
        "content":[
            {
               "writer":{
                    "name": "sqlserverwriter",
                    "parameter": {
                        "username": "sa",
                        "password": "Breeze_0806",
                        "writeMode": "insert",
                        "column": ["*"],
                        "preSql": ["create table a like b"],
                        "postSql": ["drop table a"],
                        "connection":  {
                                "url": "sqlserver://192.168.15.130:1433?database=test&encrypt=disable",
                                "table": {
                                    "db":"test",
                                    "schema":"dest",
                                    "name":"mytable"
                                }
                         },
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

#### url

- Description: Primarily used to configure the connection information for the remote end. The basic configuration format is: "sqlserver://ip:port?database=db&encrypt=disable". Here, ip:port represents the IP address and port of the SQL Server database, and db indicates the default database to connect to. For detailed connection configuration information, see [go-mssqldb](https://github.com/microsoft/go-mssqldb).
- Required: Yes
- Default: None

#### username

- Description: Primarily used to configure the SQL Server database username.
- Required: Yes
- Default: None

#### password

- Description: Primarily used to configure the SQL Server database password.
- Required: Yes
- Default: None

#### table

Describes the SQL Server table information.

##### db

- Description: Primarily used to configure the database name of the SQL Server table.
- Required: Yes
- Default: None

##### schema

- Description: Primarily used to configure the schema name of the SQL Server table.
- Required: Yes
- Default: None

##### name

- Description: Primarily used to configure the table name of the SQL Server table.
- Required: Yes
- Default: None

#### column

- Description: The set of column names that need to be synchronized in the configured table, described using a JSON array. Users can use "*" to represent all columns by default, e.g., ["*"]. Column trimming is supported, meaning users can select a subset of columns for export. Column reordering is also supported, meaning columns can be exported in an order different from the table schema. Constant configuration is supported, where users need to follow the SQL Server syntax format: ["id", "true", "power(2,3)"] where id is a regular column name, 'hello'::varchar is a string constant, true is a boolean value, 2.5 is a floating-point number, and power(2,3) is a function.
- Required: Yes
- Default: None

#### writeMode

- Description: Write mode. "insert" represents writing data using the insert into method, while "copyIn" represents bulk copy insertion.
- Required: No
- Default: insert

#### batchTimeout

- Description: Primarily used to configure the timeout interval for each batch write operation. The format is: number + unit, where the unit can be s for seconds, ms for milliseconds, or us for microseconds. If the specified time interval is exceeded, the data will be written directly. This parameter, along with batchSize, helps adjust write performance.
- Required: No
- Default: 1s

#### batchSize

- Description: Primarily used to configure the batch write size. If the specified size is exceeded, the data will be written directly. This parameter, along with batchTimeout, helps adjust write performance.
- Required: No
- Default: 1000

#### preSql

- Description: Primarily used for the set of SQL statements to be executed before writing data. Do not use select statements as they will cause errors.
- Required: No
- Default: None

#### postSql

- Description: Primarily used for the set of SQL statements to be executed after writing data. Do not use select statements as they will cause errors.
- Required: No
- Default: None

### Type Conversion

Currently, SQLServerReader supports most SQL Server types, but there are some individual types that are not supported. Please check your data types accordingly.

Below is a conversion table for SQLServerReader against SQL Server data types:

| go-etl Type | SQL Server Data Type                                          |
| ----------- | ----------------------------------------------------------- |
| bool        | bit                                                         |
| bigInt      | bigint, int, smallint, tinyint                              |
| decimal     | numeric, real, float                                         |
| string      | char, varchar, text, nchar, nvarchar, ntext                 |
| time        | date, time, datetimeoffset, datetime2, smalldatetime, datetime |
| bytes       | binary, varbinary, varbinary(max)                           |

## Performance Report

Pending testing.

## Constraints and Limitations


## FAQ
