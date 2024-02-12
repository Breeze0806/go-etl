# OracleReader Plugin Documentation

## Quick Introduction

The OracleReader plugin enables data reading from Oracle databases. Under the hood, OracleReader connects to remote Oracle databases using github.com/godror/godror and database/sql, executing corresponding SQL statements to query data from Oracle. Unlike other databases, Oracle's interaction protocol is not publicly available, so the Golang driver for Oracle is based on [ODPI-C](https://oracle.github.io/odpi/doc/installation.html), which requires the use of [Oracle Instant Client](https://www.oracle.com/database/technologies/instant-client/downloads.html) for connectivity. For example, connecting to Oracle 11g requires version 12.x.

## Implementation Principles

OracleReader connects to remote Oracle databases using Oracle Instant Client via github.com/godror/godror. It generates SQL queries based on user-provided configurations, sends them to the remote Oracle database, and assembles the returned results into an abstract dataset using go-etl's custom data types, which are then passed to downstream Writer processes.

OracleReader implements specific queries by invoking the query process defined in dbmsreader, using go-etl's custom storage/database DBWrapper. DBWrapper encapsulates many interfaces of database/sql and abstracts the database dialect, Dialect. In this case, Oracle utilizes the Dialect implemented in storage/database/oracle.

## Functionality Description

### Configuration Example

Configuring a job to synchronize data from an Oracle database to a local destination:


```json
{
    "job":{
        "content":[
            {
                "reader":{
                    "name": "oraclereader",
                    "parameter": {
                        "connection":  {
                            "url": "connectString=\"192.168.15.130:1521/xe\" heterogeneousPool=false standaloneConnection=true",
                            "table": {
                                "schema":"TEST",
                                "name":"SRC"
                            }
                        },
                        "username": "system",
                        "password": "oracle",
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
### Parameter Explanation

#### url

* Description: Primarily used to configure the connection information for the remote Oracle database. The basic configuration format for connecting to an Oracle database is: `connectString="192.168.15.130:1521/xe" heterogeneousPool=false standaloneConnection=true`. The `connectString` represents the connection information for the Oracle database. If using a server name for the connection, use `ip:port/servername` or `(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=ip)(PORT=port))(CONNECT_DATA=(SERVICE_NAME=servername)))`. If using a SID for the connection, use `(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=ip)(PORT=port))(CONNECT_DATA=(SID=sid)))`. This is similar to the connection configuration information in the [Godror User Guide](https://godror.github.io/godror/doc/contents.html), except that the username and password are extracted from the connection configuration information for easier encryption.
* Required: Yes
* Default: None

#### username

* Description: Primarily used to configure the username for the Oracle database.
* Required: Yes
* Default: None

#### password

* Description: Primarily used to configure the password for the Oracle database.
* Required: Yes
* Default: None

#### table

Describes the Oracle table information.

##### schema

* Description: Primarily used to configure the schema name for the Oracle table.
* Required: Yes
* Default: None

##### name

* Description: Primarily used to configure the table name for the Oracle table.
* Required: Yes
* Default: None

#### column

* Description: An array of column names to be synchronized from the configured table. Users can use the asterisk (*) to indicate that all columns should be used by default, for example, ["*"]. Column pruning is supported, meaning that only a subset of columns can be exported. Column reordering is also supported, meaning that columns do not need to be exported in the same order as the table schema. Constant configuration is supported, where users need to follow the Oracle SQL syntax format: ["id", "`table`", "1", "'bazhen.csy'", "null", "left(a,10)", "2.3", "true"]. In this example, "id" is a regular column name, "`table`" is a column name that contains reserved words, "1" is an integer constant, "'bazhen.csy'" is a string constant, "null" is a null pointer, "left(a,10)" is an expression, "2.3" is a floating-point number, and "true" is a boolean value.
* Required: Yes
* Default: None

#### split

##### key

* Description: Primarily used to configure the splitting key for the Oracle table. The splitting key must be of type bigInt/string/time, assuming that the data is evenly distributed based on the splitting key.
* Required: No
* Default: None

##### timeAccuracy

* Description: Primarily used to configure the time splitting key for the Oracle table. It is mainly used to describe the smallest unit of time, such as day, minute, second, millisecond, microsecond, nanosecond.
* Required: No
* Default: None

##### range

###### type

* Description: Primarily used to configure the default type for the splitting key of the Oracle table. The value can be bigInt/string/time. The system does not check the type in the table splitting key, but it is important to ensure the correct type.
* Required: No
* Default: None

###### left

* Description: Primarily used to configure the default maximum value for the splitting key of the Oracle table.
* Required: No
* Default: None

###### right

* Description: Primarily used to configure the default minimum value for the splitting key of the Oracle table.
* Required: No
* Default: None

#### where

* Description: Primarily used to configure the WHERE condition for the SELECT statement.
* Required: No
* Default: None

#### querySql

* Description: In some business scenarios, the `where` configuration item may not be sufficient to describe the filtering conditions. Users can use this configuration item to define custom SQL queries for filtering. When users configure this item, the DataX system will ignore the `table`, `column`, and other related configurations and directly use the content of this configuration item for data filtering. For example, it can be used for data synchronization after performing a join operation on multiple tables, such as `select a,b from table_a join table_b on table_a.id = table_b.id`. When `querySql` is configured in OracleReader, it directly ignores the configuration of `table`, `column`, and `where` conditions, and the priority of `querySql` is higher than that of `table`, `column`, and `where` options.
* Required: No
* Default: None

#### trimChar

* Description: Specifies whether to remove leading and trailing spaces for Oracle's char and nchar types.
* Required: No
* Default: false

### Type Conversion

Currently, OracleReader supports most Oracle types, but there are some individual types that are not supported. Please check your data types carefully.

Below is a conversion table for OracleReader regarding Oracle types:



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

1. How to configure Oracle Instant Client for Oracle?

Example configurations:

* Note that on Linux, you may need to set an environment variable like `export LD_LIBRARY_PATH=/opt/oracle/instantclient_21_1:$LD_LIBRARY_PATH`. Additionally, you may need to install `libaio`.
* On Windows, you may need to set a path variable like `set path=%path%;%GOPATH%\oracle\instantclient_21_1`. Please note that Oracle Instant Client 19 no longer supports Windows 7.