// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
