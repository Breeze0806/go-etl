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

package file

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/stream/file"
)

//Task 任务
type Task struct {
	*plugin.BaseTask

	streamer *file.InStreamer
}

//NewTask 新建任务
func NewTask() *Task {
	return &Task{
		BaseTask: plugin.NewBaseTask(),
	}
}

//Init 初始化任务
func (t *Task) Init(ctx context.Context) (err error) {
	var name string
	if name, err = t.PluginConf().GetString("opener"); err != nil {
		return
	}
	var filename string
	if filename, err = t.PluginJobConf().GetString("path"); err != nil {
		return
	}

	if t.streamer, err = file.NewInStreamer(name, filename); err != nil {
		return
	}
	return
}

//Destroy 销毁
func (t *Task) Destroy(ctx context.Context) (err error) {
	if t.streamer != nil {
		err = t.streamer.Close()
	}
	return
}

type handler struct {
	sender plugin.RecordSender
}

func newHander(sender plugin.RecordSender) *handler {
	return &handler{
		sender: sender,
	}
}

func (h *handler) CreateRecord() (element.Record, error) {
	return h.sender.CreateRecord()
}

func (h *handler) OnRecord(r element.Record) error {
	return h.sender.SendWriter(r)
}

//StartRead 开启读取数据发往sender
func (t *Task) StartRead(ctx context.Context, sender plugin.RecordSender) (err error) {
	handler := newHander(sender)

	log.Infof("jobid %v taskgroupid %v taskid %v startRead begin", t.JobID(), t.TaskGroupID(), t.TaskID())
	defer func() {
		sender.Terminate()
		log.Infof("jobid %v taskgroupid %v taskid %v startRead end", t.JobID(), t.TaskGroupID(), t.TaskID())
	}()
	var configs []*config.JSON
	configs, err = t.PluginJobConf().GetConfigArray("content")
	if err != nil {
		return err
	}

	for _, conf := range configs {
		if err = t.streamer.Read(ctx, conf, handler); err != nil {
			return
		}
	}
	return nil
}
