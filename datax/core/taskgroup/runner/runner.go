package runner

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

//Runner 运行器
type Runner interface {
	Plugin() plugin.Task           //插件任务
	Shutdown() error               //关闭
	Run(ctx context.Context) error //运行
}

type baseRunner struct {
}
