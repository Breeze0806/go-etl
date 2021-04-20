package plugin

import (
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/core/statistics/communication"
)

//DefaultJobCollector 默认工作收集器
type DefaultJobCollector struct{}

//NewDefaultJobCollector 创建默认工作收集器
func NewDefaultJobCollector(*communication.Communication) plugin.JobCollector {
	return &DefaultJobCollector{}
}

//MessageMap 空方法
func (d *DefaultJobCollector) MessageMap() map[string][]string {
	return nil
}

//MessageByKey 空方法
func (d *DefaultJobCollector) MessageByKey(key string) []string {
	return nil
}
