package plugin

import "github.com/Breeze0806/go-etl/element"

//TaskCollector 任务收集器
//todo 当前未使用
type TaskCollector interface {
	CollectDirtyRecordWithError(record element.Record, err error)
	CollectDirtyRecordWithMsg(record element.Record, msgErr string)
	CollectDirtyRecord(record element.Record, err error, msgErr string)
	CollectMessage(key string, value string)
}
