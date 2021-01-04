package datax

import (
	"context"
	"fmt"

	"github.com/Breeze0806/go-etl/datax/common/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/core"
	"github.com/Breeze0806/go-etl/datax/core/job"
	"github.com/Breeze0806/go-etl/datax/core/taskgroup"
)

type Model string

var (
	ModelJob       Model = "job"
	ModelTaskGroup Model = "taskGroup"
)

func (m Model) IsJob() bool {
	return m == ModelJob
}

func (m Model) IsTaskGroup() bool {
	return m == ModelTaskGroup
}

type Engine struct {
	core.Container
	ctx  context.Context
	conf *config.Json
}

func NewEngine(ctx context.Context, conf *config.Json) *Engine {
	return &Engine{
		ctx:  ctx,
		conf: conf,
	}
}

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
