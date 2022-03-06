package postgres

import (
	"github.com/Breeze0806/go-etl/config"
	spireader "github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
	"github.com/Breeze0806/go-etl/storage/database"

	//postgres storage
	_ "github.com/Breeze0806/go-etl/storage/database/postgres"
)

//Reader 读取器
type Reader struct {
	pluginConf *config.JSON
}

//ResourcesConfig 插件资源配置
func (r *Reader) ResourcesConfig() *config.JSON {
	return r.pluginConf
}

//Job 工作
func (r *Reader) Job() spireader.Job {
	job := &Job{
		Job: rdbm.NewJob(
			rdbm.NewBaseDbHandler(func(name string, conf *config.JSON) (q rdbm.Querier, err error) {
				if q, err = database.Open(name, conf); err != nil {
					return nil, err
				}
				return
			}, nil)),
	}
	job.SetPluginConf(r.pluginConf)
	return job
}

//Task 任务
func (r *Reader) Task() spireader.Task {
	task := &Task{
		Task: rdbm.NewTask(rdbm.NewBaseDbHandler(func(name string, conf *config.JSON) (q rdbm.Querier, err error) {
			if q, err = database.Open(name, conf); err != nil {
				return nil, err
			}
			return
		}, nil)),
	}
	task.SetPluginConf(r.pluginConf)
	return task
}
