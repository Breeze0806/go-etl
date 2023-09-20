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

package plugin

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
)

// Pluggable  可插件化接口
type Pluggable interface {
	//插件开发者,一般写入插件配置中
	Developer() (string, error)
	//插件描述,一般写入插件配置中
	Description() (string, error)
	//插件名称,一般写入插件配置中
	PluginName() (string, error)
	/*插件配置，基础配置如下，其余可以根据个性化定制
	{
		"name" : "mysqlreader",
		"developer":"Breeze0806",
		"description":"use github.com/go-sql-driver/mysql. database/sql DB execute select sql, retrieve data from the ResultSet. warn: The more you know about the database, the less problems you encounter."
	}
	*/
	PluginConf() *config.JSON
	//插件工作配置
	PluginJobConf() *config.JSON
	//对应插件名（对于Writer来说就是Reader，对应Reader来说就是Wirter）
	PeerPluginName() string
	//对应插件配置（对于Writer来说就是Reader，对应Reader来说就是Wirter）
	PeerPluginJobConf() *config.JSON
	//设置工作插件
	SetPluginJobConf(conf *config.JSON)
	//设置对应插件配置（对于Writer来说就是Reader，对应Reader来说就是Wirter）
	SetPeerPluginJobConf(conf *config.JSON)
	//设置对应插件名（对于Writer来说就是Reader，对应Reader来说就是Wirter）
	SetPeerPluginName(name string)
	//设置插件配置
	SetPluginConf(conf *config.JSON)
	//初始化插件，需要实现者个性化实现
	Init(ctx context.Context) error
	//销毁插件，需要实现者个性化实现
	Destroy(ctx context.Context) error
}

// BasePluggable 基础可插件化
// 用于辅助各类可插件化接口实现，简化其实现
type BasePluggable struct {
	pluginConf        *config.JSON
	pluginJobConf     *config.JSON
	peerPluginName    string
	peerPluginJobConf *config.JSON
}

// NewBasePluggable 创建可插件化插件
func NewBasePluggable() *BasePluggable {
	return &BasePluggable{}
}

// SetPluginConf 设置插件配置
func (b *BasePluggable) SetPluginConf(conf *config.JSON) {
	b.pluginConf = conf
}

// SetPluginJobConf 设置插件工作配置
func (b *BasePluggable) SetPluginJobConf(conf *config.JSON) {
	b.pluginJobConf = conf
}

// SetPeerPluginName 设置对应工作名
func (b *BasePluggable) SetPeerPluginName(name string) {
	b.peerPluginName = name
}

// SetPeerPluginJobConf 设置对应工作配置
func (b *BasePluggable) SetPeerPluginJobConf(conf *config.JSON) {
	b.peerPluginJobConf = conf
}

// Developer 插件开发者,当developer不存在或者不是字符串时会返回错误
func (b *BasePluggable) Developer() (string, error) {
	return b.pluginConf.GetString("developer")
}

// Description 插件描述,当description不存在或者不是字符串时会返回错误
func (b *BasePluggable) Description() (string, error) {
	return b.pluginConf.GetString("description")
}

// PluginName 插件名称,当name不存在或者不是字符串时会返回错误
func (b *BasePluggable) PluginName() (string, error) {
	return b.pluginConf.GetString("name")
}

// PluginConf 插件配置
func (b *BasePluggable) PluginConf() *config.JSON {
	return b.pluginConf
}

// PluginJobConf 工作配置
func (b *BasePluggable) PluginJobConf() *config.JSON {
	return b.pluginJobConf
}

// PeerPluginName 对应插件名称
func (b *BasePluggable) PeerPluginName() string {
	return b.peerPluginName
}

// PeerPluginJobConf 设置个性化配置
func (b *BasePluggable) PeerPluginJobConf() *config.JSON {
	return b.peerPluginJobConf
}
