package container

import "github.com/Breeze0806/go-etl/datax/common/config"

type State int

type Communicator interface {
	RegisterCommunication(configs []*config.Json)

	Collect() Communicator

	Report(communication Communicator)

	CollectState() State

	GetCommunication(id int64) Communicator

	GetCommunicationMap() map[int64]Communicator
}
