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

//Writer 写入器
type Writer interface {
	spi.Writer

	//资源插件配置
	ResourcesConfig() *config.JSON
}

//Maker 写入生成器
type Maker interface {
	Default() (Writer, error)
}

//RegisterWriter 注册创建函数new写入器,返回的是资源插件配置文件地州，出错时会返回error
//目前未在代码中实际使用，而是通过tools/datax/build的tools/datax/build命令自动将resources/plugin.json
//中的内容放入到新生成的代码文件中，用以注册Writer
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
