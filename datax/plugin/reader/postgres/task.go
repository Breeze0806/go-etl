package postgres

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
)

//Task 任务
type Task struct {
	*rdbm.Task
}

//StartRead 开始读
func (t *Task) StartRead(ctx context.Context, sender plugin.RecordSender) (err error) {
	return rdbm.StartRead(ctx, rdbm.NewBaseBatchReader(t.Task, "", nil), sender)
}
