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

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/storage/database"
)

type Config interface {
	GetUsername() string
	GetPassword() string
	GetURL() string
	GetColumns() []Column
	GetBaseTable() *database.BaseTable
	GetWhere() string
}

type Column interface {
	GetName() string
}

type BaseColumn struct {
	Name string
}

func (b *BaseColumn) GetName() string {
	return b.Name
}

type BaseConfig struct {
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	Column     []string   `json:"column"`
	Connection ConnConfig `json:"connection"`
	Where      string     `json:"where"`
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

func (b *BaseConfig) GetColumns() (columns []Column) {
	for _, v := range b.Column {
		columns = append(columns, &BaseColumn{
			Name: v,
		})
	}
	return
}

func (b *BaseConfig) GetBaseTable() *database.BaseTable {
	return database.NewBaseTable(b.Connection.Table.Db, b.Connection.Table.Schema, b.Connection.Table.Name)
}

func (b *BaseConfig) GetWhere() string {
	return b.Where
}

type ConnConfig struct {
	URL   string      `json:"url"`
	Table TableConfig `json:"table"`
}

type TableConfig struct {
	Db     string `json:"db"`
	Schema string `json:"schema"`
	Name   string `json:"name"`
}
