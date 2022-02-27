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
	Index  string `json:"index"`
	Type   string `json:"type"`
	Format string `json:"format"`
}

func (c *Column) validate() (err error) {
	switch element.ColumnType(c.Type) {
	case element.TypeBool, element.TypeBigInt,
		element.TypeDecimal, element.TypeString,
		element.TypeTime:
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
	return jodaTime.GetLayout(c.Format)
}

func NewConfig(conf *config.JSON) (c *Config, err error) {
	c = &Config{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}

	for _, v := range c.Columns {
		if err = v.validate(); err != nil {
			return nil, err
		}
	}
	return
}
