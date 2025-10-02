# DMReader Plugin Documentation

## Quick Introduction

The DMReader plugin enables data reading from a DM (Dameng) database. Internally, DMReader connects to a remote DM database using `gitee.com/chunanyong/dm` and `database/sql`, executing corresponding SQL statements to retrieve data from the DM server.

## Implementation Principles

DMReader connects to a remote DM database using `gitee.com/chunanyong/dm`. Based on user-provided configuration information, it generates SQL queries and sends them to the remote DM server. The returned results from these SQL executions are assembled into an abstract dataset using go-etl's custom data types and passed to downstream Writers for processing.

DMReader utilizes the query processes defined in `dbmsreader` and calls go-etl's custom `storage/database` DBWrapper for specific queries. DBWrapper encapsulates numerous interfaces from `database/sql` and abstracts out a database dialect. For DM, it adopts the dialect implemented in `storage/database/dm`.

## Functionality Description

### Configuration Example

Configuring a job to synchronize data from a DM database to a local destination:

```json
{
  "job":{
    "content":[
      {
        "reader":{
          "name": "dmreader",
          "parameter": {
            "username": "",
            "password": "",
            "column": [],
            "connection": {
              "url": "",
              "table": {
                "db": "",
                "name": ""
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

- Description: Primarily used to configure the connection information to the remote server. The basic configuration format is: `ip:port`, where `ip:port` represents the IP address and port of the DM server.
- Required: Yes
- Default: None

#### username

- Description: Primarily used to configure the DM database username.
- Required: Yes
- Default: None

#### password

- Description: Primarily used to configure the DM database password.
- Required: Yes
- Default: None

#### table

Describes the DM table information.

##### db

- Description: Primarily used to configure the database name of the DM table.
- Required: Yes
- Default: None

##### name

- Description: Primarily used to configure the table name of the DM table.
- Required: Yes
- Default: None

#### column

- Description: An array of column names from the configured table that need to be synchronized. Users can use the JSON array format to describe the field information. Using "*" represents selecting all columns by default, e.g., `["*"]`.

  Supports column pruning, meaning users can choose specific columns for export.

  Supports column reordering, allowing columns to be exported in an order different from the table schema.

- Required: Yes
- Default: None

#### where

- Description: Primarily used to configure the WHERE condition for the SELECT statement.
- Required: No
- Default: None

### Type Conversion

Currently, DMReader supports most DM types, but there are still some individual types that are not supported. Please check your types carefully.

## Performance Report

To be tested.

## Constraints and Limitations

### Database Encoding Issues

Currently, only the utf8 character set is supported.