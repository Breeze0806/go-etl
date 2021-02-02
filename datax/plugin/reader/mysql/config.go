package mysql

import (
	"encoding/json"

	"github.com/Breeze0806/go-etl/config"
)

type paramConfig struct {
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	Column     []string   `json:"column"`
	Connection connConfig `json:"connection"`
	Where      string     `json:"where"`
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
