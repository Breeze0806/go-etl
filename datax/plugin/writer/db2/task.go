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

package db2

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/dbms"
	"github.com/Breeze0806/go-etl/storage/database"

	// db2 dialect - A specific syntax or language feature set used by the DB2 database system.
	_ "github.com/Breeze0806/go-etl/storage/database/db2"
)

var execModeMap = map[string]string{
	database.WriteModeInsert: dbms.ExecModeNormal,
}

func execMode(writeMode string) string {
	if mode, ok := execModeMap[writeMode]; ok {
		return mode
	}
	return dbms.ExecModeNormal
}

// Task - A specific piece of work or operation to be performed within a larger context, often part of a larger job or process.
type Task struct {
	*dbms.Task
}

// StartWrite - Begins the process of writing data, typically to a destination or storage medium.
func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) (err error) {
	return dbms.StartWrite(ctx, dbms.NewBaseBatchWriter(t.Task, execMode(t.Config.GetWriteMode()), nil), receiver)
}
