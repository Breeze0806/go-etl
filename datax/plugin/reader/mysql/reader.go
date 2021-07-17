package mysql

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
	"github.com/Breeze0806/go-etl/storage/database"
	_ "github.com/Breeze0806/go-etl/storage/database/mysql"
)

var _pluginConfig string

func init() {
	var err error
	if _pluginConfig, err = rdbm.RegisterReader(
		func(filename string) (rdbm.Reader, error) {
			return NewReader(filename)
		}); err != nil {
		panic(err)
	}

}

//Reader 读取器
type Reader struct {
	pluginConf *config.JSON
}

//NewReader 创建读取器
func NewReader(filename string) (r *Reader, err error) {
	r = &Reader{}
	r.pluginConf, err = config.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	return
}

func (r *Reader) ResourcesConfig() *config.JSON {
	return r.pluginConf
}

//Job 工作
func (r *Reader) Job() reader.Job {
	job := &Job{
		BaseJob: plugin.NewBaseJob(),
		newQuerier: func(name string, conf *config.JSON) (rdbm.Querier, error) {
			return database.Open(name, conf)
		},
	}
	job.SetPluginConf(r.pluginConf)
	return job
}

//Task 任务
func (r *Reader) Task() reader.Task {
	task := &Task{
		BaseTask: plugin.NewBaseTask(),
		newQuerier: func(name string, conf *config.JSON) (rdbm.Querier, error) {
			return database.Open(name, conf)
		},
	}
	task.SetPluginConf(r.pluginConf)
	return task
}
