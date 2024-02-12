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

// Writer: Write Runner
type Writer struct {
	*baseRunner
	receiver plugin.RecordReceiver
	task     writer.Task
	describe string
}

// NewWriter: Creates a write runner by reading the task, the recorder, and the task keyword
func NewWriter(task writer.Task, receiver plugin.RecordReceiver, taskKey string) *Writer {
	return &Writer{
		baseRunner: &baseRunner{},
		receiver:   receiver,
		task:       task,
		describe:   taskKey,
	}
}

// Plugin: Plugin task
func (w *Writer) Plugin() plugin.Task {
	return w.task
}

// Run: Runs in the following order: Init->Prepare->StartWrite->Post->Destroy
func (w *Writer) Run(ctx context.Context) (err error) {
	defer func() {
		log.Debugf("datax writer runner %v starts to destroy", w.describe)
		if destroyErr := w.task.Destroy(ctx); destroyErr != nil {
			log.Errorf("task destroy fail, err: %v", destroyErr)
		}
	}()
	log.Debugf("datax writer runner %v starts to init", w.describe)
	if err = w.task.Init(ctx); err != nil {
		return err
	}

	log.Debugf("datax writer runner %v starts to prepare", w.describe)
	if err = w.task.Prepare(ctx); err != nil {
		return err
	}

	log.Debugf("datax writer runner %v starts to StartWrite", w.describe)
	if err = w.task.StartWrite(ctx, w.receiver); err != nil {
		return err
	}

	log.Debugf("datax writer runner %v starts to post", w.describe)
	if err = w.task.Post(ctx); err != nil {
		return err
	}
	return
}

// Shutdown: Closes down
func (w *Writer) Shutdown() error {
	return w.receiver.Shutdown()
}
