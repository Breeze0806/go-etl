package plugin

import (
	"context"

	"github.com/Breeze0806/go-etl/storage/database"
)

//TableDiffer 表不同
type TableDiffer struct {
	MasterTable database.Table
	SlaveTable  database.Table
	Differ      Differ
}

//DifferStorage 差异存储
type DifferStorage interface {
	Write(ctx context.Context,
		fetchDiffer func() (TableDiffer, error)) error
	Read(ctx context.Context,
		onDiffer func(differ TableDiffer) error) error
}
