package runner

import (
	"context"
	"fmt"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
)

type Reader struct {
	*baseRunner
	sender plugin.RecordSender
	task   reader.Task
}

func NewReader(task reader.Task, sender plugin.RecordSender) *Reader {
	return &Reader{
		baseRunner: &baseRunner{},
		sender:     sender,
		task:       task,
	}
}

func (r *Reader) Plugin() plugin.Task {
	return r.task
}

func (r *Reader) Run(ctx context.Context) (err error) {
	defer func() {
		if destroyErr := r.task.Destroy(ctx); destroyErr != nil {
			log.Errorf("task destroy fail, err: %v", destroyErr)
		}
	}()

	if err = r.task.Init(ctx); err != nil {
		return fmt.Errorf("task init fail, err: %v", err)
	}

	if err = r.task.Prepare(ctx); err != nil {
		return fmt.Errorf("task prepare fail, err: %v", err)
	}

	if err = r.task.StartRead(ctx, r.sender); err != nil {
		return fmt.Errorf("task startRead fail, err: %v", err)
	}

	if err = r.task.Post(ctx); err != nil {
		return fmt.Errorf("task post fail, err: %v", err)
	}
	return
}

func (r *Reader) Shutdown() error {
	return r.sender.Shutdown()
}
