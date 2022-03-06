package mysql

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/reader"
)

var _pluginConfig string

func init() {
	var err error
	maker := &maker{}
	if _pluginConfig, err = reader.RegisterReader(maker); err != nil {
		panic(err)
	}
}

//NewReaderFromFile 创建读取器
func NewReaderFromFile(filename string) (rd reader.Reader, err error) {
	r := &Reader{}
	r.pluginConf, err = config.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	rd = r
	return
}

//NewReaderFromString 创建读取器
func NewReaderFromString(filename string) (rd reader.Reader, err error) {
	r := &Reader{}
	r.pluginConf, err = config.NewJSONFromString(filename)
	if err != nil {
		return nil, err
	}
	rd = r
	return
}

type maker struct{}

func (m *maker) FromFile(filename string) (reader.Reader, error) {
	return NewReaderFromFile(filename)
}

func (m *maker) Default() (reader.Reader, error) {
	return nil, nil
}
