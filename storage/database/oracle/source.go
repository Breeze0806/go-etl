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
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/godror/godror"
)

func init() {
	var d Dialect
	database.RegisterDialect(d.Name(), d)
}

// Dialect represents the database dialect for Oracle
type Dialect struct{}

// Source generates an Oracle data source
func (d Dialect) Source(bs *database.BaseSource) (database.Source, error) {
	return NewSource(bs)
}

// Name is the registered name of the database dialect
func (d Dialect) Name() string {
	return "oracle"
}

// Source oracle refers to the Oracle data source
type Source struct {
	*database.BaseSource // Basic data source

	dsn string
}

// NewSource generates an Oracle data source and will report an error if there's an issue with the configuration file
func NewSource(bs *database.BaseSource) (s database.Source, err error) {
	source := &Source{
		BaseSource: bs,
	}
	var c *Config
	if c, err = NewConfig(source.Config()); err != nil {
		return
	}
	var connParam godror.ConnectionParams
	if connParam, err = c.FetchConnectionParams(); err != nil {
		return
	}
	source.dsn = connParam.StringWithPassword()
	return source, nil
}

// DriverName is the driver name for github.com/godror/godror
func (s *Source) DriverName() string {
	return "godror"
}

// ConnectName is the connection information for the Oracle data source using github.com/godror/godror
func (s *Source) ConnectName() string {
	return s.dsn
}

// Key is a keyword for the data source, used for reuse by DBWrapper
func (s *Source) Key() string {
	return s.dsn
}

// Table generates a table for Oracle
func (s *Source) Table(b *database.BaseTable) database.Table {
	return NewTable(b)
}

// Quoted is the quoting function for Oracle
func Quoted(s string) string {
	return `"` + s + `"`
}
