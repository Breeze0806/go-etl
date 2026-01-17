# PostgresReader Plugin Documentation

## Quick Introduction

The PostgresReader plugin enables data reading from Postgres/Greenplum databases. Under the hood, PostgresReader connects to remote Postgres/Greenplum databases using `github.com/lib/pq` and executes corresponding SQL statements to query data from the database.

## Implementation Principles

PostgresReader connects to remote Postgres/Greenplum databases using `github.com/lib/pq` and generates SQL queries based on user-provided configuration information. These queries are then sent to the remote Postgres/Greenplum database, and the returned results are assembled into an abstract dataset using go-etl's custom data types. This dataset is then passed to downstream Writer processing. Unlike directly using `github.com/lib/pq` to connect to the database, here we use `github.com/Breeze0806/go/database/pqto` to set read and write timeouts.

PostgresReader implements specific queries by calling go-etl's custom `storage/database` DBWrapper, which is defined in the dbmsreader's query process. DBWrapper encapsulates many interfaces of `database/sql` and abstracts the database dialect, Dialect. For Postgres, the implementation of Dialect provided by `storage/database/postgres` is used.

## Functionality Description

### Configuration Example

Configuring a job to synchronize data from a Postgres/Greenplum database to a local system:

```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "postgresreader",
                    "parameter": {
                        "username": "postgres",
                        "password": "123456",
                        "column": ["*"],
                        "connection":  {
                                "url": "postgres://192.168.15.130:5432/postgres?sslmode=disable",
                                "table": {
                                    "schema":"source",
                                    "name":"type_table"
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

- Description: Mainly used to configure the connection information for the remote database. The basic configuration format is: `postgres://ip:port/db`, where `ip:port` represents the IP address and port of the Postgres database, and `db` represents the default database to connect to. It is basically the same as the connection configuration information of [pq](https://pkg.go.dev/github.com/lib/pq), except that the username and password are extracted from the connection configuration information to facilitate subsequent encryption of these information. Unlike [pq](https://pkg.go.dev/github.com/lib/pq), you can use `readTimeout/writeTimeout` to configure read/write timeouts, with the same format as `batchTimeout`.
- Required: Yes
- Default: None

#### username

- Description: Mainly used to configure the Postgres database username.
- Required: Yes
- Default: None

#### password

- Description: Mainly used to configure the Postgres database password.
- Required: Yes
- Default: None

#### table

Describes the Postgres table information.

##### schema

- Description: Mainly used to configure the schema name of the Postgres table.
- Required: Yes
- Default: None

##### name

- Description: Mainly used to configure the table name of the Postgres table.
- Required: Yes
- Default: None

#### column

- Description: The set of column names that need to be synchronized from the configured table. JSON array syntax is used to describe the column information. Using "*" represents that all columns are used by default, for example, `["*"]`.

  Supports column pruning, which means users can select specific columns for export.

  Supports column reordering, meaning the columns can be exported in an order different from the table schema.

  Supports constant configuration. Users need to follow the PostgreSQL syntax format: `["id", "'hello'::varchar", "true", "2.5::real", "power(2,3)"]`. Here, "id" is a regular column name, `'hello'::varchar` is a string constant, "true" is a boolean value, "2.5" is a floating-point number, and `power(2,3)` is a function.

- Required: Yes
- Default: None

#### split

##### key

- Description: Mainly used to configure the splitting key for the Postgres table. The splitting key must be of type bigInt/string/time, assuming that the data is evenly distributed based on the splitting key.
- Required: No
- Default: None

##### timeAccuracy

- Description: Mainly used to configure the time splitting key for the Postgres table, mainly to describe the smallest unit of time, such as day (for dates), min (for minutes), s (for seconds), ms (for milliseconds), us (for microseconds), ns (for nanoseconds).
- Required: No
- Default: None

##### range

###### type
- Description: Mainly used to configure the default value type of the splitting key for the Postgres table, with values being bigInt/string/time. Here, it will check the type of the splitting key in the table, so please make sure the type is correct.
- Required: No
- Default: None

###### left
- Description: Mainly used to configure the default minimum value of the splitting key for the Postgres table.
- Required: No
- Default: None

###### right
- Description: Mainly used to configure the default maximum value of the splitting key for the Postgres table.
- Required: No
- Default: None

#### where

- Description: Mainly used to configure the where condition for the select statement.
- Required: No
- Default: None

#### querySql

- Description: In some business scenarios, the `where` configuration item is not sufficient to describe the filtering conditions, so users can use this configuration item to customize the filtering SQL. When users configure this item, the DataX system will ignore the `table`, `column`, and other configuration items, and directly use the content of this configuration item for data filtering. For example, if you need to perform a join operation on multiple tables before synchronizing the data, you can use `select a,b from table_a join table_b on table_a.id = table_b.id`.
When the user configures `querySql`, PostgresReader directly ignores the configuration of `table`, `column`, and `where` conditions. The priority of `querySql` is higher than that of `table`, `column`, and `where` options.
- Required: No
- Default: None

#### trimChar

- Description: Whether to remove leading and trailing spaces for the char type in Postgres.
- Required: No
- Default: false

### Type Conversion

Currently, PostgresReader supports most Postgres types, but there are still some individual types that are not supported. Please check your types carefully.

Below is a list of type conversions that PostgresReader performs for Postgres types:

| go-etl Type | Postgres Data Type                                         |
| ----------- | -------------------------------------------------------- |
| bool        | boolean                                                   |
| bigInt      | bigint, bigserial, integer, smallint, serial, smallserial |
| decimal     | double precision, decimal, numeric, real                 |
| string      | varchar, text, uuid                                   |
| time        | date, time, timestamp                                    |
| bytes       | char                                                     |

## Performance Report

To be tested.

## Constraints and Limitations

### Database Encoding Issues
Currently, only the utf8 character set is supported.

## FAQ
