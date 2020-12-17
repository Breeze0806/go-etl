package runner

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

type Writer struct {
	*baseRunner
	receiver plugin.RecordReceiver
	task     writer.Task
}

func NewWriter(ctx context.Context, task writer.Task, receiver plugin.RecordReceiver) *Writer {
	return &Writer{
		baseRunner: &baseRunner{
			ctx: ctx,
		},
		receiver: receiver,
		task:     task,
	}
}

func (w *Writer) Plugin() plugin.Task {
	return w.task
}

func (w *Writer) Run() (err error) {
	defer func() {
		if err = w.task.Destroy(w.ctx); err != nil {
			log.Errorf("task destroy fail, err: %v", err)
		}
	}()
	if err = w.task.Init(w.ctx); err != nil {
		return err
	}

	if err = w.task.Prepare(w.ctx); err != nil {
		return err
	}

	if err = w.task.StartWrite(w.ctx, w.receiver); err != nil {
		return err
	}

	if err = w.task.Post(w.ctx); err != nil {
		return err
	}
	return
}

func (w *Writer) Shutdown() error {
	return w.receiver.Shutdown()
}
