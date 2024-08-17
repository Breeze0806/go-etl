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

package sqlite3

import (
	"github.com/Breeze0806/go-etl/config"
	spireader "github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/dbms"
	"github.com/Breeze0806/go-etl/storage/database"

	// sqlite3 storage - Sqlite3 database storage
	_ "github.com/Breeze0806/go-etl/storage/database/sqlite3"
)

// Reader is uesed to extract data from data source
type Reader struct {
	pluginConf *config.JSON
}

// ResourcesConfig returns the configuration of the data source to initiate the reader.
func (r *Reader) ResourcesConfig() *config.JSON {
	return r.pluginConf
}

// Job returns a description of how the reader extracts data from the data source.
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

// Task returns the smallest execution unit obtained by maximizing the split of a Job
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
