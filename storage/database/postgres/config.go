package postgres

import (
	"encoding/json"
	"net/url"

	"github.com/Breeze0806/go-etl/config"
)

//Config postgres配置
type Config struct {
	URL      string `json:"url"`      //数据库url，包含数据库地址，数据库其他参数
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
}

//NewConfig 创建postgres配置，如果格式不符合要求，就会报错
func NewConfig(conf *config.JSON) (c *Config, err error) {
	c = &Config{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	return
}

//FormatDSN 生成数据源连接信息，url有错会报错
func (c *Config) FormatDSN() (dsn string, err error) {
	var URL *url.URL
	URL, err = url.Parse(c.URL)
	if err != nil {
		return
	}

	URL.User = url.User(c.Username)
	if c.Password != "" {
		URL.User = url.UserPassword(c.Username, c.Password)
	}

	return URL.String(), nil
}
