package reader

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

//Task 任务
type Task interface {
	plugin.Task

	//StartRead 开始从sender中读取
	StartRead(ctx context.Context, sender plugin.RecordSender) error
}
