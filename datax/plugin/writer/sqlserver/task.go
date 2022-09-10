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

package sqlserver

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/rdbm"
	"github.com/Breeze0806/go-etl/storage/database"
)

var execModeMap = map[string]string{
	database.WriteModeInsert: rdbm.ExecModeNormal,
}

func execMode(writeMode string) string {
	if mode, ok := execModeMap[writeMode]; ok {
		return mode
	}
	return rdbm.ExecModeNormal
}

//Task 任务
type Task struct {
	*rdbm.Task
}

//StartWrite 开始写
func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) (err error) {
	return rdbm.StartWrite(ctx, rdbm.NewBaseBatchWriter(t.Task, execMode(t.Config.GetWriteMode()), nil), receiver)
}
