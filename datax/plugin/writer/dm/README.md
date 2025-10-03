# DMWriter Plugin Documentation

## Quick Introduction

The DMWriter plugin enables writing data to DM (Dameng) databases. Under the hood, DMWriter connects to remote DM databases using gitee.com/chunanyong/dm and database/sql, executing corresponding SQL statements to write data into the DM database.

## Implementation Principles

DMWriter connects to remote DM databases using gitee.com/chunanyong/dm and generates SQL write statements based on user-configured information and go-etl's custom data types from the Reader. These statements are then sent to the remote DM database for execution.

DMWriter implements specific queries by invoking go-etl's custom DBWrapper from storage/database, using the query process defined in dbmswriter. DBWrapper encapsulates many interfaces from database/sql and abstracts the database dialect, Dialect. For DM, the Dialect implemented by storage/database/dm is used.

Based on the configured `writeMode`, DMWriter generates:
an `insert into...` statement (which will not insert conflicting rows in case of primary key/unique index conflicts)

## Functionality Description

### Configuration Example

Configuring a job to write data from memory to a DM database:

```json
{
  "job":{
    "content":[
      {
        "writer":{
          "name": "dmwriter",
          "parameter": {
            "username": "",
            "password": "",
            "column": ["*"],
            "preSql": [],
            "writeMode": "insert",
            "connection":  {
              "url": "",
              "table": {
                "db":"",
                "name":""
              }
            },
            "batchTimeout": "1s",
            "batchSize":1000
          }
        },
        "transformer": []
      }
    ]
  }
}
```

### Parameter Explanation

#### url

- Description: Used to configure the connection information for the remote end. The basic format is: ip:port, where ip:port represents the IP address and port of the DM database.
- Required: Yes
- Default: None

#### username

- Description: Used to configure the username for the DM database.
- Required: Yes
- Default: None

#### password

- Description: Used to configure the password for the DM database.
- Required: Yes
- Default: None

#### table

Describes the DM table information.

##### db

- Description: Used to configure the database name for the DM table.
- Required: Yes
- Default: None

##### name

- Description: Used to configure the table name for the DM table.
- Required: Yes
- Default: None

#### writeMode

- Description: Specifies the write mode. "insert" represents writing data using the "insert into" method.
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

Currently, DMWriter supports most DM data types, but there may be some unsupported types. Please check your data types carefully.
| go-etl Type | Dameng Data Type |
| --- | --- |
| bool  |bit |
| bigInt | BIGINT, INT, INTEGER, SMALLINT,TINYINT,BYTE |
| decimal | DOUBLE，DOUBLE PRECISION, REAL, DECIMAL，NUMERIC,NUMBER  |
| string | VARCHAR, CHAR, TEXT, LONG,CLOB，LONGVARCHAR   |
| time | DATE, TIME, DATETIME,TIMESTAMP |
| bytes | BLOB, VARBINARY, IMAGE, LONGVARBINARY  |
## Performance Report

To be tested.

## Constraints and Limitations

### Database Encoding Issues

Currently, only the utf8 character set is supported.

## FAQ