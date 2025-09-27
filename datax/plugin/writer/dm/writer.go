package dm

import (
	"github.com/Breeze0806/go-etl/config"
	spiwriter "github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/dbms"
	"github.com/Breeze0806/go-etl/storage/database"
)

// Writer Writer
type Writer struct {
	pluginConf *config.JSON
}

// ResourcesConfig Plugin Resource Configuration
func (w *Writer) ResourcesConfig() *config.JSON {
	return w.pluginConf
}

// Job Job
func (w *Writer) Job() spiwriter.Job {
	job := &Job{
		Job: dbms.NewJob(dbms.NewBaseDbHandler(
			func(name string, conf *config.JSON) (e dbms.Execer, err error) {
				if e, err = database.Open(name, conf); err != nil {
					return nil, err
				}
				return
			}, nil)),
	}
	job.SetPluginConf(w.pluginConf)
	return job
}

// Task Task
func (w *Writer) Task() spiwriter.Task {
	task := &Task{
		Task: dbms.NewTask(dbms.NewBaseDbHandler(
			func(name string, conf *config.JSON) (e dbms.Execer, err error) {
				if e, err = database.Open(name, conf); err != nil {
					return nil, err
				}
				return
			}, nil)),
	}
	task.SetPluginConf(w.pluginConf)
	return task
}
