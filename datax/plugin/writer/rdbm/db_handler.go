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
	rdbmreader "github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
	"github.com/Breeze0806/go-etl/storage/database"
)

type DbHandler interface {
	Execer(name string, conf *config.JSON) (Execer, error)
	Config(conf *config.JSON) (Config, error)
	TableParam(config Config, execer Execer) database.Parameter
}

type BaseDbHandler struct {
	newExecer func(name string, conf *config.JSON) (Execer, error)
	opts      *sql.TxOptions
}

func NewBaseDbHandler(newExecer func(name string, conf *config.JSON) (Execer, error), opts *sql.TxOptions) *BaseDbHandler {
	return &BaseDbHandler{
		newExecer: newExecer,
		opts:      opts,
	}
}

func (d *BaseDbHandler) Execer(name string, conf *config.JSON) (Execer, error) {
	return d.newExecer(name, conf)
}

func (d *BaseDbHandler) Config(conf *config.JSON) (Config, error) {
	return NewBaseConfig(conf)
}

func (d *BaseDbHandler) TableParam(config Config, execer Execer) database.Parameter {
	return rdbmreader.NewTableParam(config, execer, d.opts)
}
