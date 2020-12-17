package runner

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
)

type Reader struct {
	*baseRunner
	sender plugin.RecordSender
	task   reader.Task
}

func NewReader(ctx context.Context, task reader.Task, sender plugin.RecordSender) *Reader {
	return &Reader{
		baseRunner: &baseRunner{ctx: ctx},
		sender:     sender,
		task:       task,
	}
}

func (r *Reader) Plugin() plugin.Task {
	return r.task
}

func (r *Reader) Run() (err error) {
	defer func() {
		if err = r.task.Destroy(r.ctx); err != nil {
			log.Errorf("task destroy fail, err: %v", err)
		}
	}()
	if err = r.task.Init(r.ctx); err != nil {
		return err
	}

	if err = r.task.Prepare(r.ctx); err != nil {
		return err
	}

	if err = r.task.StartRead(r.ctx, r.sender); err != nil {
		return err
	}

	if err = r.task.Post(r.ctx); err != nil {
		return err
	}
	return
}

func (r *Reader) Shutdown() error {
	return r.sender.Shutdown()
}
