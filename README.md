# go-etl
[![LICENSE][license-img]][license]
[![Language][lang-img]][lang]
[![Build][ci-img]][ci]
[![Go Report Card][report-img]][report]
[![GitHub release][release-img]][release]
[![GitHub release date][release-date-img]][release-date]
[![Coverage Status][cov-img]][cov]
[![GoDoc][doc-img]][doc]

go-etl是一个集数据源抽取，转化，加载，同步校验的工具集，提供强大的数据同步，数据转化甚至数据校验的功能。

go-etl将提供的etl能力如下：

1. 主流数据库的数据抽取以及数据加载的能力，在storage包中实现
2. 类二维表的数据流的数据抽取以及数据加载的能力，在stream包中实现
3. 类似datax的数据同步能力，在datax包中实现
4. 以mysql sql语法为基础的数据筛选、转化能力，在transform包中实现
5. 数据库间的数据校验能力，在libra包中实现 (目前计划中)

鉴于本人实在精力有限，欢迎大家来提交issue或者***加QQ群185188648***来讨论go-etl，让我们一起进步!

## 数据同步工具

本数据数据同步工具以下数据源的同步能力

| 类型         | 数据源             | Reader（读） | Writer(写) | 文档                                                         |
| ------------ | ------------------ | ------------ | ---------- | ------------------------------------------------------------ |
| 关系型数据库 | MySQL/Mariadb/Tidb | √            | √          | [读](datax/plugin/reader/mysql/README.md)、[写](datax/plugin/writer/mysql/README.md) |
|              | Postgres/Greenplum | √            | √          | [读](datax/plugin/reader/postgres/README.md)、[写](datax/plugin/writer/postgres/README.md) |
|              | DB2 LUW            | √            | √          | [读](datax/plugin/reader/db2/README.md)、[写](datax/plugin/writer/db2/README.md) |
|              | SQL Server            | √            | √          | [读](datax/plugin/reader/sqlserver/README.md)、[写](datax/plugin/writer/sqlserver/README.md) |
|              | Oracle            | √            | √          | [读](datax/plugin/reader/oracle/README.md)、[写](datax/plugin/writer/oracle/README.md) |
| 无结构流     | CSV                | √            | √          | [读](datax/plugin/reader/csv/README.md)、[写](datax/plugin/writer/csv/README.md) |
|              | XLSX（excel）      | √            | √          | [读](datax/plugin/reader/xlsx/README.md)、[写](datax/plugin/writer/xlsx/README.md) |

### 用户手册

使用[go-etl用户手册](README_USER.md)开始数据同步

### 开发宝典

参考[go-etl开发者文档](datax/README.md)来帮助开发

## 模块简介
### datax

本包将提供类似于阿里巴巴[DataX](https://github.com/alibaba/DataX)的接口去实现go的etl框架，目前主要实现了job框架内的数据同步能力，监控等功能还未实现.

#### plan

- [x] 实现关系型数据库的任务切分
- [x] 实现监控模块
- [x] 实现流控模块
- [ ] 实现关系型数据库入库断点续传

### storage

#### database

目前已经实现了数据库的基础集成，已有mysql和postgresql的实现，如何实现可以查看godoc文档，利用它能非常方便地实现datax数据库间的同步，欢迎大家来提交新的数据同步方式，可以在下面选择新的数据库来同步

#### stream

主要用于字节流的解析，如文件，消息队列，elasticsearch等，字节流格式可以是cvs，json, xml等

#### file

#### mq

##### plan

暂无时间安排计划，欢迎来实现

#### elasticsearch

##### plan

暂无时间安排计划，欢迎来实现

### transform

主要用于类sql数据转化

#### plan

- [ ] 引入tidb数据库的mysql解析能力
- [ ] 引入tidb数据库的mysql函数计算能力
- [ ] 运用mysql解析能力和mysql函数计算能力实现数据转化能力

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

#### license

用于自动新增go代码文件中许可证

```bash
go run tools/license/main.go
```

### libra

主要用于数据库间数据校验

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