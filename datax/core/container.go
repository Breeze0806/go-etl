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

// Container - A component that holds or encapsulates other components
type Container interface {
	Start() error
	Metrics() *container.Metrics
}

// BaseContainer - A fundamental container that provides basic functionalities
type BaseCotainer struct {
	conf    *config.JSON
	metrics *container.Metrics
}

// NewBaseContainer - Creates a new instance of the base container
func NewBaseCotainer() *BaseCotainer {
	return &BaseCotainer{}
}

// SetMetrics - Sets the metrics for the container
func (b *BaseCotainer) SetMetrics(metrics *container.Metrics) {
	b.metrics = metrics
}

// Metrics - Represents the measurements or quantifications of some aspect of the container's operation
func (b *BaseCotainer) Metrics() *container.Metrics {
	return b.metrics
}

// SetConfig - Sets the JSON configuration for the container
func (b *BaseCotainer) SetConfig(conf *config.JSON) {
	b.conf = conf
}

// Config - Represents the JSON-formatted configuration settings for the container
func (b *BaseCotainer) Config() *config.JSON {
	return b.conf
}
