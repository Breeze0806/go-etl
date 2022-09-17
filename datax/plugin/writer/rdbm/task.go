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
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/storage/database"
)

//Task 任务
type Task struct {
	*writer.BaseTask

	Handler DbHandler
	Execer  Execer
	Config  Config
	Table   database.Table
}

//NewTask 通过数据库句柄handler创建任务
func NewTask(handler DbHandler) *Task {
	return &Task{
		BaseTask: writer.NewBaseTask(),
		Handler:  handler,
	}
}

//Init 初始化
func (t *Task) Init(ctx context.Context) (err error) {
	var name string
	if name, err = t.PluginConf().GetString("dialect"); err != nil {
		return t.Wrapf(err, "GetString fail")
	}

	if t.Config, err = t.Handler.Config(t.PluginJobConf()); err != nil {
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

	if t.Execer, err = t.Handler.Execer(name, jobSettingConf); err != nil {
		return t.Wrapf(err, "Execer fail")
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = t.Execer.PingContext(timeoutCtx)
	if err != nil {
		return t.Wrapf(err, "PingContext fail")
	}

	param := t.Handler.TableParam(t.Config, t.Execer)
	if t.Table, err = t.Execer.FetchTableWithParam(ctx, param); err != nil {
		return t.Wrapf(err, "FetchTableWithParam fail")
	}

	if setter, ok := t.Table.(database.TableConfigSetter); ok {
		setter.SetConfig(t.PluginJobConf())
	}
	return
}

//Destroy 销毁
func (t *Task) Destroy(ctx context.Context) (err error) {
	if t.Execer != nil {
		err = t.Execer.Close()
	}
	return t.Wrapf(err, "Close fail")
}
