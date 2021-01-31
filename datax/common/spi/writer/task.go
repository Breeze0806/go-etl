package writer

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

//Task 写入任务
type Task interface {
	plugin.Task

	//开始从receiver中读取记录写入
	StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error
	//是否支持故障转移，就是是否在写入后失败重试
	SupportFailOver() bool
}

//BaseTask 基础写入任务，辅助和简化写入任务接口的实现
type BaseTask struct {
	*plugin.BaseTask
}

//NewBaseTask 创建基础任务
func NewBaseTask() *BaseTask {
	return &BaseTask{
		BaseTask: plugin.NewBaseTask(),
	}
}

//SupportFailOver 是否支持故障转移，就是是否在写入后失败重试
func (b *BaseTask) SupportFailOver() bool {
	return false
}
