package plugin

import (
	"context"

	"github.com/Breeze0806/go-etl/storage/database"
)

//OffsetTracker 位移追踪器
type OffsetTracker interface {
	Read(ctx context.Context) error
	Offset() Offset
	SetOffset(offset Offset)
	Write(ctx context.Context) error
	Close() error
}

//PageParamTracker 页查询参数追踪器
type PageParamTracker interface {
	Read(ctx context.Context) ([]PageParam, error)
	Add(pageParam PageParam)
	Delete(pageParam PageParam)
	Write(ctx context.Context) error
	Close() error
}

//Tracker 追踪器
type Tracker interface {
	OffsetTracker(master database.Table) (OffsetTracker, error)
	PageParamTracker(master database.Table) (PageParamTracker, error)
}

