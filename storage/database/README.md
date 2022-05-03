# 数据库存储开发者指南

数据库存储是数据库查询和执行SQL的框架，用于关系型数据库的抽象

## 数据库存储简介

数据库存储通过封装golang标准库的database/sql的DB结构体用于查询和执行数据库SQL，在db.go封装了较为丰富的方法，不仅提供database/sql的DB原有的方法BeginTx，PingContext，QueryContext和ExecContext，而且还提供了FetchTable和FetchTableWithParam用于获取表结构，FetchRecord和FetchRecordWithTx用于获取表中的记录，BatchExec，BatchExecWithTx和BatchExecStmtWithTx用于执行写入记录语句。

但是对于不同的数据库使用数据库存储可以实现不同的数据库方言dialect，本文档会介绍如何实现数据库方言dialect接口。

## 数据库方言dialect接口介绍

实现dialect接口的前提是**数据库对应的驱动库可以实现golang标准库的database/sql的接口**。

具体实现时，可以参考以下目录去实现，这里以mysql为例

```go
storage--database--mysql----+--config.go        
                            |--doc.go
							|--field.go
                            |--source.go
                            |--table.go

```

通过这种方式，目前已经实现了mysql，postgres，db2

### 数据源接口

```golang
//Dialect 数据库方言
type Dialect interface {
	Source(*BaseSource) (Source, error) //数据源
}

//Source 数据源,包含驱动信息，包信息，配置文件以及连接信息
type Source interface {
	Config() *config.JSON   //配置信息
	Key() string            //一般是连接信息
	DriverName() string     //驱动名，用于sql.Open的第1个参数
	ConnectName() string    //连接信息，用于sql.Open的第2个参数
	Table(*BaseTable) Table //获取具体表
}
```

具体实现Source接口时，可以组合BaseSource以简化具体实现Source接口的实现Table方法可以返回具体的表结构接口。可以看mysql包source.go的实现。

另外，连接信息依赖Config的依赖。目前Config需要用下面的方式定义，否则无法使用rdbm包来实现datax的插件，可以看mysql包config.go的实现。

```go
type Config struct {
	URL      string `json:"url"`      //数据库url，包含数据库地址，数据库其他参数
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
}
```

此外，需要使用init函数将具体Dialect注册

```go
func init() {
	var d Dialect
	database.RegisterDialect(d.Name(), d)
}
```

### 表结构接口

```go
//Table 表结构
type Table interface {
	fmt.Stringer

	Quoted() string   //引用的表名全称
	Instance() string //实例名，例如对于mysql就是数据库
	Schema() string   //模式名，例如对于oracle就是用户名（模式名）
	Name() string     //表名，例如对于mysql就是表
	Fields() []Field  //显示所有列
}

//FieldsFetcher Table的补充方法，用于特殊获取表的所有列
type FieldsFetcher interface {
	FetchFields(ctx context.Context, db *DB) error //获取具体列
}

//FieldAdder Table的补充方法，用于新增表的列
type FieldAdder interface {
	AddField(*BaseField) //新增具体列
}

//ExecParameter Table的补充方法，用于写模式获取生成sql语句的方法
type ExecParameter interface {
	ExecParam(string, *sql.TxOptions) (Parameter, bool)
}
```

具体实现Table接口时，可以组合BaseTable以简化具体Table接口的实现，其中Fields方法必须返回相应数据库的具体字段接口集合。具体可以看mysql包table.go的实现。

FetchFields和FieldAdder只可以实现一个，一般选择FieldAdder接口。ExecParameter可以用于实现批量入表的SQL语句，如对于mysql，可以实现replace into的方法插入，目前默认实现了普遍适用的insert方法，但是对于如实现oracle使用gora驱动的话，insert方法不适用于这种情况。

```go
//Parameter 带有表，事务模式，sql语句的执行参数
type Parameter interface {
	Table() Table                                 //表或者视图
	TxOptions() *sql.TxOptions                    //事务模式
	Query([]element.Record) (string, error)       //sql prepare语句
	Agrs([]element.Record) ([]interface{}, error) //prepare参数
}
```

如上对于实现replace into的方法插入需要实现Parameter，可以组合BaseParam简化具体实现Parameter接口的实现，可以参考mysql包table.go的实现。

### 字段接口

```go
//Field 数据库字段
type Field interface {
	fmt.Stringer

	Index() int                   //索引
	Name() string                 //字段名
	Quoted() string               //引用字段名
	BindVar(int) string           //占位符号
	Select() string               //select字段名
	Type() FieldType              //字段类型
	Scanner() Scanner             //扫描器,用于将数据库对应数据转化成列
	Valuer(element.Column) Valuer //赋值器,用于将列转化成数据库对应数据
}
```

具体实现Field接口时，可以组合BaseField以简化具体实现Field接口的实现，Type()方法必须返回相应数据库的列类型，Scanner必须返回相应数据库的扫描器，Valuer须返回相应数据库的扫描器。具体可以看mysql包中field.go的实现。

```go
//ColumnType 列类型,抽象 sql.ColumnType，也方便自行实现对应函数
type ColumnType interface {
	Name() string                                   //列名
	ScanType() reflect.Type                         //扫描类型
	Length() (length int64, ok bool)                //长度
	DecimalSize() (precision, scale int64, ok bool) //精度
	Nullable() (nullable, ok bool)                  //是否为空
	DatabaseTypeName() string                       //列数据库类型名
}

//FieldType 字段类型
type FieldType interface {
	ColumnType

	IsSupportted() bool //是否支持该类型
}
```

具体实现FieldType接口时，可以组合BaseFieldType以简化具体实现FieldFieldType接口的实现，ColumnType在事实上是sql.ColumnType的抽象。具体可以看mysql包field.go的实现。

```go
//Scanner 列数据扫描器 数据库驱动的值扫描成列数据
type Scanner interface {
	sql.Scanner

	Column() element.Column //获取列数据
}
```

具体实现Scanner的接口时,可以组合BaseFieldType以简化具体实现FieldType接口的实现,Scanner的作用将数据库驱动读取到的数据转化为单列的数据。具体可以看mysql包field.go的实现。

```go
//Valuer 赋值器 将对应数据转化成数据库驱动的值
type Valuer interface {
	driver.Valuer
}

//ValuerGoType 用于赋值器的golang类型判定,是Field的可选功能，
//就是对对应驱动的值返回相应的值，方便GoValuer进行判定
type ValuerGoType interface {
	GoType() GoType
}
```

具体实现Valuer的接口时,可以组合GoValuer以简化具体实现Valuer接口的实现，使用GoValuer需要在数据库层Field实现ValuerGoType接口，Valuer的作用将单列的数据转化为数据库驱动写入的数据类型。具体可以看mysql包field.go的实现。