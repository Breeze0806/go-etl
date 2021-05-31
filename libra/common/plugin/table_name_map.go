package plugin

import "github.com/Breeze0806/go-etl/storage/database"

//TableNameMap 表名映射
type TableNameMap interface {
	//通过主数据库表名获取master获取对应的从数据库表名slave,有错误时返回err
	SlaveTableName(master *database.BaseTable) (slave *database.BaseTable, err error)
	Close() error
}

type TableNameMapMaker interface {
	TableNameMap() TableNameMap
}
