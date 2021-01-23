package database

import (
	"encoding/json"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go/time2"
)

type Config struct {
	Pool PoolConfig `json:"pool"`
}

func NewConfig(conf *config.Json) (c *Config, err error) {
	c = &Config{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal(%v) error: %v", conf.String(), err)
	}
	return
}

type PoolConfig struct {
	MaxOpenConns    int            `json:"maxOpenConns"`
	MaxIdleConns    int            `json:"maxIdleConns"`
	ConnMaxIdleTime time2.Duration `json:"connMaxIdleTime"`
	ConnMaxLifetime time2.Duration `json:"connMaxLifetime"`
}

func (c *PoolConfig) GetMaxOpenConns() int {
	if c.MaxOpenConns <= 0 {
		return DefaultMaxOpenConns
	}
	return c.MaxOpenConns
}

func (c *PoolConfig) GetMaxIdleConns() int {
	if c.MaxIdleConns <= 0 {
		return DefaultMaxIdleConns
	}
	return c.MaxIdleConns
}
