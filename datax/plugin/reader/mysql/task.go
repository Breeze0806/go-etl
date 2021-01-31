package mysql

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

//Task 任务
type Task struct {
	*plugin.BaseTask
}

//Init 初始化
func (t *Task) Init(ctx context.Context) (err error) {
	return
}

//Destroy 销毁
func (t *Task) Destroy(ctx context.Context) (err error) {
	return
}

//StartRead 开始读
func (t *Task) StartRead(ctx context.Context, sender plugin.RecordSender) error {
	return nil
}
