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

func NewWriter(task writer.Task, receiver plugin.RecordReceiver) *Writer {
	return &Writer{
		baseRunner: &baseRunner{},
		receiver:   receiver,
		task:       task,
	}
}

func (w *Writer) Plugin() plugin.Task {
	return w.task
}

func (w *Writer) Run(ctx context.Context) (err error) {
	defer func() {
		if destroyErr := w.task.Destroy(ctx); destroyErr != nil {
			log.Errorf("task destroy fail, err: %v", destroyErr)
		}
	}()
	if err = w.task.Init(ctx); err != nil {
		log.Errorf("task init fail, err: %v", err)
		return
	}

	if err = w.task.Prepare(ctx); err != nil {
		log.Errorf("task prepare fail, err: %v", err)
		return
	}

	if err = w.task.StartWrite(ctx, w.receiver); err != nil {
		log.Errorf("task startWrite fail, err: %v", err)
		return
	}

	if err = w.task.Post(ctx); err != nil {
		log.Errorf("task post fail, err: %v", err)
		return
	}
	return
}

func (w *Writer) Shutdown() error {
	return w.receiver.Shutdown()
}
