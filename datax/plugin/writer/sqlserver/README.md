# SQLServerWriter Plugin Documentation

## Quick Introduction

The SQLServerWriter plugin enables writing data to SQL Server databases. Under the hood, SQLServerWriter connects to remote SQL Server databases using github.com/microsoft/go-mssqldb and database/sql, executing corresponding SQL statements to write data into the SQL Server database.

## Implementation Principles

SQLServerWriter connects to remote SQL Server databases using github.com/microsoft/go-mssqldb. It generates write SQL statements based on user-configured information and go-etl's custom data types from the Reader. These statements are then sent to the remote SQL Server database for execution.

SQLServerWriter implements specific queries by invoking go-etl's custom storage/database DBWrapper, which is defined in the dbmswriter query process. DBWrapper encapsulates many interfaces of database/sql and abstracts the database dialect, Dialect. For SQL Server, it adopts the Dialect implemented by storage/database/sqlserver.

Based on your configured `writeMode`, it generates:

- `insert into...` (If there is a conflict with the primary key/unique index, the conflicting row will not be inserted.)

**Or**

- bulk copy, i.e., `BULK INSERT ...` which behaves similarly to insert into but is much faster. Compared to `insert into...`, we recommend this writing mode more.

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

#### bulkOption

- Description: Primarily used for configuring bulk write operations in `copyIn`, affecting the settings for `BULK INSERT`. For details, reference to [BULK INSERT](https://learn.microsoft.com/en-us/sql/t-sql/statements/bulk-insert-transact-sql?view=sql-server-ver16).
- Required: No
- Default Value: None

##### CheckConstraints

+ Description: Means `CHECK_CONSTRAINTS `. Specifies that all constraints on the target table or view must be checked during the bulk-import operation. Without the `CHECK_CONSTRAINTS` option, any CHECK and FOREIGN KEY constraints are ignored, and after the operation, the constraint on the table is marked as not-trusted.
+ Required: No
+ Default Value: None

##### FireTriggers

+ Description: Means `FIRE_TRIGGERS`. Specifies that any insert triggers defined on the destination table execute during the bulk-import operation. If triggers are defined for INSERT operations on the target table, they're fired for every completed batch.If FIRE_TRIGGERS isn't specified, no insert triggers execute.
+ Required: No
+ Default Value: None

##### KeepNulls

+ Description: Means `KEEPNULLS`. Specifies that empty columns should retain a null value during the bulk-import operation, instead of having any default values for the columns inserted. 
+ Required: No
+ Default Value: None

##### KilobytesPerBatch

+ Description: Means `KILOBYTES_PER_BATCH`. Specifies the approximate number of kilobytes (KB) of data per batch as *kilobytes_per_batch*. By default, `KILOBYTES_PER_BATCH` is unknown.
+ Required: No
+ Default Value: None

##### RowsPerBatch

+ Description: Means `ROWS_PER_BATCH `. Indicates the approximate number of rows of data in the data file.By default, all the data in the data file is sent to the server as a single transaction, and the number of rows in the batch is unknown to the query optimizer. If you specify `ROWS_PER_BATCH` (with a value > 0) the server uses this value to optimize the bulk-import operation. The value specified for `ROWS_PER_BATCH` should approximately the same as the actual number of rows.
+ Required: No
+ Default Value: None

##### Order

+ Description: Means `ORDER`.Specifies how the data in the data file is sorted. Bulk import performance is improved if the data being imported is sorted according to the clustered index on the table, if any. If the data file is sorted in a different order, that is other than the order of a clustered index key or if there's no clustered index on the table, the `ORDER` clause is ignored. The column names supplied must be valid column names in the destination table. By default, the bulk insert operation assumes the data file is unordered. For optimized bulk import, SQL Server also validates that the imported data is sorted.
+ Required: No
+ Default Value: None

##### Tablock

+ Description: Means `TABLOCK `.Specifies that a table-level lock is acquired for the duration of the bulk-import operation. A table can be loaded concurrently by multiple clients if the table has no indexes and TABLOCK is specified. By default, locking behavior is determined by the table option **table lock on bulk load**. Holding a lock for the duration of the bulk-import operation reduces lock contention on the table, in some cases can significantly improve performance. 
+ Required: No
+ Default Value: None

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
