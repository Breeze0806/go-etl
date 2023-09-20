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

// Dialect oracle数据库方言
type Dialect struct{}

// Source 生成oracle数据源
func (d Dialect) Source(bs *database.BaseSource) (database.Source, error) {
	return NewSource(bs)
}

// Name 数据库方言的注册名
func (d Dialect) Name() string {
	return "oracle"
}

// Source oracle数据源
type Source struct {
	*database.BaseSource //基础数据源

	dsn string
}

// NewSource 生成oracle数据源，在配置文件错误时会报错
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

// DriverName github.com/godror/godror的驱动名
func (s *Source) DriverName() string {
	return "godror"
}

// ConnectName github.com/godror/godror的数据源连接信息
func (s *Source) ConnectName() string {
	return s.dsn
}

// Key 数据源的关键字，用于DBWrapper的复用
func (s *Source) Key() string {
	return s.dsn
}

// Table 生成oracle的表
func (s *Source) Table(b *database.BaseTable) database.Table {
	return NewTable(b)
}

// Quoted db2引用函数
func Quoted(s string) string {
	return `"` + s + `"`
}
