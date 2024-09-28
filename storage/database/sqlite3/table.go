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

package sqlite3

import (
	"database/sql"
	"database/sql/driver"

	"github.com/Breeze0806/go-etl/storage/database"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/pingcap/errors"
)

// Table represents a Sqlite3 table.
type Table struct {
	*database.BaseTable
	database.BaseConfigSetter
}

// NewTable creates a new Sqlite3 table. Note that at this point, the schema parameter in BaseTable refers to the schema name, instance is the database name, and name is the table name.
func NewTable(b *database.BaseTable) *Table {
	return &Table{
		BaseTable: b,
	}
}

// Quoted refers to the fully qualified name of the table.
func (t *Table) Quoted() string {
	return Quoted(t.Name())
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

// ExecParam retrieves execution parameters, where the copy in parameter mode has been registered.
func (t *Table) ExecParam(mode string, txOpts *sql.TxOptions) (database.Parameter, bool) {
	return nil, false
}

// ShouldRetry determines whether a retry is necessary.
func (t *Table) ShouldRetry(err error) bool {
	switch cause := errors.Cause(err).(type) {
	case sqlite3.Error:
		return true
	default:
		return cause == driver.ErrBadConn
	}
}

// ShouldOneByOne specifies whether to retry one operation at a time.
func (t *Table) ShouldOneByOne(err error) bool {
	switch errors.Cause(err).(type) {
	case sqlite3.Error:
		return true
	}
	return false
}
