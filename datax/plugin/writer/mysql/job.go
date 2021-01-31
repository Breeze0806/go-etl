package mysql

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

//Job 工作
type Job struct {
	*plugin.BaseJob
}

//Init 初始化
func (j *Job) Init(ctx context.Context) (err error) {
	return
}

//Destroy 销毁
func (j *Job) Destroy(ctx context.Context) (err error) {
	return
}

//Split 切分任务
func (j *Job) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return nil, nil
}
