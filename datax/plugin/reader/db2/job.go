package db2

import (
	"github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
)

//Job 工作
type Job struct {
	*rdbm.Job
}