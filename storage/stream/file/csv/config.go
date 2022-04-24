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
	"fmt"
	"strconv"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/jodaTime"
)

//Config csv配置
type Config struct {
	Columns    []Column `json:"column"`     // 列信息
	Encoding   string   `json:"encoding"`   // 编码
	Delimiter  string   `json:"delimiter"`  // 分割符
	NullFormat string   `json:"nullFormat"` // null文本
}

//Column 列信息
type Column struct {
	Index    string `json:"index"`  // 索引 从1开始，代表第几列
	Type     string `json:"type"`   // 类型 bool bigInt decimal string time
	Format   string `json:"format"` // joda时间格式
	indexNum int
	goLayout string
}

//validate 校验
func (c *Column) validate() (err error) {
	switch element.ColumnType(c.Type) {
	case element.TypeBool, element.TypeBigInt,
		element.TypeDecimal, element.TypeString:
	case element.TypeTime:
		if c.Format == "" {
			return fmt.Errorf("type %v format %v is empty", c.Type, c.Format)
		}
	default:
		return fmt.Errorf("type %v is not valid", c.Type)
	}
	var i int
	if i, err = strconv.Atoi(c.Index); err != nil {
		return
	}
	if i < 1 {
		return fmt.Errorf("index is less than 1")
	}

	return
}

//index 列索引
func (c *Column) index() (i int) {
	if c.indexNum > 0 {
		return c.indexNum - 1
	}
	c.indexNum, _ = strconv.Atoi(c.Index)
	return c.indexNum - 1
}

//layout 变为golang 时间格式
func (c *Column) layout() string {
	if c.goLayout != "" {
		return c.goLayout
	}
	c.goLayout = jodaTime.GetLayout(c.Format)
	return c.goLayout
}

//NewConfig 通过conf获取csv配置
func NewConfig(conf *config.JSON) (c *Config, err error) {
	c = &Config{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}

	if c.Delimiter == "" {
		c.Delimiter = ","
	}

	if len(c.Delimiter) != 1 {
		return nil, fmt.Errorf("Delimiter is not valid")
	}

	if c.Encoding == "" {
		c.Encoding = "utf-8"
	}

	switch c.Encoding {
	case "utf-8":
	default:
		return nil, fmt.Errorf("encoding %v does not support", c.Encoding)
	}

	for _, v := range c.Columns {
		if err = v.validate(); err != nil {
			return nil, err
		}
	}
	return
}
