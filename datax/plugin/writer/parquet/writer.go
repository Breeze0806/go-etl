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

package parquet

import (
	"github.com/Breeze0806/go-etl/config"
	spiwriter "github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/file"
)

// Writer parquet写入器
type Writer struct {
	pluginConf *config.JSON
}

// ResourcesConfig 返回初始化写入器的数据源配置。
func (w *Writer) ResourcesConfig() *config.JSON {
	return w.pluginConf
}

// Job 返回描述写入器如何从数据源提取数据的Job。
func (w *Writer) Job() spiwriter.Job {
	job := NewJob()
	job.SetPluginConf(w.pluginConf)
	return job
}

// Task 返回通过最大化拆分Job获得的最小执行单元。
func (w *Writer) Task() spiwriter.Task {
	task := file.NewTask(func(conf *config.JSON) (file.Config, error) {
		c, err := file.NewBaseConfig(conf)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
	task.SetPluginConf(w.pluginConf)
	return task
}
