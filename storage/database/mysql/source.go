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
	"database/sql/driver"

	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/go-sql-driver/mysql"
)

func init() {
	var d Dialect
	database.RegisterDialect(d.Name(), d)
}

// Dialect represents the database dialect for MySQL
type Dialect struct{}

// Source refers to the production data source
func (d Dialect) Source(bs *database.BaseSource) (database.Source, error) {
	return NewSource(bs)
}

// Name is the registered name of the database dialect
func (d Dialect) Name() string {
	return "mysql"
}

// Source mysql represents the MySQL data source
type Source struct {
	*database.BaseSource // Basic data source

	dsn       string
	mysqlConf *mysql.Config
}

// NewSource generates a MySQL data source and will report an error if there's an issue with the configuration file
func NewSource(bs *database.BaseSource) (s database.Source, err error) {
	source := &Source{
		BaseSource: bs,
	}
	var c *Config
	if c, err = NewConfig(source.Config()); err != nil {
		return
	}

	if source.mysqlConf, err = c.FetchMysqlConfig(); err != nil {
		return
	}
	source.dsn = source.mysqlConf.FormatDSN()
	return source, nil
}

// DriverName is the driver name for github.com/go-sql-driver/mysql
func (s *Source) DriverName() string {
	return "mysql"
}

// ConnectName is the connection information for the MySQL data source using github.com/go-sql-driver/mysql
func (s *Source) ConnectName() string {
	return s.dsn
}

// Key is a keyword for the data source, used for reuse by DBWrapper
func (s *Source) Key() string {
	return s.dsn
}

// Table generates a table for MySQL
func (s *Source) Table(b *database.BaseTable) database.Table {
	return NewTable(b)
}

// Connector is the data source connector for github.com/go-sql-driver/mysql
func (s *Source) Connector() (driver.Connector, error) {
	return mysql.NewConnector(s.mysqlConf)
}

// Quoted is the quoting function for MySQL
func Quoted(s string) string {
	return "`" + s + "`"
}
