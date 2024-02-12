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

package reader

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi"
	"github.com/pingcap/errors"
)

// Reader is a database reader.
type Reader interface {
	spi.Reader

	ResourcesConfig() *config.JSON // Plugin resource configuration.
}

// Maker is a write generator.
type Maker interface {
	Default() (Reader, error)
}

// RegisterReader registers a new database reader function and returns the address of the plugin resource configuration file. If an error occurs, it will be wrapped in err.
// Currently, it is not used directly in the code, but instead, the go generate command in tools/datax/build automatically inserts the content from resources/plugin.json into a newly generated code file to register the Reader.
// This approach is used to register the Reader without manually editing the code.
func RegisterReader(maker Maker) (err error) {
	var reader Reader
	if reader, err = maker.Default(); err != nil {
		return errors.Wrap(err, "Default fail")
	}

	name := ""
	name, err = reader.ResourcesConfig().GetString("name")
	if err != nil {
		return errors.Wrap(err, "GetString fail")
	}
	if name == "" {
		return errors.New("name is empty")
	}
	loader.RegisterReader(name, reader)
	return nil
}
