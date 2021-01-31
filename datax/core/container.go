package core

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/core/statistics/communication"
)

//Container 容器
type Container interface {
	Start() error
}

//BaseCotainer 基础容器
type BaseCotainer struct {
	conf *config.JSON
	com  communication.Communication
}

//NewBaseCotainer 创建基础容器
func NewBaseCotainer() *BaseCotainer {
	return &BaseCotainer{}
}

//SetConfig 设置JSON配置
func (b *BaseCotainer) SetConfig(conf *config.JSON) {
	b.conf = conf
}

//Config JSON配置
func (b *BaseCotainer) Config() *config.JSON {
	return b.conf
}

//SetCommunication 未真正使用
func (b *BaseCotainer) SetCommunication(com communication.Communication) {
	b.com = com
}

//Communication  未真正使用
func (b *BaseCotainer) Communication() communication.Communication {
	return b.com
}
