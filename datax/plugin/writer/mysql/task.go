package mysql

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

//Task 任务
type Task struct {
	*writer.BaseTask
}

//Init 初始化
func (t *Task) Init(ctx context.Context) (err error) {
	return
}

//Destroy 销毁
func (t *Task) Destroy(ctx context.Context) (err error) {
	return
}

//StartWrite 开始写
func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error {
	return nil
}
