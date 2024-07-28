package sqlite3

import (
	"github.com/Breeze0806/go-etl/config"
	spireader "github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/dbms"
	"github.com/Breeze0806/go-etl/storage/database"

	// sqlite3 storage - Sqlite3 database storage
	_ "github.com/Breeze0806/go-etl/storage/database/sqlite3"
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
	job := &Job{
		Job: dbms.NewJob(dbms.NewBaseDbHandler(func(name string, conf *config.JSON) (q dbms.Querier, err error) {
			if q, err = database.Open(name, conf); err != nil {
				return nil, err
			}
			return
		}, nil)),
	}
	job.SetPluginConf(r.pluginConf)
	return job
}

// Task returns the smallest execution unit obtained by maximizing the split of a Job
func (r *Reader) Task() spireader.Task {
	task := &Task{
		Task: dbms.NewTask(dbms.NewBaseDbHandler(func(name string, conf *config.JSON) (q dbms.Querier, err error) {
			if q, err = database.Open(name, conf); err != nil {
				return nil, err
			}
			return
		}, nil)),
	}
	task.SetPluginConf(r.pluginConf)
	return task
}
