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

package plugin

import (
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/core/statistics/container"
	"github.com/Breeze0806/go/encoding"
)

// DefaultJobCollector - A default job information collector
type DefaultJobCollector struct {
	metrics *container.Metrics
}

// NewDefaultJobCollector - Creates a new instance of the default job information collector
func NewDefaultJobCollector(metrics *container.Metrics) plugin.JobCollector {
	return &DefaultJobCollector{metrics: metrics}
}

// JSON - Retrieves metrics in JSON format
func (d *DefaultJobCollector) JSON() *encoding.JSON {
	return d.metrics.JSON()
}

// JSONByKey - Retrieves metrics in JSON format based on the given key
func (d *DefaultJobCollector) JSONByKey(key string) *encoding.JSON {
	return d.metrics.Get(key)
}
