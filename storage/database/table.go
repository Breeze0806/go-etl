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

package database

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

//Table 表结构
type Table interface {
	fmt.Stringer

	Quoted() string   //引用的表名全称
	Instance() string //实例名，例如对于mysql就是数据库
	Schema() string   //模式名，例如对于oracle就是用户名（模式名）
	Name() string     //表名，例如对于mysql就是表
	Fields() []Field  //显示所有列
}

//Parameter 带有表，事务模式，sql语句的执行参数
type Parameter interface {
	Table() Table                                 //表或者视图
	TxOptions() *sql.TxOptions                    //事务模式
	Query([]element.Record) (string, error)       //sql prepare语句
	Agrs([]element.Record) ([]interface{}, error) //prepare参数
}

//ParameterOptions 参数选项
type ParameterOptions struct {
	Table     Table            //表或者视图
	Mode      string           //写入模式，例如mysql
	TxOptions *sql.TxOptions   //事务模式
	Records   []element.Record //写入行
}

//FieldsFetcher Table的补充方法，用于特殊获取表的所有列
type FieldsFetcher interface {
	FetchFields(ctx context.Context, db *DB) error //获取具体列
}

//FieldAdder Table的补充方法，用于新增表的列
type FieldAdder interface {
	AddField(*BaseField) //新增具体列
}

//TableConfigSetter Table的补充方法，用于设置json配置文件
type TableConfigSetter interface {
	SetConfig(conf *config.JSON)
}

//ExecParameter Table的补充方法，用于写模式获取生成sql语句的方法
type ExecParameter interface {
	ExecParam(string, *sql.TxOptions) (Parameter, bool)
}

//BaseTable 基本表，用于嵌入各种数据库Table的实现
type BaseTable struct {
	instance string
	schema   string
	name     string
	fields   []Field
}

//NewBaseTable ，通过实例名，模式名，表明获取基本表
func NewBaseTable(instance, schema, name string) *BaseTable {
	return &BaseTable{
		instance: instance,
		schema:   schema,
		name:     name,
	}
}

//Instance 实例名，例如对于mysql就是数据库，对于oracle就是实例
func (b *BaseTable) Instance() string {
	return b.instance
}

//Schema 模式名，例如对于mysql就是数据库，对于oracle就是用户名
func (b *BaseTable) Schema() string {
	return b.schema
}

//Name 表名，例如对于mysql就是表
func (b *BaseTable) Name() string {
	return b.name
}

//String 用于打印的显示字符串
func (b *BaseTable) String() string {
	return b.instance + "." + b.schema + "." + b.name
}

//Fields 显示所有列
func (b *BaseTable) Fields() []Field {
	return b.fields
}

//AppendField 追加列
func (b *BaseTable) AppendField(f Field) {
	b.fields = append(b.fields, f)
}

//BaseParam 基础参数，用于嵌入各类数据库sql参数的
type BaseParam struct {
	table  Table
	txOpts *sql.TxOptions
}

//NewBaseParam 通过表table和事务参数txOps生成基础参数
func NewBaseParam(table Table, txOpts *sql.TxOptions) *BaseParam {
	return &BaseParam{
		table:  table,
		txOpts: txOpts,
	}
}

//Table 获取表
func (b *BaseParam) Table() Table {
	return b.table
}

//TxOptions 获取事务参数
func (b *BaseParam) TxOptions() *sql.TxOptions {
	return b.txOpts
}

//InsertParam 插入参数
type InsertParam struct {
	*BaseParam
}

//NewInsertParam  通过表table和事务参数txOps插入参数
func NewInsertParam(t Table, txOps *sql.TxOptions) *InsertParam {
	return &InsertParam{
		BaseParam: NewBaseParam(t, txOps),
	}
}

//Query 通过多条记录 records生成批量插入sql语句
func (i *InsertParam) Query(records []element.Record) (query string, err error) {
	buf := bytes.NewBufferString("insert into ")
	buf.WriteString(i.Table().Quoted())
	buf.WriteString("(")
	for fi, f := range i.Table().Fields() {
		if fi > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(f.Quoted())
	}
	buf.WriteString(") values")

	for ri := range records {
		if ri > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("(")
		for fi, f := range i.Table().Fields() {
			if fi > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(
				f.BindVar(ri*len(i.Table().Fields()) + fi + 1))
		}
		buf.WriteString(")")
	}
	return buf.String(), nil
}

//Agrs 通过多条记录 records生成批量插入参数
func (i *InsertParam) Agrs(records []element.Record) (valuers []interface{}, err error) {
	for _, r := range records {
		for fi, f := range i.Table().Fields() {
			var c element.Column
			if c, err = r.GetByIndex(fi); err != nil {
				return nil, fmt.Errorf("GetByIndex(%v) err: %v", fi, err)
			}
			var v driver.Value
			if v, err = f.Valuer(c).Value(); err != nil {
				return nil, err
			}

			valuers = append(valuers, interface{}(v))
		}
	}
	return
}

//TableQueryParam 表结构查询参数
type TableQueryParam struct {
	*BaseParam
}

//NewTableQueryParam 通过表Table生成表结构查询参数
func NewTableQueryParam(table Table) *TableQueryParam {
	return &TableQueryParam{
		BaseParam: NewBaseParam(table, nil),
	}
}

//Query 生成select * from table where 1=2来获取表结构
func (t *TableQueryParam) Query(_ []element.Record) (s string, err error) {
	s = "select * from "
	s += t.table.Quoted() + " where 1 = 2"
	return s, nil
}

//Agrs  生成参数，不过为空
func (t *TableQueryParam) Agrs(_ []element.Record) (a []interface{}, err error) {
	return nil, nil
}
