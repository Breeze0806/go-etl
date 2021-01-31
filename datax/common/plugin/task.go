package plugin

//Task 任务接口
type Task interface {
	Plugin

	//任务信息收集器，todo 未使用
	TaskCollector() TaskCollector
	//设置任务信息收集器，todo 未使用
	SetTaskCollector(collector TaskCollector)
	//任务ID
	TaskID() int
	//设置任务ID
	SetTaskID(taskID int)
	//任务组ID
	TaskGroupID() int
	//设置任务组ID
	SetTaskGroupID(taskGroupID int)
}

//BaseTask 基础任务，用于辅助和简化任务接口的实现
type BaseTask struct {
	*BasePlugin

	taskID      int
	taskGroupID int
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
func (b *BaseTask) TaskID() int {
	return b.taskID
}

//SetTaskID 设置任务ID
func (b *BaseTask) SetTaskID(taskID int) {
	b.taskID = taskID
}

//TaskGroupID 任务组ID
func (b *BaseTask) TaskGroupID() int {
	return b.taskGroupID
}

//SetTaskGroupID 设置任务组ID
func (b *BaseTask) SetTaskGroupID(taskGroupID int) {
	b.taskGroupID = taskGroupID
}
