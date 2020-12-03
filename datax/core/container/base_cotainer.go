package container

import (
	"github.com/Breeze0806/go-etl/datax/common/config"
	statContainer "github.com/Breeze0806/go-etl/datax/core/statistics/container"
)

type baseCotainer struct {
	conf *config.Json
	com  statContainer.Communicator
}

func (b *baseCotainer) SetConfig(conf *config.Json) {
	b.conf = conf
}

func (b *baseCotainer) Config() *config.Json {
	return b.conf
}

func (b *baseCotainer) SetCommunication(com statContainer.Communicator) {
	b.com = com
}

func (b *baseCotainer) Communication() statContainer.Communicator {
	return b.com
}
