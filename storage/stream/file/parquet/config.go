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

package parquet

import (
	"encoding/json"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

type InConfig struct {
	Columns []Column `json:"column"`
}
type OutConfig struct {
	Columns []Column `json:"column"`
}

func NewOutConfig(conf *config.JSON) (c *OutConfig, err error) {
	c = &OutConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	for _, v := range c.Columns {
		if err = v.validate(); err != nil {
			return nil, err
		}
	}
	return c, nil
}

type Column struct {
	Type string
	Name string
}

func NewInConfig(conf *config.JSON) (c *InConfig, err error) {
	c = &InConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	for _, v := range c.Columns {
		if err = v.validate(); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Column) validate() error {
	switch element.ColumnType(c.Type) {
	case element.TypeBool, element.TypeBigInt,
		element.TypeDecimal, element.TypeString,
		element.TypeTime:
	default:
		return fmt.Errorf("unsupported column type: %s", c.Type)
	}
	return nil
}
