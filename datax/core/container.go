package core

import (
	"github.com/Breeze0806/go-etl/datax/common/config"
	statContainer "github.com/Breeze0806/go-etl/datax/core/statistics/container"
)

type Container interface {
	Start() error
}

type BaseCotainer struct {
	conf *config.Json
	com  statContainer.Communicator
}

func (b *BaseCotainer) SetConfig(conf *config.Json) {
	b.conf = conf
}

func (b *BaseCotainer) Config() *config.Json {
	return b.conf
}

func (b *BaseCotainer) SetCommunication(com statContainer.Communicator) {
	b.com = com
}

func (b *BaseCotainer) Communication() statContainer.Communicator {
	return b.com
}
