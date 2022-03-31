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
	"encoding/json"
	"time"

	"github.com/Breeze0806/go-etl/config"
	rdbmreader "github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/Breeze0806/go/time2"
)

var (
	defalutBatchSize    = 1000
	defalutBatchTimeout = 1 * time.Second
)

type Config interface {
	GetUsername() string
	GetPassword() string
	GetURL() string
	GetColumns() []rdbmreader.Column
	GetBaseTable() *database.BaseTable
	GetWriteMode() string
	GetBatchSize() int              //单次批量写入数
	GetBatchTimeout() time.Duration //单次批量写入超时时间
}

type BaseConfig struct {
	Username     string                `json:"username"`
	Password     string                `json:"password"`
	Column       []string              `json:"column"`
	Connection   rdbmreader.ConnConfig `json:"connection"`
	WriteMode    string                `json:"writeMode"`
	BatchSize    int                   `json:"batchSize"`
	BatchTimeout time2.Duration        `json:"batchTimeout"`
}

func NewBaseConfig(conf *config.JSON) (c *BaseConfig, err error) {
	c = &BaseConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	return
}

func (b *BaseConfig) GetUsername() string {
	return b.Username
}

func (b *BaseConfig) GetPassword() string {
	return b.Password
}

func (b *BaseConfig) GetURL() string {
	return b.Connection.URL
}

func (b *BaseConfig) GetColumns() (columns []rdbmreader.Column) {
	for _, v := range b.Column {
		columns = append(columns, &rdbmreader.BaseColumn{
			Name: v,
		})
	}
	return
}

func (b *BaseConfig) GetBaseTable() *database.BaseTable {
	return database.NewBaseTable(b.Connection.Table.Db, b.Connection.Table.Schema, b.Connection.Table.Name)
}

func (b *BaseConfig) GetWriteMode() string {
	return b.WriteMode
}

func (b *BaseConfig) GetBatchTimeout() time.Duration {
	if b.BatchTimeout.Duration == 0 {
		return defalutBatchTimeout
	}
	return b.BatchTimeout.Duration
}

func (b *BaseConfig) GetBatchSize() int {
	if b.BatchSize == 0 {
		return defalutBatchSize
	}

	return b.BatchSize
}
