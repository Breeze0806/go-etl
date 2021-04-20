package plugin

import (
	"context"
)

//MasterTable 主数据库表
type MasterTable interface {
	StorageTable
	//切分为splitNum个页，每个页通过页查询参数表述，切分失败时返回错误
	Split(ctx context.Context, splitNum int) ([]PageParam, error)
}
