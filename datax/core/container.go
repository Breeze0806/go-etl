package core

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/core/statistics/communication"
)

type Container interface {
	Start() error
}

type BaseCotainer struct {
	conf *config.JSON
	com  communication.Communication
}

func NewBaseCotainer() *BaseCotainer {
	return &BaseCotainer{}
}

func (b *BaseCotainer) SetConfig(conf *config.JSON) {
	b.conf = conf
}

func (b *BaseCotainer) Config() *config.JSON {
	return b.conf
}

func (b *BaseCotainer) SetCommunication(com communication.Communication) {
	b.com = com
}

func (b *BaseCotainer) Communication() communication.Communication {
	return b.com
}
