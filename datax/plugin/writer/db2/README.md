# DB2Writer Plugin Documentation

## Quick Introduction

The DB2Writer plugin enables data writing to DB2 LUW (Linux, Unix, and Windows) databases. Under the hood, DB2Writer utilizes `github.com/ibmdb/go_ibm_db` along with `database/sql` to connect to remote DB2 LUW databases and executes corresponding SQL statements to write data into the DB2 database. Unlike other databases, DB2 does not publicly expose its interaction protocol, so the Golang driver for DB2 leverages DB2's ODBC library for database connectivity.

## Implementation Principles

DB2Writer connects to remote DB2 LUW databases using the ODBC library through `github.com/ibmdb/go_ibm_db`. It generates write SQL statements based on user-configured information and go-etl's custom data types from the Reader. These statements are then sent to the remote DB2 database for execution.

DB2Writer implements specific queries by invoking go-etl's custom `DBWrapper` from the query process defined in `dbmswriter`. The `DBWrapper` encapsulates many interfaces from `database/sql` and abstracts the database dialect, known as `Dialect`. In this case, DB2 utilizes the `Dialect` implementation from `storage/database/db2`.

Based on the configured `writeMode`, DB2Writer can generate:

- `insert into...` (which may fail to insert conflicting rows if there are primary key or uniqueness index violations)

## Functional Description

### Configuration Example

Configuring a job to write data from memory to a DB2 database:

```json
{
    "job":{
        "content":[
            {
               "writer":{
                    "name": "db2writer",
                    "parameter": {
                        "connection":  {
                            "url": "HOSTNAME=127.0.0.1;PORT=50000;DATABASE=db",
                            "table": {
                                "schema":"SOURCE",
                                "name":"TEST"
                            }
                        },
                        "username": "root",
                        "password": "12345678",
                        "writeMode": "insert",
                        "column": ["*"],
                        "preSql": ["create table a like b"],
                        "postSql": ["drop table a"],
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

- Description: Primarily used to configure the connection information for the remote DB2 database. The basic format is: `HOSTNAME=ip;PORT=port;DATABASE=db`, where `ip` represents the IP address and `port` of the DB2 database, and `db` denotes the default database to connect to. This configuration is similar to that of [ibm db2](https://github.com/ibmdb/go_ibm_db), except that the username and password are extracted from the connection configuration for easier encryption in the future.
- Required: Yes
- Default: None

#### username

- Description: Used to specify the DB2 database username.
- Required: Yes
- Default: None

#### password

- Description: Used to specify the DB2 database password.
- Required: Yes
- Default: None

#### table

Describes the DB2 table information.

##### schema

- Description: Specifies the schema name of the DB2 table.
- Required: Yes
- Default: None

##### name

- Description: Specifies the table name of the DB2 table.
- Required: Yes
- Default: None

#### writeMode

- Description: Write mode. `insert` represents writing data using the `insert into` method.
- Required: No
- Default: `insert`

#### column

- Description: An array of column names from the configured table that need to be synchronized. Using `["*"]` selects all columns by default. Supports column pruning (selecting only a subset of columns for insertion) and column reordering (inserting columns in an order different from the table schema).
- Required: Yes
- Default: None

#### batchTimeout

- Description: Specifies the timeout interval for each batch write operation. The format is `number+unit`, where the unit can be `s` for seconds, `ms` for milliseconds, or `us` for microseconds. If the specified time interval elapses, the data is written immediately. This parameter, along with `batchSize`, helps adjust write performance.
- Required: No
- Default: `1s`

#### batchSize

- Description: Specifies the size of each batch write operation. If the specified size is reached, the data is written immediately. This parameter, along with `batchTimeout`, helps adjust write performance.
- Required: No
- Default: 1000

#### preSql

- Description: An array of SQL statements to be executed before writing data. Avoid using `select` statements as they may cause errors.
- Required: No
- Default: None

#### postSql

- Description: An array of SQL statements to be executed after writing data. Avoid using `select` statements as they may cause errors.
- Required: No
- Default: None

### Type Conversion

Currently, DB2Reader supports most DB2 data types, but there may be some unsupported individual types. Please check your data types carefully.

The following table lists the type conversion between go-etl types and DB2 data types supported by DB2Reader:

| go-etl Type | DB2 Data Type |
| --- | --- |
| bool | BOOLEAN |
| bigInt | BIGINT, INTEGER, SMALLINT |
| decimal | DOUBLE, REAL, DECIMAL |
| string | VARCHAR, CHAR |
| time | DATE, TIME, TIMESTAMP |
| bytes | BLOB, CLOB |

## Performance Report

Pending testing.

## Constraints and Limitations

### Database Encoding Issues

Currently, only the UTF-8 character set is supported.

## Frequently Asked Questions (FAQ)

1. How to configure the DB2 ODBC library?

   - On Linux, set the environment variable `LD_LIBRARY_PATH` to include the path to the DB2 ODBC library, as shown in the Makefile: `export LD_LIBRARY_PATH=${DB2HOME}/lib`.
   - On Windows, update the system path to include the path to the DB2 ODBC library, as shown in `release.bat`: `set path=%path%;%GOPATH%\src\github.com\ibmdb\go_ibm_db\clidriver\bin`.