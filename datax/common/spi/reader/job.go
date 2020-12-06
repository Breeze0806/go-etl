package reader

import (
	"github.com/Breeze0806/go-etl/datax/common/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

type Job interface {
	plugin.Job
	Split(int) ([]*config.Json, error)
}
