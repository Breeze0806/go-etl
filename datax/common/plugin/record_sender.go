package plugin

import "github.com/Breeze0806/go-etl/datax/common/element"

type RecordSender interface {
	CreateRecord() element.Record
	SendWriter(record element.Record)
	Flush() error
	Terminate() error
	Shutdown() error
}
