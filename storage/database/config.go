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
	"bytes"
	"encoding/json"
	"strings"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go/time2"
)

// Config is the basic configuration for database connections, typically used for sql.DB configurations
type Config struct {
	Pool PoolConfig `json:"pool"`
}

// NewConfig retrieves the database connection configuration 'c' from a JSON configuration
// 'err' refers to an error where the JSON configuration cannot be converted into a database connection configuration
func NewConfig(conf *config.JSON) (c *Config, err error) {
	c = &Config{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	return
}

// PoolConfig is the configuration for the database connection pool
// Generally, the maximum number of open connections should be the same as the maximum number of idle connections, otherwise it can lead to insufficient file resources due to unreleased connections
type PoolConfig struct {
	MaxOpenConns    int            `json:"maxOpenConns"`    // Maximum number of open connections
	MaxIdleConns    int            `json:"maxIdleConns"`    // Maximum number of idle connections
	ConnMaxIdleTime time2.Duration `json:"connMaxIdleTime"` // Maximum idle time for connections
	ConnMaxLifetime time2.Duration `json:"connMaxLifetime"` // Maximum lifetime for connections
}

// GetMaxOpenConns retrieves the maximum number of open connections, with a default return value of 4
func (c *PoolConfig) GetMaxOpenConns() int {
	if c.MaxOpenConns <= 0 {
		return DefaultMaxOpenConns
	}
	return c.MaxOpenConns
}

// GetMaxIdleConns retrieves the maximum number of idle connections, with a default return value of 4
func (c *PoolConfig) GetMaxIdleConns() int {
	if c.MaxIdleConns <= 0 {
		return DefaultMaxIdleConns
	}
	return c.MaxIdleConns
}

// ConfigSetter is an additional method for Table, used to set the JSON configuration file
type ConfigSetter interface {
	SetConfig(conf *config.JSON)
}

// BaseConfig is the configuration for the base table
type BaseConfig struct {
	TrimChar bool `json:"trimChar"`
}

// BaseConfigSetter is the setter for the base table configuration
type BaseConfigSetter struct {
	BaseConfig

	conf *config.JSON
}

// SetConfig sets the table configuration
func (b *BaseConfigSetter) SetConfig(conf *config.JSON) {
	b.conf = conf
	if b.conf != nil {
		json.Unmarshal([]byte(b.conf.String()), &b.BaseConfig)
	}
}

// Config retrieves the table configuration
func (b *BaseConfigSetter) Config() *config.JSON {
	return b.conf
}

// TrimStringChar removes leading and trailing spaces from a string character
func (b *BaseConfigSetter) TrimStringChar(char string) string {
	if b.TrimChar {
		return strings.TrimSpace(char)
	}
	return char
}

// TrimByteChar removes leading and trailing spaces from a byte array character
func (b *BaseConfigSetter) TrimByteChar(char []byte) []byte {
	if b.TrimChar {
		return bytes.TrimSpace(char)
	}
	return char
}
