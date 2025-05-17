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

[English](README.md) | 简体中文

go-etl是一个集数据源抽取，转化，加载的工具集，提供强大的数据同步能力。

go-etl将提供的etl能力如下：

1. 主流数据库的数据抽取以及数据加载的能力，在storage包中实现
2. 类二维表的数据流的数据抽取以及数据加载的能力，在stream包中实现
3. 类似datax的数据同步能力，在datax包中实现

鉴于本人实在精力有限，欢迎大家来提交issue或者***加QQ群185188648***来讨论go-etl，让我们一起进步!

## 数据同步工具

本数据数据同步工具以下数据源的同步能力

| 类型         | 数据源             | Reader（读） | Writer(写) | 文档                                                         |
| ------------ | ------------------ | ------------ | ---------- | ------------------------------------------------------------ |
| 关系型数据库 | MySQL/Mariadb/Tidb | √            | √          | [读](datax/plugin/reader/mysql/README_zh-CN.md)、[写](datax/plugin/writer/mysql/README_zh-CN.md) |
|              | Postgres/Greenplum | √            | √          | [读](datax/plugin/reader/postgres/README_zh-CN.md)、[写](datax/plugin/writer/postgres/README_zh-CN.md) |
|              | DB2 LUW            | √            | √          | [读](datax/plugin/reader/db2/README_zh-CN.md)、[写](datax/plugin/writer/db2/README_zh-CN.md) |
|              | SQL Server            | √            | √          | [读](datax/plugin/reader/sqlserver/README_zh-CN.md)、[写](datax/plugin/writer/sqlserver/README_zh-CN.md) |
|              | Oracle            | √            | √          | [读](datax/plugin/reader/oracle/README_zh-CN.md)、[写](datax/plugin/writer/oracle/README_zh-CN.md) |
|              | Sqlite3            | √            | √          | [读](datax/plugin/reader/sqlite3/README.md)、[写](datax/plugin/writer/sqlite3/README.md) |
| 无结构流     | CSV                | √            | √          | [读](datax/plugin/reader/csv/README_zh-CN.md)、[写](datax/plugin/writer/csv/README_zh-CN.md) |
|              | XLSX（excel）      | √            | √          | [读](datax/plugin/reader/xlsx/README_zh-CN.md)、[写](datax/plugin/writer/xlsx/README_zh-CN.md) |

### 快速开始

#### 从测试性能开始

如果你想直接想获取性能相关的信息，可以直接参考[prometheus监控部署手册](docker/README_zh-CN.md)去部署,来获取相关的性能参数，甚至性能图。

#### 从获取二进制程序开始

可以在[最新发布版本](https://github.com/Breeze0806/go-etl/releases)下载到windows或者linux操作系统的64位版本二进制程序。

使用方式参考[go-etl数据同步用户手册](README_USER_zh-CN.md)开始数据同步。

#### 从获取docker镜像开始

获取docker镜像

```bash
docker pull go-etl:v0.2.2
```

启动容器

```bash
docker run -d -p 6080:6080 --name etl -v /data:/usr/local/go-etl/data go-etl:v0.2.2
```

进入容器

```bash
docker exec -it etl bash
```
但当前还不直接支持sqlite3，oracle的直接使用，这需要下载对应的odbc并配置环境变量

#### 从源代码开始

##### linux

###### 编译依赖

1. golang 1.20以及以上版本

###### 构建

```bash
cd ${GO_PATH}\src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
make dependencies
make release
```

###### 去掉db2依赖

在编译前需要export IGNORE_PACKAGES=db2 

```bash
export IGNORE_PACKAGES=db2
cd ${GO_PATH}\src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
make dependencies
make release
```

##### windows

###### 编译依赖

1. 需要mingw-w64 with gcc 7.2.0以上的环境进行编译
2. golang 1.20以及以上
3. 最小编译环境为win7 

###### 构建

```bash
cd ${GO_PATH}\src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
release.bat
```

###### 去掉db2依赖

在编译前需要set IGNORE_PACKAGES=db2

```bash
cd ${GO_PATH}\src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
set IGNORE_PACKAGES=db2
release.bat
```

##### 编译产物

```
    +---go-etl---|---plugin---+---reader--mysql---|--README_zh-CN.md
    |                        | .......
    |                        |
    |                        |---writer--mysql---|--README_zh-CN.md
    |                        | .......
    |
    +---bin----go-etl
    +---exampales---+---csvpostgres----config.json
    |               |---db2------------config.json
    |               | .......
    |
    +---README_USER_zh-CN.md

```

+ go-etl/plugin下是各插件的文档
+ bin下的是数据同步程序go-etl
+ exampales下是各场景的数据同步的配置文档
+ README_USER_zh-CN.md是中文用户使用手册

#### 从编译docker镜像开始

使用命令获取`go-etl`项目,版本是`go-etl`的`v0.2.2`,
```bash
git clone https://github.com/Breeze0806/go-etl.git
cd go-etl
git describe --abbrev=0 --tags
```

通过以下命令可以编译成镜像
```bash
docker build . -t go-etl:v0.2.2
```

启动容器

```bash
docker run -d -p 6080:6080 --name etl -v /data:/usr/local/go-etl/data go-etl:v0.2.2
```

进入容器

```bash
docker exec -it etl bash
```

但当前还不直接支持sqlite3，DB2，oracle的直接使用，这需要下载对应的odbc并配置环境变量

### 数据同步开发宝典

参考[go-etl数据同步开发者文档](datax/README_zh-CN.md)来帮助开发

## 模块简介
### datax

本包将提供类似于阿里巴巴[DataX](https://github.com/alibaba/DataX)的接口去实现go语言的离线数据同步框架框架，

```
readerPlugin(reader)—> Framework(Exchanger+Transformer) ->writerPlugin(riter)  
```

采用Framework + plugin架构构建。将数据源读取和写入抽象成为Reader/Writer插件，纳入到整个同步框架中。

+ Reader：Reader为数据采集模块，负责采集数据源的数据，将数据发送给Framework。 
+ Writer：Writer为数据写入模块，负责不断向Framework取数据，并将数据写入到目的端。
+ Framework：Framework用于连接reader和writer，作为两者的数据传输通道，并处理缓冲，流控，并发，数据转换等核心技术问题

具体可以参考[go-etl数据同步开发者文档](datax/README_zh-CN.md)。

### element

目前已经实现了go-etl中的数据类型以及数据类型转换，可以参考[go-etl数据类型说明](element\README_zh-CN.md)。

### storage

#### database

目前已经实现了数据库的基础集成，抽象了数据库方言(Dialect)接口，具体实现可以参考[数据库存储开发者指南](storage/database/README_zh-CN.md)。

#### stream

主要用于字节流的解析，如文件，消息队列，elasticsearch等，字节流格式可以是cvs，json, xml等。

##### file
主要用于文件的解析，如cvs，excel等，抽象了输入流（InputStream）和输出流（OutputStream）接口，具体实现可以参考[类二维表文件存储开发者指南](storage/stream/file/README_zh-CN.md)。
### tools

工具集用于编译，新增许可证等

#### datax

##### build

```bash
go generate ./...
```
发布命令，用于将由开发者开发的reader和writer插件注册到程序中的代码

另外，该命令也会把编译信息如软件版本，git版本，go编译版本和编译时间写入命令行中

##### plugin

数据源插件模板新增工具，用于新增一个reader或writer模板，配合发布命令使用，减少开发者负担

##### release

数据同步程序和用户使用文档打包工具

#### license

用于自动新增go代码文件中许可证并使用gofmt -s -w格式化代码

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