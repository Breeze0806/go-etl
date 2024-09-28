# Sqlite3Reader Plugin Documentation

## Quick Introduction

The Sqlite3Reader plugin enables data reading from Sqlite3 databases. Under the hood, Sqlite3Reader connects to remote Sqlite3 databases using `github.com/mattn/go-sqlite3` and executes corresponding SQL statements to query data from the database.

## Implementation Principles

Sqlite3Reader connects to remote Sqlite3 databases using `github.com/mattn/go-sqlite3` and generates SQL queries based on user-provided configuration information. These queries are then sent to the remote Sqlite3 database, and the returned results are assembled into an abstract dataset using go-etl's custom data types. This dataset is then passed to downstream Writer processing.
Sqlite3Reader implements specific queries by calling go-etl's custom `storage/database` DBWrapper, which is defined in the dbmsreader's query process. DBWrapper encapsulates many interfaces of `database/sql` and abstracts the database dialect, Dialect. For sqlite3, the implementation of Dialect provided by `storage/database/sqlite3` is used.

## Functionality Description

### Configuration Example

Configuring a job to synchronize data from a Sqlite3 database to a local system:

```json
{
  "job": {
    "content": [
      {
        "reader": {
          "name": "sqlite3reader",
          "parameter": {
            "column": [
              "*"
            ],
            "connection": {
              "url": "E:\\Sqlite3\\test.db",
              "table": {
                "db": "main",
                "name": "type_table"
              }
            },
            "where": ""
          }
        }
      }
    ]
  }
}
```

### Parameter Explanation

#### url

- Description: It is mainly used to configure the path of sqlite3 database files
- Required: Yes
- Default: None

#### table

Describes the sqlite3 table information.

##### name

- Description: Mainly used to configure the table name of the sqlite3 table.
- Required: Yes
- Default: None

#### column

- Description: The set of column names that need to be synchronized from the configured table. JSON array syntax is used to describe the column information. Using "*" represents that all columns are used by default, for example, `["*"]`.

  Supports column pruning, which means users can select specific columns for export.

  Supports column reordering, meaning the columns can be exported in an order different from the table schema.

  Supports constant configuration. Users need to follow the sqlite3 syntax format.

- Required: Yes
- Default: None

#### split

##### key

- Description: Mainly used to configure the splitting key for the sqlite3 table. The splitting key must be of type bigInt/string/time, assuming that the data is evenly distributed based on the splitting key.
- Required: No
- Default: None

##### timeAccuracy

- Description: Mainly used to configure the time splitting key for the sqlite3 table, mainly to describe the smallest unit of time, such as day (for dates), min (for minutes), s (for seconds), ms (for milliseconds), us (for microseconds), ns (for nanoseconds).
- Required: No
- Default: None

##### range

###### type
- Description: Mainly used to configure the default value type of the splitting key for the sqlite3 table, with values being bigInt/string/time. Here, it will check the type of the splitting key in the table, so please make sure the type is correct.
- Required: No
- Default: None

###### left
- Description: Mainly used to configure the default maximum value of the splitting key for the sqlite3 table.
- Required: No
- Default: None

###### right
- Description: Mainly used to configure the default minimum value of the splitting key for the sqlite3 table.
- Required: No
- Default: None

#### where

- Description: Mainly used to configure the where condition for the select statement.
- Required: No
- Default: None

#### querySql

- Description: In some business scenarios, the `where` configuration item is not sufficient to describe the filtering conditions, so users can use this configuration item to customize the filtering SQL. When users configure this item, the DataX system will ignore the `table`, `column`, and other configuration items, and directly use the content of this configuration item for data filtering. For example, if you need to perform a join operation on multiple tables before synchronizing the data, you can use `select a,b from table_a join table_b on table_a.id = table_b.id`.
When the user configures `querySql`, Sqlite3Reader directly ignores the configuration of `table`, `column`, and `where` conditions. The priority of `querySql` is higher than that of `table`, `column`, and `where` options.
- Required: No
- Default: None

#### trimChar

- Description: Whether to remove leading and trailing spaces for the char type in sqlite3.
- Required: No
- Default: false

### Type Conversion

Currently, Sqlite3Reader supports most sqlite3 types, but there are still some individual types that are not supported. Please check your types carefully.

Below is a list of type conversions that Sqlite3Reader performs for sqlite3 types:

| go-etl的类型 | sqlite3数据类型        |
| ------------ |--------------------|
| string       | INTEGER、TEXT、REAL、BLOB |

## Performance Report

To be tested.

## Constraints and Limitations

### Database Encoding Issues
Currently, only the utf8 character set is supported.

## FAQ
