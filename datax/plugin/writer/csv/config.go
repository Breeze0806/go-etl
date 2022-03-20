package csv

import (
	"encoding/json"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/file"
	"github.com/Breeze0806/go-etl/storage/stream/file/csv"
)

type CsvConfig struct {
	csv.Config
	file.BaseConfig
}

type JobConfig struct {
	CsvConfig

	Path []string `json:"path"`
}

func NewJobConfig(conf *config.JSON) (*JobConfig, error) {
	c := &JobConfig{}
	if err := json.Unmarshal([]byte(conf.String()), c); err != nil {
		return nil, err
	}
	return c, nil
}
