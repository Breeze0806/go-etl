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

package sqlserver

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/pingcap/errors"
)

// WriteModeCopyIn represents the copy in write mode.
const WriteModeCopyIn = "copyIn"

// Table represents an MSSQL table.
type Table struct {
	database.BaseConfigSetter
	*database.BaseTable
}

// NewTable creates a new MSSQL table. Note that at this point, the schema parameter in BaseTable is empty, instance is the database name, and name is the table name.
func NewTable(b *database.BaseTable) *Table {
	return &Table{
		BaseTable: b,
	}
}

// Quoted refers to the fully qualified name of the table.
func (t *Table) Quoted() string {
	return Quoted(t.Instance()) + "." + Quoted(t.Schema()) + "." + Quoted(t.Name())
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
	case WriteModeCopyIn:
		return NewCopyInParam(t, txOpts), true
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
	_, ok := errors.Cause(err).(*mssql.Error)
	return ok
}

// CopyInParam represents the parameters for the copy in operation.
type CopyInParam struct {
	*database.BaseParam
}

// NewCopyInParam creates copy-in parameters based on the table and transaction options (txOpts).
func NewCopyInParam(t database.Table, txOpts *sql.TxOptions) *CopyInParam {
	return &CopyInParam{
		BaseParam: database.NewBaseParam(t, txOpts),
	}
}

// Query generates a batch of copy in SQL statements for insertion.
func (ci *CopyInParam) Query(_ []element.Record) (query string, err error) {
	var conf *config.JSON
	conf, err = ci.Table().(*Table).Config().GetConfig("bulkOption")
	if err != nil {
		err = nil
		conf, _ = config.NewJSONFromString("{}")
	}

	opt := mssql.BulkOptions{}
	err = json.Unmarshal([]byte(conf.String()), &opt)
	if err != nil {
		return
	}

	var columns []string
	for _, f := range ci.Table().Fields() {
		columns = append(columns, f.Name())
	}
	return mssql.CopyIn(ci.Table().Quoted(), opt,
		columns...), nil
}

// Agrs generates a batch of copy in parameters based on multiple records.
func (ci *CopyInParam) Agrs(records []element.Record) (valuers []interface{}, err error) {
	for _, r := range records {
		for fi, f := range ci.Table().Fields() {
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
