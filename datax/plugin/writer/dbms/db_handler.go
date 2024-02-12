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
	"database/sql"

	"github.com/Breeze0806/go-etl/config"
	dbmsreader "github.com/Breeze0806/go-etl/datax/plugin/reader/dbms"
	"github.com/Breeze0806/go-etl/storage/database"
)

// DbHandler Database Execution Handler Encapsulation
type DbHandler interface {
	Execer(name string, conf *config.JSON) (Execer, error)      // Obtain an executor through the database name and configuration
	Config(conf *config.JSON) (Config, error)                   // Obtain relational database configuration through configuration
	TableParam(config Config, execer Execer) database.Parameter // Obtain table parameters through relational database configuration and executor
}

// BaseDbHandler Basic Database Execution Handler Encapsulation
type BaseDbHandler struct {
	newExecer func(name string, conf *config.JSON) (Execer, error)
	opts      *sql.TxOptions
}

// NewBaseDbHandler Create a database execution handler encapsulation using the executor function newExecer and database transaction execution options opts
func NewBaseDbHandler(newExecer func(name string, conf *config.JSON) (Execer, error), opts *sql.TxOptions) *BaseDbHandler {
	return &BaseDbHandler{
		newExecer: newExecer,
		opts:      opts,
	}
}

// Execer Obtain an executor through the database name and configuration
func (d *BaseDbHandler) Execer(name string, conf *config.JSON) (Execer, error) {
	return d.newExecer(name, conf)
}

// Config Obtain relational database configuration through configuration
func (d *BaseDbHandler) Config(conf *config.JSON) (Config, error) {
	return NewBaseConfig(conf)
}

// TableParam Obtain table parameters through relational database configuration and executor
func (d *BaseDbHandler) TableParam(config Config, execer Execer) database.Parameter {
	return dbmsreader.NewTableParam(config, execer, d.opts)
}
