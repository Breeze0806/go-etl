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
	"github.com/Breeze0806/go-etl/storage/stream/file/compress"
	"github.com/Breeze0806/jodaTime"
)

//InConfig csv配置
type InConfig struct {
	Columns    []Column `json:"column"`     // 列信息
	Encoding   string   `json:"encoding"`   // 编码
	Delimiter  string   `json:"delimiter"`  // 分割符
	NullFormat string   `json:"nullFormat"` // null文本
	StartRow   int      `json:"startRow"`   // 读取开始行数，从1开始
	Comment    string   `json:"comment"`    // 注释
	Compress   string   `json:"compress"`   // 压缩
}

//NewInConfig 通过conf获取csv配置
func NewInConfig(conf *config.JSON) (c *InConfig, err error) {
	c = &InConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}

	if c.startRow() < 1 {
		return nil, fmt.Errorf("startRow is not valid")
	}

	if len([]rune(c.Delimiter)) > 1 {
		return nil, fmt.Errorf("delimiter is not valid")
	}

	if len([]rune(c.Comment)) > 1 {
		return nil, fmt.Errorf("comment is not valid")
	}

	switch c.encoding() {
	case "utf-8", "gbk":
	default:
		return nil, fmt.Errorf("encoding %v does not support", c.Encoding)
	}

	switch compress.Type(c.Compress) {
	case compress.TypeNone, compress.TypeGzip, compress.TypeZip:
	default:
		return nil, fmt.Errorf("compress %v does not support", c.Encoding)
	}

	for _, v := range c.Columns {
		if err = v.validate(); err != nil {
			return nil, err
		}
	}
	return
}

func (c *InConfig) startRow() int {
	if c.StartRow == 0 {
		return 1
	}
	return c.StartRow
}

func (c *InConfig) encoding() string {
	if c.Encoding == "" {
		return "utf-8"
	}
	return c.Encoding
}

func (c *InConfig) comma() rune {
	if len([]rune(c.Delimiter)) == 1 {
		return []rune(c.Delimiter)[0]
	}
	return rune(',')
}

func (c *InConfig) comment() rune {
	if len([]rune(c.Comment)) == 1 {
		return []rune(c.Comment)[0]
	}
	return rune(0)
}

//OutConfig csv配置
type OutConfig struct {
	Columns    []Column `json:"column"`     // 列信息
	Encoding   string   `json:"encoding"`   // 编码
	Delimiter  string   `json:"delimiter"`  // 分割符
	NullFormat string   `json:"nullFormat"` // null文本
	HasHeader  bool     `json:"hasHeader"`  // 是否有列头
	Header     []string `json:"header"`     // 列头
	Compress   string   `json:"compress"`   // 压缩
}

//NewOutConfig 通过conf获取csv配置
func NewOutConfig(conf *config.JSON) (c *OutConfig, err error) {
	c = &OutConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}

	if len([]rune(c.Delimiter)) > 1 {
		return nil, fmt.Errorf("delimiter is not valid")
	}

	switch c.encoding() {
	case "utf-8", "gbk":
	default:
		return nil, fmt.Errorf("encoding %v does not support", c.Encoding)
	}

	switch compress.Type(c.Compress) {
	case compress.TypeNone, compress.TypeGzip, compress.TypeZip:
	default:
		return nil, fmt.Errorf("compress %v does not support", c.Encoding)
	}
	for _, v := range c.Columns {
		if err = v.validate(); err != nil {
			return nil, err
		}
	}
	return
}

func (c *OutConfig) encoding() string {
	if c.Encoding == "" {
		return "utf-8"
	}
	return c.Encoding
}

func (c *OutConfig) comma() rune {
	if c.Delimiter == "" {
		return rune(',')
	}
	return []rune(c.Delimiter)[0]
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
