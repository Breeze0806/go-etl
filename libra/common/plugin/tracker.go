package plugin

import (
	"context"

	"github.com/Breeze0806/go-etl/storage/database"
)

//OffsetTracker 位移追踪器
type OffsetTracker interface {
	Write(ctx context.Context, offset Offset) error
	Read(ctx context.Context) (Offset, error)
	Close() error
}

//PageParamTracker 页查询参数追踪器
type PageParamTracker interface {
	Write(ctx context.Context, params []PageParam) error
	Read(ctx context.Context) ([]PageParam, error)
	Close() error
}

//Tracker 追踪器
type Tracker interface {
	OffsetTracker(master database.Table) (OffsetTracker, error)
	PageParamTracker(master database.Table) (PageParamTracker, error)
}
