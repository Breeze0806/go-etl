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

// Default Parameters
const (
	DefaultMaxOpenConns = 4
	DefaultMaxIdleConns = 4
)

// Source Data Source, containing driver information, package information, configuration files, and connection information
type Source interface {
	Config() *config.JSON   // Configuration Information
	Key() string            // Typically connection information
	DriverName() string     // Driver Name, used as the first parameter for sql.Open
	ConnectName() string    // Connection Information, used as the second parameter for sql.Open
	Table(*BaseTable) Table // Get Specific Table
}

// WithConnector Data Source with Connection, the data source prefers to call this method to generate a data connection pool DB
type WithConnector interface {
	Connector() (driver.Connector, error) // go 1.10 Get Connection
}

// NewSource Obtain the corresponding data source by the name of the database dialect
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

// BaseSource Basic data source for storing JSON configuration files
// Used to embed Source, facilitating the implementation of various database Fields
type BaseSource struct {
	conf *config.JSON
}

// NewBaseSource Generate a basic data source from the JSON configuration file conf
func NewBaseSource(conf *config.JSON) *BaseSource {
	return &BaseSource{
		conf: conf.CloneConfig(),
	}
}

// Config Configuration file for the basic data source
func (b *BaseSource) Config() *config.JSON {
	return b.conf
}
