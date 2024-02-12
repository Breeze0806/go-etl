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
	"github.com/Breeze0806/go-etl/storage/database"
)

func init() {
	var d Dialect
	database.RegisterDialect(d.Name(), d)
}

// Dialect is the dialect for the DB2 database
type Dialect struct{}

// Source generates a DB2 data source
func (d Dialect) Source(bs *database.BaseSource) (database.Source, error) {
	return NewSource(bs)
}

// Name is the registered name for the database dialect
func (d Dialect) Name() string {
	return "db2"
}

// Source is the DB2 data source
type Source struct {
	*database.BaseSource // Basic data source

	dsn string
}

// NewSource generates a DB2 data source and will report an error if the configuration file is incorrect
func NewSource(bs *database.BaseSource) (s database.Source, err error) {
	source := &Source{
		BaseSource: bs,
	}
	var c *Config
	if c, err = NewConfig(source.Config()); err != nil {
		return
	}

	if source.dsn, err = c.FormatDSN(); err != nil {
		return
	}
	return source, nil
}

// DriverName is the driver name for github.com/ibmdb/go_ibm_db
func (s *Source) DriverName() string {
	return "go_ibm_db"
}

// ConnectName is the connection information for the data source from github.com/ibmdb/go_ibm_db
func (s *Source) ConnectName() string {
	return s.dsn
}

// Key is the keyword for the data source, used for reuse by DBWrapper
func (s *Source) Key() string {
	return s.dsn
}

// Table generates a DB2 table
func (s *Source) Table(b *database.BaseTable) database.Table {
	return NewTable(b)
}

// Quoted is the quoting function for DB2
func Quoted(s string) string {
	return `"` + s + `"`
}
