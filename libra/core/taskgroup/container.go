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

package taskgroup

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
)

//Option 容器选项
type Option func(c *Container) error

//Container 容器
type Container struct {
	jobID string
	conf  *config.JSON
	ctx   context.Context
}

//NewContainer 容器
func NewContainer(ctx context.Context, options ...Option) (c *Container, err error) {
	c = &Container{
		ctx: ctx,
	}
	for _, option := range options {
		if err = option(c); err != nil {
			return
		}
	}
	return
}

func setConfig(conf *config.JSON) Option {
	return func(c *Container) error {
		c.conf = conf
		return nil
	}
}

func setJobID(jobID string) Option {
	return func(c *Container) error {
		c.jobID = jobID
		return nil
	}
}

//Do 执行
func (c *Container) Do() (err error) {
	return
}
