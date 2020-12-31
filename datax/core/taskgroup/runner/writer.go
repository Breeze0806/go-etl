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
	describe string
}

func NewWriter(task writer.Task, receiver plugin.RecordReceiver, taskKey string) *Writer {
	return &Writer{
		baseRunner: &baseRunner{},
		receiver:   receiver,
		task:       task,
		describe:   taskKey,
	}
}

func (w *Writer) Plugin() plugin.Task {
	return w.task
}

func (w *Writer) Run(ctx context.Context) (err error) {
	defer func() {
		log.Debugf("datax writer runner %v starts to destroy", w.describe)
		if destroyErr := w.task.Destroy(ctx); destroyErr != nil {
			log.Errorf("task destroy fail, err: %v", destroyErr)
		}
	}()
	log.Debugf("datax writer runner %v starts to init", w.describe)
	if err = w.task.Init(ctx); err != nil {
		log.Errorf("task init fail, err: %v", err)
		return
	}

	log.Debugf("datax writer runner %v starts to prepare", w.describe)
	if err = w.task.Prepare(ctx); err != nil {
		log.Errorf("task prepare fail, err: %v", err)
		return
	}

	log.Debugf("datax writer runner %v starts to StartWrite", w.describe)
	if err = w.task.StartWrite(ctx, w.receiver); err != nil {
		log.Errorf("task startWrite fail, err: %v", err)
		return
	}

	log.Debugf("datax writer runner %v starts to post", w.describe)
	if err = w.task.Post(ctx); err != nil {
		log.Errorf("task post fail, err: %v", err)
		return
	}
	return
}

func (w *Writer) Shutdown() error {
	return w.receiver.Shutdown()
}
