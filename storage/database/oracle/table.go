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

package oracle

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/godror/godror"
	"github.com/pingcap/errors"
)

// WriteModeInsert represents the insert into write mode.
const WriteModeInsert = "insert"

// Table represents an Oracle table.
type Table struct {
	*database.BaseTable
	database.BaseConfigSetter
}

// NewTable creates a new Oracle table. Note that at this point, the schema parameter in BaseTable is empty, instance is the database name, and name is the table name.
func NewTable(b *database.BaseTable) *Table {
	return &Table{
		BaseTable: b,
	}
}

// Quoted refers to the fully qualified name of the table.
func (t *Table) Quoted() string {
	return Quoted(t.Schema()) + "." + Quoted(t.Name())
}

func (t *Table) String() string {
	return t.Quoted()
}

// AddField adds a new column to the table.
func (t *Table) AddField(baseField *database.BaseField) {
	f := NewField(baseField)
	f.SetConfig(t.Config())
	t.AppendField(f)
}

// ExecParam retrieves execution parameters, where the replace into parameter mode has been registered.
func (t *Table) ExecParam(mode string, txOpts *sql.TxOptions) (database.Parameter, bool) {
	switch mode {
	case WriteModeInsert:
		return NewInsertParam(t, txOpts), true
	}
	return nil, false
}

// ShouldRetry determines whether a retry is necessary.
func (t *Table) ShouldRetry(err error) bool {
	return godror.IsBadConn(errors.Cause(err))
}

// ShouldOneByOne specifies whether to retry one operation at a time.
func (t *Table) ShouldOneByOne(err error) bool {
	_, ok := errors.Cause(err).(*godror.OraErr)
	return ok && !godror.IsBadConn(err)
}

// InsertParam represents the parameters for the insert into operation.
type InsertParam struct {
	*database.BaseParam
}

// NewInsertParam creates insert parameters based on the table and transaction options (txOpts).
func NewInsertParam(t database.Table, txOpts *sql.TxOptions) *InsertParam {
	return &InsertParam{
		BaseParam: database.NewBaseParam(t, txOpts),
	}
}

// Query generates a batch of insert into SQL statements for insertion based on multiple records.
func (ip *InsertParam) Query(_ []element.Record) (query string, err error) {
	buf := bytes.NewBufferString("insert into ")
	buf.WriteString(ip.Table().Quoted())
	buf.WriteString("(")
	for fi, f := range ip.Table().Fields() {
		if fi > 0 {
			buf.WriteString(",")
		}
		_, err = buf.WriteString(f.Quoted())
	}
	buf.WriteString(") values (")

	for fi, f := range ip.Table().Fields() {
		if fi > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(
			f.BindVar(fi + 1))
	}
	buf.WriteString(")")

	return buf.String(), nil
}

// Agrs generates a batch of insert into parameters based on multiple records.
func (ip *InsertParam) Agrs(records []element.Record) (valuers []interface{}, err error) {
	for fi, f := range ip.Table().Fields() {
		var ba [][]byte
		var sa []string
		for _, r := range records {
			var c element.Column
			if c, err = r.GetByIndex(fi); err != nil {
				return nil, fmt.Errorf("GetByIndex(%v) err: %v", fi, err)
			}
			var v driver.Value
			if v, err = f.Valuer(c).Value(); err != nil {
				return nil, err
			}
			switch data := v.(type) {
			case nil:
				ba = append(ba, nil)
			case []byte:
				ba = append(ba, data)
			case string:
				sa = append(sa, data)
			}
		}
		var a interface{}

		if len(ba) > 0 {
			a = ba
		}
		if len(sa) > 0 {
			a = sa
		}
		valuers = append(valuers, a)
	}
	return
}
