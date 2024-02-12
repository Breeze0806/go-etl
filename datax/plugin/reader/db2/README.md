# DB2Reader Plugin Documentation

## Quick Introduction

The DB2Reader plugin enables data extraction from DB2 LUW (Linux, Unix, and Windows) databases. Under the hood, it utilizes the `github.com/ibmdb/go_ibm_db` package along with `database/sql` to connect to remote DB2 LUW databases, execute SQL queries, and retrieve data. Unlike other databases, DB2 does not have a publicly available interaction protocol, so the GoLang driver for DB2 leverages the DB2 ODBC library for database connectivity.

## Implementation Principles

DB2Reader establishes a connection to the remote DB2 LUW database using the ODBC library via `github.com/ibmdb/go_ibm_db`. It generates SQL queries based on user-provided configuration information, sends them to the remote DB2 LUW database, and receives the results. These results are then packaged into an abstract dataset using go-etl's custom data types and passed downstream to the Writer for processing.

DB2Reader accomplishes the specific querying by invoking the query process defined in `dbmsreader` and utilizing go-etl's custom `storage/database` DBWrapper. The DBWrapper encapsulates many interfaces from `database/sql` and abstracts the database dialect. In the case of DB2, it implements the dialect defined in `storage/database/db2`.

## Functionality Description

### Configuration Example

Configuring a job to synchronize data from a DB2 database to a local system:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "db2reader",
                    "parameter": {
                        "connection":  {
                            "url": "HOSTNAME=127.0.0.1;PORT=50000;DATABASE=db",
                            "table": {
                                "schema":"SOURCE",
                                "name":"TEST"
                            }
                        },
                        "username": "user",
                        "password": "passwd",
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

- Description: Specifies the connection information for the remote DB2 database. The basic format is: `HOSTNAME=ip;PORT=port;DATABASE=db`, where `ip` represents the IP address and `port` represents the port number of the DB2 database, and `db` represents the default database to connect to. This configuration is similar to the connection information used by [ibm db2](https://github.com/ibmdb/go_ibm_db) but separates the username and password for easier encryption.
- Required: Yes
- Default: None

#### username

- Description: The username for the DB2 database.
- Required: Yes
- Default: None

#### password

- Description: The password for the DB2 database.
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

#### split

##### key

- Description: The split key for the DB2 table. The split key must be of type bigInt, string, or time, assuming the data is evenly distributed based on the split key.
- Required: No
- Default: None

##### timeAccuracy

- Description: Specifies the time precision for the split key of the DB2 table. Valid values are day, min, s, ms, us, ns. This setting is required if a range is specified.
- Required: No
- Default: None

##### range

###### type
- Description: Specifies the data type of the default value for the split key of the DB2 table. Valid values are bigInt, string, time. Please ensure the correct type is selected.
- Required: No
- Default: None

###### left
- Description: Specifies the default maximum value for the split key of the DB2 table.
- Required: No
- Default: None

###### right
- Description: Specifies the default minimum value for the split key of the DB2 table.
- Required: No
- Default: None

#### column

- Description: An array of column names to be synchronized from the configured table. The user can specify "*" to select all columns by default, e.g., ["*"]. Column pruning (selecting only specific columns) and column reordering (not following the table schema order) are supported. Constant values can also be configured using DB2 SQL syntax, e.g., ["id", "`table`", "1", "'bazhen.csy'", "null", "left(a,10)", "2.3", "true"].
- Required: Yes
- Default: None

#### where

- Description: Specifies the WHERE condition for the SELECT query.
- Required: No
- Default: None

#### querySql

- Description: Allows the user to define a custom SQL query for data filtering. When this option is configured, the system ignores the table, column, and where settings, and directly uses the content of this configuration for data filtering. This is useful for scenarios that require joining multiple tables or complex filtering conditions.
- Required: No
- Default: None

#### trimChar

- Description: Specifies whether to trim leading and trailing spaces from char type columns in DB2.
- Required: No
- Default: false

### Type Conversions

DB2Reader supports most DB2 data types, but there may be some unsupported types. Please check your data types carefully.

The following table lists the type conversion mappings supported by DB2Reader:

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

## FAQ

1. How to configure the DB2 ODBC library?

   - For Linux, set the environment variable `LD_LIBRARY_PATH` to point to the DB2 ODBC library path, e.g., `export LD_LIBRARY_PATH=${DB2HOME}/lib`.
   - For Windows, update the system PATH to include the path to the DB2 ODBC library, e.g., `set path=%path%;%GOPATH%\src\github.com\ibmdb\go_ibm_db\clidriver\bin`.