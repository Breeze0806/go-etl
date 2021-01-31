package spi

import "github.com/Breeze0806/go-etl/datax/common/spi/writer"

//Writer 写入器
type Writer interface {
	Job() writer.Job   //获取写入工作,一般不能为空
	Task() writer.Task //获取写入任务,一般不能为空
}
