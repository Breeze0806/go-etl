package database

import (
	"fmt"

	"github.com/Breeze0806/go-etl/config"
)

//默认参数
const (
	DefaultMaxOpenConns = 4
	DefaultMaxIdleConns = 4
)

//Source 数据源,包含驱动信息，包信息，配置文件以及连接信息
type Source interface {
	Config() *config.JSON   //配置信息
	Key() string            //dbMap Key
	DriverName() string     //驱动名，用于sql.Open
	ConnectName() string    //连接信息，用于sql.Open
	Table(*BaseTable) Table //获取具体表
}

//NewSource 通过数据库方言的名字获取对应数据源
func NewSource(name string, conf *config.JSON) (source Source, err error) {
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

//BaseSource 基础数据源，用于存储json配置文件
//用于嵌入Source，方便实现各个数据库的Field
type BaseSource struct {
	conf *config.JSON
}

//NewBaseSource 通过json配置文件conf生成基础数据源
func NewBaseSource(conf *config.JSON) *BaseSource {
	return &BaseSource{
		conf: conf.CloneConfig(),
	}
}

//Config 基础数据源的配置文件
func (b *BaseSource) Config() *config.JSON {
	return b.conf
}
