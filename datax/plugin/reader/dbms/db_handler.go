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

// DbHandler 数据库句柄
type DbHandler interface {
	Querier(name string, conf *config.JSON) (Querier, error)         //通过数据库名name和json配置conf获取查询器
	Config(conf *config.JSON) (Config, error)                        //通过json配置conf获取关系型数据库输入配置
	TableParam(config Config, querier Querier) database.Parameter    //通过关系型数据库输入配置config和查询器querier获取表参数
	SplitParam(config Config, querier Querier) database.Parameter    //通过关系型数据库输入配置config和查询器querier获取切分表参数
	MinParam(config Config, table database.Table) database.Parameter //通过关系型数据库输入配置config和表Table获取切分最小值参数
	MaxParam(config Config, table database.Table) database.Parameter //通过关系型数据库输入配置config和表询器Table获取切分最大值参数
}

// BaseDbHandler 基础数据库句柄
type BaseDbHandler struct {
	newQuerier func(name string, conf *config.JSON) (Querier, error)
	opts       *sql.TxOptions
}

// NewBaseDbHandler 通过获取查询器函数newQuerier和事务选项opts获取基础数据库句柄
func NewBaseDbHandler(newQuerier func(name string, conf *config.JSON) (Querier, error), opts *sql.TxOptions) *BaseDbHandler {
	return &BaseDbHandler{
		newQuerier: newQuerier,
		opts:       opts,
	}
}

// Querier 通过数据库名name和json配置conf获取查询器
func (d *BaseDbHandler) Querier(name string, conf *config.JSON) (Querier, error) {
	return d.newQuerier(name, conf)
}

// Config 通过json配置conf获取关系型数据库输入配置
func (d *BaseDbHandler) Config(conf *config.JSON) (Config, error) {
	return NewBaseConfig(conf)
}

// TableParam 通过关系型数据库输入配置config和查询器querier获取表参数
func (d *BaseDbHandler) TableParam(config Config, querier Querier) database.Parameter {
	return NewTableParam(config, querier, d.opts)
}

// SplitParam 通过关系型数据库输入配置config和表Table获取切分最小值参数
func (d *BaseDbHandler) SplitParam(config Config, querier Querier) database.Parameter {
	return NewSplitParam(config, querier, d.opts)
}

// MinParam 通过关系型数据库输入配置config和查询器querier获取切分表参数
func (d *BaseDbHandler) MinParam(config Config, table database.Table) database.Parameter {
	return NewMinParam(config, table, d.opts)
}

// MaxParam 通过关系型数据库输入配置config和表询器Table获取切分最大值参数
func (d *BaseDbHandler) MaxParam(config Config, table database.Table) database.Parameter {
	return NewMaxParam(config, table, d.opts)
}
