package plugin

import (
	"github.com/Breeze0806/go-etl/element"
)

//RecordReceiver 记录接收器
type RecordReceiver interface {
	GetFromReader() (element.Record, error) //从reader中读取记录
	Shutdown() error                        // 关闭
}
