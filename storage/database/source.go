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

package database

import (
	"database/sql/driver"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
)

// 默认参数
const (
	DefaultMaxOpenConns = 4
	DefaultMaxIdleConns = 4
)

// Source 数据源,包含驱动信息，包信息，配置文件以及连接信息
type Source interface {
	Config() *config.JSON   //配置信息
	Key() string            //一般是连接信息
	DriverName() string     //驱动名，用于sql.Open的第1个参数
	ConnectName() string    //连接信息，用于sql.Open的第2个参数
	Table(*BaseTable) Table //获取具体表
}

// WithConnector 带有连接的数据源, 数据源优先调用该方法生成数据连接池DB
type WithConnector interface {
	Connector() (driver.Connector, error) //go 1.10 获取连接
}

// NewSource 通过数据库方言的名字获取对应数据源
func NewSource(name string, conf *config.JSON) (source Source, err error) {
	d, ok := dialects.dialect(name)
	if !ok {
		return nil, fmt.Errorf("dialect %v does not exsit", name)
	}
	source, err = d.Source(NewBaseSource(conf))
	if err != nil {
		return nil, fmt.Errorf("dialect %v Source() err: %v", name, err)
	}
	return
}

// BaseSource 基础数据源，用于存储json配置文件
// 用于嵌入Source，方便实现各个数据库的Field
type BaseSource struct {
	conf *config.JSON
}

// NewBaseSource 通过json配置文件conf生成基础数据源
func NewBaseSource(conf *config.JSON) *BaseSource {
	return &BaseSource{
		conf: conf.CloneConfig(),
	}
}

// Config 基础数据源的配置文件
func (b *BaseSource) Config() *config.JSON {
	return b.conf
}
