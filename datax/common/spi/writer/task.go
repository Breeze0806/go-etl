package writer

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

type Task interface {
	plugin.Task
	StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error
	SupportFailOver() bool
}
