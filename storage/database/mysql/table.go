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

package mysql

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/go-sql-driver/mysql"
	"github.com/pingcap/errors"
)

// WriteModeReplace represents the replace into write mode.
const WriteModeReplace = "replace"

// Table represents a MySQL table.
type Table struct {
	*database.BaseTable
	database.BaseConfigSetter
}

// NewTable creates a new MySQL table. Note that at this point, the schema parameter in BaseTable is empty, instance is the database name, and name is the table name.
func NewTable(b *database.BaseTable) *Table {
	return &Table{
		BaseTable: b,
	}
}

// Quoted refers to the fully qualified name of the table.
func (t *Table) Quoted() string {
	return Quoted(t.Instance()) + "." + Quoted(t.Name())
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
	case "replace":
		return NewReplaceParam(t, txOpts), true
	}
	return nil, false
}

// ShouldRetry determines whether a retry is necessary.
func (t *Table) ShouldRetry(err error) bool {
	switch cause := errors.Cause(err).(type) {
	case net.Error:
		return true
	default:
		return cause == driver.ErrBadConn
	}
}

// ShouldOneByOne specifies whether to retry one operation at a time.
func (t *Table) ShouldOneByOne(err error) bool {
	_, ok := errors.Cause(err).(*mysql.MySQLError)
	return ok
}

// ReplaceParam represents the parameters for the replace into operation.
type ReplaceParam struct {
	*database.BaseParam
}

// NewReplaceParam creates replace parameters based on the table and transaction options (txOpts).
func NewReplaceParam(t database.Table, txOpts *sql.TxOptions) *ReplaceParam {
	return &ReplaceParam{
		BaseParam: database.NewBaseParam(t, txOpts),
	}
}

// Query generates a batch of replace into SQL statements for insertion based on multiple records.
func (rp *ReplaceParam) Query(records []element.Record) (query string, err error) {
	buf := bytes.NewBufferString("replace into ")
	buf.WriteString(rp.Table().Quoted())
	buf.WriteString("(")
	for fi, f := range rp.Table().Fields() {
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
		for fi, f := range rp.Table().Fields() {
			if fi > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(
				f.BindVar(ri*len(rp.Table().Fields()) + fi + 1))
		}
		buf.WriteString(")")
	}
	return buf.String(), nil
}

// Agrs generates a batch of replace into parameters based on multiple records.
func (rp *ReplaceParam) Agrs(records []element.Record) (valuers []interface{}, err error) {
	for _, r := range records {
		for fi, f := range rp.Table().Fields() {
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
