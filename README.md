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
| Relational Database | MySQL/Mariadb/Tidb | √            | √          | [Read](datax/plugin/reader/mysql/README.md)、[Write](datax/plugin/writer/mysql/README.md) |
|              | Postgres/Greenplum | √            | √          | [Read](datax/plugin/reader/postgres/README.md)、[Write](datax/plugin/writer/postgres/README.md) |
|              | DB2 LUW            | √            | √          | [Read](datax/plugin/reader/db2/README.md)、[Write](datax/plugin/writer/db2/README.md) |
|              | SQL Server            | √            | √          | [Read](datax/plugin/reader/sqlserver/README.md)、[Write](datax/plugin/writer/sqlserver/README.md) |
|              | Oracle            | √            | √          | [Read](datax/plugin/reader/oracle/README.md)、[Write](datax/plugin/writer/oracle/README.md) |
|              | Sqlite3            | √            | √          | [Read](datax/plugin/reader/sqlite3/README.md)、[Write](datax/plugin/writer/sqlite3/README.md) |
| Unstructured Data Stream    | CSV                | √            | √          | [Read](datax/plugin/reader/csv/README.md)、[Write](datax/plugin/writer/csv/README.md) |
|              | XLSX（excel）      | √            | √          | [Read](datax/plugin/reader/xlsx/README.md)、[Write](datax/plugin/writer/xlsx/README.md) |

### Getting Started

#### Run by testing performance

If you want to directly obtain performance-related information, you can deploy according to the [Prometheus Monitoring Deployment Manual](docker/README.md) to acquire relevant performance metrics and even performance visualization charts.

#### Run by obtaining the binary program

You can download the 64-bit binary executable for Windows or Linux operating systems from the 
[latest releases](https://github.com/Breeze0806/go-etl/releases)

Start data synchronization with the [go-etl Data Synchronization User Manual](README_USER.md)

#### Run by obtaining the Docker image

**Pull Docker Image**
```bash
docker pull go-etl:v0.2.2
```

**Start Container**
```bash
docker run -d -p 6080:6080 --name etl -v /data:/usr/local/go-etl/data go-etl:v0.2.2
```

**Execute Command in Container**
```bash
docker exec -it etl bash
```

Important Note:  
Current version doesn't support direct usage of SQLite3 or Oracle databases. To enable these databases, you need to:  
1. Download the corresponding ODBC drivers  
2. Configure environment variables for database connections  

#### Run by source code

##### Linux

###### Compilation dependencies

1. golang 1.20 and later versions
2. gcc 4.8 and later versions

###### build

```bash
cd ${GO_PATH}/src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
make dependencies
make release
```

###### Removing DB2 dependency

Before compilation, it is necessary to use `export IGNORE_PACKAGES=db2`

```bash
cd ${GO_PATH}/src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
export IGNORE_PACKAGES=db2
make dependencies
make release
```

##### windows

###### Compilation Dependencies:

1. Mingw-w64 with gcc 7.2.0 or higher is required for compilation.
2. Golang version 1.20 or above is necessary.
3. The minimum compilation environment is Windows 7.

###### build

```bash
cd ${GO_PATH}\src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
release.bat
```

###### Removing DB2 dependency

Before compilation, it is necessary to use `set IGNORE_PACKAGES=db2`

```bash
cd ${GO_PATH}\src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
set IGNORE_PACKAGES=db2
release.bat
```

#### Compilation output

```
    +---go-etl---|---plugin---+---reader--mysql---|--README.md
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
+ README_USER is the user manual or guide in English.

#### Run by docker

To retrieve the `go-etl` project (version v0.2.1), follow these steps:

```bash
# Clone the repository
git clone https://github.com/Breeze0806/go-etl.git
cd go-etl

# Verify the tag version
git describe --abbrev=0 --tags
```

For Docker image compilation:

```bash
# Build Docker image
docker build . -t go-etl:v0.2.2
```

To start the container:

```bash
# Run container in detached mode
docker run -d -p 6080:6080 --name etl -v /data:/usr/local/go-etl/data go-etl:v0.2.2
```

To access the container shell:

```bash
# Enter the running container
docker exec -it etl bash
```

### Data Synchronization Development Handbook

Refer to the [go-etl Data Synchronization Developer Documentation](datax/README.md) to assist with your development.

## Module Introduction
### datax

This package will provide an interface similar to Alibaba's [DataX](https://github.com/alibaba/DataX) to implement an offline data synchronization framework in the Go programming language. The framework will enable users to perform data synchronization tasks efficiently and reliably, leveraging the power and flexibility of the Go language. It may include features such as pluggable data sources and destinations, support for various data formats, and robust error handling mechanisms.

```
readerPlugin(reader)—> Framework(Exchanger+Transformer) ->writerPlugin(riter)  
```

The system is built using a Framework + plugin architecture. In this design, the reading and writing of data sources are abstracted into Reader/Writer plugins, which are integrated into the overall synchronization framework.

+ Reader: The Reader module is responsible for data acquisition. It collects data from the source and sends it to the Framework.
+ Writer: The Writer module handles data writing. It continuously retrieves data from the Framework and writes it to the destination.
+ Framework: The Framework serves as the connection between the Reader and Writer. It functions as a data transmission channel, handling core technical aspects such as buffering, flow control, concurrency, and data transformation.
This architecture allows for flexibility and scalability, as new data sources and destinations can be easily added by developing new Reader and Writer plugins, respectively.

For detailed information, please refer to the [go-etl Data Synchronization Developer Documentation](datax/README.md). This documentation provides guidance on how to use the go-etl framework for data synchronization, including information on its architecture, plugin system, and how to develop custom Reader and Writer plugins.

### element
Currently, the data types and data type conversions in go-etl have been implemented. For more information, please refer to the [go-etl Data Type Descriptions](element/README.md). This documentation provides details on the supported data types, their usage, and how to perform conversions between different types within the go-etl framework.

### storage

#### database

We have now implemented basic integration for databases, abstracting the database dialect (Dialect) interface. For specific implementation details, please refer to the [Database Storage Developer Guide](storage/database/README.md). This guide provides information on how to work with different database dialects within the framework, allowing for flexible and extensible database support.

#### Stream

Primarily used for parsing byte streams, such as files, message queues, Elasticsearch, etc. The byte stream format can be CSV, JSON, XML, etc.

##### File

Focused on file parsing, including CSV, Excel, etc. It abstracts the InputStream and OutputStream interfaces. For specific implementation details, refer to the [Developer Guide for Tabular File Storage](storage/stream/file/README.md).

### Tools

A collection of utilities for compilation, adding licenses, etc.

#### DataX

##### Build

```bash
go generate ./...
```
This is the build command used to register developer-created reader and writer plugins into the program's code. Additionally, this command inserts compilation information, such as software version, Git version, Go compiler version, and compilation time, into the command line tool.

##### Plugin

A plugin template creation tool for data sources. It's used to create a new reader or writer template, in conjunction with the build command, to reduce the developer's workload.

##### Release

A packaging tool for the data synchronization program and user documentation.

#### License

Automatically adds a license to Go code files and formats the code using `gofmt -s -w`.

```bash
go run tools/license/main.go
```

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