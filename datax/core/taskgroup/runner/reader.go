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
	"fmt"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
)

//Reader 读取运行器
type Reader struct {
	*baseRunner

	sender   plugin.RecordSender
	task     reader.Task
	describe string
}

//NewReader 通过读取任务task、记录发送器sender以及任务关键字taskKey创建读取运行器
func NewReader(task reader.Task, sender plugin.RecordSender, taskKey string) *Reader {
	return &Reader{
		baseRunner: &baseRunner{},
		sender:     sender,
		task:       task,
		describe:   taskKey,
	}
}

//Plugin 插件任务
func (r *Reader) Plugin() plugin.Task {
	return r.task
}

//Run 运行，运行顺序：Init->Prepare->StartRead->Post->Destroy
func (r *Reader) Run(ctx context.Context) (err error) {
	defer func() {
		log.Debugf("datax reader runner %v starts to destroy", r.describe)
		if destroyErr := r.task.Destroy(ctx); destroyErr != nil {
			log.Errorf("task destroy fail, err: %v", destroyErr)
		}
	}()

	log.Debugf("datax reader runner %v starts to init", r.describe)
	if err = r.task.Init(ctx); err != nil {
		return fmt.Errorf("task init fail, err: %v", err)
	}

	log.Debugf("datax reader runner %v starts to prepare", r.describe)
	if err = r.task.Prepare(ctx); err != nil {
		return fmt.Errorf("task prepare fail, err: %v", err)
	}

	log.Debugf("datax reader runner %v starts to startRead", r.describe)
	if err = r.task.StartRead(ctx, r.sender); err != nil {
		return fmt.Errorf("task startRead fail, err: %v", err)
	}

	log.Debugf("datax reader runner %v starts to post", r.describe)
	if err = r.task.Post(ctx); err != nil {
		return fmt.Errorf("task post fail, err: %v", err)
	}
	return
}

//Shutdown 关闭
func (r *Reader) Shutdown() error {
	return r.sender.Shutdown()
}
