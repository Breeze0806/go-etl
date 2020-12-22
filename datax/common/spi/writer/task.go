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

type BaseTask struct {
	*plugin.BaseTask
}

func NewBaseTask() *BaseTask {
	return &BaseTask{
		BaseTask: plugin.NewBaseTask(),
	}
}

func (b *BaseTask) SupportFailOver() bool {
	return false
}
