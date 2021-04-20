package plugin

import (
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

//PageParam 页查询参数
type PageParam struct {
	Start Offset             //页面开始位移
	End   Offset             //页面结束位移
	Param database.Parameter //页面查询语句
}

//Page 页查询结果
type Page struct {
	Min     Offset                //页查询结果最小位移
	Max     Offset                //页查询结果最大位移
	Records map[string]PageRecord //页查询结果映射
}

//PageRecord 页记录
type PageRecord struct {
	Now    Offset         //当前页记录位移
	Record element.Record //页记录查询结果
}
