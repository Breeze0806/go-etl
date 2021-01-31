package plugin

//JobCollector 工作信息采集器，用于统计整个工作的进度，错误信息等
//toto 当前未实现监控模块，为此需要在后面来实现这个接口的结构体
type JobCollector interface {
	MessageMap() map[string][]string
	MessageByKey(key string) []string
}
