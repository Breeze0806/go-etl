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

//Task 写入任务
type Task interface {
	plugin.Task

	//开始从receiver中读取记录写入
	StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error
	//是否支持故障转移，就是是否在写入后失败重试
	SupportFailOver() bool
}

//BaseTask 基础写入任务，辅助和简化写入任务接口的实现
type BaseTask struct {
	*plugin.BaseTask
}

//NewBaseTask 创建基础任务
func NewBaseTask() *BaseTask {
	return &BaseTask{
		BaseTask: plugin.NewBaseTask(),
	}
}

//SupportFailOver 是否支持故障转移，就是是否在写入后失败重试
func (b *BaseTask) SupportFailOver() bool {
	return false
}
