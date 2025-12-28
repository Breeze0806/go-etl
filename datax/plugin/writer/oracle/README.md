# OracleWriter Plugin Documentation

## Quick Introduction

The OracleWriter plugin enables data writing to an Oracle database. Under the hood, OracleReader connects to a remote Oracle database using `github.com/godror/godror` and `database/sql`. Unlike other databases, Oracle's interaction protocol is not publicly available. Therefore, the Golang driver for Oracle is based on [ODPI-C](https://oracle.github.io/odpi/doc/installation.html) and requires [Oracle Instant Client](https://www.oracle.com/database/technologies/instant-client/downloads.html) for the connection. For instance, connecting to Oracle 11g requires version 12.x of the client.

## Implementation Principles

OracleReader connects to a remote Oracle database using Oracle Instant Client via `github.com/godror/godror`. It generates SQL statements for writing based on user-configured information and go-etl's custom data types from Reader. These statements are then sent to the remote Oracle database for execution.

OracleReader implements specific queries by invoking go-etl's custom `storage/database` DBWrapper, defined in the query process of `dbmswriter`. DBWrapper encapsulates numerous interfaces of `database/sql` and abstracts the database dialect, `Dialect`. In this case, Oracle uses the `Dialect` implementation from `storage/database/oracle`.

Based on your configured `writeMode`, it generates:

* `insert into...` (if there's a primary key/unique index conflict, the conflicting row won't be inserted).

Note that the insert method here is not the usual `storage/database` insert implementation but a specific Oracle approach. In this implementation, the query might be `insert into a(x,y,x) values(:1,:2,:3)`, where the args for x, y, and z are arrays consisting of column values.

## Functionality Description

### Configuration Example

Configuring a job to write data from memory to an Oracle database:


```json
{
    "job":{
        "content":[
            {
                 "writer":{
                    "name": "oraclewriter",
                    "parameter": {
                        "connection":  {
                            "url": "connectString=\"192.168.15.130:1521/xe\" heterogeneousPool=false standaloneConnection=true",
                            "table": {
                                "schema":"TEST",
                                "name":"DEST"
                            }
                        },
                        "username": "system",
                        "password": "oracle",
                        "writeMode": "insert",
                        "column": ["*"],
                        "preSql": ["create table a like b"],
                        "postSql": ["drop table a"],
                        "batchTimeout": "1s",
                        "batchSize":1000
                    }
                },
            }
        ]
    }
}
```
### Parameter Description

#### url

* Description: Mainly used to configure the connection information for the target database. The basic format for Oracle database connections is: `connectString="192.168.15.130:1521/xe" heterogeneousPool=false standaloneConnection=true`. The `connectString` represents the connection information for the Oracle database. If using a server name for the connection, please use `ip:port/servername` or `(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=ip)(PORT=port))(CONNECT_DATA=(SERVICE_NAME=servername)))`. If using a SID for the connection, use `(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=ip)(PORT=port))(CONNECT_DATA=(SID=sid)))`. This is similar to the connection configuration information in the [Godror User Guide](https://godror.github.io/godror/doc/contents.html), except that the username and password are extracted from the connection configuration for easier encryption in the future.
* Required: Yes
* Default: None

#### username

* Description: Mainly used to configure the Oracle database username.
* Required: Yes
* Default: None

#### password

* Description: Mainly used to configure the Oracle database password.
* Required: Yes
* Default: None

#### table

Describes the Oracle table information.

##### schema

* Description: Mainly used to configure the schema name of the Oracle table.
* Required: Yes
* Default: None

##### name

* Description: Mainly used to configure the table name of the Oracle table.
* Required: Yes
* Default: None

#### writeMode

* Description: Write mode. `insert` represents writing data using the `insert into` method.
* Required: No
* Default: insert

#### column

* Description: An array of column names from the configured table that need to be synchronized, described using JSON array notation. Users can use `*` to represent all columns by default, e.g., `["*"]`. Supports column pruning, which means you can select specific columns for insertion. Supports column reordering, which means columns can be inserted in a different order from the table schema.
* Required: Yes
* Default: None

#### batchTimeout

* Description: Mainly used to configure the timeout interval for each batch write operation. Format: number + unit. Units: s for seconds, ms for milliseconds, us for microseconds. If the timeout interval is exceeded, the data will be written directly. This parameter, along with `batchSize`, helps adjust write performance.
* Required: No
* Default: 1s

#### batchSize

* Description: Mainly used to configure the batch write size. If the size is exceeded, the data will be written directly. This parameter, along with `batchTimeout`, helps adjust write performance.
* Required: No
* Default: 1000

#### preSql

* Description: Mainly used for SQL statement groups executed before writing data. Do not use select statements as they will cause errors.
* Required: No
* Default: None

#### postSql

* Description: Mainly used for SQL statement groups executed after writing data. Do not use select statements as they will cause errors.
* Required: No
* Default: None

### Type Conversion

Currently, OracleWriter supports most Oracle types, but there may be some individual types that are not supported. Please check your types carefully.

Below is a list of OracleWriter type conversions for Oracle types:



| go-etl Type | Oracle Data Type |
| --- | --- |
| bool | BOOLEAN |
| bigInt | NUMBER, INTEGER, SMALLINT |
| decimal | BINARY_FLOAT, FLOAT, BINARY_DOUBLE, REAL, DECIMAL, NUMERIC |
| string | VARCHAR, CHAR, NCHAR, VARCHAR2, NVARCHAR2, CLOB, NCLOB |
| time | DATE, TIMESTAMP |
| bytes | BLOB, RAW, LONG RAW, LONG |

## Performance Report

To be tested.

## Constraints and Limitations

### Database Encoding Issues

Currently, only the UTF-8 character set is supported.

## FAQ

1. How to configure Oracle's Oracle Instant Client?

Here's an example:

* On Linux, set the environment variable like this: `export LD_LIBRARY_PATH=/opt/oracle/instantclient_21_1:$LD_LIBRARY_PATH`. Also, note that you need to install `libaio`.
* Note that on Windows, you can use the command: set path=%path%;%GOPATH%\oracle\instantclient_21_1.Oracle Instant Client 19 no longer supports Windows 7. In addition, you need to install [Oracle Instant Client and the corresponding Visual Studio Redistributable](https://odpi-c.readthedocs.io/en/latest/user_guide/installation.html#windows).


2. How to eliminate `godor WARNING: discrepancy between SESSIONTIMEZONE and SYSTIMESTAMP`

Either speak with your DBA to synchronize the DB's time zone (DBTIMEZONE) with the underlying OS' time zone, or use

```sql
ALTER SESSION SET TIME_ZONE='Europe/Berlin'
```


or set one chosen timezone in the [./connection.md](connection string):

```ini
timezone="Europe/Berlin"
```


(it is parsed with time.LoadLocation, so such names can be used, or local, or a numeric +0500 fixed zone).

WARNING: time zone altered with ALTER SESSION may not be read each and every time, so either always ALTER SESSION consistently to the same timezone, or use the

```ini
perSessionTimezone=1
```


connection parameter, to force checking the time zone for each session (and not cache it per DB).

3. Why is the writing of time formats so slow?
Currently, the OracleWriter converts time formats by using the to_date or to_timestamp functions instead of the bind variable method, which results in relatively slow writing of time formats.