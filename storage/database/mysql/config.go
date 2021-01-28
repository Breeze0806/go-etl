package mysql

import (
	"encoding/json"

	"github.com/Breeze0806/go-etl/config"
	"github.com/go-sql-driver/mysql"
)

//Config mysql配置，读入的时间都需要解析即parseTime=true
type Config struct {
	URL      string `json:"url"`      //数据库url，包含数据库地址，数据库其他参数
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
}

//NewConfig 创建mysql配置，如果格式不符合要求，就会报错
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
	var mysqlConf *mysql.Config
	if mysqlConf, err = mysql.ParseDSN(c.URL); err != nil {
		return
	}
	mysqlConf.User = c.Username
	mysqlConf.Passwd = c.Password
	mysqlConf.ParseTime = true
	dsn = mysqlConf.FormatDSN()
	return
}
