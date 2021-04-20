package plugin

//SlaveTable 从数据库表
type SlaveTable interface {
	StorageTable
	//通过[start, end)生成页查询参数PageParam, 在转化失败时返回错误
	TransformPage(start Offset, end Offset) (PageParam, error)
}
