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

// Job: a unit of work
type Job interface {
	Plugin
	// Job ID: a unique identifier for a job
	JobID() int64
	// Set Job ID: a function or method to set the ID of a job
	SetJobID(jobID int64)
	Collector() JobCollector   // todo: The job collector is currently not in use
	SetCollector(JobCollector) // todo: The function or method to set the job collector is currently not in use
}

// BaseJob: a fundamental job class that assists and simplifies the implementation of job interfaces
type BaseJob struct {
	*BasePlugin

	id        int64
	collector JobCollector
}

// NewBaseJob: a function or method to acquire a new instance of BaseJob
func NewBaseJob() *BaseJob {
	return &BaseJob{
		BasePlugin: NewBasePlugin(),
	}
}

// JobID: the identifier for a job
func (b *BaseJob) JobID() int64 {
	return b.id
}

// SetJobID: a function or method to set the ID of a job
func (b *BaseJob) SetJobID(jobID int64) {
	b.id = jobID
}

// Collector: a component or system that collects data or information
func (b *BaseJob) Collector() JobCollector {
	return b.collector
}

// SetCollector: a function or method to set or configure a collector
func (b *BaseJob) SetCollector(collector JobCollector) {
	b.collector = collector
}
