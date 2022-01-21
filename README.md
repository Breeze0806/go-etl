# go-etl
[![Go Report Card][report-img]][report][![GoDoc][doc-img]][doc][![LICENSE][license-img]][license]

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

### 安装

#### linux

```bash
export GO111MODULE=on
cd cmd/datax
go build
```

#### windows

```bash
set GO111MODULE=on
cd cmd/datax
go build
```

### 使用方式

#### 试用mysql同步

- 可以使用cmd/datax/mysql/init.sql初始化数据库
- 开启同步mysql命令

```
datax -c mysql/config.json
```

#### 试用postgres同步

- 可以使用cmd/datax/postgres/init.sql初始化数据库
- 开启同步postgres命令

```
datax -c postgres/config.json
```

#### 其他

你也可以编写任意支持数据源之间的同步，欢迎大家来提交新的数据同步方式，可以在下面选择新的数据库来同步

### Support Data Channels

| 类型         | 数据源   | Reader（读） | Writer(写) | 文档                                                         |
| ------------ | -------- | ------------ | ---------- | ------------------------------------------------------------ |
| 关系型数据库 | MySQL    | √            | √          | [读](datax/plugin/reader/mysql/README.md)、[写](datax/plugin/writer/mysql/README.md) |
|              | Postgres | √            | √          | [读](datax/plugin/reader/postgres/README.md)、[写](datax/plugin/writer/postgres/README.md) |

### plan

- [ ] 实现db2数据库reader/writer插件
- [ ] 实现sql server数据库reader/writer插件
- [ ] 实现oracle数据库reader/writer插件
- [ ] 实现cvs文件reader/writer插件
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

#### plan

- [ ] 实现stream的数据传输框架
- [ ] 单元测试stream的数据传输框架
- [ ] 实现cvs文件字节流的数据传输框架并单元测试

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
