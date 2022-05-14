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
	"os"
	"path/filepath"
	"runtime"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi"
	"github.com/pingcap/errors"
)

//Reader 数据库读取器
type Reader interface {
	spi.Reader

	ResourcesConfig() *config.JSON //插件资源配置
}

//Maker 写入生成器
type Maker interface {
	Default() (Reader, error)
	FromFile(filename string) (Reader, error)
}

//RegisterReader 通过生成数据库读取器函数new注册到读取器，返回插件资源配置文件地址，在出错时会包err
//目前未在代码中实际使用，而是通过tools/datax/build的go generate命令自动将resources/plugin.json
//中的内容放入到新生成的代码文件中，用以注册Reader
func RegisterReader(maker Maker) (pluginConfig string, err error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("fail to get filename")
	}
	path := filepath.Dir(file)
	pluginConfig = filepath.Join(path, "resources", "plugin.json")
	var reader Reader

	if reader, err = maker.FromFile(pluginConfig); err != nil {
		if !os.IsNotExist(errors.Cause(err)) {
			return
		}
		if reader, err = maker.Default(); err != nil {
			return "", errors.Wrap(err, "Default fail")
		}
	}
	name := ""
	name, err = reader.ResourcesConfig().GetString("name")
	if err != nil {
		return "", errors.Wrap(err, "GetString fail")
	}
	if name == "" {
		return "", errors.New("name is empty")
	}
	loader.RegisterReader(name, reader)
	return pluginConfig, nil
}
