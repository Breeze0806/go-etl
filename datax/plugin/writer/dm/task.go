package dm

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/dbms"
	"github.com/Breeze0806/go-etl/storage/database"
)

var execModeMap = map[string]string{
	database.WriteModeInsert: dbms.ExecModeNormal,
}

func execMode(writeMode string) string {
	if mode, ok := execModeMap[writeMode]; ok {
		return mode
	}
	return dbms.ExecModeNormal
}

// Task
type Task struct {
	*dbms.Task
}

// StartWrite
func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) (err error) {
	return dbms.StartWrite(ctx, dbms.NewBaseBatchWriter(t.Task, execMode(t.Config.GetWriteMode()), nil), receiver)
}
