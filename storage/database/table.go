package database

import (
	"bytes"
	"context"
	"database/sql"
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
	ExecParam(string, *sql.TxOptions) Parameter
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

func (b *BaseTable) Fields() []Field {
	return b.fields
}

func (b *BaseTable) AppendField(f Field) {
	b.fields = append(b.fields, f)
}

type BaseParam struct {
	table Table
}

func NewBaseParam(table Table) *BaseParam {
	return &BaseParam{
		table: table,
	}
}

func (b *BaseParam) Table() Table {
	return b.table
}

type InsertParam struct {
	*BaseParam
	txOps *sql.TxOptions
}

func NewInsertParam(t Table, txOps *sql.TxOptions) *InsertParam {
	return &InsertParam{
		BaseParam: NewBaseParam(t),
		txOps:     txOps,
	}
}

func (i *InsertParam) TxOptions() *sql.TxOptions {
	return i.txOps
}

func (i *InsertParam) Query(records []element.Record) (s string, err error) {
	buf := bytes.NewBufferString("insert into ")
	if _, err = buf.WriteString(i.table.Quoted()); err != nil {
		return
	}
	if _, err = buf.WriteString("("); err != nil {
		return
	}
	for fi, f := range i.Table().Fields() {
		if fi > 0 {
			if _, err = buf.WriteString(","); err != nil {
				return
			}
		}
		if _, err = buf.WriteString(f.Quoted()); err != nil {
			return
		}
	}
	if _, err = buf.WriteString(") values"); err != nil {
		return
	}

	for ri := range records {
		if ri > 0 {
			if _, err = buf.WriteString(","); err != nil {
				return
			}
		}
		if _, err = buf.WriteString("("); err != nil {
			return
		}
		for fi, f := range i.Table().Fields() {
			if fi > 0 {
				if _, err = buf.WriteString(","); err != nil {
					return
				}
			}
			if _, err = buf.WriteString(
				f.BindVar(ri*len(i.table.Fields()) + ri + 1)); err != nil {
				return
			}
		}
		if _, err = buf.WriteString(")"); err != nil {
			return
		}
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
			valuers = append(valuers, f.Valuer(c))
		}
	}
	return
}

type TableQueryParam struct {
	*BaseParam
}

func NewTableQueryParam(table Table) *TableQueryParam {
	return &TableQueryParam{
		BaseParam: NewBaseParam(table),
	}
}

func (t *TableQueryParam) TxOptions() *sql.TxOptions {
	return nil
}

func (t *TableQueryParam) Query(records []element.Record) (s string, err error) {
	s = "select * from "
	s += t.table.Quoted() + "where 1 = 2"
	return s, nil
}

func (t *TableQueryParam) Agrs(records []element.Record) (a []interface{}, err error) {
	return nil, nil
}
