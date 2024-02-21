# MysqlWriter Plugin Documentation

## Quick Introduction

The MysqlWriter plugin enables writing data to Postgres/Greenplum databases. Under the hood, MysqlWriter connects to remote Mysql databases using github.com/go-sql-driver/mysql and database/sql, executing corresponding SQL statements to write data into the Mysql database.

## Implementation Principles

MysqlWriter connects to remote Mysql databases using github.com/go-sql-driver/mysql and generates SQL write statements based on user-configured information and go-etl's custom data types from the Reader. These statements are then sent to the remote Mysql database for execution.

MysqlWriter implements specific queries by invoking go-etl's custom DBWrapper from storage/database, using the query process defined in dbmswriter. DBWrapper encapsulates many interfaces from database/sql and abstracts the database dialect, Dialect. For Mysql, the Dialect implemented by storage/database/mysql is used.

Based on the configured `writeMode`, MysqlWriter generates either an `insert into...` statement (which will not insert conflicting rows in case of primary key/unique index conflicts) or a `replace into...` statement (which behaves like `insert into` when no conflicts occur, but replaces the entire row with new values when conflicts arise). Data is buffered in memory and written in batches to optimize performance.

## Functionality Description

### Configuration Example

Configuring a job to write data from memory to a Mysql database:

```json
{
    "job":{
        "content":[
            {
               "writer":{
                    "name": "mysqlwriter",
                    "parameter": {
                        "username": "root",
                        "password": "123456",
                        "writeMode": "insert",
                        "column": ["*"],
                        "preSql": ["create table a like b"],
                        "postSql": ["drop table a"],
                        "connection":  {
                                "url": "tcp(192.168.0.1:3306)/mysql?parseTime=false",
                                "table": {
                                    "db":"destination",
                                    "name":"type_table"
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

### Parameter Explanation

#### url

- Description: Used to configure the connection information for the remote end. The basic format is: tcp(ip:port)/db, where ip:port represents the IP address and port of the Mysql database, and db indicates the default database to connect to. This is similar to the connection configuration for [mysql](https://github.com/go-sql-driver/mysql), except that the username and password are extracted for easier encryption.
- Required: Yes
- Default: None

#### username

- Description: Used to configure the username for the Mysql database.
- Required: Yes
- Default: None

#### password

- Description: Used to configure the password for the Mysql database.
- Required: Yes
- Default: None

#### table

Describes the Mysql table information.

##### db

- Description: Used to configure the database name for the Mysql table.
- Required: Yes
- Default: None

##### name

- Description: Used to configure the table name for the Mysql table.
- Required: Yes
- Default: None

#### writeMode

- Description: Specifies the write mode. "insert" represents writing data using the "insert into" method, while "replace" represents writing data using the "replace into" method.
- Required: No
- Default: insert

#### column

- Description: Specifies the set of column names that need to be synchronized in the configured table. JSON array format is used to describe the column information. Using "*" represents including all columns by default, e.g., ["*"]. Column pruning is supported, meaning only selected columns can be inserted. Column reordering is also supported, meaning columns can be inserted in any order, not necessarily following the table schema.
- Required: Yes
- Default: None

#### batchTimeout

- Description: Configures the timeout interval for each batch write operation. The format is: number + unit, where the unit can be s for seconds, ms for milliseconds, or us for microseconds. If the specified time interval is exceeded, the data will be written immediately. This parameter, along with batchSize, can be adjusted to optimize write performance.
- Required: No
- Default: 1s

#### batchSize

- Description: Configures the size of each batch write operation. If the specified size is exceeded, the data will be written immediately. This parameter, along with batchTimeout, can be adjusted to optimize write performance.
- Required: No
- Default: 1000

#### preSql

- Description: Specifies a set of SQL statements to be executed before writing data. Do not use select statements as they will cause errors.
- Required: No
- Default: None

#### postSql

- Description: Specifies a set of SQL statements to be executed after writing data. Do not use select statements as they will cause errors.
- Required: No
- Default: None

### Type Conversion

Currently, MysqlWriter supports most Mysql data types, but there may be some unsupported types. Please check your data types carefully.

Below is a conversion table for MysqlWriter and Mysql data types:

| go-etl Type | Mysql Data Type                                        |
| ----------  | --------------------------------------------------- |
| bigInt      | int, tinyint, smallint, mediumint, bigint, year, unsigned int, unsigned bigint, unsigned smallint, unsigned tinyint       |
| decimal     | float, double, decimal                                 |
| string      | varchar, char, tinytext, text, mediumtext, longtext   |
| time        | date, datetime, timestamp, time                        |
| bytes       | tinyblob, mediumblob, blob, longblob, varbinary, bit  |

## Performance Report

To be tested.

## Constraints and Limitations

### Database Encoding Issues

Currently, only the utf8 character set is supported.

## FAQ
