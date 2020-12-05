package writer

import "github.com/Breeze0806/go-etl/datax/common/plugin"

type Task interface {
	StartWrite(receiver plugin.RecordReceiver) error
	SupportFailOver() bool
}
