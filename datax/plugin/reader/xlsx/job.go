package xlsx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/file"
	"github.com/Breeze0806/go-etl/storage/stream/file/xlsx"
)

type Job struct {
	*file.Job

	conf *Config
}

func NewJob() *Job {
	return &Job{
		Job: file.NewJob(),
	}
}

func (j *Job) Init(ctx context.Context) (err error) {
	j.conf, err = NewConfig(j.PluginJobConf())
	return
}

func (j *Job) Split(ctx context.Context, number int) (configs []*config.JSON, err error) {
	for _, x := range j.conf.Xlsxs {
		conf, _ := config.NewJSONFromString("{}")
		if err = conf.Set("path", x.Path); err != nil {
			return
		}
		for i, v := range x.Sheets {
			xlsxConfig := xlsx.InConfig{
				Sheet:   v,
				Columns: j.conf.Columns,
			}
			if err = conf.Set("content."+strconv.Itoa(i), xlsxConfig); err != nil {
				return
			}
		}
		fmt.Println(conf.String())
		configs = append(configs, conf)
	}
	return
}
