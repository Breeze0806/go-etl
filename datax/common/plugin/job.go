package plugin

type Job interface {
	Plugin
	Collector() JobCollector
	SetCollector(JobCollector)
}

type BaseJob struct {
	*BasePlugin
	collector JobCollector
}

func (b *BaseJob) Collector() JobCollector {
	return b.collector
}

func (b *BaseJob) SetCollector(collector JobCollector) {
	b.collector = collector
}
