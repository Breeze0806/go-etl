package postgres

import (
	"github.com/Breeze0806/go-etl/config"
	spiwriter "github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/rdbm"
	"github.com/Breeze0806/go-etl/storage/database"

	//postgres storage
	_ "github.com/Breeze0806/go-etl/storage/database/postgres"
)

//Writer 写入器
type Writer struct {
	pluginConf *config.JSON
}

//ResourcesConfig 插件资源配置
func (w *Writer) ResourcesConfig() *config.JSON {
	return w.pluginConf
}

//Job 工作
func (w *Writer) Job() spiwriter.Job {
	job := &Job{
		Job: rdbm.NewJob(rdbm.NewBaseDbHandler(
			func(name string, conf *config.JSON) (e rdbm.Execer, err error) {
				if e, err = database.Open(name, conf); err != nil {
					return nil, err
				}
				return
			}, nil)),
	}
	job.SetPluginConf(w.pluginConf)
	return job
}

//Task 任务
func (w *Writer) Task() spiwriter.Task {
	task := &Task{
		Task: rdbm.NewTask(rdbm.NewBaseDbHandler(
			func(name string, conf *config.JSON) (e rdbm.Execer, err error) {
				if e, err = database.Open(name, conf); err != nil {
					return nil, err
				}
				return
			}, nil)),
	}
	task.SetPluginConf(w.pluginConf)
	return task
}
