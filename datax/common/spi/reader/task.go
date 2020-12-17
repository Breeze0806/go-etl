package reader

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

type Task interface {
	plugin.Task
	StartRead(ctx context.Context, sender plugin.RecordSender) error
}
