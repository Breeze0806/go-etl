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

type Reader struct {
	pluginConf *config.JSON
}

func NewReader(filename string) (r *Reader, err error) {
	r = &Reader{}
	r.pluginConf, err = config.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	return
}

func (r *Reader) Job() reader.Job {
	job := &Job{
		BaseJob: plugin.NewBaseJob(),
	}
	job.SetPluginConf(r.pluginConf)
	return job
}

func (r *Reader) Task() reader.Task {
	task := &Task{
		BaseTask: plugin.NewBaseTask(),
	}
	task.SetPluginConf(r.pluginConf)
	return task
}
