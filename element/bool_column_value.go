package element

import (
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

//NilBoolColumnValue 空值布尔列值
type NilBoolColumnValue struct {
	nilColumnValue
}

//NewNilBoolColumnValue 生成空值布尔列值
func NewNilBoolColumnValue() ColumnValue {
	return &NilBoolColumnValue{}
}

//Type 返回列类型
func (n *NilBoolColumnValue) Type() ColumnType {
	return TypeBool
}

//Clone 克隆空值布尔列值
func (n *NilBoolColumnValue) Clone() ColumnValue {
	return NewNilBoolColumnValue()
}

//BoolColumnValue 布尔列值
type BoolColumnValue struct {
	notNilColumnValue

	val bool //布尔值
}

//NewBoolColumnValue 从布尔值v生成布尔列值
func NewBoolColumnValue(v bool) ColumnValue {
	return &BoolColumnValue{
		val: v,
	}
}

//Type 返回列类型
func (b *BoolColumnValue) Type() ColumnType {
	return TypeBool
}

//AsBool 转化成布尔值
func (b *BoolColumnValue) AsBool() (bool, error) {
	return b.val, nil
}

//AsBigInt 转化成整数，true转化为1，false转化为0
func (b *BoolColumnValue) AsBigInt() (*big.Int, error) {
	if b.val {
		return big.NewInt(1), nil
	}
	return big.NewInt(0), nil
}

//AsDecimal 转化成高精度实数，true转化为1.0，false转化为0.0
func (b *BoolColumnValue) AsDecimal() (decimal.Decimal, error) {
	if b.val {
		return decimal.New(1, 0), nil
	}
	return decimal.New(0, 1), nil
}

//AsString 转化成字符串，true转化为"true"，false转化为"false"
func (b *BoolColumnValue) AsString() (string, error) {
	if b.val {
		return b.String(), nil
	}
	return b.String(), nil
}

//AsBytes 转化成字节流，true转化为"true"，false转化为"false"
func (b *BoolColumnValue) AsBytes() ([]byte, error) {
	if b.val {
		return []byte(b.String()), nil
	}
	return []byte(b.String()), nil
}

//AsTime 目前布尔无法转化成时间
func (b *BoolColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformErrorFormColumnTypes(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
}

func (b *BoolColumnValue) String() string {
	if b.val {
		return "true"
	}
	return "false"
}

//Clone 克隆布尔列值
func (b *BoolColumnValue) Clone() ColumnValue {
	return NewBoolColumnValue(b.val)
}
