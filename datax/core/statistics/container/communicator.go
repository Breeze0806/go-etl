package container

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/core/statistics/communication"
)

//State 状态
type State int

//Communicator 交换器 todo 未使用
type Communicator interface {
	RegisterCommunication(configs []*config.JSON)

	Collect() Communicator

	Report(communication communication.Communication)

	CollectState() State

	GetCommunication(id int64) communication.Communication

	GetCommunicationMap() map[int64]communication.Communication
}
