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

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/schedule"
)

// Table Table structure
type Table interface {
	fmt.Stringer

	Quoted() string   // Full name of the referenced table
	Instance() string // Instance name, e.g., for MySQL, it's the database name
	Schema() string   // Schema name, e.g., for Oracle, it's the username (schema name)
	Name() string     // Table name, e.g., for MySQL, it's the table name
	Fields() []Field  // Show all columns
}

// Parameter Execution parameters for SQL statements with table, transaction mode, and SQL
type Parameter interface {
	SetTable(Table)                         // Set table or view
	Table() Table                           // Table or view
	TxOptions() *sql.TxOptions              // Transaction mode
	Query([]element.Record) (string, error) // SQL prepare statement
	Agrs([]element.Record) ([]any, error)   // Prepare parameters
}

// ParameterOptions Options for parameters
type ParameterOptions struct {
	Table     Table            // Table or view
	Mode      string           // Write mode, e.g., for MySQL
	TxOptions *sql.TxOptions   // Transaction mode
	Records   []element.Record // Write row
}

// FieldsFetcher Supplementary method for Table, used to specially fetch all columns of a table
type FieldsFetcher interface {
	FetchFields(ctx context.Context, db *DB) error // Get specific column
}

// FieldAdder Supplementary method for Table, used to add columns to a table
type FieldAdder interface {
	AddField(*BaseField) // Add specific column
}

// ExecParameter Supplementary method for Table, used to get the method to generate SQL statements for write mode
type ExecParameter interface {
	ExecParam(string, *sql.TxOptions) (Parameter, bool)
}

// Judger Error evaluator
type Judger interface {
	schedule.RetryJudger

	ShouldOneByOne(err error) bool
}

// BaseTable Basic table, used to embed implementations of various database tables
type BaseTable struct {
	instance string
	schema   string
	name     string
	fields   []Field
}

// NewBaseTable, acquire the basic table through instance name, schema name, and table name
func NewBaseTable(instance, schema, name string) *BaseTable {
	return &BaseTable{
		instance: instance,
		schema:   schema,
		name:     name,
	}
}

// Instance Instance name, e.g., for MySQL, it's the database name; for Oracle, it's the instance name
func (b *BaseTable) Instance() string {
	return b.instance
}

// Schema Schema name, e.g., for MySQL, it's the database name; for Oracle, it's the username
func (b *BaseTable) Schema() string {
	return b.schema
}

// Name Table name, e.g., for MySQL, it's the table name
func (b *BaseTable) Name() string {
	return b.name
}

// String Display string for printing
func (b *BaseTable) String() string {
	return b.instance + "." + b.schema + "." + b.name
}

// Fields Show all columns
func (b *BaseTable) Fields() []Field {
	return b.fields
}

// AppendField Append column
func (b *BaseTable) AppendField(f Field) {
	b.fields = append(b.fields, f)
}

// BaseParam Basic parameters, used to embed SQL parameters for various databases
type BaseParam struct {
	table  Table
	txOpts *sql.TxOptions
}

// NewBaseParam Generate basic parameters through table and transaction parameters txOps
func NewBaseParam(table Table, txOpts *sql.TxOptions) *BaseParam {
	return &BaseParam{
		table:  table,
		txOpts: txOpts,
	}
}

// SetTable Set table
func (b *BaseParam) SetTable(table Table) {
	b.table = table
}

// Table Get table
func (b *BaseParam) Table() Table {
	return b.table
}

// TxOptions Get transaction parameters
func (b *BaseParam) TxOptions() *sql.TxOptions {
	return b.txOpts
}

// InsertParam Insert parameters
type InsertParam struct {
	*BaseParam
}

// NewInsertParam Generate insert parameters through table and transaction parameters txOps
func NewInsertParam(t Table, txOps *sql.TxOptions) *InsertParam {
	return &InsertParam{
		BaseParam: NewBaseParam(t, txOps),
	}
}

// Query Generate a bulk insert SQL statement from multiple records
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

// Args Generate bulk insert parameters from multiple records
func (i *InsertParam) Agrs(records []element.Record) (valuers []any, err error) {
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

			valuers = append(valuers, any(v))
		}
	}
	return
}

// TableQueryParam Table structure query parameters
type TableQueryParam struct {
	*BaseParam
}

// NewTableQueryParam Generate table structure query parameters from Table
func NewTableQueryParam(table Table) *TableQueryParam {
	return &TableQueryParam{
		BaseParam: NewBaseParam(table, nil),
	}
}

// Query Generate a select * from table where 1=2 to acquire the table structure
func (t *TableQueryParam) Query(_ []element.Record) (s string, err error) {
	s = "select * from "
	s += t.table.Quoted() + " where 1 = 2"
	return s, nil
}

// Args Generate parameters, but they are empty
func (t *TableQueryParam) Agrs(_ []element.Record) (a []any, err error) {
	return nil, nil
}
