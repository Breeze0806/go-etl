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

package csv

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/file"
)

//Job 工作
type Job struct {
	*file.Job

	conf *Config
}

//NewJob 创建工作
func NewJob() *Job {
	return &Job{
		Job: file.NewJob(),
	}
}

//Init 初始化
func (j *Job) Init(ctx context.Context) (err error) {
	j.conf, err = NewConfig(j.PluginJobConf())
	return
}

//Split 切分
func (j *Job) Split(ctx context.Context, number int) (configs []*config.JSON, err error) {
	for _, v := range j.conf.Path {
		conf, _ := config.NewJSONFromString("{}")
		if err = conf.Set("path", v); err != nil {
			return
		}
		if err = conf.Set("content.0", j.conf.Config); err != nil {
			return
		}
		configs = append(configs, conf)
	}
	return
}
