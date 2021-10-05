package rdbm

import (
	"context"
	"database/sql"
	"time"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/plugin/writer"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

const (
	ExecModeNormal = "Normal"
	ExecModeTx     = "Tx"
	ExecModeStmt   = "Stmt"
)

type BaseBatchWriter struct {
	task     *Task
	execMode string
	opts     *database.ParameterOptions
}

func NewBaseBatchWriter(task *Task, execMode string, opts *sql.TxOptions) *BaseBatchWriter {
	w := &BaseBatchWriter{
		task:     task,
		execMode: execMode,
	}
	w.opts = &database.ParameterOptions{
		Table:     task.Table,
		Mode:      task.Config.GetWriteMode(),
		TxOptions: opts,
	}
	return w
}

func (b *BaseBatchWriter) JobID() int64 {
	return b.task.JobID()
}

func (b *BaseBatchWriter) TaskGroupID() int64 {
	return b.task.TaskGroupID()
}

func (b *BaseBatchWriter) TaskID() int64 {
	return b.task.TaskID()
}

func (b *BaseBatchWriter) BatchSize() int {
	return b.task.Config.GetBatchSize()
}

func (b *BaseBatchWriter) BatchTimeout() time.Duration {
	return b.task.Config.GetBatchTimeout()
}

func (b *BaseBatchWriter) BatchWrite(ctx context.Context, records []element.Record) error {
	b.opts.Records = records
	defer func() {
		b.opts.Records = nil
	}()
	switch b.execMode {
	case ExecModeTx:
		return b.task.Execer.BatchExecWithTx(ctx, b.opts)
	case ExecModeStmt:
		return b.task.Execer.BatchExecStmtWithTx(ctx, b.opts)
	}
	return b.task.Execer.BatchExec(ctx, b.opts)
}

//StartWrite 通过批量写入器writer和记录接受器receiver将记录写入数据库
func StartWrite(ctx context.Context, w writer.BatchWriter,
	receiver plugin.RecordReceiver) error {
	return writer.StartWrite(ctx, w, receiver)
}
