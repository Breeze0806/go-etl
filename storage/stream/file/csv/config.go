package csv

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/jodaTime"
)

type Config struct {
	Columns   []Column `json:"column"`
	Encoding  string   `json:"encoding"`
	Delimiter string   `json:"delimiter"`
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

	if _, err = strconv.Atoi(c.Index); err != nil {
		return
	}
	return
}

func (c *Column) index() (i int) {
	i, _ = strconv.Atoi(c.Index)
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

	if c.Delimiter == "" {
		c.Delimiter = ","
	}

	if len(c.Delimiter) != 1 {
		return nil, fmt.Errorf("Delimiter is not valid")
	}

	switch c.Encoding {
	case "", "utf-8":
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
