package mysql

import (
	"github.com/Breeze0806/go-etl/datax/plugin/writer/rdbm"
)

//Job 工作
type Job struct {
	*rdbm.Job
}
