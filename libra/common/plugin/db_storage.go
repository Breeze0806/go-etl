package plugin

import (
	"context"

	"github.com/Breeze0806/go-etl/storage/database"
)

//StorageTable 数据库表
type StorageTable interface {
	Table() database.Table        //数据库表
	ReadPage(ctx context.Context, //读取一页数据库
		param PageParam) (Page, error)
	Close() error //关闭
}

//DBStorage 数据库
type DBStorage interface {
	AllTable(ctx context.Context) ([]*database.BaseTable, error) //查询所有表
	MasterTable(*database.BaseTable) (MasterTable, error)        //主数据库表
	SlaveTable(*database.BaseTable) (SlaveTable, error)          //从数据库表
}
