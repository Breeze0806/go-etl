package table

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/libra/common/plugin"
)

type taskExecer struct {
	page plugin.PageParam
}

func newTaskExecer(conf *config.JSON, page plugin.PageParam) (*taskExecer, error) {
	return nil, nil
}

func (t *taskExecer) Do() error {
	
	return nil
}
