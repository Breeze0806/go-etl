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

package oracle

import (
	"github.com/Breeze0806/go-etl/config"
	spiwriter "github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/dbms"
	"github.com/Breeze0806/go-etl/storage/database"
)

// Writer Writer
type Writer struct {
	pluginConf *config.JSON
}

// ResourcesConfig Plugin Resource Configuration
func (w *Writer) ResourcesConfig() *config.JSON {
	return w.pluginConf
}

// Job Job
func (w *Writer) Job() spiwriter.Job {
	job := &Job{
		Job: dbms.NewJob(dbms.NewBaseDbHandler(
			func(name string, conf *config.JSON) (e dbms.Execer, err error) {
				if e, err = database.Open(name, conf); err != nil {
					return nil, err
				}
				return
			}, nil)),
	}
	job.SetPluginConf(w.pluginConf)
	return job
}

// Task Task
func (w *Writer) Task() spiwriter.Task {
	task := &Task{
		Task: dbms.NewTask(dbms.NewBaseDbHandler(
			func(name string, conf *config.JSON) (e dbms.Execer, err error) {
				if e, err = database.Open(name, conf); err != nil {
					return nil, err
				}
				return
			}, nil)),
	}
	task.SetPluginConf(w.pluginConf)
	return task
}
