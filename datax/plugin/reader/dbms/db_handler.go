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
	"github.com/Breeze0806/go-etl/storage/database"
)

// DbHandler - Database Handler
type DbHandler interface {
	Querier(name string, conf *config.JSON) (Querier, error)         // Obtain a querier based on the database name (name) and JSON configuration (conf)
	Config(conf *config.JSON) (Config, error)                        // Acquire the relational database input configuration using the JSON configuration (conf)
	TableParam(config Config, querier Querier) database.Parameter    // Retrieve table parameters using the relational database input configuration (config) and querier
	SplitParam(config Config, querier Querier) database.Parameter    // Obtain split table parameters using the relational database input configuration (config) and querier
	MinParam(config Config, table database.Table) database.Parameter // Get the minimum split value parameter based on the relational database input configuration (config) and table (Table)
	MaxParam(config Config, table database.Table) database.Parameter // Retrieve the maximum split value parameter using the relational database input configuration (config) and table querier (Table)
}

// BaseDbHandler - Basic Database Handler
type BaseDbHandler struct {
	newQuerier func(name string, conf *config.JSON) (Querier, error)
	opts       *sql.TxOptions
}

// Create a new instance of the BasicDbHandler using the function to obtain a querier (newQuerier) and transaction options (opts)
func NewBaseDbHandler(newQuerier func(name string, conf *config.JSON) (Querier, error), opts *sql.TxOptions) *BaseDbHandler {
	return &BaseDbHandler{
		newQuerier: newQuerier,
		opts:       opts,
	}
}

// Querier - Acquire a querier based on the database name (name) and JSON configuration (conf)
func (d *BaseDbHandler) Querier(name string, conf *config.JSON) (Querier, error) {
	return d.newQuerier(name, conf)
}

// Config - Retrieve the relational database input configuration using the JSON configuration (conf)
func (d *BaseDbHandler) Config(conf *config.JSON) (Config, error) {
	return NewBaseConfig(conf)
}

// TableParam - Get table parameters using the relational database input configuration (config) and querier
func (d *BaseDbHandler) TableParam(config Config, querier Querier) database.Parameter {
	return NewTableParam(config, querier, d.opts)
}

// SplitParam - Obtain the minimum split value parameter based on the relational database input configuration (config) and table (Table)
func (d *BaseDbHandler) SplitParam(config Config, querier Querier) database.Parameter {
	return NewSplitParam(config, querier, d.opts)
}

// MinParam - Retrieve split table parameters using the relational database input configuration (config) and querier
func (d *BaseDbHandler) MinParam(config Config, table database.Table) database.Parameter {
	return NewMinParam(config, table, d.opts)
}

// MaxParam - Get the maximum split value parameter based on the relational database input configuration (config) and table querier (Table)
func (d *BaseDbHandler) MaxParam(config Config, table database.Table) database.Parameter {
	return NewMaxParam(config, table, d.opts)
}
