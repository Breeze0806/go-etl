// Package mysql 实现了mysql的数据库方言Dialect，支持mysql 5.6+ 对应数据库
// 驱动为github.com/go-sql-driver/mysql
// 数据源Source使用BaseSource来简化实现, 对github.com/go-sql-driver/mysql
// 驱动进行包装.对于数据库配置，需要和Config一致
// 表Table使用BaseTable来简化实现,也是基于github.com/go-sql-driver/mysql的
// 封装,Table实现了FieldAdder的方式去获取列,在ExecParameter中实现写入模式为
// replace的repalce into批量数据处理模式,写入模式为insert的插入模式复用
// 已有的database.InsertParam
// 列Field使用BaseField来简化实现,其中FieldType采用了原来的sql.ColumnType，
// 并实现了ValuerGoType
// 扫描器Scanner使用BaseScanner来简化实现
// 赋值器Valuer 使用了GoValuer的实现方式
package mysql
