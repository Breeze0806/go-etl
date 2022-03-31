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

package runner

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

//Writer 写入运行器
type Writer struct {
	*baseRunner
	receiver plugin.RecordReceiver
	task     writer.Task
	describe string
}

//NewWriter 通过读取任务task、记录接受器receiver以及任务关键字taskKey创建写入运行器
func NewWriter(task writer.Task, receiver plugin.RecordReceiver, taskKey string) *Writer {
	return &Writer{
		baseRunner: &baseRunner{},
		receiver:   receiver,
		task:       task,
		describe:   taskKey,
	}
}

//Plugin 插件任务
func (w *Writer) Plugin() plugin.Task {
	return w.task
}

//Run 运行，运行顺序：Init->Prepare->StartWrite->Post->Destroy
func (w *Writer) Run(ctx context.Context) (err error) {
	defer func() {
		log.Debugf("datax writer runner %v starts to destroy", w.describe)
		if destroyErr := w.task.Destroy(ctx); destroyErr != nil {
			log.Errorf("task destroy fail, err: %v", destroyErr)
		}
	}()
	log.Debugf("datax writer runner %v starts to init", w.describe)
	if err = w.task.Init(ctx); err != nil {
		log.Errorf("task init fail, err: %v", err)
		return
	}

	log.Debugf("datax writer runner %v starts to prepare", w.describe)
	if err = w.task.Prepare(ctx); err != nil {
		log.Errorf("task prepare fail, err: %v", err)
		return
	}

	log.Debugf("datax writer runner %v starts to StartWrite", w.describe)
	if err = w.task.StartWrite(ctx, w.receiver); err != nil {
		log.Errorf("task startWrite fail, err: %v", err)
		return
	}

	log.Debugf("datax writer runner %v starts to post", w.describe)
	if err = w.task.Post(ctx); err != nil {
		log.Errorf("task post fail, err: %v", err)
		return
	}
	return
}

// Shutdown 关闭
func (w *Writer) Shutdown() error {
	return w.receiver.Shutdown()
}
