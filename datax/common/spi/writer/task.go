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

package writer

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

// Task - a task related to writing operations
type Task interface {
	plugin.Task

	// Start reading records from the receiver and write them
	StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error
	// SupportFailOver - whether fault tolerance is supported, i.e., whether to retry after a failed write
	SupportFailOver() bool
}

// BaseTask - a fundamental task class that assists and simplifies the implementation of writing task interfaces
type BaseTask struct {
	*plugin.BaseTask
}

// NewBaseTask - a function or method to create a new instance of BaseTask
func NewBaseTask() *BaseTask {
	return &BaseTask{
		BaseTask: plugin.NewBaseTask(),
	}
}

// SupportFailOver - whether fault tolerance is supported, i.e., whether to retry after a failed write
func (b *BaseTask) SupportFailOver() bool {
	return false
}
