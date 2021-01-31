package plugin

//Job 工作
type Job interface {
	Plugin
	Collector() JobCollector   //todo 工作采集器目前未使用
	SetCollector(JobCollector) //todo  设置工作采集器目前未使用
}

//BaseJob 基础工作，用于辅助和简化工作接口的实现
type BaseJob struct {
	*BasePlugin

	collector JobCollector
}

//NewBaseJob 获取NewBaseJob
func NewBaseJob() *BaseJob {
	return &BaseJob{
		BasePlugin: NewBasePlugin(),
	}
}

//Collector 采集器
func (b *BaseJob) Collector() JobCollector {
	return b.collector
}

//SetCollector 设置采集器
func (b *BaseJob) SetCollector(collector JobCollector) {
	b.collector = collector
}
