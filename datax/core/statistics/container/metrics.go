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

package container

import (
	"sync"

	"github.com/Breeze0806/go/encoding"
)

// Metrics - Metrics in JSON format
type Metrics struct {
	sync.RWMutex

	metricJSON *encoding.JSON
}

// NewMetrics - New instance of metrics in JSON format
func NewMetrics() *Metrics {
	j, _ := encoding.NewJSONFromString("{}")
	return &Metrics{
		metricJSON: j,
	}
}

// JSON - Represents data in JSON format
func (m *Metrics) JSON() *encoding.JSON {
	m.RLock()
	defer m.RUnlock()
	return m.metricJSON
}

// Set - Sets the value at the specified path
func (m *Metrics) Set(path string, value any) error {
	m.Lock()
	defer m.Unlock()
	return m.metricJSON.Set(path, value)
}

// Get - Gets the value at the specified path
func (m *Metrics) Get(key string) *encoding.JSON {
	m.RLock()
	defer m.RUnlock()
	j, err := m.metricJSON.GetJSON(key)
	if err != nil {
		return nil
	}
	return j
}
