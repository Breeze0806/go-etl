package mysql

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/writer"
)

var _pluginConfig string

func init() {
	var err error
	maker := &maker{}
	if _pluginConfig, err = writer.RegisterWriter(maker); err != nil {
		panic(err)
	}
}

//NewWriterFromFile 创建写入器
func NewWriterFromFile(filename string) (wr writer.Writer, err error) {
	w := &Writer{}
	w.pluginConf, err = config.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	wr = w
	return
}

//NewWriterFromString 创建写入器
func NewWriterFromString(filename string) (wr writer.Writer, err error) {
	w := &Writer{}
	w.pluginConf, err = config.NewJSONFromString(filename)
	if err != nil {
		return nil, err
	}
	wr = w
	return
}

type maker struct{}

func (m *maker) FromFile(filename string) (writer.Writer, error) {
	return NewWriterFromFile(filename)
}

func (m *maker) Default() (writer.Writer, error) {
	return nil, nil
}
