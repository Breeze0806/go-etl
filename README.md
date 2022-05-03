# go-etl
[![Go Report Card][report-img]][report][![GoDoc][doc-img]][doc][![LICENSE][license-img]][license][![Build Status][ci-img]][ci][![Coverage Status][cov-img]][cov]

go-etl是一个集数据源抽取，转化，加载，同步校验的工具集，提供强大的数据同步，数据转化甚至数据校验的功能。

go-etl将提供的etl能力如下：

1. 主流数据库的数据抽取以及数据加载的能力，在storage包中实现
2. 类二维表的数据流的数据抽取以及数据加载的能力，在stream包中实现
3. 类似datax的数据同步能力，在datax包中实现
4. 以mysql sql语法为基础的数据筛选、转化能力，在transform包中实现
5. 数据库间的数据校验能力，在libra包中实现 

鉴于本人实在精力有限，欢迎大家来提交issue或者***加QQ群185188648***来讨论go-etl，让我们一起进步!

## 数据同步工具

本数据数据同步工具以下数据源的同步能力

| 类型         | 数据源             | Reader（读） | Writer(写) | 文档                                                         |
| ------------ | ------------------ | ------------ | ---------- | ------------------------------------------------------------ |
| 关系型数据库 | MySQL/Mariadb/Tidb | √            | √          | [读](datax/plugin/reader/mysql/README.md)、[写](datax/plugin/writer/mysql/README.md) |
|              | Postgres/Greenplum | √            | √          | [读](datax/plugin/reader/postgres/README.md)、[写](datax/plugin/writer/postgres/README.md) |
|              | DB2 LUW            | √            | √          | [读](datax/plugin/reader/db2/README.md)、[写](datax/plugin/writer/db2/README.md) |
| 无结构流     | CSV                | √            | √          | [读](datax/plugin/reader/csv/README.md)、[写](datax/plugin/writer/csv/README.md) |
|              | XLSX（excel）      | √            | √          | [读](datax/plugin/reader/xlsx/README.md)、[写](datax/plugin/writer/xlsx/README.md) |

### 安装和使用

在使用下列命令前请确保你已经安装go的编译环境并且设置好了GOPATH

#### linux

```bash
make dependencies
make release
```

#### windows

```bash
release.bat
```
### 使用方式

```bash
Usage of datax:
  -c string
        config (default "config.json")
```

调用datax十分简单，只要直接调用它即可

```bash
data -c config.json
```

当返回值是0，并且显示run success,表示执行成功

当返回值是1，并且显示run fail,并告知执行失败的原因

### 使用示例

#### 使用mysql同步

- 可以使用cmd/datax/mysql/init.sql初始化数据库
- 开启同步mysql命令

```bash
datax -c mysql/config.json
```

#### 使用postgres同步

- 可以使用cmd/datax/postgres/init.sql初始化数据库
- 开启同步postgres命令

```bash
datax -c postgres/config.json
```

#### 使用db2同步

- 注意使用前请下载相应的db2的odbc库，如linux的make dependencies和release.bat
- 注意在linux下如Makefile所示export LD_LIBRARY_PATH=${DB2HOME}/lib
- 注意在windows下如release.bat所示set path=%path%;%GOPATH%\src\github.com\ibmdb\go_ibm_db\clidriver\bin
- 可以使用cmd/datax/db2/init.sql初始化数据库
- 开启同步命令

```bash
datax -c db2/config.json
```

#### 使用csv同步到postgres

- 可以使用cmd/datax/csvpostgres/init.sql初始化数据库
- 开启同步命令

```bash
datax -c csvpostgres/config.json
```

#### 使用xlsx同步到postgres

- 可以使用cmd/datax/csvpostgres/init.sql初始化数据库
- 开启同步命令

```bash
datax -c xlsxpostgres/config.json
```

#### 使用postgres同步csv

- 可以使用cmd/datax/csvpostgres/init.sql初始化数据库
- 开启同步命令

```bash
datax -c postgrescsv/config.json
```

#### 使用postgres同步xlsx

- 可以使用cmd/datax/csvpostgres/init.sql初始化数据库
- 开启同步命令

```bash
datax -c postgresxlsx/config.json
```

#### 其他数据源

你可以选择：
1. 提交issue让别人帮你开发来支持新的数据源
2. 参考开发宝典自己开发来支持新的数据源

### 开发宝典

可以参考[Datax开发者文档](datax/README.md)来帮助开发

## 模块简介
### datax

本包将提供类似于阿里巴巴[DataX](https://github.com/alibaba/DataX)的接口去实现go的etl框架，目前主要实现了job框架内的数据同步能力，监控等功能还未实现.

#### plan

- [x] 实现db2数据库reader/writer插件
- [ ] 实现sql server数据库reader/writer插件
- [ ] 实现oracle数据库reader/writer插件
- [x] 实现cvs文件reader/writer插件
- [x] 实现xlsx文件reader/writer插件
- [ ] 实现监控模块
- [ ] 实现流控模块
- [ ] 实现关系型数据库入库保证数据不丢失功能
- [ ] 实现关系型数据库入库断点续传

### storage

#### database

目前已经实现了数据库的基础集成，已有mysql和postgresql的实现，如何实现可以查看godoc文档，利用它能非常方便地实现datax数据库间的同步，欢迎大家来提交新的数据同步方式，可以在下面选择新的数据库来同步

##### plan

- [x] 实现db2数据库的dialect 
- [ ] 实现sql server数据库的dialect
- [ ] 实现oracle数据库的dialect

#### stream

主要用于字节流的解析，如文件，消息队列，elasticsearch等，字节流格式可以是cvs，json, xml等

#### file
##### plan

- [x] 实现文件流的数据传输框架
- [x] 单元测试文件流的数据传输框架
- [x] 实现cvs文件字节流的数据传输框架
- [x] 单元测试cvs文件字节流的数据传输框架
- [x] 实现xlsx文件字节流的数据传输框架
- [x] 单元测试xlsx文件字节流的数据传输框架

#### mq

##### plan

暂无时间安排计划，欢迎来实现

#### elasticsearch

##### plan

暂无时间安排计划，欢迎来实现

### libra

主要用于数据库间数据校验

#### plan

- [ ] 实现libra的数据校验框架
- [ ] 单元测试libra的数据校验框架
- [ ] 实现MySQL数据库的libra接口并单元测试
- [ ] 系统测试MySQL数据库间校验
- [ ] 完善相关文档，包含代码注释（通过go lint 检查）

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

发布命令，用于将由开发者开发的reader和writer插件注册到程序中的代码

##### plugin

数据源插件模板新增工具，用于新增一个reader或writer模板，配合发布命令使用，减少开发者负担

#### license

用于自动新增go代码文件中许可证

[report-img]:https://goreportcard.com/badge/github.com/Breeze0806/go-etl
[report]:https://goreportcard.com/report/github.com/Breeze0806/go-etl
[doc-img]:https://godoc.org/github.com/Breeze0806/go-etl?status.svg
[doc]:https://godoc.org/github.com/Breeze0806/go-etl
[license-img]: https://img.shields.io/badge/License-Apache%202.0-blue.svg
[license]: https://github.com/Breeze0806/go-etl/blob/main/LICENSE
[ci-img]: https://app.travis-ci.com/Breeze0806/go-etl.svg?branch=main
[ci]: https://app.travis-ci.com/Breeze0806/go-etl
[cov-img]: https://codecov.io/gh/Breeze0806/go-etl/branch/main/graph/badge.svg?token=UGb27Nysga
[cov]: https://codecov.io/gh/Breeze0806/go-etl