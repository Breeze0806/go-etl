package mysql

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

type Task struct {
	*plugin.BaseTask
}

func (t *Task) Init(ctx context.Context) (err error) {
	return
}

func (t *Task) Destroy(ctx context.Context) (err error) {
	return
}

func (t *Task) StartRead(ctx context.Context, sender plugin.RecordSender) error {
	return nil
}
