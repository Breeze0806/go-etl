package plugin

import "github.com/Breeze0806/go-etl/datax/common/element"

type TaskCollector interface {
	CollectDirtyRecordWithError(record element.Record, err error)
	CollectDirtyRecordWithMsg(record element.Record, msgErr string)
	CollectDirtyRecord(record element.Record, err error, msgErr string)
	CollectMessage(key string, value string)
}
