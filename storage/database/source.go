package database

import (
	"fmt"

	"github.com/Breeze0806/go-etl/config"
)

const (
	DefaultMaxOpenConns = 4
	DefaultMaxIdleConns = 4
)

//Source 数据源
type Source interface {
	Config() *config.Json   //配置信息
	Key() string            //dbMap Key
	DriverName() string     //驱动名，用于sql.Open
	ConnectName() string    //连接信息，用于sql.Open
	Table(*BaseTable) Table //获取具体表
}

func NewSource(name string, conf *config.Json) (source Source, err error) {
	d, ok := dialects.dialect(name)
	if !ok {
		return nil, fmt.Errorf("dialect %v does not exsit", name)
	}
	source, err = d.Source(NewBaseSource(conf))
	if err != nil {
		return nil, fmt.Errorf("dialect %v Source() err: %v", name, err)
	}
	return
}

type BaseSource struct {
	conf *config.Json
}

func NewBaseSource(conf *config.Json) *BaseSource {
	return &BaseSource{
		conf: conf.CloneConfig(),
	}
}

func (b *BaseSource) Config() *config.Json {
	return b.conf
}
