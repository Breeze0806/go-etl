package database

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go/time2"
)

//Source 数据源
type Source interface {
	Config() *config.Json   //配置信息
	DriverName() string     //驱动名，用于sql.Open
	ConnectName() string    //连接信息，用于sql.Open
	Table(*BaseTable) Table //获取具体表
}

type BaseSource struct {
	conf *config.Json
}

func NewBaseSource(conf *config.Json) *BaseSource {
	return &BaseSource{
		conf: conf,
	}
}

func (b *BaseSource) Config() *config.Json {
	return b.conf
}

type Config struct {
	MaxOpenConns    int            `json:"maxOpenConns"`
	MaxIdleConns    int            `json:"maxIdleConns"`
	ConnMaxIdleTime time2.Duration `json:"connMaxIdleTime"`
	ConnMaxLifetime time2.Duration `json:"connMaxLifetime"`
}

func (c *Config) GetMaxOpenConns() int {
	if c.MaxOpenConns == 0 {
		return 4
	}
	return c.MaxOpenConns
}

func (c *Config) GetMaxIdleConns() int {
	if c.MaxIdleConns == 0 {
		return 2
	}
	return c.MaxIdleConns
}
