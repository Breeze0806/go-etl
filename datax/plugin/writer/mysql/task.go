package mysql

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

type Task struct {
	*writer.BaseTask
}

func (t *Task) Init(ctx context.Context) (err error) {
	return
}

func (t *Task) Destroy(ctx context.Context) (err error) {
	return
}

func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error {
	return nil
}
