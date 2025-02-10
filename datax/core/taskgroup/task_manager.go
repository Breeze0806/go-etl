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

// taskManager - Task Manager
type taskManager struct {
	manager *schedule.MappedTaskManager
}

// newTaskManager - Create a new Task Manager
func newTaskManager() *taskManager {
	return &taskManager{
		manager: schedule.NewTaskManager(),
	}
}

// isEmpty - Check if the Task Manager is empty
func (t *taskManager) isEmpty() bool {
	return t.manager.IsEmpty()
}

// size - The number of tasks, including pending and running tasks
func (t *taskManager) size() int {
	return t.manager.Size()
}

// removeRunAndPushRemain - Move a task from the running queue to the pending queue
func (t *taskManager) removeRunAndPushRemain(te *taskExecer) {
	t.manager.RemoveRunAndPushRemain(te)
}

// pushRemain - Add a task to the pending queue
func (t *taskManager) pushRemain(te *taskExecer) {
	t.manager.PushRemain(te)
}

// removeRun - Remove a task from the running queue
func (t *taskManager) removeRun(te *taskExecer) {
	t.manager.RemoveRun(te)
}

// popRemainAndAddRun - Move a task from the pending queue to the running queue
func (t *taskManager) popRemainAndAddRun() (te *taskExecer, ok bool) {
	var task schedule.MappedTask
	task, ok = t.manager.PopRemainAndAddRun()
	if ok {
		return task.(*taskExecer), ok
	}
	return nil, ok
}

func (t *taskManager) isRunsEmpty() bool {
	return len(t.manager.Runs()) == 0
}
