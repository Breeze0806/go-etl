# go-etl
[![Go Report Card][report-img]][report][![GoDoc][doc-img]][doc][![LICENSE][license-img]][license][![Build Status][ci-img]][ci][![Coverage Status][cov-img]][cov]

go-etl是一个集数据源抽取，转化，加载，同步校验的工具集，提供强大的数据同步，数据校验甚至数据转储的功能。

go-etl将提供的etl能力如下：

1. 主流数据库的数据抽取以及数据加载的能力，在storage包中实现
2. 类二维表的数据流的数据抽取以及数据加载的能力，在stream包中实现
3. 类似datax的数据同步能力，在datax包中实现
4. 数据库间的数据校验能力，在libra包中实现
5. 以mysql sql语法为基础的数据筛选、转化能力，在transform包中实现

鉴于本人实在精力有限，欢迎大家来提交issue或者***加QQ群185188648***来讨论go-etl，让我们一起进步!

## datax

目前已经基本完成数据同步框架，已有类mysql和类postgresql的数据库的同步能力

### 开发者文档

#### 新增许可证（license）
当你开发完一个功能后在提交前，请运行go run license.go用于自动加入许可证

#### 关系型数据库

如果你想帮忙实现关系型数据库的数据源，根据以下方式去实现你的数据源将更加方便
1. 先实现storage/database的接口，更多信息使用 go doc storage/database/doc.go，可以参考storage/database/mysql和storage/database/postgres的实现。
2. 再利用datax/plugin/reader/rdbm和datax/plugin/writer/rdbm可以更加快速地实现对应功能，实现reader/writer，可以参考storage/database/mysql和storage/database/postgres的实现。
3. 使用 go doc datax/doc.go即可获取datax以及插件开发的要点。

#### 二维表文件流

如果你想帮忙实现二维表文件流的数据源，根据以下方式去实现你的数据源将更加方便

1. 先实现storage/stream/file的接口，更多信息使用 go doc storage/stream/file/doc.go，可以参考storage/stream/file/csv和storage/stream/file/xlsx的实现。
2. 再利用datax/plugin/reader/file和datax/plugin/writer/file可以更加快速地实现对应功能，实现reader/writer，可以参考storage/stream/file/csv和storage/stream/file/xlsx的实现。
3. 使用 go doc datax/doc.go即可获取datax以及插件开发的要点。

#### 其他数据源

- 如果你想实现其他数据源，使用 go doc datax/doc.go即可获取datax以及插件开发的要点。
- 提交issue让其他人帮助你实现。

### 安装和发布

#### linux

```bash
export GO111MODULE=on
go mod download
go mod vendor
go generate ./...
cd cmd/datax
go build
```

#### windows

```bash
set GO111MODULE=on
go mod download
go mod vendor   
go generate ./...
cd cmd/datax
go build
```

### 发布命令解析

```bash
go generate ./...
```
本命令生成将这些reader和writer注册到程序中的代码

主要的原理如下会将对应datax/plugin插件中的reader和writer的resources的plugin.json生成plugin.go，同时在datax目录下生成plugin.go用于导入这些插件， 具体在datax/build实现。

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

如上数据可以在各个数据源之间流转，如MySQL到Postgres

### Support Data Channels

| 类型         | 数据源        | Reader（读） | Writer(写) | 文档                                                         |
| ------------ | ------------- | ------------ | ---------- | ------------------------------------------------------------ |
| 关系型数据库 | MySQL         | √            | √          | [读](datax/plugin/reader/mysql/README.md)、[写](datax/plugin/writer/mysql/README.md) |
|              | Postgres      | √            | √          | [读](datax/plugin/reader/postgres/README.md)、[写](datax/plugin/writer/postgres/README.md) |
| 无结构流     | CVS           | √            | √          | [读](datax/plugin/reader/csv/README.md)、[写](datax/plugin/writer/csv/README.md) |
|              | XLSX（excel） | √            | √           | [读](datax/plugin/reader/xlsx/README.md)、[写](datax/plugin/writer/xlsx/README.md) |

### plan

- [ ] 实现db2数据库reader/writer插件
- [ ] 实现sql server数据库reader/writer插件
- [ ] 实现oracle数据库reader/writer插件
- [x] 实现cvs文件reader/writer插件
- [x] 实现xlsx文件reader/writer插件
- [ ] 实现监控模块
- [ ] 实现流控模块
- [ ] 实现关系型数据库入库保证数据不丢失功能
- [ ] 实现关系型数据库入库断点续传

## storage

### database

目前已经实现了数据库的基础集成，已有mysql和postgresql的实现，如何实现可以查看godoc文档，利用它能非常方便地实现datax数据库间的同步，，欢迎大家来提交新的数据同步方式，可以在下面选择新的数据库来同步

#### plan

- [ ] 实现db2数据库的dialect 
- [ ] 实现sql server数据库的dialect
- [ ] 实现oracle数据库的dialect

### stream

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

## libra

主要用于数据库间数据校验

### plan

- [ ] 实现libra的数据校验框架
- [ ] 单元测试libra的数据校验框架
- [ ] 实现MySQL数据库的libra接口并单元测试
- [ ] 系统测试MySQL数据库间校验
- [ ] 完善相关文档，包含代码注释（通过go lint 检查）

## transform

主要用于类sql数据转化

### plan

- [ ] 引入tidb数据库的mysql解析能力
- [ ] 引入tidb数据库的mysql函数计算能力
- [ ] 运用mysql解析能力和mysql函数计算能力实现数据转化能力

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