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

package sqlite3

import (
	"encoding/json"
	"github.com/Breeze0806/go-etl/config"
	"github.com/pingcap/errors"
	"os"
)

// Config is the Sqlite3 configuration
type Config struct {
	URL string `json:"url"` // Database URL, including the database address and other database parameters
}

// NewConfig creates a Sqlite3 configuration and will report an error if the format does not meet the requirements
func NewConfig(conf *config.JSON) (c *Config, err error) {
	c = &Config{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	return
}

// FormatDSN generates data source connection information and will report an error if the URL is incorrect
func (c *Config) FormatDSN() (dsn string, err error) {
	if c.isValidPath(c.URL) {
		err = errors.New("configure a url that is not a valid file path")
		return
	}
	return c.URL, nil
}

// IsValidPath to check whether a given path points to an existing file or directory.
func (c *Config) isValidPath(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	return true
}
