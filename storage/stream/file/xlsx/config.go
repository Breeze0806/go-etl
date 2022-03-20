package xlsx

import (
	"encoding/json"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/jodaTime"
	"github.com/xuri/excelize/v2"
)

type Config struct {
	Columns []Column `json:"column"`
	Sheet   string   `json:"sheet"`
}

type Column struct {
	Index    string `json:"index"`
	Type     string `json:"type"`
	Format   string `json:"format"`
	goLayout string
}

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
		return
	}
	return
}

func (c *Column) index() (i int) {
	i, _ = excelize.ColumnNameToNumber(c.Index)
	return i - 1
}

func (c *Column) layout() string {
	if c.goLayout != "" {
		return c.goLayout
	}
	c.goLayout = jodaTime.GetLayout(c.Format)
	return c.goLayout
}

func NewConfig(conf *config.JSON) (c *Config, err error) {
	c = &Config{}
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
