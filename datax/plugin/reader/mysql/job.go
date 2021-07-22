package postgres

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
)

//Job 工作
type Job struct {
	*rdbm.Job
}

type DbHandler struct {
	newQuerier func(name string, conf *config.JSON) (rdbm.Querier, error)
}
