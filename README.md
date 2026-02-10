# go-etl
[![LICENSE][license-img]][license]
[![Language][lang-img]][lang]
[![Build][ci-img]][ci]
[![Go Report Card][report-img]][report]
[![GitHub release][release-img]][release]
[![GitHub release date][release-date-img]][release-date]
[![Coverage Status][cov-img]][cov]
[![GoDoc][doc-img]][doc]
[![Mentioned in Awesome Go][awesome-img]][awesome]
[![Docker Version][docker-img]][docker]

English | [简体中文](README_zh-CN.md)

go-etl is a toolset for extracting, transforming, and loading data sources, providing powerful data synchronization capabilities.

go-etl will provide the following ETL capabilities:

- The ability to extract and load data from mainstream databases is implemented in the storage package
- The ability to extract and load data from data streams in a two-dimensional table-like format is implemented in the stream package
- Similar data synchronization capabilities to datax, implemented in the datax package

Since I have limited energy, everyone is welcome to submit issues to discuss go-etl, let's make progress together!


## Data Synchronization Tool

This data synchronization tool has the synchronization capability for the following data sources.

| Type         | Data Source             | Reader | Writer| Document                                                        |
| ------------ | ------------------ | ------------ | ---------- | ------------------------------------------------------------ |
| Relational Database | MySQL/Mariadb/Tidb/TDSQL MySQL | √            | √          | [Read](datax/plugin/reader/mysql/README.md)、[Write](datax/plugin/writer/mysql/README.md) |
|              | Postgres/Greenplum | √            | √          | [Read](datax/plugin/reader/postgres/README.md)、[Write](datax/plugin/writer/postgres/README.md) |
|              | DB2 LUW            | √            | √          | [Read](datax/plugin/reader/db2/README.md)、[Write](datax/plugin/writer/db2/README.md) |
|              | SQL Server            | √            | √          | [Read](datax/plugin/reader/sqlserver/README.md)、[Write](datax/plugin/writer/sqlserver/README.md) |
|              | Oracle            | √            | √          | [Read](datax/plugin/reader/oracle/README.md)、[Write](datax/plugin/writer/oracle/README.md) |
|              | Sqlite3            | √            | √          | [Read](datax/plugin/reader/sqlite3/README.md)、[Write](datax/plugin/writer/sqlite3/README.md) |
|              | Dameng            | √            | √          | [Read](datax/plugin/reader/dm/README.md)、[Write](datax/plugin/writer/dm/README.md) |
| Unstructured Data Stream    | CSV                | √            | √          | [Read](datax/plugin/reader/csv/README.md)、[Write](datax/plugin/writer/csv/README.md) |
|              | XLSX（excel）      | √            | √          | [Read](datax/plugin/reader/xlsx/README.md)、[Write](datax/plugin/writer/xlsx/README.md) |

### Getting Started

#### Quick Start (3 Minutes)

##### Start from Binary Program

