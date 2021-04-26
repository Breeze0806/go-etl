package plugin

//Task 任务接口
type Task interface {
	Plugin

	//任务信息收集器，todo 未使用
	TaskCollector() TaskCollector
	//设置任务信息收集器，todo 未使用
	SetTaskCollector(collector TaskCollector)

	//工作ID
	JobID() int64
	//设置工作ID
	SetJobID(jobID int64)
	//任务组ID
	TaskGroupID() int64
	//设置任务组ID
	SetTaskGroupID(taskGroupID int64)
	//任务ID
	TaskID() int64
	//设置任务ID
	SetTaskID(taskID int64)
}

//BaseTask 基础任务，用于辅助和简化任务接口的实现
type BaseTask struct {
	*BasePlugin

	jobID       int64
	taskID      int64
	taskGroupID int64
	collector   TaskCollector
}

//NewBaseTask 创建基础任务
func NewBaseTask() *BaseTask {
	return &BaseTask{
		BasePlugin: NewBasePlugin(),
	}
}

//TaskCollector 任务信息收集器
func (b *BaseTask) TaskCollector() TaskCollector {
	return b.collector
}

//SetTaskCollector 设置任务信息收集器
func (b *BaseTask) SetTaskCollector(collector TaskCollector) {
	b.collector = collector
}

//TaskID 任务ID
func (b *BaseTask) TaskID() int64 {
	return b.taskID
}

//SetTaskID 设置任务ID
func (b *BaseTask) SetTaskID(taskID int64) {
	b.taskID = taskID
}

//TaskGroupID 任务组ID
func (b *BaseTask) TaskGroupID() int64 {
	return b.taskGroupID
}

//SetTaskGroupID 设置任务组ID
func (b *BaseTask) SetTaskGroupID(taskGroupID int64) {
	b.taskGroupID = taskGroupID
}

//JobID 工作ID
func (b *BaseTask) JobID() int64 {
	return b.jobID
}

//SetJobID 设置工作ID
func (b *BaseTask) SetJobID(jobID int64) {
	b.jobID = jobID
}
