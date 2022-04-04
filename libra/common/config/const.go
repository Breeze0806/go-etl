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

package config

//libra参数
const (
	LibraCoreContainerJobMaster          = "core.container.job.master"
	LibraCoreContainerJobSlave           = "core.container.job.slave"
	LibraCoreContainerJobMasterName      = "core.container.job.master.name"
	LibraCoreContainerJobMasterParameter = "core.container.job.master.parameter"
	LibraCoreContainerJobSlaveName       = "core.container.job.slave.name"
	LibraCoreContainerJobSlaveParameter  = "core.container.job.slave.parameter"

	LibraJobContent                = "job.content"
	LibraJobContentMasterName      = "job.content.0.master.name"
	LibraJobContentMasterParameter = "job.content.0.master.parameter"
	LibraJobContentSlaveName       = "job.content.0.slave.name"
	LibraJobContentSlaveParameter  = "job.content.0.slave.parameter"
	LibraJobSetting                = "job.setting"
)
