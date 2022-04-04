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

package xlsx

import (
	"encoding/json"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/file"
	"github.com/Breeze0806/go-etl/storage/stream/file/xlsx"
)

//Config xlsx输出配置
type Config struct {
	file.BaseConfig

	Columns []xlsx.Column `json:"column"`
	Xlsxs   []Xlsx        `json:"xlsxs"`
}

//Xlsx xlsx内容配置
type Xlsx struct {
	Path   string   `json:"path"`
	Sheets []string `json:"sheets"`
}

//NewConfig 通过json配置conf创建xlsx输出配置
func NewConfig(conf *config.JSON) (c *Config, err error) {
	c = &Config{}
	if err = json.Unmarshal([]byte(conf.String()), c); err != nil {
		return
	}
	return
}
