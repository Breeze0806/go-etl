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

package dbms

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

// Task normal dbms task
type Task struct {
	*plugin.BaseTask

	handler DbHandler
	Querier Querier
	Config  Config
	Table   database.Table
}

// NewTask Get task through database handler
func NewTask(handler DbHandler) *Task {
	return &Task{
		BaseTask: plugin.NewBaseTask(),

		handler: handler,
	}
}

// Init Initialization
func (t *Task) Init(ctx context.Context) (err error) {
	var name string
	if name, err = t.PluginConf().GetString("dialect"); err != nil {
		return t.Wrapf(err, "GetString fail")
	}

	if t.Config, err = t.handler.Config(t.PluginJobConf()); err != nil {
		return t.Wrapf(err, "Config fail")
	}

	var jobSettingConf *config.JSON
	if jobSettingConf, err = t.PluginJobConf().GetConfig(coreconst.DataxJobSetting); err != nil {
		jobSettingConf, _ = config.NewJSONFromString("{}")
		err = nil
	}
	jobSettingConf.Set("username", t.Config.GetUsername())
	jobSettingConf.Set("password", t.Config.GetPassword())
	jobSettingConf.Set("url", t.Config.GetURL())

	if t.Querier, err = t.handler.Querier(name, jobSettingConf); err != nil {
		return t.Wrapf(err, "Querier fail")
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = t.Querier.PingContext(timeoutCtx)
	if err != nil {
		return
	}

	param := t.handler.TableParam(t.Config, t.Querier)
	if setter, ok := param.Table().(database.ConfigSetter); ok {
		setter.SetConfig(t.PluginJobConf())
	}

	if len(t.Config.GetQuerySQL()) == 0 {
		if t.Table, err = t.Querier.FetchTableWithParam(ctx, param); err != nil {
			return t.Wrapf(err, "FetchTableWithParam fail")
		}
		return
	}
	t.Table = param.Table()
	return
}

// Destroy Destruction
func (t *Task) Destroy(ctx context.Context) (err error) {
	if t.Querier != nil {
		err = t.Querier.Close()
	}
	return t.Wrapf(err, "Close fail")
}

// BatchReader Batch reader
type BatchReader interface {
	JobID() int64       // Job number
	TaskGroupID() int64 // Task group number
	TaskID() int64      // Task number
	Read(ctx context.Context, param database.Parameter,
		handler database.FetchHandler) (err error) // Query through context ctx, description, and database handler
	Parameter() database.Parameter // Query parameters
}

// BaseBatchReader Basic batch reader
type BaseBatchReader struct {
	task *Task
	mode string
	opts *sql.TxOptions
}

// NewBaseBatchReader Get basic batch reader through task, query mode, and transaction options
func NewBaseBatchReader(task *Task, mode string, opts *sql.TxOptions) *BaseBatchReader {
	return &BaseBatchReader{
		task: task,
		mode: mode,
		opts: opts,
	}
}

// JobID Job number
func (b *BaseBatchReader) JobID() int64 {
	return b.task.JobID()
}

// TaskID Task number
func (b *BaseBatchReader) TaskID() int64 {
	return b.task.TaskID()
}

// TaskGroupID Task group number
func (b *BaseBatchReader) TaskGroupID() int64 {
	return b.task.TaskGroupID()
}

// Parameter Query parameters
func (b *BaseBatchReader) Parameter() database.Parameter {
	return NewQueryParam(b.task.Config, b.task.Table, b.opts)
}

// Query through context ctx, description, and database handler
func (b *BaseBatchReader) Read(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error) {
	if b.mode == "Tx" {
		return b.task.Querier.FetchRecordWithTx(ctx, param, handler)
	}
	return b.task.Querier.FetchRecord(ctx, param, handler)
}

// StartRead Start reading
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
