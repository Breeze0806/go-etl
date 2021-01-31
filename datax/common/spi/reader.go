package spi

import "github.com/Breeze0806/go-etl/datax/common/spi/reader"

//Reader 读取器
type Reader interface {
	Job() reader.Job   //获取读取工作,一般不能为空
	Task() reader.Task //获取读取任务,一般不能为空
}
