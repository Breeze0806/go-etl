package transform

import "github.com/Breeze0806/go-etl/element"

//Transformer 转化器
type Transformer interface {
	DoTransform(element.Record) (element.Record, error)
}

//NilTransformer 空转化器
type NilTransformer struct{}

//DoTransform 转化
func (n *NilTransformer) DoTransform(record element.Record) (element.Record, error) {
	return record, nil
}
