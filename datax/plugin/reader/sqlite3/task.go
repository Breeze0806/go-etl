package sqlite3

import (
	"context"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/dbms"
)

// Task
type Task struct {
	*dbms.Task
}

// StartRead
func (t *Task) StartRead(ctx context.Context, sender plugin.RecordSender) (err error) {
	return dbms.StartRead(ctx, dbms.NewBaseBatchReader(t.Task, "", nil), sender)
}
