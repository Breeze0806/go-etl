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

package writer

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi"
	"github.com/pingcap/errors"
)

// Writer
type Writer interface {
	spi.Writer

	// Resource Plugin Configuration
	ResourcesConfig() *config.JSON
}

// Maker Writer Generator
type Maker interface {
	Default() (Writer, error)
}

// RegisterWriter Registers the creation function for a new writer, returning the path to the resource plugin configuration file. If an error occurs, an error is returned.
// This is currently not used in the code directly, but the content of resources/plugin.json is automatically placed into the newly generated code file through the tools/datax/build command, for the purpose of registering the Writer.
// The content in resources/plugin.json is used to register the Writer by being automatically placed into the newly generated code file through the tools/datax/build command, without being explicitly used in the code.
func RegisterWriter(maker Maker) (err error) {
	var writer Writer

	if writer, err = maker.Default(); err != nil {
		return err
	}

	name, err := writer.ResourcesConfig().GetString("name")
	if err != nil {
		return err
	}
	if name == "" {
		return errors.New("name is empty")
	}
	loader.RegisterWriter(name, writer)
	return nil
}
