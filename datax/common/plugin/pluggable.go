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

// Pluggable - A pluggable interface
type Pluggable interface {
	// Plugin Developer, generally written in the plugin configuration
	Developer() (string, error)
	// Plugin Description, generally written in the plugin configuration
	Description() (string, error)
	// Plugin Name, generally written in the plugin configuration
	PluginName() (string, error)
	/* Plugin Configuration, basic configuration is as follows, the rest can be customized according to individual needs
	{
		"name" : "mysqlreader",
		"developer":"Breeze0806",
		"description":"use github.com/go-sql-driver/mysql. database/sql DB execute select sql, retrieve data from the ResultSet. warn: The more you know about the database, the less problems you encounter."
	}
	*/
	PluginConf() *config.JSON
	// Plugin Working Configuration
	PluginJobConf() *config.JSON
	// Corresponding Plugin Name (for Writer, it's Reader; for Reader, it's Writer)
	PeerPluginName() string
	// Corresponding Plugin Configuration (for Writer, it's Reader; for Reader, it's Writer)
	PeerPluginJobConf() *config.JSON
	// Set Working Plugin
	SetPluginJobConf(conf *config.JSON)
	// Set Corresponding Plugin Configuration (for Writer, it's Reader; for Reader, it's Writer)
	SetPeerPluginJobConf(conf *config.JSON)
	// Set Corresponding Plugin Name (for Writer, it's Reader; for Reader, it's Writer)
	SetPeerPluginName(name string)
	// Set Plugin Configuration
	SetPluginConf(conf *config.JSON)
	// Initialize Plugin, needs to be implemented by the implementer according to their needs
	Init(ctx context.Context) error
	// Destroy Plugin, needs to be implemented by the implementer according to their needs
	Destroy(ctx context.Context) error
}

// BasePluggable - A basic pluggable interface
// Used to assist in the implementation of various pluggable interfaces, simplifying their implementation
type BasePluggable struct {
	pluginConf        *config.JSON
	pluginJobConf     *config.JSON
	peerPluginName    string
	peerPluginJobConf *config.JSON
}

// NewBasePluggable - Creates a pluggable plugin
func NewBasePluggable() *BasePluggable {
	return &BasePluggable{}
}

// SetPluginConf - Sets the plugin configuration
func (b *BasePluggable) SetPluginConf(conf *config.JSON) {
	b.pluginConf = conf
}

// SetPluginJobConf - Sets the plugin's working configuration
func (b *BasePluggable) SetPluginJobConf(conf *config.JSON) {
	b.pluginJobConf = conf
}

// SetPeerPluginName - Sets the corresponding peer plugin name
func (b *BasePluggable) SetPeerPluginName(name string) {
	b.peerPluginName = name
}

// SetPeerPluginJobConf - Sets the corresponding peer plugin's working configuration
func (b *BasePluggable) SetPeerPluginJobConf(conf *config.JSON) {
	b.peerPluginJobConf = conf
}

// Developer - Plugin Developer, will return an error if developer is not present or not a string
func (b *BasePluggable) Developer() (string, error) {
	return b.pluginConf.GetString("developer")
}

// Description - Plugin Description, will return an error if description is not present or not a string
func (b *BasePluggable) Description() (string, error) {
	return b.pluginConf.GetString("description")
}

// PluginName - Plugin Name, will return an error if name is not present or not a string
func (b *BasePluggable) PluginName() (string, error) {
	return b.pluginConf.GetString("name")
}

// PluginConf - Plugin Configuration
func (b *BasePluggable) PluginConf() *config.JSON {
	return b.pluginConf
}

// PluginJobConf - Working Configuration
func (b *BasePluggable) PluginJobConf() *config.JSON {
	return b.pluginJobConf
}

// PeerPluginName - Corresponding Plugin Name
func (b *BasePluggable) PeerPluginName() string {
	return b.peerPluginName
}

// PeerPluginJobConf - Set personalized configuration
func (b *BasePluggable) PeerPluginJobConf() *config.JSON {
	return b.peerPluginJobConf
}
