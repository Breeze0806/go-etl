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

package xlsx

import (
	"context"
	"strconv"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/file"
	"github.com/pingcap/errors"
)

// Job - A unit of work or task to be performed
type Job struct {
	*file.Job

	conf *Config
}

// NewJob - Creates a new instance of a Job
func NewJob() *Job {
	return &Job{
		Job: file.NewJob(),
	}
}

// Init - Initializes or sets up the Job for execution
func (j *Job) Init(ctx context.Context) (err error) {
	j.conf, err = NewConfig(j.PluginJobConf())
	return errors.Wrapf(err, "NewConfig fail. val: %v", j.PluginJobConf())
}

// Split - To divide or separate into smaller parts or segments
func (j *Job) Split(ctx context.Context, number int) (configs []*config.JSON, err error) {
	for _, x := range j.conf.Xlsxs {
		conf, _ := config.NewJSONFromString("{}")
		conf.Set("path", x.Path)

		for i, v := range x.Sheets {
			xlsxConfig := j.conf.InConfig
			xlsxConfig.Sheet = v
			conf.Set("content."+strconv.Itoa(i), xlsxConfig)
		}
		configs = append(configs, conf)
	}
	return
}
