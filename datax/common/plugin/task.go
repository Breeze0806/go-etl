package plugin

type Task interface {
	Plugin
	TaskCollector() TaskCollector
	SetTaskCollector(collector TaskCollector)
	TaskId() int
	SetTaskId(taskId int)
	TaskGroupId() int
	SetTaskGroupId(taskGroupId int)
}

type BaseTask struct {
	*BasePlugin

	taskId      int
	taskGroupId int
	collector   TaskCollector
}

func NewBaseTask() *BaseTask {
	return &BaseTask{
		BasePlugin: NewBasePlugin(),
	}
}

func (b *BaseTask) TaskCollector() TaskCollector {
	return b.collector
}

func (b *BaseTask) SetTaskCollector(collector TaskCollector) {
	b.collector = collector
}

func (b *BaseTask) TaskId() int {
	return b.taskId
}

func (b *BaseTask) SetTaskId(taskId int) {
	b.taskId = taskId
}

func (b *BaseTask) TaskGroupId() int {
	return b.taskGroupId
}

func (b *BaseTask) SetTaskGroupId(taskGroupId int) {
	b.taskGroupId = taskGroupId
}
