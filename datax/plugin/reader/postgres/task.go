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

package postgres

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
)

//Task 任务
type Task struct {
	*rdbm.Task
}

//StartRead 开始读
func (t *Task) StartRead(ctx context.Context, sender plugin.RecordSender) (err error) {
	return rdbm.StartRead(ctx, rdbm.NewBaseBatchReader(t.Task, "", nil), sender)
}
