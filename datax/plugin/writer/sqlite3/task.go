package sqlite3

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/dbms"
	"github.com/Breeze0806/go-etl/storage/database"
)

const maxNumPlaceholder = 65535

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

type batchWriter struct {
	*dbms.BaseBatchWriter
}

func (b *batchWriter) BatchSize() (size int) {
	size = maxNumPlaceholder / len(b.Task.Table.Fields())
	if b.Task.Config.GetBatchSize() < size {
		size = b.Task.Config.GetBatchSize()
	}
	return
}

// StartWrite
func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) (err error) {
	return dbms.StartWrite(ctx, &batchWriter{BaseBatchWriter: dbms.NewBaseBatchWriter(t.Task, execMode(t.Config.GetWriteMode()), nil)}, receiver)
}
