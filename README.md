# go-etl
[![Go Report Card][report-img]][report][![GoDoc][doc-img]][doc][![LICENSE][license-img]][license]

go-etl是一个集数据源抽取，转化，加载，同步校验的工具集，提供强大的数据同步，数据校验甚至数据转储的功能。

go-etl将提供的etl能力如下：

1. 主流数据库的数据抽取以及数据加载的能力，这个计划在storage包中实现
2. 类似datax的数据同步能力，这个计划在datax包中实现
3. 数据库间的数据校验能力，这个计划在libra包中实现
4. 以mysql sql语法为基础的数据筛选、转化能力，这个计划在transform包中实现（计划中）

## plan

### datax

- [x] 实现datax的同步框架，不包含监控以及流控模块
- [x] 单元测试datax的同步框架，不包含监控以及流控模块
- [ ] 实现MySQL基于datax的同步接口，并单元测试
- [ ] 系统测试MySQL数据库间的同步
- [ ] 完善相关文档，包含代码注释（通过go lint 检查）
- [ ] 实现监控以及流控模块,并单元测试（延后实现）

### storage

- [x] 实现数据库的数据抽取以及数据加载框架，并单元测试
- [x] 实现MySQL数据库数据抽取以及数据加载的相应接口，并单元测试
- [ ] 结合MySQL测试系统测试数据库的数据抽取以及数据加载框架
- [x] 完善相关文档，包含代码注释（通过go lint 检查）

### libra

- [ ] 实现libra的数据校验框架
- [ ] 单元测试libra的数据校验框架
- [ ] 实现MySQL数据库的libra接口并单元测试
- [ ] 系统测试MySQL数据库间校验
- [ ] 完善相关文档，包含代码注释（通过go lint 检查）

### transform

目前计划中

[report-img]:https://goreportcard.com/badge/github.com/Breeze0806/go-etl
[report]:https://goreportcard.com/report/github.com/Breeze0806/go-etl
[doc-img]:https://godoc.org/github.com/Breeze0806/go-etl?status.svg
[doc]:https://godoc.org/github.com/Breeze0806/go-etl
[license-img]: https://img.shields.io/badge/License-Apache%202.0-blue.svg
[license]: https://github.com/Breeze0806/go-etl/blob/main/LICENSE












