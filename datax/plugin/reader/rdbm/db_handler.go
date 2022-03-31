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

package rdbm

import (
	"database/sql"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/storage/database"
)

type DbHandler interface {
	Querier(name string, conf *config.JSON) (Querier, error)
	Config(conf *config.JSON) (Config, error)
	TableParam(config Config, querier Querier) database.Parameter
}

type BaseDbHandler struct {
	newQuerier func(name string, conf *config.JSON) (Querier, error)
	opts       *sql.TxOptions
}

func NewBaseDbHandler(newQuerier func(name string, conf *config.JSON) (Querier, error), opts *sql.TxOptions) *BaseDbHandler {
	return &BaseDbHandler{
		newQuerier: newQuerier,
		opts:       opts,
	}
}

func (d *BaseDbHandler) Querier(name string, conf *config.JSON) (Querier, error) {
	return d.newQuerier(name, conf)
}

func (d *BaseDbHandler) Config(conf *config.JSON) (Config, error) {
	return NewBaseConfig(conf)
}

func (d *BaseDbHandler) TableParam(config Config, querier Querier) database.Parameter {
	return NewTableParam(config, querier, d.opts)
}
