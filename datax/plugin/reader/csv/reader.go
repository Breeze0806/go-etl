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
	"github.com/Breeze0806/go-etl/config"
	spireader "github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/file"
)

// Reader reader
type Reader struct {
	pluginConf *config.JSON
}

// ResourcesConfig plugin resource configuration
func (r *Reader) ResourcesConfig() *config.JSON {
	return r.pluginConf
}

// Job job
func (r *Reader) Job() spireader.Job {
	job := NewJob()
	job.SetPluginConf(r.pluginConf)
	return job
}

// Task task
func (r *Reader) Task() spireader.Task {
	task := file.NewTask()
	task.SetPluginConf(r.pluginConf)
	return task
}
