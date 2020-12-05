package plugin

type Job interface {
	Defalut
	Collector() JobCollector
	SetCollector(JobCollector)
}

type BaseJob struct {
	*BaseDefalut
	collector JobCollector
}

func (b *BaseJob) Collector() JobCollector {
	return b.collector
}

func (b *BaseJob) SetCollector(collector JobCollector) {
	b.collector = collector
}
