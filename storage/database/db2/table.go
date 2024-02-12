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

package db2

import (
	"database/sql"
	"database/sql/driver"

	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/ibmdb/go_ibm_db"
	"github.com/pingcap/errors"
)

// Table - A db2 table
type Table struct {
	*database.BaseTable
	database.BaseConfigSetter
}

// NewTable - Creates a db2 table. Note that the schema parameter in BaseTable is empty at this point, instance is the database name, and name is the table name.
func NewTable(b *database.BaseTable) *Table {
	return &Table{
		BaseTable: b,
	}
}

// Quoted - Refer to the table's full name
func (t *Table) Quoted() string {
	return Quoted(t.Schema()) + "." + Quoted(t.Name())
}

func (t *Table) String() string {
	return t.Quoted()
}

// AddField - Add a new column
func (t *Table) AddField(baseField *database.BaseField) {
	f := NewField(baseField)
	f.SetConfig(t.Config())
	t.AppendField(f)
}

// ExecParam - Acquire execution parameter
func (t *Table) ExecParam(mode string, txOpts *sql.TxOptions) (database.Parameter, bool) {
	return nil, false
}

// ShouldRetry - Retry
func (t *Table) ShouldRetry(err error) bool {
	return errors.Cause(err) == driver.ErrBadConn
}

// ShouldOneByOne - Retry individually
func (t *Table) ShouldOneByOne(err error) bool {
	_, ok := errors.Cause(err).(*go_ibm_db.Error)
	return ok
}
