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

package schedule

import "sync"

// MappedTaskManager task manager
// toto I don't know why len(remain) + len(run) can't accurately represent the number of real-time tasks, mainly because len(run) is not accurate
type MappedTaskManager struct {
	sync.Mutex

	remain []MappedTask          // Pending queue
	run    map[string]MappedTask // Running queue
	num    int                   // Number of tasks
}

// NewTaskManager create task manager
func NewTaskManager() *MappedTaskManager {
	return &MappedTaskManager{
		run: make(map[string]MappedTask),
	}
}

// IsEmpty check if the task manager is empty
func (t *MappedTaskManager) IsEmpty() bool {
	t.Lock()
	defer t.Unlock()
	return t.lockedSize() == 0
}

// Size number of tasks
func (t *MappedTaskManager) Size() int {
	t.Lock()
	defer t.Unlock()
	return t.lockedSize()
}

// Runs get currently running tasks
func (t *MappedTaskManager) Runs() (tasks []MappedTask) {
	t.Lock()
	for _, v := range t.run {
		tasks = append(tasks, v)
	}
	t.Unlock()
	return
}

// lockedSize number of unlocked tasks
func (t *MappedTaskManager) lockedSize() int {
	return t.num
}

// RemoveRunAndPushRemain move task from running queue to pending queue
func (t *MappedTaskManager) RemoveRunAndPushRemain(task MappedTask) {
	t.Lock()
	defer t.Unlock()
	t.lockedRemoveRun(task)
	t.lockedPushRemain(task)
}

// PushRemain add task to pending queue
func (t *MappedTaskManager) PushRemain(task MappedTask) {
	t.Lock()
	defer t.Unlock()
	t.lockedPushRemain(task)
}

// RemoveRun remove task from running queue
func (t *MappedTaskManager) RemoveRun(task MappedTask) {
	t.Lock()
	defer t.Unlock()
	t.lockedRemoveRun(task)
}

// PopRemainAndAddRun move task from pending queue to running queue
func (t *MappedTaskManager) PopRemainAndAddRun() (task MappedTask, ok bool) {
	t.Lock()
	defer t.Unlock()
	task, ok = t.lockedPopRemain()
	if ok {
		t.lockedAddRun(task)
	}
	return
}

// lockedRemoveRun remove task from running queue with lock
func (t *MappedTaskManager) lockedRemoveRun(task MappedTask) {
	t.run[task.Key()] = nil
	delete(t.run, task.Key())
	t.num--
}

// lockedPushRemain add task to pending queue with lock
func (t *MappedTaskManager) lockedPushRemain(task MappedTask) {
	t.remain = append(t.remain, task)
	t.num++
}

// lockedPushRemain add task to running queue with lock
func (t *MappedTaskManager) lockedAddRun(task MappedTask) {
	t.run[task.Key()] = task
}

// lockedPopRemain retrieve task from pending queue. If there are no tasks in the pending queue
func (t *MappedTaskManager) lockedPopRemain() (task MappedTask, ok bool) {
	if len(t.remain) == 0 {
		return nil, false
	}
	task, ok = t.remain[0], true
	t.remain, t.remain[0] = t.remain[1:], nil
	return
}
