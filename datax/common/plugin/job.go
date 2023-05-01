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

//Job 工作
type Job interface {
	Plugin
	//工作ID
	JobID() int64
	//设置工作ID
	SetJobID(jobID int64)
	Collector() JobCollector   //todo 工作采集器目前未使用
	SetCollector(JobCollector) //todo  设置工作采集器目前未使用
}

//BaseJob 基础工作，用于辅助和简化工作接口的实现
type BaseJob struct {
	*BasePlugin

	id        int64
	collector JobCollector
}

//NewBaseJob 获取NewBaseJob
func NewBaseJob() *BaseJob {
	return &BaseJob{
		BasePlugin: NewBasePlugin(),
	}
}

//JobID 工作ID
func (b *BaseJob) JobID() int64 {
	return b.id
}

//SetJobID 设置工作ID
func (b *BaseJob) SetJobID(jobID int64) {
	b.id = jobID
}

//Collector 采集器
func (b *BaseJob) Collector() JobCollector {
	return b.collector
}

//SetCollector 设置采集器
func (b *BaseJob) SetCollector(collector JobCollector) {
	b.collector = collector
}
