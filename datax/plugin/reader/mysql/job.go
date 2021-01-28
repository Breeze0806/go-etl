package mysql

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

type Job struct {
	*plugin.BaseJob
}

func (j *Job) Init(ctx context.Context) (err error) {
	return
}

func (j *Job) Destroy(ctx context.Context) (err error) {
	return
}

func (j *Job) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return nil, nil
}
