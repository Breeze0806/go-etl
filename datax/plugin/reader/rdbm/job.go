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
	"fmt"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/pingcap/errors"
)

//Job 工作
type Job struct {
	*plugin.BaseJob

	Querier Querier
	Config  Config
	handler DbHandler
}

//NewJob 通过数据库句柄handler获取工作
func NewJob(handler DbHandler) *Job {
	return &Job{
		BaseJob: plugin.NewBaseJob(),
		handler: handler,
	}
}

//Init 初始化
func (j *Job) Init(ctx context.Context) (err error) {
	var name string
	if name, err = j.PluginConf().GetString("dialect"); err != nil {
		return errors.Wrapf(err, "GetString fail")
	}

	if j.Config, err = j.handler.Config(j.PluginJobConf()); err != nil {
		return errors.Wrapf(err, "Config fail")
	}

	var jobSettingConf *config.JSON
	if jobSettingConf, err = j.PluginJobConf().GetConfig(coreconst.DataxJobSetting); err != nil {
		jobSettingConf, _ = config.NewJSONFromString("{}")
		err = nil
	}

	jobSettingConf.Set("username", j.Config.GetUsername())
	jobSettingConf.Set("password", j.Config.GetPassword())
	jobSettingConf.Set("url", j.Config.GetURL())

	if j.Querier, err = j.handler.Querier(name, jobSettingConf); err != nil {
		return errors.Wrapf(err, "Querier fail")
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = j.Querier.PingContext(timeoutCtx)
	if err != nil {
		return errors.Wrapf(err, "PingContext fail")
	}
	return
}

//Destroy 销毁
func (j *Job) Destroy(ctx context.Context) (err error) {
	if j.Querier != nil {
		err = j.Querier.Close()
	}
	return errors.Wrapf(err, "Close fail")
}

//Split 切分
func (j *Job) Split(ctx context.Context, number int) (configs []*config.JSON, err error) {
	if j.Config.GetSplitConfig().Key == "" || number == 1 {
		return []*config.JSON{j.PluginJobConf().CloneConfig()}, nil
	}

	var splitTable database.Table
	param := j.handler.SplitParam(j.Config, j.Querier)
	if splitTable, err = j.Querier.FetchTableWithParam(ctx, param); err != nil {
		err = errors.Wrapf(err, "FetchTableWithParam fail")
		return
	}

	var minColumn element.Column
	minHandler := database.NewBaseFetchHandler(func() (element.Record, error) {
		return element.NewDefaultRecord(), nil
	}, func(r element.Record) (err error) {
		minColumn, err = r.GetByIndex(0)
		return nil
	})
	minParam := j.handler.MinParam(j.Config, splitTable)
	if err = j.Querier.FetchRecord(ctx, minParam, minHandler); err != nil {
		err = errors.Wrapf(err, "FetchTableWithParam fail")
		return
	}

	var maxColumn element.Column
	maxHandler := database.NewBaseFetchHandler(func() (element.Record, error) {
		return element.NewDefaultRecord(), nil
	}, func(r element.Record) error {
		maxColumn, err = r.GetByIndex(0)
		return nil
	})

	maxParam := j.handler.MaxParam(j.Config, splitTable)
	if err = j.Querier.FetchRecord(ctx, maxParam, maxHandler); err != nil {
		err = errors.Wrapf(err, "FetchTableWithParam fail")
		return
	}

	ranges, err := split(minColumn, maxColumn, number,
		j.Config.GetSplitConfig().TimeAccuracy, splitTable.Fields()[0])
	if err != nil {
		err = errors.Wrapf(err, "split fail")
		return
	}

	for _, r := range ranges {
		clone := j.PluginJobConf().CloneConfig()
		_ = clone.Set("split.range", r)
		where := r.where
		if j.Config.GetWhere() != "" {
			where = fmt.Sprintf("(%s) and (%s)", j.Config.GetWhere(), r.where)
		}
		_ = clone.Set("where", where)
		configs = append(configs, clone)
	}

	return
}
