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

package csv

import (
	"encoding/json"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/file"

	// csv storage - Storage or handling of data in the CSV format.
	"github.com/Breeze0806/go-etl/storage/stream/file/csv"
)

// SingleConfig csv single input setting - Configuration settings for a single CSV input.
type SingleConfig struct {
	csv.OutConfig
	file.BaseConfig
}

// Config csv input configuration - Configuration settings for reading or processing data from a CSV file.
type Config struct {
	SingleConfig

	Path []string `json:"path"`
}

// NewConfig - A function or method that retrieves CSV input configuration settings from a JSON configuration file (conf).
func NewConfig(conf *config.JSON) (*Config, error) {
	c := &Config{}
	if err := json.Unmarshal([]byte(conf.String()), c); err != nil {
		return nil, err
	}
	return c, nil
}