You can download the 64-bit binary executable for Windows or Linux operating systems from the [latest releases](https://github.com/Breeze0806/go-etl/releases).

Start data synchronization with the [go-etl Data Synchronization User Manual](README_USER.md).

##### Start from Docker Image

**Pull Docker Image**
```bash
docker pull go-etl:v0.2.3
```

**Start Container**
```bash
docker run -d -p 6080:6080 --name etl -v /data:/usr/local/go-etl/data go-etl:v0.2.3
```

**Enter Container**
```bash
docker exec -it etl bash
```

**Execute Sync in Container**
```bash
docker exec -it etl release/bin/go-etl -c data/config.json
```

#### Start from Performance Testing

If you want to directly obtain performance-related information, you can deploy according to the [Prometheus Monitoring Deployment Manual](docker/README.md) to acquire relevant performance metrics and even performance visualization charts.

#### Start from Source Code

##### Linux

###### Compilation Dependencies

1. golang 1.20 and later versions
2. gcc 4.8 and later versions

###### Build

```bash
cd ${GO_PATH}/src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
make dependencies
make release
```

###### Removing DB2 Dependency

Before compilation, it is necessary to use `export IGNORE_PACKAGES=db2` 

```bash
export IGNORE_PACKAGES=db2
cd ${GO_PATH}/src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
make dependencies
make release
```

##### Windows

###### Compilation Dependencies

1. A MinGW-w64 environment with GCC 7.2.0 or higher is required for compilation.
2. golang 1.20 and later versions
3. The minimum compilation environment is Windows 7.

###### Build

```bash
cd ${GO_PATH}\src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
release.bat
```

###### Removing DB2 Dependency

Before compilation, it is necessary to use `set IGNORE_PACKAGES=db2`

```bash
cd ${GO_PATH}\src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
set IGNORE_PACKAGES=db2
release.bat
```

##### Compilation Output

```
    +---datax---|---plugin---+---reader--mysql---|--README.md
    |                        | .......
    |                        |
    |                        |---writer--mysql---|--README.md
    |                        | .......
    |
    +---bin----go-etl
    +---exampales---+---csvpostgres----config.json
    |               |---db2------------config.json
    |               | .......
    |
    +---README_USER.md
```

+ The datax/plugin directory contains the documentation for various plugins.
+ The bin directory houses the data synchronization program, named go-etl.
+ The examples directory includes configuration files for data synchronization in different scenarios.
+ README_USER.md is the user manual or guide in English.

#### Start from Compiled Docker Image

Use the following commands to get the `go-etl` project (version `v0.2.3`):

```bash
git clone https://github.com/Breeze0806/go-etl.git
cd go-etl
git describe --abbrev=0 --tags
```

Build the Docker image with the following command:

```bash
docker build . -t go-etl:v0.2.3
```

Start the container:

```bash
docker run -d -p 6080:6080 --name etl -v /data:/usr/local/go-etl/data go-etl:v0.2.3
```

Enter the container:

```bash
docker exec -it etl bash
```

Note that currently, sqlite3, DB2, and Oracle are not directly supported and require downloading the corresponding ODBC and configuring environment variables.

#### Batch Sync

Use a wizard CSV file to batch sync multiple tables.

**1. Create data source config `config.json`** - same as single sync

**2. Create wizard file `wizard.csv`** - each row defines a source-target table pair:
```csv
source_table,target_table
table1,table1_copy
table2,table2_copy
```

**3. Generate batch configs and run script:**

**Linux:**
```bash
./go-etl -c config.json -w wizard.csv; ./run.sh
```

**Windows:**
```powershell
.\go-etl.exe -c config.json -w wizard.csv; run.bat
```

**Docker:**
```bash
docker exec -it etl release/bin/go-etl -c data/config.json -w data/wizard.csv; docker exec -it etl bash run.sh
```

### Data Synchronization Development Guide

Refer to the [go-etl Data Synchronization Developer Documentation](datax/README.md) to assist with your development.

## Module Introduction

### datax

This package provides an interface similar to Alibaba's [DataX](https://github.com/alibaba/DataX) to implement an offline data synchronization framework in Go.

```
readerPlugin(reader)—> Framework(Exchanger+Transformer) ->writerPlugin(writer)  
```

Built using a Framework + plugin architecture. Data source reading and writing are abstracted into Reader/Writer plugins and integrated into the overall synchronization framework.

+ Reader: The Reader is the data acquisition module, responsible for collecting data from the data source and sending it to the Framework.
+ Writer: The Writer is the data writing module, responsible for continuously fetching data from the Framework and writing it to the destination.
+ Framework: The Framework connects the reader and writer, serving as a data transmission channel, and handles core technical aspects such as buffering, flow control, concurrency, and data transformation.

For detailed information, please refer to the [go-etl Data Synchronization Developer Documentation](datax/README.md).

### element

The data types and data type conversions in go-etl have been implemented. For more information, please refer to the [go-etl Data Type Descriptions](element/README.md).

### storage

#### database

Basic integration for databases has been implemented, abstracting the database dialect (Dialect) interface. For specific implementation details, please refer to the [Database Storage Developer Guide](storage/database/README.md).

#### stream

Primarily used for parsing byte streams, such as files, message queues, Elasticsearch, etc. The byte stream format can be CSV, JSON, XML, etc.

##### file

Focused on file parsing, including CSV, Excel, etc. It abstracts the InputStream and OutputStream interfaces. For specific implementation details, refer to the [Developer Guide for Tabular File Storage](storage/stream/file/README.md).

### tools

A collection of utilities for compilation, adding licenses, etc.

#### datax

##### build

```bash
go generate ./...
```

Release command used to register developer-created reader and writer plugins into the program's code.

Additionally, this command inserts compilation information such as software version, Git version, Go compiler version, and compilation time into the command line.

##### plugin

A plugin template creation tool for data sources. It's used to create a new reader or writer template, in conjunction with the release command, to reduce the developer's workload.

##### release

A packaging tool for the data synchronization program and user documentation.

#### license

Automatically adds a license to Go code files and formats the code using `gofmt -s -w`.

```bash
go run tools/license/main.go
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on how to contribute to this project.

### Ways to Contribute

- Report bugs and issues
- Suggest new features
- Submit pull requests
- Improve documentation
- Share your use cases

### Getting Help

- Check the [Documentation](README.md)
- Check the [User Manual](README_USER.md)
- Submit a GitHub Issue for discussion

[lang-img]:https://img.shields.io/badge/Language-Go-blue.svg
[lang]:https://golang.org/
[report-img]:https://goreportcard.com/badge/github.com/Breeze0806/go-etl
[report]:https://goreportcard.com/report/github.com/Breeze0806/go-etl
[doc-img]:https://godoc.org/github.com/Breeze0806/go-etl?status.svg
[doc]:https://godoc.org/github.com/Breeze0806/go-etl
[license-img]: https://img.shields.io/badge/License-Apache%202.0-blue.svg
[license]: https://github.com/Breeze0806/go-etl/blob/main/LICENSE
[ci-img]: https://github.com/Breeze0806/go-etl/actions/workflows/Build.yml/badge.svg
[ci]: https://github.com/Breeze0806/go-etl/actions/workflows/Build.yml
[release-img]: https://img.shields.io/github/tag/Breeze0806/go-etl.svg?label=release
[release]: https://github.com/Breeze0806/go-etl/releases
[release-date-img]: https://img.shields.io/github/release-date/Breeze0806/go-etl.svg
[release-date]: https://github.com/Breeze0806/go-etl/releases
[cov-img]: https://codecov.io/gh/Breeze0806/go-etl/branch/main/graph/badge.svg?token=UGb27Nysga
[cov]: https://codecov.io/gh/Breeze0806/go-etl
[awesome-img]:https://awesome.re/mentioned-badge.svg
[awesome]:https://github.com/avelino/awesome-go
[docker-img]:https://img.shields.io/docker/v/breeze0806/go-etl?sort=semver&logo=docker&logoColor=white&label=Docker&color=blue
[docker]:https://hub.docker.com/r/breeze0806/go-etl