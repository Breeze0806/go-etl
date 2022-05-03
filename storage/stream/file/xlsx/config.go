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
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/jodaTime"
	"github.com/xuri/excelize/v2"
)

//InConfig 输入xlsx配置
type InConfig struct {
	Columns    []Column `json:"column"`     //列信息数组
	Sheet      string   `json:"sheet"`      //表格名
	NullFormat string   `json:"nullFormat"` //null文本
}

//OutConfig 输出xlsx配置
type OutConfig struct {
	Columns    []Column `json:"column"`     //列信息数组
	Sheets     []string `json:"sheets"`     //表格名
	NullFormat string   `json:"nullFormat"` //null文本
}

//Column 列信息
type Column struct {
	Index    string `json:"index"`  //列索引，A,B,C....AA.....
	Type     string `json:"type"`   //类型 类型 bool bigInt decimal string time
	Format   string `json:"format"` //joda时间格式
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

	if _, err = excelize.ColumnNameToNumber(c.Index); err != nil {
		return fmt.Errorf("index %v err: %v", c.Type, err)
	}
	return
}

//index 列索引
func (c *Column) index() (i int) {
	if c.indexNum > 0 {
		return c.indexNum - 1
	}
	c.indexNum, _ = excelize.ColumnNameToNumber(c.Index)
	return c.indexNum - 1
}

//layout go时间格式
func (c *Column) layout() string {
	if c.goLayout != "" {
		return c.goLayout
	}
	c.goLayout = jodaTime.GetLayout(c.Format)
	return c.goLayout
}

//NewInConfig 新建以json配置conf的输入xlsx配置
func NewInConfig(conf *config.JSON) (c *InConfig, err error) {
	c = &InConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}

	if c.Sheet == "" {
		return nil, fmt.Errorf("sheet should not be empty")
	}

	for _, v := range c.Columns {
		if err = v.validate(); err != nil {
			return nil, err
		}
	}
	return
}

//NewOutConfig 新建以json配置conf的输出xlsx配置
func NewOutConfig(conf *config.JSON) (c *OutConfig, err error) {
	c = &OutConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	if len(c.Sheets) == 0 {
		return nil, fmt.Errorf("sheets should not be empty")
	}

	for _, v := range c.Columns {
		if err = v.validate(); err != nil {
			return nil, err
		}
	}
	return
}
