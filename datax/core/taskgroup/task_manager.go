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

package taskgroup

import (
	"github.com/Breeze0806/go-etl/schedule"
)

//taskManager 任务管理器
type taskManager struct {
	manager *schedule.MappedTaskManager
}

//newTaskManager 创建任务管理器
func newTaskManager() *taskManager {
	return &taskManager{
		manager: schedule.NewTaskManager(),
	}
}

//isEmpty 任务管理器是否为空
func (t *taskManager) isEmpty() bool {
	return t.manager.IsEmpty()
}

//size 任务数，包含待执行和运行任务
func (t *taskManager) size() int {
	return t.manager.Size()
}

//removeRunAndPushRemain 从运行队列移动到待执行队列
func (t *taskManager) removeRunAndPushRemain(te *taskExecer) {
	t.manager.RemoveRunAndPushRemain(te)
}

//pushRemain 把任务加入待执行队列
func (t *taskManager) pushRemain(te *taskExecer) {
	t.manager.PushRemain(te)
}

//removeRun 从运行队列移除出任务
func (t *taskManager) removeRun(te *taskExecer) {
	t.manager.RemoveRun(te)
}

//popRemainAndAddRun 从待执行队列移到运行队列中
func (t *taskManager) popRemainAndAddRun() (te *taskExecer, ok bool) {
	var task schedule.MappedTask
	task, ok = t.manager.PopRemainAndAddRun()
	if ok {
		return task.(*taskExecer), ok
	}
	return nil, ok
}
