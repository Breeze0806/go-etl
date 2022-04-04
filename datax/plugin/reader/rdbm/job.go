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
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

//Job 工作
type Job struct {
	*plugin.BaseJob

	Querier Querier
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
		return
	}

	var conf Config
	if conf, err = j.handler.Config(j.PluginJobConf()); err != nil {
		return
	}

	var jobSettingConf *config.JSON
	if jobSettingConf, err = j.PluginJobConf().GetConfig(coreconst.DataxJobSetting); err != nil {
		jobSettingConf, _ = config.NewJSONFromString("{}")
		err = nil
	}
	if err = jobSettingConf.Set("username", conf.GetUsername()); err != nil {
		return
	}

	if err = jobSettingConf.Set("password", conf.GetPassword()); err != nil {
		return
	}

	if err = jobSettingConf.Set("url", conf.GetURL()); err != nil {
		return
	}

	if j.Querier, err = j.handler.Querier(name, jobSettingConf); err != nil {
		return
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = j.Querier.PingContext(timeoutCtx)
	if err != nil {
		return
	}
	return
}

//Destroy 销毁
func (j *Job) Destroy(ctx context.Context) (err error) {
	if j.Querier != nil {
		err = j.Querier.Close()
	}
	return
}

//Split 切分
func (j *Job) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return []*config.JSON{j.PluginJobConf().CloneConfig()}, nil
}
