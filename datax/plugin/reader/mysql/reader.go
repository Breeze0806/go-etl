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

package mysql

import (
	"github.com/Breeze0806/go-etl/config"
	spireader "github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/dbms"
	"github.com/Breeze0806/go-etl/storage/database"

	//mysql storage
	_ "github.com/Breeze0806/go-etl/storage/database/mysql"
)

// Reader 读取器
type Reader struct {
	pluginConf *config.JSON
}

// ResourcesConfig 插件资源配置
func (r *Reader) ResourcesConfig() *config.JSON {
	return r.pluginConf
}

// Job 工作
func (r *Reader) Job() spireader.Job {
	job := &Job{
		Job: dbms.NewJob(dbms.NewBaseDbHandler(func(name string, conf *config.JSON) (q dbms.Querier, err error) {
			if q, err = database.Open(name, conf); err != nil {
				return nil, err
			}
			return
		}, nil)),
	}
	job.SetPluginConf(r.pluginConf)
	return job
}

// Task 任务
func (r *Reader) Task() spireader.Task {
	task := &Task{
		Task: dbms.NewTask(dbms.NewBaseDbHandler(func(name string, conf *config.JSON) (q dbms.Querier, err error) {
			if q, err = database.Open(name, conf); err != nil {
				return nil, err
			}
			return
		}, nil)),
	}
	task.SetPluginConf(r.pluginConf)
	return task
}
