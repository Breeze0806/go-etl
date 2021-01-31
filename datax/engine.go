package datax

import (
	"context"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/core"
	"github.com/Breeze0806/go-etl/datax/core/job"
	"github.com/Breeze0806/go-etl/datax/core/taskgroup"
)

//Model 模式
type Model string

//容器工作模式
var (
	ModelJob       Model = "job"       //以工作为单位工作
	ModelTaskGroup Model = "taskGroup" //以任务组为单位工作
)

//IsJob 是否以工作为单位工作
func (m Model) IsJob() bool {
	return m == ModelJob
}

//IsTaskGroup 以任务组为单位工作
func (m Model) IsTaskGroup() bool {
	return m == ModelTaskGroup
}

//Engine 执行引擎
type Engine struct {
	core.Container
	ctx  context.Context
	conf *config.JSON
}

//NewEngine 通过上下文ctx以及JSON配置conf创建新执行引擎
func NewEngine(ctx context.Context, conf *config.JSON) *Engine {
	return &Engine{
		ctx:  ctx,
		conf: conf,
	}
}

//Start 启动
func (e *Engine) Start() (err error) {
	model := Model(e.conf.GetStringOrDefaullt(coreconst.DataxCoreContainerModel, string(ModelJob)))
	switch {
	case model.IsJob():
		e.Container, err = job.NewContainer(e.ctx, e.conf)
		if err != nil {
			return
		}
	case model.IsTaskGroup():
		e.Container, err = taskgroup.NewContainer(e.ctx, e.conf)
		if err != nil {
			return
		}
	default:
		return fmt.Errorf("model is %v", model)
	}

	return e.Container.Start()
}
