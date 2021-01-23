package mysql

import (
	"encoding/json"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/go-sql-driver/mysql"
)

type Config struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewConfig(conf *config.Json) (c *Config, err error) {
	c = &Config{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	return
}

func (c *Config) FormatDSN() (dsn string, err error) {
	var mysqlConf *mysql.Config
	if mysqlConf, err = mysql.ParseDSN(c.URL); err != nil {
		return
	}
	fmt.Printf("%+v", mysqlConf)
	mysqlConf.User = c.Username
	mysqlConf.Passwd = c.Password
	mysqlConf.ParseTime = true
	dsn = mysqlConf.FormatDSN()
	return
}
