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

package plugin

import (
	"fmt"

	"github.com/pingcap/errors"
)

// Task - An interface for representing tasks
type Task interface {
	Plugin

	// Task Information Collector, todo: not currently used
	TaskCollector() TaskCollector
	// Set Task Information Collector, todo: not currently used
	SetTaskCollector(collector TaskCollector)

	// Job ID
	JobID() int64
	// Set Job ID
	SetJobID(jobID int64)
	// Task Group ID
	TaskGroupID() int64
	// Set Task Group ID
	SetTaskGroupID(taskGroupID int64)
	// Task ID
	TaskID() int64
	// Set Task ID
	SetTaskID(taskID int64)
	// Wrap Error
	Wrapf(err error, format string, args ...interface{}) error
	// Format - Log format
	Format(format string) string
}

// BaseTask - A basic task that assists and simplifies the implementation of task interfaces
type BaseTask struct {
	*BasePlugin

	jobID       int64
	taskID      int64
	taskGroupID int64
	collector   TaskCollector
}

// NewBaseTask - Creates a new instance of a base task
func NewBaseTask() *BaseTask {
	return &BaseTask{
		BasePlugin: NewBasePlugin(),
	}
}

// TaskCollector - Collects information related to tasks
func (b *BaseTask) TaskCollector() TaskCollector {
	return b.collector
}

// SetTaskCollector - Sets the task information collector
func (b *BaseTask) SetTaskCollector(collector TaskCollector) {
	b.collector = collector
}

// TaskID - The unique identifier for a task
func (b *BaseTask) TaskID() int64 {
	return b.taskID
}

// SetTaskID - Sets the unique identifier for a task
func (b *BaseTask) SetTaskID(taskID int64) {
	b.taskID = taskID
}

// TaskGroupID - The unique identifier for a group of tasks
func (b *BaseTask) TaskGroupID() int64 {
	return b.taskGroupID
}

// SetTaskGroupID - Sets the unique identifier for a group of tasks
func (b *BaseTask) SetTaskGroupID(taskGroupID int64) {
	b.taskGroupID = taskGroupID
}

// JobID - The unique identifier for a job
func (b *BaseTask) JobID() int64 {
	return b.jobID
}

// SetJobID - Sets the unique identifier for a job
func (b *BaseTask) SetJobID(jobID int64) {
	b.jobID = jobID
}

// Wrapf - Wraps an error with additional context
func (b *BaseTask) Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, b.Format(format), args...)
}

// Format - The format for logging messages
func (b *BaseTask) Format(format string) string {
	return fmt.Sprintf("jobId : %v taskgroupId: %v taskId: %v %v", b.jobID, b.taskGroupID, b.taskID, format)
}
