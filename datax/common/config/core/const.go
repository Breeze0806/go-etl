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

//datax 配置文件路径
var (
	//datax全局配置路劲
	DataxCoreContainerTaskgroupChannel                = "core.container.taskGroup.channel"
	DataxCoreContainerModel                           = "core.container.model"
	DataxCoreContainerJobID                           = "core.container.job.id"
	DataxCoreContainerTraceEnable                     = "core.container.trace.enable"
	DataxCoreContainerJobMode                         = "core.container.job.mode"
	DataxCoreContainerJobReportinterval               = "core.container.job.reportInterval"
	DataxCoreContainerJobSleepinterval                = "core.container.job.sleepInterval"
	DataxCoreContainerJobMaxWorkerNumber              = "core.container.job.maxWorkerNumber"
	DataxCoreContainerTaskGroupID                     = "core.container.taskGroup.id"
	DataxCoreContainerTaskGroupSleepinterval          = "core.container.taskGroup.sleepInterval"
	DataxCoreContainerTaskGroupReportinterval         = "core.container.taskGroup.reportInterval"
	DataxCoreContainerTaskGroupMaxWorkerNumber        = "core.container.taskGroup.maxWorkerNumber"
	DataxCoreContainerTaskFailoverMaxretrytimes       = "core.container.task.failover.maxRetryTimes"
	DataxCoreContainerTaskFailoverRetryintervalinmsec = "core.container.task.failover.retryIntervalInMsec"
	DataxCoreContainerTaskFailoverMaxwaitinmsec       = "core.container.task.failover.maxWaitInMsec"
	DataxCoreDataxserverAddress                       = "core.dataXServer.address"
	DataxCoreDscAddress                               = "core.dsc.address"
	DataxCoreDataxserverTimeout                       = "core.dataXServer.timeout"
	DataxCoreReportDataxLog                           = "core.dataXServer.reportDataxLog"
	DataxCoreReportDataxPerflog                       = "core.dataXServer.reportPerfLog"
	DataxCoreTransportChannelClass                    = "core.transport.channel.class"
	DataxCoreTransportChannelCapacity                 = "core.transport.channel.capacity"
	DataxCoreTransportChannelCapacityByte             = "core.transport.channel.byteCapacity"
	DataxCoreTransportChannelSpeed                    = "core.transport.channel.speed"
	DataxCoreTransportChannelSpeedByte                = "core.transport.channel.speed.byte"
	DataxCoreTransportChannelSpeedRecord              = "core.transport.channel.speed.record"
	DataxCoreTransportChannelFlowcontrolinterval      = "core.transport.channel.flowControlInterval"
	DataxCoreTransportExchangerBuffersize             = "core.transport.exchanger.bufferSize"
	DataxCoreTransportRecordClass                     = "core.transport.record.class"
	DataxCoreStatisticsCollectorPluginTaskclass       = "core.statistics.collector.plugin.taskClass"
	DataxCoreStatisticsCollectorPluginMaxdirtynum     = "core.statistics.collector.plugin.maxDirtyNumber"
	DataxJobContentReaderName                         = "job.content.0.reader.name"
	DataxJobContentReaderParameter                    = "job.content.0.reader.parameter"
	DataxJobContentWriterName                         = "job.content.0.writer.name"
	DataxJobContentWriterParameter                    = "job.content.0.writer.parameter"
	DataxJobJobinfo                                   = "job.jobInfo"
	DataxJobContent                                   = "job.content"
	DataxJobContentTransformer                        = "job.content.0.transformer"
	DataxJobSetting                                   = "job.setting"
	DataxJobSettingKeyversion                         = "job.setting.keyVersion"
	DataxJobSettingSpeedByte                          = "job.setting.speed.byte"
	DataxJobSettingSpeedRecord                        = "job.setting.speed.record"
	DataxJobSettingSpeedChannel                       = "job.setting.speed.channel"
	DataxJobSettingSpeed                              = "job.setting.speed"
	DataxJobSettingErrorlimit                         = "job.setting.errorLimit"
	DataxJobSettingErrorlimitRecord                   = "job.setting.errorLimit.record"
	DataxJobSettingErrorlimitPercent                  = "job.setting.errorLimit.percentage"
	DataxJobSettingDryrun                             = "job.setting.dryRun"
	DataxJobPreHandlerPluginType                      = "job.preHandler.pluginType"
	DataxJobPreHandlerPluginName                      = "job.preHandler.pluginName"
	DataxJobPostHandlerPluginType                     = "job.postHandler.pluginType"
	DataxJobPostHandlerPluginName                     = "job.postHandler.pluginName"
	//datax局部配置路径
	JobWriter                                 = "writer"
	JobReader                                 = "reader"
	JobTransformer                            = "transformer"
	JobReaderName                             = "reader.name"
	JobReaderParameter                        = "reader.parameter"
	JobWriterName                             = "writer.name"
	JobWriterParameter                        = "writer.parameter"
	TransformerParameterColumnindex           = "parameter.columnIndex"
	TransformerParameterParas                 = "parameter.paras"
	TransformerParameterContext               = "parameter.context"
	TransformerParameterCode                  = "parameter.code"
	TransformerParameterExtrapackage          = "parameter.extraPackage"
	TaskID                                    = "taskId"
	JobReaderParameterLoadBalanceResourceMark = "reader.parameter.loadBalanceResourceMark"
	JobWriterParameterLoadBalanceResourceMark = "writer.parameter.loadBalanceResourceMark"
)
