package csv

import (
	"encoding/json"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/file"

	//csv storage
	"github.com/Breeze0806/go-etl/storage/stream/file/csv"
)

type CsvConfig struct {
	csv.Config
	file.BaseConfig
}

type Config struct {
	CsvConfig

	Path []string `json:"path"`
}

func NewConfig(conf *config.JSON) (*Config, error) {
	c := &Config{}
	if err := json.Unmarshal([]byte(conf.String()), c); err != nil {
		return nil, err
	}
	return c, nil
}
