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
	"encoding/json"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/storage/database"
)

//Config 关系型数据读入器配置
type Config interface {
	GetUsername() string               //获取用户名
	GetPassword() string               //获取密码
	GetURL() string                    //获取连接url
	GetColumns() []Column              //获取列信息
	GetBaseTable() *database.BaseTable //获取表信息
	GetWhere() string                  //获取查询条件
	GetSplitConfig() SplitConfig       //获取切分配置
}

//Column 列信息
type Column interface {
	GetName() string //获取表名
}

//BaseColumn 基础列信息
type BaseColumn struct {
	Name string
}

//GetName 获取列名
func (b *BaseColumn) GetName() string {
	return b.Name
}

//BaseConfig 基础关系型数据读入器配置
type BaseConfig struct {
	Username   string      `json:"username"`   //用户名
	Password   string      `json:"password"`   //密码
	Column     []string    `json:"column"`     //列信息
	Connection ConnConfig  `json:"connection"` //连接信息
	Where      string      `json:"where"`      //查询条件
	Split      SplitConfig `json:"split"`      //切分键
}

//NewBaseConfig 通过json配置conf获取基础关系型数据读入器配置
func NewBaseConfig(conf *config.JSON) (c *BaseConfig, err error) {
	c = &BaseConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	return
}

//GetUsername 获取用户名
func (b *BaseConfig) GetUsername() string {
	return b.Username
}

//GetPassword 获取密码
func (b *BaseConfig) GetPassword() string {
	return b.Password
}

//GetURL 获取关系型数据库连接url
func (b *BaseConfig) GetURL() string {
	return b.Connection.URL
}

//GetColumns 获取列信息
func (b *BaseConfig) GetColumns() (columns []Column) {
	for _, v := range b.Column {
		columns = append(columns, &BaseColumn{
			Name: v,
		})
	}
	return
}

//GetBaseTable 获取表信息
func (b *BaseConfig) GetBaseTable() *database.BaseTable {
	return database.NewBaseTable(b.Connection.Table.Db, b.Connection.Table.Schema, b.Connection.Table.Name)
}

//GetWhere 获取查询条件
func (b *BaseConfig) GetWhere() string {
	return b.Where
}

//GetSplitConfig 获取切分配置
func (b *BaseConfig) GetSplitConfig() SplitConfig {
	return b.Split
}

//ConnConfig 连接配置
type ConnConfig struct {
	URL   string      `json:"url"`   //连接数据库
	Table TableConfig `json:"table"` //表配置
}

//TableConfig 表配置
type TableConfig struct {
	Db     string `json:"db"`     //库
	Schema string `json:"schema"` //模式
	Name   string `json:"name"`   //表名
}
