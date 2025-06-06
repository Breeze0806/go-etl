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

package file

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

// Job normal file job
type Job struct {
	*plugin.BaseJob
}

// NewJob create normal file job
func NewJob() *Job {
	return &Job{
		plugin.NewBaseJob(),
	}
}

// Destroy destroy normal file job
func (j *Job) Destroy(ctx context.Context) (err error) {
	return
}
