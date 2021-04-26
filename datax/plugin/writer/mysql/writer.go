package mysql

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/storage/database"
	_ "github.com/Breeze0806/go-etl/storage/database/mysql"
)

var _pluginConfig string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic(fmt.Errorf("fail to get filename"))
	}
	path := filepath.Dir(file)
	_pluginConfig = filepath.Join(path, "resources", "plugin.json")
	writer, err := NewWriter(_pluginConfig)
	if err != nil {
		panic(err)
	}
	name, err := writer.pluginConf.GetString("name")
	if err != nil {
		panic(err)
	}
	if name == "" {
		panic("name is empty")
	}
	loader.RegisterWriter(name, writer)
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

//Job 工作
func (w *Writer) Job() writer.Job {
	job := &Job{
		BaseJob: plugin.NewBaseJob(),
		newExecer: func(name string, conf *config.JSON) (Execer, error) {
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
		newExecer: func(name string, conf *config.JSON) (Execer, error) {
			return database.Open(name, conf)
		},
	}
	task.SetPluginConf(w.pluginConf)
	return task
}
