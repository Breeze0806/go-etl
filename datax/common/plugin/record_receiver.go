package plugin

import (
	"github.com/Breeze0806/go-etl/element"
)

type RecordReceiver interface {
	GetFromReader() (element.Record, error)
	Shutdown() error
}
