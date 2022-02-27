package xlsx

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/file"
)

type Job struct {
	*file.Job
}

func NewJob() *Job {
	return &Job{
		Job: file.NewJob(),
	}
}

func (j *Job) Init(ctx context.Context) (err error) {
	_, err = NewConfig(j.PluginJobConf())
	return
}

func (j *Job) Split(ctx context.Context, number int) (configs []*config.JSON, err error) {
	configs, err = j.PluginJobConf().GetConfigArray("xlsxs")
	if err != nil {
		return nil, err
	}
	return
}
