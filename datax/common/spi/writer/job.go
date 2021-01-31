package writer

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

//Job 工作
type Job interface {
	plugin.Job

	//根据Job进行切分，将原有任务尽量切分成number个任务，主要以配置文件的形式传递给每个任务
	Split(ctx context.Context, number int) ([]*config.JSON, error)
}
