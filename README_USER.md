# go-etl Data Synchronization User Manual

go-etl is a data synchronization tool that currently supports data synchronization between mainstream relational databases such as MySQL, Postgres, Oracle, SQL Server, DB2, Sqlite3, and file formats like CSV and XLSX.

## 1 How to Obtain

Refer to the [project documentation](README.md) for instructions on obtaining the binary program, starting from source code compilation, or building from Docker images.

## 2 Getting Started

Acquire the go-etl binary program or Docker image to begin usage. On Linux, as shown in the Makefile, export LD_LIBRARY_PATH=/home/ibmdb/clidriver/lib. This library can be downloaded from [ibm db2](https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli), otherwise it will not run.

Additionally, for Oracle, you need to download the corresponding 64-bit version of the ODBC dependency from [Oracle](https://www.oracle.com/database/technologies/instant-client/downloads.html).

### 2.1 Single-Task Data Synchronization

Invoking go-etl is straightforward; you simply call it directly.

```bash
./go-etl -c config.json
```

- The `-c` flag specifies the data source configuration file.

When the return value is `0` and the message `run success` is displayed at the end, it indicates that the execution was successful.

When the return value is `1` and the message `run fail` is displayed at the end, along with the reason for the failure, it indicates that the execution failed.

#### 2.1.1 Data Source Configuration File

The data source configuration file is a JSON file that combines data sources. For example, to synchronize data from MySQL to Postgres:

```json
{
    "core" : {
        "container": {
            "job":{
                "id": 1,
                "sleepInterval":100
            }
        }
    },
    "job":{
        "content":[
            {
                "reader":{
                    "name": "mysqlreader",
                    "parameter": {
                        "username": "test:",
                        "password": "test:",
                        "column": ["*"],
                        "connection":  {
                                "url": "tcp(192.168.15.130:3306)/source?parseTime=false",
                                "table": {
                                    "db":"source",
                                    "name":"type_table"
                                }
                            },
                        "where": ""
                    }
                },
                "writer":{
                    "name": "postgreswriter",
                    "parameter": {
                        "username": "postgres",
                        "password": "123456",
                        "writeMode": "insert",
                        "column": ["*"],
                        "preSql": [],
                        "connection":  {
                                "url": "postgres://192.168.15.130:5432/postgres?sslmode=disable&connect_timeout=2",
                                "table": {
                                    "schema":"destination",
                                    "name":"type_table"
                                }
                         },
                        "batchTimeout": "1s",
                        "batchSize":1000
                    }
                },
               "transformer":[]
            }
        ]
    }
}
```
The configurations for `reader` and `writer` are as follows:

| Type                | Data Source        | Reader (Read) | Writer (Write) | Documentation                                                |
| ------------------- | ------------------ | ------------- | -------------- | ------------------------------------------------------------ |
| Relational Database | MySQL/Mariadb/Tidb | √             | √              | [Read](datax/plugin/reader/mysql/README.md), [Write](datax/plugin/writer/mysql/README.md) |
|                     | Postgres/Greenplum | √             | √              | [Read](datax/plugin/reader/postgres/README.md), [Write](datax/plugin/writer/postgres/README.md) |
|                     | DB2 LUW            | √             | √              | [Read](datax/plugin/reader/db2/README.md), [Write](datax/plugin/writer/db2/README.md) |
|                     | SQL Server         | √             | √              | [Read](datax/plugin/reader/sqlserver/README.md), [Write](datax/plugin/writer/sqlserver/README.md) |
|                     | Oracle             | √             | √              | [Read](datax/plugin/reader/oracle/README.md), [Write](datax/plugin/writer/oracle/README.md) |
|              | Sqlite3            | √            | √          | [Read](datax/plugin/reader/sqlite3/README.md)、[Write](datax/plugin/writer/sqlite3/README.md) |
| Unstructured Stream | CSV                | √             | √              | [Read](datax/plugin/reader/csv/README.md), [Write](datax/plugin/writer/csv/README.md) |
|                     | XLSX (excel)       | √             | √              | [Read](datax/plugin/reader/xlsx/README.md), [Write](datax/plugin/writer/xlsx/README.md) |

#### 2.1.2 Usage Examples

Note: On Linux, as indicated in the Makefile, export `LD_LIBRARY_PATH=${DB2HOME}/lib`

##### 2.1.2.1 Synchronizing with MySQL

* Initialize the database using `cmd/datax/examples/mysql/init.sql` **for testing purposes**
* Start the MySQL synchronization command:

```bash
./go-etl -c examples/mysql/config.json
```

##### 2.1.2.2 Synchronizing with PostgreSQL

* Initialize the database using `cmd/datax/examples/postgres/init.sql` **for testing purposes**
* Start the PostgreSQL synchronization command:

```bash
./go-etl -c examples/postgres/config.json
```

##### 2.1.2.3 Synchronizing with DB2

* Before use, download the corresponding DB2 ODBC library, e.g., `make dependencies` and `release.bat` for Linux
* Note: On Linux, as indicated in the Makefile, export `LD_LIBRARY_PATH=${DB2HOME}/lib`
* Note: On Windows, as indicated in `release.bat`, set `path=%path%;%GOPATH%\src\github.com\ibmdb\go_ibm_db\clidriver\bin`
* Initialize the database using `cmd/datax/examples/db2/init.sql` **for testing purposes**
* Start the synchronization command:

```bash
./go-etl -c examples/db2/config.json
```

##### 2.1.2.4 Synchronizing with Oracle

* Before use, download the corresponding [Oracle Instant Client](https://www.oracle.com/database/technologies/instant-client/downloads.html). For example, it is recommended to download version 12.x to connect to Oracle 11g.
* Note: On Linux, export `LD_LIBRARY_PATH=/opt/oracle/instantclient_21_1:$LD_LIBRARY_PATH`. Additionally, install `libaio`.
* Note: On Windows, set `path=%path%;%GOPATH%\oracle\instantclient_21_1`. Oracle Instant Client 19 no longer supports Windows 7.
* Initialize the database using `cmd/datax/examples/oracle/init.sql` **for testing purposes**
* Start the synchronization command:

```bash
./go-etl -c examples/oracle/config.json
```

##### 2.1.2.5 Synchronizing with SQL Server

* Initialize the database using `cmd/datax/examples/sqlserver/init.sql` **for testing purposes**
* Start the SQL Server synchronization command:

```bash
./go-etl -c examples/sqlserver/config.json
```

##### 2.1.2.6 Synchronizing CSV to PostgreSQL

* Initialize the database using `cmd/datax/examples/csvpostgres/init.sql` **for testing purposes**
* Start the synchronization command:

```bash
./go-etl -c examples/csvpostgres/config.json
```

##### 2.1.2.7 Synchronizing XLSX to PostgreSQL

* Initialize the database using `cmd/datax/examples/csvpostgres/init.sql` **for testing purposes** (Note: The path may need correction as it seems inconsistent with other examples)
* Start the synchronization command:

```bash
./go-etl -c examples/xlsxpostgres/config.json
```

##### 2.1.2.8 Synchronizing PostgreSQL to CSV

* Initialize the database using `cmd/datax/examples/csvpostgres/init.sql` **for testing purposes**
* Start the synchronization command:

```bash
./go-etl -c examples/postgrescsv/config.json
```

##### 2.1.2.9 Synchronizing PostgreSQL to XLSX

* Initialize the database using `cmd/datax/examples/csvpostgres/init.sql` **for testing purposes** (Note: The initialization script may not be specific to XLSX synchronization)
* Start the synchronization command:

```bash
./go-etl -c examples/postgresxlsx/config.json
```

##### 2.1.2.10 Synchronizing with sqlite3

* Before use, download the corresponding [SQLite Download Page](https://www.sqlite.org/download.html). 
* Note: On Windows, set `path=%path%;/opt/sqlite/sqlite3.dll`. 
* Initialize the database using `cmd/datax/examples/sqlite3/init.sql` **for testing purposes**
* In `examples/sqlite3/config.json`, `url` is the path of sqlite3 database files. On Windows, it can be `E:\sqlite3\test.db`, meanwhile, on Linux, it can be `/sqlite3/test.db`,
* Start the sqlite3 synchronization command:

```bash
./go-etl -c examples/sqlite3/config.json
```

##### 2.1.2.11 Other Synchronization Examples

In addition to the above examples, all data sources listed in the go-etl features can be used interchangeably. Configurations can be set up for data sources such as MySQL to PostgreSQL, MySQL to Oracle, Oracle to DB2, etc.

#### 2.1.3 Global Database Configuration

```json
{
    "job":{
        "setting":{
            "pool":{
              "maxOpenConns":8,
              "maxIdleConns":8,
              "connMaxIdleTime":"40m",
              "connMaxLifetime":"40m"
            },
            "retry":{
              "type":"ntimes",
              "strategy":{
                "n":3,
                "wait":"1s"
              },
              "ignoreOneByOneError":true
            }
        }
    }
}
```

##### 2.1.3.1 Connection Pool (pool)

+ `maxOpenConns`: Maximum number of open connections
+ `maxIdleConns`: Maximum number of idle connections
+ `connMaxIdleTime`: Maximum idle time for a connection
+ `connMaxLifetime`: Maximum lifetime of a connection

##### 2.1.3.2 Retry (retry)

`ignoreOneByOneError`: Whether to ignore individual retry errors

+ Retry types and strategies:
	1. Type `ntimes`: Retry a fixed number of times. Strategy: `"strategy":{"n":3,"wait":"1s"}`, where `n` is the number of retries and `wait` is the waiting time between retries.
	2. Type `forever`: Retry indefinitely. Strategy: `"strategy":{"wait":"1s"}`, where `wait` is the waiting time between retries.
	3. Type `exponential`: Exponential backoff retry. Strategy: `"strategy":{"init":"100ms","max":"4s"}`, where `init` is the initial waiting time and `max` is the maximum waiting time.

##### 2.1.3.3 Test Data

```bash
./go-etl -c examples/global/config.json
```

#### 2.1.4 Using Split Keys

It is assumed that data is evenly distributed based on the split key. Proper use of such a split key can make synchronization faster. Additionally, to speed up queries for maximum and minimum values, preset maximum and minimum values can be used for large tables.

##### 2.1.4.1 Testing Approach

* Generate MySQL data using a program and create `split.csv`

```bash
cd cmd/datax/examples/split
go run main.go
```

* Use `init.sql` to create the table
* Synchronize to the MySQL database

```bash
cd ../..
./go-etl -c examples/split/csv.json
```

* Modify `examples/split/config.json` to set the split key as `id,dt,str`
* Synchronize MySQL data with integer, date, and string types using the split key

```bash
./go-etl -c examples/split/config.json
```

#### 2.1.5 Using preSql and postSql

`preSql` and `postSql` are sets of SQL statements executed before and after writing data, respectively.

##### 2.1.5.1 Testing Approach

In this example, a full import is used:
1. Before writing data, a temporary table is created.
2. After writing data, the original table is deleted, and the temporary table is renamed to the new table.

```bash
./go-etl -c examples/prePostSql/config.json
```

#### 2.1.6 Flow Control Configuration

Previously, the `byte` and `record` configurations for speed did not take effect. Now, with the introduction of flow control, `byte` and `record` will be effective. `byte` limits the size of cached messages in bytes, while `record` limits the number of cached messages. If `byte` is set too low, it can cause the cache to be too small, resulting in failed data synchronization. When `byte` is 0 or negative, the limiter will not work. For example, `byte` can be set to 10485760, which is equivalent to 10Mb (10x1024x1024).

```json
{
    "job":{
        "setting":{
            "speed":{
                "byte":10485760,
                "record":1024,
                "channel":4
            }
        }
    }
}
```

##### 2.1.6.1 Flow Control Testing

* Generate `src.csv` using a program and initiate the flow control test

```bash
cd cmd/datax/examples/limit
go run main.go
cd ../..
./go-etl -c examples/limit/config.json
```

#### 2.1.7 querySql Configuration

The database reader uses `querySql` to query the database.2.2 多任务数据同步

- #### 2.2.1 Usage

  ##### 2.2.1.1 Data Source Configuration File

  Configure the data source configuration file, such as syncing from MySQL to PostgreSQL:

  ```json
  {
      "core" : {
          "container": {
              "job":{
                  "id": 1,
                  "sleepInterval":100
              }
          }
      },
      "job":{
          "content":[
              {
                  "reader":{
                      "name": "mysqlreader",
                      "parameter": {
                          "username": "test:",
                          "password": "test:",
                          "column": ["*"],
                          "connection":  {
                                  "url": "tcp(192.168.15.130:3306)/source?parseTime=false",
                                  "table": {
                                      "db":"source",
                                      "name":"type_table"
                                  }
                              },
                          "where": ""
                      }
                  },
                  "writer":{
                      "name": "postgreswriter",
                      "parameter": {
                          "username": "postgres",
                          "password": "123456",
                          "writeMode": "insert",
                          "column": ["*"],
                          "preSql": [],
                          "connection":  {
                                  "url": "postgres://192.168.15.130:5432/postgres?sslmode=disable&connect_timeout=2",
                                  "table": {
                                      "schema":"destination",
                                      "name":"type_table"
                                  }
                           },
                          "batchTimeout": "1s",
                          "batchSize":1000
                      }
                  },
                 "transformer":[]
              }
          ]
      }
  }
  ```

  ##### 2.2.1.2 Source-Destination Configuration Wizard File

  The source-destination configuration wizard file is a CSV file. Each row of configuration can be set as follows:

  ```csv
  path[table],path[table]
  ```

  Each column can be a path or table name. Note that all tables should have a configured database name or schema name, which needs to be configured in the data source configuration file.

  ##### 2.2.1.3 Batch Generation of Data Configuration Sets and Execution Scripts

  ```bash
  ./go-etl -c tools/testData/xlsx.json -w tools/testData/wizard.csv 
  ```

  -c specifies the data source configuration file, and -w specifies the source-destination configuration wizard file.

  The execution result will generate a set of configuration files in the data source configuration file directory, with the number of rows in the source-destination configuration wizard file. The configuration sets will be named as specified_data_source_config_file1.json, specified_data_source_config_file2.json, ..., specified_data_source_config_file[n].json.

  Additionally, an execution script named run.bat or run.sh will be generated in the current directory.

  ##### 2.2.1.4 Batch Execution of Generated Data Configuration Sets

  ###### Windows

  ```bash
  run.bat
  ```

  Linux

  ```bash
  run.sh
  ```

  #### 2.2.2 Test Results

  You can run the test data in cmd/datax/testData:

  ```bash
  cd cmd/datax
  ./go-etl -c testData/xlsx.json -w testData/wizard.csv 
  ```

  The result will generate a set of configuration files in the testData directory, with the number of rows in the wizard.csv file. The configuration sets will be named as xlsx1.json, xlsx2.json, ..., xlsx[n].json.

  ### 2.3 Data Synchronization Help Manual

  #### 2.3.1 Help Command

  ```
 ./go-etl -h
  ```

  Help display:

  ```bash
  Usage of go-etl:
    -c string
          config (default "config.json")
    -http string
          http
    -w string
          wizard
  ```

  -http adds a listening port, such as 6080. After enabling, access 127.0.0.1:6080/metrics to get real-time throughput.

  #### 2.3.2 View Version

  ```bash
  ./go-etl version
  ```

  Display: `version number` (git commit: `git commit number`) compiled by go version `go version number`

  ```bash
  v0.1.0 (git commit: c82eb302218f38cd3851df4b425256e93f85160d) compiled by go version go1.16.5 windows/amd64
  ```

  #### 2.3.3 Start Monitoring Port

  ```bash
./go-etl -http :6080 -c examples/limit/config.json
  ```

  ##### 2.3.3.1 Retrieve Current Monitoring Data

  Access `http://127.0.0.1:6080/metrics` using a browser to obtain monitoring data similar to that provided by [Prometheus exporters](https://prometheus.io/docs/instrumenting/writing_exporters/).

  ```bash
  # HELP datax_channel_byte the number of bytes currently being synchronized in the channel
  # TYPE datax_channel_byte gauge
  datax_channel_byte{job_id="1",task_group_id="0",task_id="0"} 20480
  datax_channel_byte{job_id="1",task_group_id="0",task_id="1"} 20500
  # HELP datax_channel_record the number of records currently being synchronized in the channel
  # TYPE datax_channel_record gauge
  datax_channel_record{job_id="1",task_group_id="0",task_id="0"} 1024
  datax_channel_record{job_id="1",task_group_id="0",task_id="1"} 1025
  # HELP datax_channel_total_byte the total number of bytes synchronized
  # TYPE datax_channel_total_byte counter
  datax_channel_total_byte{job_id="1",task_group_id="0",task_id="0"} 2.75355e+06
  datax_channel_total_byte{job_id="1",task_group_id="0",task_id="1"} 5.29381e+06
  # HELP datax_channel_total_record the total number of records synchronized
  # TYPE datax_channel_total_record counter
  datax_channel_total_record{job_id="1",task_group_id="0",task_id="0"} 143233
  datax_channel_total_record{job_id="1",task_group_id="0",task_id="1"} 270246
  ```
  
  Additionally, you can access `http://127.0.0.1:6080/metrics?t=json` to obtain monitoring data in `JSON` format.
  
  ```json
  {
      "jobID": 1,
      "metrics": [
          {
              "taskGroupID": 0,
              "metrics": [
                  {
                      "taskID": 0,
                      "channel": {
                          "totalByte": 7069190,
                          "totalRecord": 359015,
                          "byte": 20500,
                          "record": 1025
                      }
                  },
                  {
                      "taskID": 1,
                      "channel": {
                          "totalByte": 13245910,
                          "totalRecord": 667851,
                          "byte": 20460,
                          "record": 1023
                      }
                  }
              ]
          }
      ]
  }
  ```
  
  - `totalByte` corresponds to `datax_channel_total_byte`, representing the total number of bytes synchronized.
  - `totalRecord` corresponds to `datax_channel_total_record`, representing the total number of records synchronized.
  - `byte` corresponds to `datax_channel_byte`, representing the number of bytes currently being synchronized in the channel.
  - `record` corresponds to `datax_channel_record`, representing the number of records currently being synchronized in the channel.
