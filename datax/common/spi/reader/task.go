package reader

import "github.com/Breeze0806/go-etl/datax/common/plugin"

type Task interface {
	plugin.Task
	StartRead(sender plugin.RecordSender) error
}
