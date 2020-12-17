package plugin

import (
	"github.com/Breeze0806/go-etl/datax/common/element"
)

type RecordReceiver interface {
	GetFromReader() (element.Record, error)
	Shutdown() error
}
