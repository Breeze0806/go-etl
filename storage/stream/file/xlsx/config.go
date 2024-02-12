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

// InConfig represents the input XLSX configuration
type InConfig struct {
	Columns    []Column `json:"column"`     // Column information array
	Sheet      string   `json:"sheet"`      // Sheet name
	NullFormat string   `json:"nullFormat"` // Null text
	StartRow   int      `json:"startRow"`   // Starting row for reading, starting from the 1st row
}

// NewInConfig creates a new input XLSX configuration based on the JSON configuration conf
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

func (c *InConfig) startRow() int {
	if c.StartRow == 0 {
		return 1
	}
	return c.StartRow
}

// OutConfig represents the output XLSX configuration
type OutConfig struct {
	Columns    []Column `json:"column"`     // Column information array
	Sheets     []string `json:"sheets"`     // Sheet name
	NullFormat string   `json:"nullFormat"` // Null text
	HasHeader  bool     `json:"hasHeader"`  // Whether there is a column header
	Header     []string `json:"header"`     // Column header
	SheetRow   int      `json:"sheetRow"`   // Maximum number of rows in the sheet
}

// NewOutConfig creates a new output XLSX configuration based on the JSON configuration conf
func NewOutConfig(conf *config.JSON) (c *OutConfig, err error) {
	c = &OutConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	if len(c.Sheets) == 0 {
		return nil, fmt.Errorf("sheets should not be empty")
	}

	if c.SheetRow > excelize.TotalRows || c.SheetRow < 0 {
		return nil, fmt.Errorf("sheetRow should be not less than %v and positive", excelize.TotalRows)
	}

	for _, v := range c.Columns {
		if err = v.validate(); err != nil {
			return nil, err
		}
	}
	return
}

func (c *OutConfig) sheetRow() int {
	if c.SheetRow == 0 {
		return excelize.TotalRows
	}
	return c.SheetRow
}

// Column represents column information
type Column struct {
	Index    string `json:"index"`  // Column index, e.g., A, B, C, ..., AA, ...
	Type     string `json:"type"`   // Type (bool, bigInt, decimal, string, time)
	Format   string `json:"format"` // Joda time format
	indexNum int
	goLayout string
}

// Validate performs validation
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

// index represents the column index
func (c *Column) index() (i int) {
	if c.indexNum > 0 {
		return c.indexNum - 1
	}
	c.indexNum, _ = excelize.ColumnNameToNumber(c.Index)
	return c.indexNum - 1
}

// layout converts to the Go time format
func (c *Column) layout() string {
	if c.goLayout != "" {
		return c.goLayout
	}
	c.goLayout = jodaTime.GetLayout(c.Format)
	return c.goLayout
}
