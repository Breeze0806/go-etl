package xlsx

import (
	"encoding/json"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/storage/stream/file/xlsx"
)

type Config struct {
	Columns []xlsx.Column `json:"column"`
	Xlsxs   []Xlsx        `json:"xlsxs"`
}

type Xlsx struct {
	Path   string   `json:"path"`
	Sheets []string `json:"sheets"`
}

func NewConfig(conf *config.JSON) (c *Config, err error) {
	c = &Config{}
	if err = json.Unmarshal([]byte(conf.String()), c); err != nil {
		return
	}
	return
}
