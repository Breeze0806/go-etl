package mysql

import (
	"encoding/json"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go/time2"
)

var (
	defalutBatchSize    = 1000
	defalutBatchTimeout = 1 * time.Second
)

type paramConfig struct {
	Username     string         `json:"username"`
	Password     string         `json:"password"`
	Column       []string       `json:"column"`
	Connection   connConfig     `json:"connection"`
	WriteMode    string         `json:"writeMode"`
	BatchSize    int            `json:"batchSize"`
	BatchTimeout time2.Duration `json:"batchTimeout"`
}

type connConfig struct {
	URL   string      `json:"url"`
	Table tableConfig `json:"table"`
}

type tableConfig struct {
	Db   string `json:"db"`
	Name string `json:"name"`
}

func newParamConfig(conf *config.JSON) (c *paramConfig, err error) {
	c = &paramConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	return
}

func (p *paramConfig) getBatchSize() int {
	if p.BatchSize <= 0 {
		return defalutBatchSize
	}
	return p.BatchSize
}

func (p *paramConfig) getBatchTimeout() time.Duration {
	if p.BatchTimeout.Duration == 0 {
		return defalutBatchTimeout
	}
	return p.BatchTimeout.Duration
}
