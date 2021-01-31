package mysql

import (
	"path/filepath"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
)

func init() {
	reader, err := NewReader(filepath.Join("resources", "plugin.json"))
	if err != nil {
		panic(err)
	}
	name, err := reader.pluginConf.GetString("name")
	if err != nil {
		panic(err)
	}
	if name != "" {
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
	}
	job.SetPluginConf(r.pluginConf)
	return job
}

//Task 任务
func (r *Reader) Task() reader.Task {
	task := &Task{
		BaseTask: plugin.NewBaseTask(),
	}
	task.SetPluginConf(r.pluginConf)
	return task
}
