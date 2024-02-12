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

import (
	"errors"
	"sync"
	goatomic "sync/atomic"

	"go.uber.org/atomic"
)

// ErrorClosed - Represents an error indicating that the operation is closed.
var (
	ErrClosed = errors.New("task schduler closed")
)

type taskWrapper struct {
	task   Task
	result chan error
}

// TaskScheduler - Represents a task scheduler.
type TaskSchduler struct {
	taskWrappers chan *taskWrapper // PendingTaskQueue - Queue of tasks waiting to be executed.
	wg           sync.WaitGroup
	stop         chan struct{} // StopSignal - Signal to stop the scheduler.
	stopped      int32         // StopFlag - Flag indicating whether the scheduler should stop.
	size         *atomic.Int32 // PendingQueueSize - Size of the pending task queue.
}

// NewTaskScheduler - Creates a new task scheduler based on the number of workers (workerNumber) and the capacity of the pending task queue (capacity).
func NewTaskSchduler(workerNumer, cap int) *TaskSchduler {
	t := &TaskSchduler{
		taskWrappers: make(chan *taskWrapper, cap),
		stop:         make(chan struct{}),
		stopped:      0,
		size:         atomic.NewInt32(0),
	}

	for i := 0; i < workerNumer; i++ {
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			t.processTask()
		}()
	}

	return t
}

// Push - Adds a task to the queue and receives a notification channel for the execution result. An error is reported if the queue is closed.
func (t *TaskSchduler) Push(task Task) (<-chan error, error) {
	if goatomic.CompareAndSwapInt32(&t.stopped, 1, 1) {
		return nil, ErrClosed
	}
	tw := &taskWrapper{
		task:   task,
		result: make(chan error, 1),
	}

	select {
	case t.taskWrappers <- tw:
		t.size.Inc()
		return tw.result, nil
	case <-t.stop:
		return nil, ErrClosed
	}
}

// Size - Represents the size of the pending task queue.
func (t *TaskSchduler) Size() int32 {
	return t.size.Load()
}

// Stop - Stops the task scheduler.
func (t *TaskSchduler) Stop() {
	if !goatomic.CompareAndSwapInt32(&t.stopped, 0, 1) {
		return
	}
	close(t.stop)
	t.wg.Wait()
}

func (t *TaskSchduler) processTask() {
	for {
		select {
		case tw, ok := <-t.taskWrappers:
			if !ok {
				return
			}
			tw.result <- tw.task.Do()
			t.size.Dec()
		case <-t.stop:
			return
		}
	}
}
