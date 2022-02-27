package file

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

type Job struct {
	*plugin.BaseJob
}

func NewJob() *Job {
	return &Job{
		plugin.NewBaseJob(),
	}
}

func (j *Job) Destroy(ctx context.Context) (err error) {
	return
}
