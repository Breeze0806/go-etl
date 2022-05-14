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

//Task 任务接口
type Task interface {
	Plugin

	//任务信息收集器，todo 未使用
	TaskCollector() TaskCollector
	//设置任务信息收集器，todo 未使用
	SetTaskCollector(collector TaskCollector)

	//工作ID
	JobID() int64
	//设置工作ID
	SetJobID(jobID int64)
	//任务组ID
	TaskGroupID() int64
	//设置任务组ID
	SetTaskGroupID(taskGroupID int64)
	//任务ID
	TaskID() int64
	//设置任务ID
	SetTaskID(taskID int64)
	//包裹错误
	Wrapf(err error, format string, args ...interface{}) error
	//Format 日志格式
	Format(format string) string
}

//BaseTask 基础任务，用于辅助和简化任务接口的实现
type BaseTask struct {
	*BasePlugin

	jobID       int64
	taskID      int64
	taskGroupID int64
	collector   TaskCollector
}

//NewBaseTask 创建基础任务
func NewBaseTask() *BaseTask {
	return &BaseTask{
		BasePlugin: NewBasePlugin(),
	}
}

//TaskCollector 任务信息收集器
func (b *BaseTask) TaskCollector() TaskCollector {
	return b.collector
}

//SetTaskCollector 设置任务信息收集器
func (b *BaseTask) SetTaskCollector(collector TaskCollector) {
	b.collector = collector
}

//TaskID 任务ID
func (b *BaseTask) TaskID() int64 {
	return b.taskID
}

//SetTaskID 设置任务ID
func (b *BaseTask) SetTaskID(taskID int64) {
	b.taskID = taskID
}

//TaskGroupID 任务组ID
func (b *BaseTask) TaskGroupID() int64 {
	return b.taskGroupID
}

//SetTaskGroupID 设置任务组ID
func (b *BaseTask) SetTaskGroupID(taskGroupID int64) {
	b.taskGroupID = taskGroupID
}

//JobID 工作ID
func (b *BaseTask) JobID() int64 {
	return b.jobID
}

//SetJobID 设置工作ID
func (b *BaseTask) SetJobID(jobID int64) {
	b.jobID = jobID
}

//Wrapf 包裹错误
func (b *BaseTask) Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, b.Format(format), args...)
}

//Format 日志格式
func (b *BaseTask) Format(format string) string {
	return fmt.Sprintf("jobId : %v taskgroupId: %v taskId: %v %v", b.jobID, b.taskGroupID, b.taskID, format)
}
