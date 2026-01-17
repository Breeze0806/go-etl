# PostgresWriter Plugin Documentation

## Quick Introduction

The PostgresWriter plugin enables writing data to Postgres/Greenplum databases. Under the hood, PostgresWriter connects to remote Postgres/Greenplum databases using github.com/lib/pq and database/sql, executing corresponding SQL statements to write data into the Postgres/Greenplum database.

## Implementation Principles

PostgresWriter connects to remote Postgres/Greenplum databases via github.com/lib/pq. It generates SQL statements for writing based on user-provided configuration information and go-etl's custom data types from the Reader. These statements are then sent to the remote Postgres/Greenplum database for execution.

Postgres/Greenplum implements specific queries by utilizing the query process defined in dbmswriter, calling go-etl's custom storage/database DBWrapper. DBWrapper encapsulates numerous interfaces from database/sql and abstracts the database dialect, Dialect. For Postgres/Greenplum, it adopts the Dialect implemented in storage/database/postgres.

Based on your configured `writeMode`, it generates either:

- `insert into...` (which may fail to insert conflicting rows in case of primary key/unique index conflicts)

**or**

- `copy in ...` which behaves similarly to insert into but offers faster performance. For optimal performance, data is buffered in memory and written only when the memory reaches a predefined threshold.

**or**

- `insert into ... on conflict ... do update set ...` which allows you to handle conflicts when inserting data. If the inserted data violates a unique constraint (such as a primary key or unique index), you can choose to update the existing record instead of throwing an error.
 
## Functionality Description

### Configuration Example

Configuring a job to synchronously write data to a Postgres/Greenplum database:

```json
{
    "job":{
        "content":[
            {
               "writer":{
                    "name": "postgreswriter",
                    "parameter": {
                        "username": "postgres",
                        "password": "123456",
                        "writeMode": "insert",
                        "column": ["*"],
                        "session": [],
                        "preSql": [],
                        "connection":  {
                                "url": "postgres://192.168.15.130:5432/postgres?sslmode=disable",
                                "table": {
                                    "schema":"destination",
                                    "name":"type_table"
                                }
                         },
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

- Description: Primarily used to configure the connection information for the remote end. The basic format is: postgres://ip:port/db, where ip:port represents the IP address and port of the Postgres database, and db indicates the default database to connect to. It is similar to the connection configuration information of [pq](https://pkg.go.dev/github.com/lib/pq), except that the username and password are extracted from the connection configuration for easier encryption in the future. Unlike [pq](https://pkg.go.dev/github.com/lib/pq), read/write timeouts can be configured using readTimeout/writeTimeout in the same format as batchTimeout.
- Required: Yes
- Default: None

#### username

- Description: Used to configure the username for the Postgres database.
- Required: Yes
- Default: None

#### password

- Description: Used to configure the password for the Postgres database.
- Required: Yes
- Default: None

#### name

Describes the Postgres table information.

##### schema

- Description: Primarily used to configure the schema name of the Postgres table.
- Required: Yes
- Default: None

##### table

- Description: Primarily used to configure the table name of the Postgres table.
- Required: Yes
- Default: None

#### column

- Description: A set of column names from the configured table that need to be synchronized, described using a JSON array. Users can use * to indicate that all columns should be used by default, for example, ["*"].

  Supports column pruning, allowing only selected columns to be exported.

  Supports column reordering, meaning columns can be exported in an order different from the table schema.

  Supports constant configuration. Users need to follow the PostgreSQL syntax format: ["id", "'hello'::varchar", "true", "2.5::real", "power(2,3)"] where id is a regular column name, 'hello'::varchar is a string constant, true is a boolean value, 2.5 is a floating-point number, and power(2,3) is a function.

- Required: Yes
- Default: None

#### writeMode

- Description: Write mode. "insert" represents writing data using the `insert into` method, "copyIn" represents writing data using the `copy in` method, "upsert" represents writing data using the `insert into ... on conflict ... do update set ... `method.
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

#### upsertSql

- Description: Primarily used to configure the `on conflict ... do update set ...` statement in the `upsert` mode.
- Required: No
- Default: None

### Type Conversion

Currently, PostgresWriter supports most Postgres types, but there may be some individual types that are not supported. Please check your types accordingly.

Below is a conversion table for PostgresWriter with regards to Postgres types:

| go-etl Type | Postgres Data Type |
| --- | --- |
| bool | boolean |
| bigInt | bigint, bigserial, integer, smallint, serial, smallserial |
| decimal | double precision, decimal, numeric, real |
| string | varchar, text, uuid  |
| time | date, time, timestamp |
| bytes | char |

## Performance Report

Pending testing.

## Constraints and Limitations

### Database Encoding Issues

Currently, only the utf8 character set is supported.

## FAQ

1. upsert mode support postgres 9.6+ version， What PostgreSQL and Greenplum versions are supported for upsert mode?
   - PostgreSQL: The upsert functionality (`INSERT ... ON CONFLICT ... DO UPDATE SET`) is supported in PostgreSQL 9.5+, but go-etl specifically requires PostgreSQL 9.6+ for stable upsert operations due to improvements made in version 9.6.
   - Greenplum: Since Greenplum is based on PostgreSQL, upsert functionality is supported starting from Greenplum 7.x and later versions that are built on PostgreSQL 12.12.
