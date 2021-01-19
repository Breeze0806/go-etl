package database

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/Breeze0806/go-etl/element"
)

type Table interface {
	fmt.Stringer

	Quoted() string   //引用的表名全称
	Instance() string //实例名
	Schema() string   //模式名
	Name() string     //表名
	Fields() []Field  //显示所有列
}

type Parameter interface {
	Table() Table
	TxOptions() *sql.TxOptions
	Query([]element.Record) (string, error)
	Agrs([]element.Record) ([]interface{}, error)
}

type ParameterOptions struct {
	Table     Table
	Mode      string
	TxOptions *sql.TxOptions
	Records   []element.Record
}

type FieldsFetcher interface {
	FetchFields(ctx context.Context, db *DB) error //获取具体列
}

type FieldAdder interface {
	AddField(*BaseField) //新增具体列
}

type ExecParameter interface {
	ExecParam(string, *sql.TxOptions) (Parameter, bool)
}

type BaseTable struct {
	instance string
	schema   string
	name     string
	fields   []Field
}

func NewBaseTable(instance, schema, name string) *BaseTable {
	return &BaseTable{
		instance: instance,
		schema:   schema,
		name:     name,
	}
}

func (b *BaseTable) Instance() string {
	return b.instance
}

func (b *BaseTable) Schema() string {
	return b.schema
}

func (b *BaseTable) Name() string {
	return b.name
}

func (b *BaseTable) String() string {
	return b.instance + "." + b.schema + "." + b.name
}

func (b *BaseTable) Fields() []Field {
	return b.fields
}

func (b *BaseTable) AppendField(f Field) {
	b.fields = append(b.fields, f)
}

type BaseParam struct {
	table Table
	txOps *sql.TxOptions
}

func NewBaseParam(table Table, txOps *sql.TxOptions) *BaseParam {
	return &BaseParam{
		table: table,
		txOps: txOps,
	}
}

func (b *BaseParam) SetTable(t Table) {
	b.table = t
}

func (b *BaseParam) SettxOps(txOps *sql.TxOptions) {
	b.txOps = txOps
}

func (b *BaseParam) Table() Table {
	return b.table
}

func (b *BaseParam) TxOptions() *sql.TxOptions {
	return b.txOps
}

type InsertParam struct {
	*BaseParam
}

func NewInsertParam(t Table, txOps *sql.TxOptions) *InsertParam {
	return &InsertParam{
		BaseParam: NewBaseParam(t, txOps),
	}
}

func (i *InsertParam) Query(records []element.Record) (query string, err error) {
	buf := bytes.NewBufferString("insert into ")
	buf.WriteString(i.table.Quoted())
	buf.WriteString("(")
	for fi, f := range i.Table().Fields() {
		if fi > 0 {
			buf.WriteString(",")
		}
		_, err = buf.WriteString(f.Quoted())
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
				f.BindVar(ri*len(i.table.Fields()) + fi + 1))
		}
		buf.WriteString(")")
	}
	return buf.String(), nil
}

func (i *InsertParam) Agrs(records []element.Record) (valuers []interface{}, err error) {
	for _, r := range records {
		for _, f := range i.Table().Fields() {
			var c element.Column
			if c, err = r.GetByName(f.Name()); err != nil {
				return nil, fmt.Errorf("GetByName(%v) err: %v", f.Name(), err)
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

type TableQueryParam struct {
	*BaseParam
}

func NewTableQueryParam(table Table) *TableQueryParam {
	return &TableQueryParam{
		BaseParam: NewBaseParam(table, nil),
	}
}

func (t *TableQueryParam) Query(records []element.Record) (s string, err error) {
	s = "select * from "
	s += t.table.Quoted() + " where 1 = 2"
	return s, nil
}

func (t *TableQueryParam) Agrs(records []element.Record) (a []interface{}, err error) {
	return nil, nil
}
