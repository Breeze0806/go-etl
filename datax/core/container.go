package core

import (
	"github.com/Breeze0806/go-etl/datax/common/config"
	"github.com/Breeze0806/go-etl/datax/core/statistics/communication"
)

type Container interface {
	Start() error
}

type BaseCotainer struct {
	conf *config.Json
	com  communication.Communication
}

func (b *BaseCotainer) SetConfig(conf *config.Json) {
	b.conf = conf
}

func (b *BaseCotainer) Config() *config.Json {
	return b.conf
}

func (b *BaseCotainer) SetCommunication(com communication.Communication) {
	b.com = com
}

func (b *BaseCotainer) Communication() communication.Communication {
	return b.com
}
