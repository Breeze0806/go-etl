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

package core

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/core/statistics/container"
)

// Container 容器
type Container interface {
	Start() error
	Metrics() *container.Metrics
}

// BaseCotainer 基础容器
type BaseCotainer struct {
	conf    *config.JSON
	metrics *container.Metrics
}

// NewBaseCotainer 创建基础容器
func NewBaseCotainer() *BaseCotainer {
	return &BaseCotainer{}
}

// SetMetrics 设置指标
func (b *BaseCotainer) SetMetrics(metrics *container.Metrics) {
	b.metrics = metrics
}

// Metrics 指标
func (b *BaseCotainer) Metrics() *container.Metrics {
	return b.metrics
}

// SetConfig 设置JSON配置
func (b *BaseCotainer) SetConfig(conf *config.JSON) {
	b.conf = conf
}

// Config JSON配置
func (b *BaseCotainer) Config() *config.JSON {
	return b.conf
}
