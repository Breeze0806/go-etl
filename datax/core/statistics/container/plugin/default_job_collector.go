package plugin

import (
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/core/statistics/communication"
)

type DefaultJobCollector struct {
}

func NewDefaultJobCollector(communication.Communication) plugin.JobCollector {
	return &DefaultJobCollector{}
}

func (d *DefaultJobCollector) MessageMap() map[string][]string {
	return nil
}

func (d *DefaultJobCollector) MessageByKey(key string) []string {
	return nil
}
