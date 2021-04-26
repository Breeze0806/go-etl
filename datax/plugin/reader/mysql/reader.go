package mysql

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
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
	reader, err := NewReader(_pluginConfig)
	if err != nil {
		panic(err)
	}
	name, err := reader.pluginConf.GetString("name")
	if err != nil {
		panic(err)
	}
	if name == "" {
		panic("name is empty")
	}
	loader.RegisterReader(name, reader)
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

//Job 工作
func (r *Reader) Job() reader.Job {
	job := &Job{
		BaseJob: plugin.NewBaseJob(),
		newQuerier: func(name string, conf *config.JSON) (Querier, error) {
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
		newQuerier: func(name string, conf *config.JSON) (Querier, error) {
			return database.Open(name, conf)
		},
	}
	task.SetPluginConf(r.pluginConf)
	return task
}
