package mysql

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/rdbm"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/Breeze0806/go-etl/storage/database/mysql"
)

const maxNumPlaceholder = 65535

var execModeMap = map[string]string{
	database.WriteModeInsert: rdbm.ExecModeNormal,
	mysql.WriteModeReplace:   rdbm.ExecModeNormal,
}

func execMode(writeMode string) string {
	if mode, ok := execModeMap[writeMode]; ok {
		return mode
	}
	return rdbm.ExecModeNormal
}

//Task 任务
type Task struct {
	*rdbm.Task
}

type batchWriter struct {
	*rdbm.BaseBatchWriter
}

func (b *batchWriter) BatchSize() (size int) {
	size = maxNumPlaceholder / len(b.Task.Table.Fields())
	if b.Task.Config.GetBatchSize() < size {
		size = b.Task.Config.GetBatchSize()
	}
	return
}

//StartWrite 开始写
func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) (err error) {
	writer := &batchWriter{
		BaseBatchWriter: rdbm.NewBaseBatchWriter(t.Task, execMode(t.Config.GetWriteMode()), nil),
	}
	return rdbm.StartWrite(ctx, writer, receiver)
}
