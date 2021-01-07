package container

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/core/statistics/communication"
)

type State int

type Communicator interface {
	RegisterCommunication(configs []*config.Json)

	Collect() Communicator

	Report(communication communication.Communication)

	CollectState() State

	GetCommunication(id int64) communication.Communication

	GetCommunicationMap() map[int64]communication.Communication
}
