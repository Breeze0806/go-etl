package csv

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/file"
)

type Job struct {
	*file.Job
	conf *JobConfig
}

func NewJob() *Job {
	return &Job{
		Job: file.NewJob(),
	}
}

func (j *Job) Init(ctx context.Context) (err error) {
	j.conf, err = NewJobConfig(j.PluginJobConf())
	return
}

func (j *Job) Split(ctx context.Context, number int) (configs []*config.JSON, err error) {
	for _, v := range j.conf.Path {
		conf, _ := config.NewJSONFromString("{}")
		if err = conf.Set("path", v); err != nil {
			return
		}
		if err = conf.Set("content", j.conf.CsvConfig); err != nil {
			return
		}
		if err = conf.Set("content.batchTimeout", j.conf.GetBatchTimeout().String()); err != nil {
			return
		}

		configs = append(configs, conf)
	}
	return
}
