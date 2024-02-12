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

package dbms

import (
	"encoding/json"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/storage/database"
)

// Config represents the configuration for a relational data reader.
type Config interface {
	GetUsername() string               // GetUsername retrieves the username.
	GetPassword() string               // GetPassword retrieves the password.
	GetURL() string                    // GetURL retrieves the connection URL.
	GetColumns() []Column              // GetColumns retrieves the column information.
	GetBaseTable() *database.BaseTable // GetBaseTable retrieves the table information.
	GetWhere() string                  // GetWhere retrieves the query conditions.
	GetSplitConfig() SplitConfig       // GetSplitConfig retrieves the splitting configuration.
	GetQuerySQL() []string             // GetQuerySQL retrieves the query SQL.
}

// Column represents column information.
type Column interface {
	GetName() string // GetTableName retrieves the table name.
}

// BaseColumn represents basic column information.
type BaseColumn struct {
	Name string
}

// GetName retrieves the column name.
func (b *BaseColumn) GetName() string {
	return b.Name
}

// BaseConfig represents the basic configuration for a relational data reader.
type BaseConfig struct {
	Username   string      `json:"username"`   // Username is the user's name.
	Password   string      `json:"password"`   // Password is the user's password.
	Column     []string    `json:"column"`     // Columns is the list of column information.
	Connection ConnConfig  `json:"connection"` // ConnectionInfo is the database connection information.
	Where      string      `json:"where"`      // Where is the query condition.
	Split      SplitConfig `json:"split"`      // SplitKey is the key used for splitting.
	QuerySQL   []string    `json:"querySql"`   // QuerySQL is the SQL query.
}

// NewBaseConfig creates a new instance of BaseConfig based on the provided JSON configuration conf.
func NewBaseConfig(conf *config.JSON) (c *BaseConfig, err error) {
	c = &BaseConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	return
}

// GetUsername retrieves the username.
func (b *BaseConfig) GetUsername() string {
	return b.Username
}

// GetPassword retrieves the password.
func (b *BaseConfig) GetPassword() string {
	return b.Password
}

// GetURL retrieves the URL for connecting to the relational database.
func (b *BaseConfig) GetURL() string {
	return b.Connection.URL
}

// GetColumns retrieves the column information.
func (b *BaseConfig) GetColumns() (columns []Column) {
	for _, v := range b.Column {
		columns = append(columns, &BaseColumn{
			Name: v,
		})
	}
	return
}

// GetBaseTable retrieves the table information.
func (b *BaseConfig) GetBaseTable() *database.BaseTable {
	return database.NewBaseTable(b.Connection.Table.Db, b.Connection.Table.Schema, b.Connection.Table.Name)
}

// GetWhere retrieves the query conditions.
func (b *BaseConfig) GetWhere() string {
	return b.Where
}

// GetSplitConfig retrieves the splitting configuration.
func (b *BaseConfig) GetSplitConfig() SplitConfig {
	return b.Split
}

// GetQuerySQL retrieves the SQL query.
func (b *BaseConfig) GetQuerySQL() []string {
	return b.QuerySQL
}

// ConnConfig represents the configuration for connecting to a database.
type ConnConfig struct {
	URL   string      `json:"url"`   // ConnectToDatabase establishes a connection to the database.
	Table TableConfig `json:"table"` // TableConfig represents the configuration for a table.
}

// TableConfig represents the configuration for a table.
type TableConfig struct {
	Db     string `json:"db"`     // Database is the name of the database.
	Schema string `json:"schema"` // Schema is the schema name.
	Name   string `json:"name"`   // TableName is the name of the table.
}
