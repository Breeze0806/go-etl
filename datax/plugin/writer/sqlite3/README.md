# Sqlite3Writer Plugin Documentation

## Quick Introduction

The Sqlite3Writer plugin enables writing data to Sqlite3 databases. Under the hood, Sqlite3Writer connects to remote Sqlite3 databases using github.com/mattn/go-sqlite3, executing corresponding SQL statements to write data into the Sqlite3 database.

## Implementation Principles

Sqlite3Writer connects to remote Sqlite3 databases via github.com/mattn/go-sqlite3. It generates SQL statements for writing based on user-provided configuration information and go-etl's custom data types from the Reader. These statements are then sent to the remote Sqlite3 database for execution.

Sqlite3 implements specific queries by utilizing the query process defined in dbmswriter, calling go-etl's custom storage/database DBWrapper. DBWrapper encapsulates numerous interfaces from database/sql and abstracts the database dialect, Dialect. For Sqlite3, it adopts the Dialect implemented in storage/database/postgres.

Based on your configured `writeMode`, it generates either:

- `insert into...` (which may fail to insert conflicting rows in case of primary key/unique index conflicts)



## Functionality Description

### Configuration Example

Configuring a job to synchronously write data to a Sqlite3 database:

```json
{
  "content": [
    {
      "writer": {
        "name": "sqlite3writer",
        "parameter": {
          "writeMode": "insert",
          "column": [
            "*"
          ],
          "connection": {
            "url": "E:\\Sqlite3\\test.db",
            "table": {
              "db": "main",
              "name": "type_table_copy"
            }
          },
          "preSql": ["create table a like b"],
          "postSql": ["drop table a"],
          "batchTimeout": "1s",
          "batchSize": 1000
        }
      }
    }
  ]
}
```

### Parameter Description

#### url

- Description: Primarily used to configure the connection information for the remote end. 
- Required: Yes
- Default: None

#### name

Describes the Sqlite3 table information.

##### table

- Description: Primarily used to configure the table name of the Sqlite3 table.
- Required: Yes
- Default: None

#### column

- Description: A set of column names from the configured table that need to be synchronized, described using a JSON array. Users can use * to indicate that all columns should be used by default, for example, ["*"].

  Supports column pruning, allowing only selected columns to be exported.

  Supports column reordering, meaning columns can be exported in an order different from the table schema.

  Supports constant configuration. Users need to follow the PostgreSQL syntax format.

- Required: Yes
- Default: None

#### writeMode

- Description: Write mode. "insert" represents writing data using the insert into method, while "copyIn" represents writing data using the copy in method.
- Required: No
- Default: insert

#### batchTimeout

- Description: Primarily used to configure the timeout interval for each batch write operation. The format is: number + unit, where the unit can be s for seconds, ms for milliseconds, or us for microseconds. If the specified time interval is exceeded, the data will be written directly. This parameter, along with batchSize, can be adjusted for optimal write performance.
- Required: No
- Default: 1s

#### batchSize

- Description: Primarily used to configure the size of each batch write operation. If the specified size is exceeded, the data will be written directly. This parameter, along with batchTimeout, can be adjusted for optimal write performance.
- Required: No
- Default: 1000

#### preSql

- Description: Primarily used for SQL statement groups executed before writing data. Do not use select statements as they will result in an error.
- Required: No
- Default: None

#### postSql

- Description: Primarily used for SQL statement groups executed after writing data. Do not use select statements as they will result in an error.
- Required: No
- Default: None

### Type Conversion

Currently, Sqlite3Writer supports most Sqlite3 types, but there may be some individual types that are not supported. Please check your types accordingly.

Below is a conversion table for Sqlite3Writer with regards to Sqlite3 types:

| go-etl的类型 | sqlite3数据类型                                         |
| ------------ | -------------------------------------------------------- |
| string         |  INTEGER、TEXT、REAL、BLOB |

## Performance Report

Pending testing.

## Constraints and Limitations

### Database Encoding Issues

Currently, only the utf8 character set is supported.

## FAQ

(Frequently Asked Questions section to be added if applicable.)