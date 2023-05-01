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

//DbHandler 数据库执行句柄封装
type DbHandler interface {
	Execer(name string, conf *config.JSON) (Execer, error)      //通过数据库名name和配置获取执行器
	Config(conf *config.JSON) (Config, error)                   //通过配置获取关系型数据库配置
	TableParam(config Config, execer Execer) database.Parameter //通过关系型数据库配置和执行器获取表参数
}

//BaseDbHandler 基础数据库执行句柄封装
type BaseDbHandler struct {
	newExecer func(name string, conf *config.JSON) (Execer, error)
	opts      *sql.TxOptions
}

//NewBaseDbHandler 通过获取执行器函数newExecer和数据库事务执行选项opts创建数据库执行句柄封装
func NewBaseDbHandler(newExecer func(name string, conf *config.JSON) (Execer, error), opts *sql.TxOptions) *BaseDbHandler {
	return &BaseDbHandler{
		newExecer: newExecer,
		opts:      opts,
	}
}

//Execer 通过数据库名name和配置获取执行器
func (d *BaseDbHandler) Execer(name string, conf *config.JSON) (Execer, error) {
	return d.newExecer(name, conf)
}

//Config 通过配置获取关系型数据库配置
func (d *BaseDbHandler) Config(conf *config.JSON) (Config, error) {
	return NewBaseConfig(conf)
}

//TableParam 通过关系型数据库配置和执行器获取表参数
func (d *BaseDbHandler) TableParam(config Config, execer Execer) database.Parameter {
	return dbmsreader.NewTableParam(config, execer, d.opts)
}
