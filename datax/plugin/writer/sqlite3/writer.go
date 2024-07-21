package sqlite3

import (
	"github.com/Breeze0806/go-etl/config"
	spiwriter "github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

// Writer is uesed to load data into data source
type Writer struct {
	pluginConf *config.JSON
}

// ResourcesConfig returns the configuration of the data source to initiate the writer.
func (w *Writer) ResourcesConfig() *config.JSON {
	return w.pluginConf
}

// Job returns a description of how the reader extracts data from the data source.
func (w *Writer) Job() spiwriter.Job {
	// todo like below
	//job := NewJob()
	//job.SetPluginConf(w.pluginConf)
	//return job
	return nil
}

// Task returns the smallest execution unit obtained by maximizing the split of a Job
func (w *Writer) Task() spiwriter.Task {
	// todo like below
	//task := NewTask()
	//task.SetPluginConf(w.pluginConf)
	//return task
	return nil
}
