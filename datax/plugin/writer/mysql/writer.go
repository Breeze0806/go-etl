package mysql

import (
	"path/filepath"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

func init() {
	writer, err := NewWriter(filepath.Join("resources", "plugin.json"))
	if err != nil {
		panic(err)
	}
	name, err := writer.pluginConf.GetString("name")
	if err != nil {
		panic(err)
	}
	if name != "" {
		panic("name is empty")
	}
	loader.RegisterWriter(name, writer)
}

type Writer struct {
	pluginConf *config.Json
}

func NewWriter(filename string) (w *Writer, err error) {
	w = &Writer{}
	w.pluginConf, err = config.NewJsonFromFile(filename)
	if err != nil {
		return nil, err
	}
	return
}

func (w *Writer) Job() writer.Job {
	job := &Job{
		BaseJob: plugin.NewBaseJob(),
	}
	job.SetPluginConf(w.pluginConf)
	return job
}

func (w *Writer) Task() writer.Task {
	task := &Task{
		BaseTask: writer.NewBaseTask(),
	}
	task.SetPluginConf(w.pluginConf)
	return task
}
