package database

import (
	"encoding/json"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go/time2"
)

//Config 数据库连接基础配置，一般用于sql.DB的配置
type Config struct {
	Pool PoolConfig `json:"pool"`
}

//NewConfig 从Json配置中获取数据库连接配置c
//err是指Json配置无法转化为数据库连接配置
func NewConfig(conf *config.JSON) (c *Config, err error) {
	c = &Config{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal(%v) error: %v", conf.String(), err)
	}
	return
}

//PoolConfig 数据库连接池配置
//一般让最大打开连接数和最大空闲时连接数一致，否则会导致释放连接不及导致文件数不足
type PoolConfig struct {
	MaxOpenConns    int            `json:"maxOpenConns"`    //最大打开连接数
	MaxIdleConns    int            `json:"maxIdleConns"`    //最大空闲时连接数
	ConnMaxIdleTime time2.Duration `json:"connMaxIdleTime"` //最大连接空闲时间
	ConnMaxLifetime time2.Duration `json:"connMaxLifetime"` //最大连接存活时间
}

//GetMaxOpenConns 获取最大连接数，默认返回值为4
func (c *PoolConfig) GetMaxOpenConns() int {
	if c.MaxOpenConns <= 0 {
		return DefaultMaxOpenConns
	}
	return c.MaxOpenConns
}

//GetMaxIdleConns 获取空闲时最大连接数，默认返回为4
func (c *PoolConfig) GetMaxIdleConns() int {
	if c.MaxIdleConns <= 0 {
		return DefaultMaxIdleConns
	}
	return c.MaxIdleConns
}
