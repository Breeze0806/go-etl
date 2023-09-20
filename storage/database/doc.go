// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package database 对实现标准库database/sql的接口的数据库进行封装
// 本包提供了DB作为数据库连接池来操作数据库
// DB可以通过FetchRecord来获得每一行的记录，例如数据库方言名为name,数据库配置文件为conf
//
//	source, err := NewSource(name, conf)
//	if  err != nil {
//		fmt.Println(err)
//		return
//	}
//	db, err:= NewDB(source)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer db.Close()
//	t, err := db.FetchTable(context.TODO(), NewBaseTable("db", "schema", "table"))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	var gotRecords []element.Record
//	if err = db.FetchRecord(context.TODO(), NewTableQueryParam(t),  NewBaseFetchHandler(
//		func() (element.Record, error) {
//			return element.NewDefaultRecord(), nil
//		},
//		func(r element.Record) error {
//			gotRecords = append(gotRecords, r)
//			return nil
//		})); err != nil {
//		fmt.Println(err)
//		return
//	}
//
// DB也可以通过BatchExec来批量处理数据
//
//	source, err := NewSource(name, conf)
//	if  err != nil {
//		fmt.Println(err)
//		return
//	}
//	db, err:= NewDB(source)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer db.Close()
//	t, err := db.FetchTable(context.TODO(), NewBaseTable("db", "schema", "table"))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	columns := [][]element.Column{
//		{
//			element.NewDefaultColumn(element.NewBoolColumnValue(false), "f1", 0),
//			element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "f2", 0),
//			element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1.0), "f3", 0),
//			element.NewDefaultColumn(element.NewStringColumnValue("1"), "f4", 0),
//		},
//		{
//			element.NewDefaultColumn(element.NewBoolColumnValue(true), "f1", 0),
//			element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(2), "f2", 0),
//			element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(2.0), "f3", 0),
//			element.NewDefaultColumn(element.NewStringColumnValue("2"), "f4", 0),
//		},
//	}
//	var wantRecords []element.Record
//	for _, row := range columns {
//		record := element.NewDefaultRecord()
//		for _, c := range row {
//			record.Add(c)
//		}
//		wantRecords = append(wantRecords, record)
//	}
//	if err = db.BatchExec(context.TODO(), &ParameterOptions{
//		Table:     gotTable,
//		TxOptions: nil,
//		Mode:      "insert",
//		Records:   wantRecords,
//	}); err != nil {
//		fmt.Println(err)
//		return
//	}
//
// DB也可以像sql.DB那样通过BeginTx，QueryContext，ExecContext去实现操作数据库的要求
// 另外database包提供了DBWrapper使得数据库连接池能够被复用
// 它实际上是DB的包装，为此它可以表现的和DB一样，例如
//
//	db, err:= Open(name, conf)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer db.Close()
//	t, err := db.FetchTable(context.TODO(), NewBaseTable("db", "schema", "table"))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
// 要使用上述特性，各类实现标准库database/sql的接口的数据库
// 通过RegisterDialect通过名字注册一下数据库方言接口
// 当然数据库配置文件要满足Config的格式
//
//	type Dialect interface {
//		Source(*BaseSource) (Source, error)
//	}
//
// 数据库方言接口主要是返回数据源接口
//
//	type Source interface {
//		Config() *config.JSON   //配置信息
//		Key() string            //dbMap Key
//		DriverName() string     //驱动名，用于sql.Open
//		ConnectName() string    //连接信息，用于sql.Open
//		Table(*BaseTable) Table //获取具体表
//	}
//
// 当然这里可以使用BaseSource来简化Source的实现
// 数据源接口返回对应的表
// Table 表结构
//
//	type Table interface {
//		fmt.Stringer
//		Quoted() string   //引用的表名全称
//		Instance() string //实例名，例如对于mysql就是数据库，对于oracle就是实例
//		Schema() string   //模式名，例如对于mysql就是数据库，对于oracle就是用户名
//		Name() string     //表名，例如对于mysql就是表
//		Fields() []Field  //显示所有列
//	}
//
// 当Table一般实现基本的添加列接口用于FetchTable
//
//	type FieldAdder interface {
//		AddField(*BaseField) //新增具体列
//	}
//
// 如果有特别的SQL获取列方式，则需要实现下列接口去获取
//
//	type FieldsFetcher interface {
//		FetchFields(ctx context.Context, db *DB) error //获取具体列
//	}
//
// 如果在批量处理数据时需要除了insert批量处理数据的语句，或者特殊的insert语句时,还需要实现下列接口
//
//	type ExecParameter interface {
//		ExecParam(string, *sql.TxOptions) (Parameter, bool)
//	}
//
// 通过传入写入模式字符串以及事务选项获取执行参数
// 当然这里也可以使用BaseTable来简化Table的实现
// 每个表包含多列Field
//
//	type Field interface {
//		fmt.Stringer
//		Name() string                 //字段名
//		Quoted() string               //引用字段名
//		BindVar(int) string           //占位符号
//		Select() string               //select字段名
//		Type() FieldType              //字段类型
//		Scanner() Scanner             //扫描器
//		Valuer(element.Column) Valuer //赋值器
//	}
//
// 当然这里可以使用BaseField来简化Field的实现
// 这里接口FieldType时sql.ColumnType抽象，一般无需自己实现
// ，如有特殊的需求也可以实现这个接口用于自己的列类型
//
//	type FieldType interface {
//		Name() string                                   //列名
//		ScanType() reflect.Type                         //扫描类型
//		Length() (length int64, ok bool)                //长度
//		DecimalSize() (precision, scale int64, ok bool) //精度
//		Nullable() (nullable, ok bool)                  //是否为空
//		DatabaseTypeName() string                       //列数据库类型名
//	}
//
// 扫描器接口Scanner是将数据库驱动底层的数据转化成列Column类型，用于读取数据
//
//	type Scanner interface {
//		sql.Scanner
//		Column() element.Column //获取列数据
//	}
//
// 当然这里可以使用BaseScanner来简化Scanner的实现
// 赋值器接口Valuer是将列Column数据库转化成驱动底层的数据，用于处理数据
//
//	type Valuer interface {
//		driver.Valuer
//	}
//
// 特别地，如果使用GoValuer来作为赋值器Valuer的简单实现
// 那么需要在FieldType的基础上简单实现下列接口就可以实现对应的Valuer
//
//	type ValuerGoType interface {
//		GoType() GoType
//	}
//
// 通过上述内容的实现，我们就可以愉快的使用数据库
package database
