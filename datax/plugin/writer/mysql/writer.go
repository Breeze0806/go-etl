package mysql

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/rdbm"
	"github.com/Breeze0806/go-etl/storage/database"

	//msyql driver
	_ "github.com/Breeze0806/go-etl/storage/database/mysql"
)

var _pluginConfig string

func init() {
	var err error
	_pluginConfig, err = rdbm.RegisterWriter(
		func(path string) (rdbm.Writer, error) {
			return NewWriter(path)
		})
	if err != nil {
		panic(err)
	}
}

//Writer 写入器
type Writer struct {
	pluginConf *config.JSON
}

//NewWriter 创建写入器
func NewWriter(filename string) (w *Writer, err error) {
	w = &Writer{}
	w.pluginConf, err = config.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	return
}

//ResourcesConfig 插件资源配置
func (w *Writer) ResourcesConfig() *config.JSON {
	return w.pluginConf
}

//Job 工作
func (w *Writer) Job() writer.Job {
	job := &Job{
		BaseJob: plugin.NewBaseJob(),
		newExecer: func(name string, conf *config.JSON) (rdbm.Execer, error) {
			return database.Open(name, conf)
		},
	}
	job.SetPluginConf(w.pluginConf)
	return job
}

//Task 任务
func (w *Writer) Task() writer.Task {
	task := &Task{
		BaseTask: writer.NewBaseTask(),
		newExecer: func(name string, conf *config.JSON) (rdbm.Execer, error) {
			return database.Open(name, conf)
		},
	}
	task.SetPluginConf(w.pluginConf)
	return task
}
