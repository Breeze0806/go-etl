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

package rdbm

import (
	"context"
	"database/sql"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

//Task 任务
type Task struct {
	*plugin.BaseTask

	handler DbHandler
	Querier Querier
	Config  Config
	Table   database.Table
}

func NewTask(handler DbHandler) *Task {
	return &Task{
		BaseTask: plugin.NewBaseTask(),

		handler: handler,
	}
}

//Init 初始化
func (t *Task) Init(ctx context.Context) (err error) {
	var name string
	if name, err = t.PluginConf().GetString("dialect"); err != nil {
		return
	}

	if t.Config, err = t.handler.Config(t.PluginJobConf()); err != nil {
		return
	}

	var jobSettingConf *config.JSON
	if jobSettingConf, err = t.PluginJobConf().GetConfig(coreconst.DataxJobSetting); err != nil {
		jobSettingConf, _ = config.NewJSONFromString("{}")
		err = nil
	}
	if err = jobSettingConf.Set("username", t.Config.GetUsername()); err != nil {
		return
	}

	if err = jobSettingConf.Set("password", t.Config.GetPassword()); err != nil {
		return
	}

	if err = jobSettingConf.Set("url", t.Config.GetURL()); err != nil {
		return
	}

	if t.Querier, err = t.handler.Querier(name, jobSettingConf); err != nil {
		return
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = t.Querier.PingContext(timeoutCtx)
	if err != nil {
		return
	}

	param := t.handler.TableParam(t.Config, t.Querier)
	if t.Table, err = t.Querier.FetchTableWithParam(ctx, param); err != nil {
		return
	}

	return
}

//Destroy 销毁
func (t *Task) Destroy(ctx context.Context) (err error) {
	if t.Querier != nil {
		err = t.Querier.Close()
	}
	return
}

type BatchReader interface {
	JobID() int64
	TaskGroupID() int64
	TaskID() int64
	Read(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error)
	Parameter() database.Parameter
}

type BaseBatchReader struct {
	task *Task
	mode string
	opts *sql.TxOptions
}

func NewBaseBatchReader(task *Task, mode string, opts *sql.TxOptions) *BaseBatchReader {
	return &BaseBatchReader{
		task: task,
		mode: mode,
		opts: opts,
	}
}

func (b *BaseBatchReader) JobID() int64 {
	return b.task.JobID()
}

func (b *BaseBatchReader) TaskID() int64 {
	return b.task.TaskID()
}

func (b *BaseBatchReader) TaskGroupID() int64 {
	return b.task.TaskGroupID()
}

func (b *BaseBatchReader) Parameter() database.Parameter {
	return NewQueryParam(b.task.Config, b.task.Table, b.opts)
}

func (b *BaseBatchReader) Read(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error) {
	if b.mode == "Tx" {
		return b.task.Querier.FetchRecordWithTx(ctx, param, handler)
	}
	return b.task.Querier.FetchRecord(ctx, param, handler)
}

//StartRead 开始读
func StartRead(ctx context.Context, reader BatchReader, sender plugin.RecordSender) (err error) {
	handler := database.NewBaseFetchHandler(func() (element.Record, error) {
		return sender.CreateRecord()
	}, func(r element.Record) error {
		return sender.SendWriter(r)
	})

	log.Infof("jobid %v taskgroupid %v taskid %v startRead begin", reader.JobID(), reader.TaskGroupID(), reader.TaskID())
	defer func() {
		sender.Terminate()
		log.Infof("jobid %v taskgroupid %v taskid %v startRead end", reader.JobID(), reader.TaskGroupID(), reader.TaskID())
	}()
	if err = reader.Read(ctx, reader.Parameter(), handler); err != nil {
		return
	}
	return nil
}
