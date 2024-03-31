package sqlite3

import (
	"github.com/Breeze0806/go-etl/config"
	spireader "github.com/Breeze0806/go-etl/datax/common/spi/reader"
)

// Reader is uesed to extract data from data source
type Reader struct {
	pluginConf *config.JSON
}

// ResourcesConfig returns the configuration of the data source to initiate the reader.
func (r *Reader) ResourcesConfig() *config.JSON {
	return r.pluginConf
}

// Job returns a description of how the reader extracts data from the data source.
func (r *Reader) Job() spireader.Job {
	// todo like below
	//job := NewJob()
	//job.SetPluginConf(r.pluginConf)
	//return job
	return nil
}

// Task returns the smallest execution unit obtained by maximizing the split of a Job
func (r *Reader) Task() spireader.Task {
	// todo like below
	//task := fNewTask()
	//task.SetPluginConf(r.pluginConf)
	//return task
	return nil
}
