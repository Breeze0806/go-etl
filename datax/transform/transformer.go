package transform

import "github.com/Breeze0806/go-etl/datax/common/element"

type Transformer interface {
	DoTransform(element.Record) (element.Record, error)
}

type NilTransformer struct{}

func (n *NilTransformer) DoTransform(record element.Record) (element.Record, error) {
	return record, nil
}
