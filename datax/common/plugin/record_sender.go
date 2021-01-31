package plugin

import (
	"github.com/Breeze0806/go-etl/element"
)

//RecordSender 记录发送器
type RecordSender interface {
	CreateRecord() (element.Record, error)  //创建记录
	SendWriter(record element.Record) error //将记录发往写入器
	Flush() error                           //将记录刷新到记录发送器
	Terminate() error                       //终止发送信号
	Shutdown() error                        //关闭
}
