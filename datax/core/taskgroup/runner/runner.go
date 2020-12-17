package runner

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

type Runner interface {
	Plugin() plugin.Task
	Shutdown() error
	Run() error
}

type baseRunner struct {
	ctx context.Context
}
