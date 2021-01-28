package element

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

//NilBytesColumnValue 空值字节流列值
type NilBytesColumnValue struct {
	nilColumnValue
}

//NewNilBytesColumnValue 创建空值字节流列值
func NewNilBytesColumnValue() ColumnValue {
	return &NilBytesColumnValue{}
}

//Type 返回列类型
func (n *NilBytesColumnValue) Type() ColumnType {
	return TypeBytes
}

//Clone 克隆空值字节流列值
func (n *NilBytesColumnValue) Clone() ColumnValue {
	return NewNilBytesColumnValue()
}

//BytesColumnValue 字节流列值
type BytesColumnValue struct {
	notNilColumnValue
	TimeEncoder //时间编码器

	val []byte //字节流值
}

//NewBytesColumnValue 从字节流v 生成字节流列值
func NewBytesColumnValue(v []byte) ColumnValue {
	return NewBytesColumnValueWithEncoder(v, NewStringTimeEncoder(time.RFC3339Nano))
}

//NewBytesColumnValueWithEncoder 从字节流v 和时间编码器e 生成字节流列值
func NewBytesColumnValueWithEncoder(v []byte, e TimeEncoder) ColumnValue {
	return &BytesColumnValue{
		val:         v,
		TimeEncoder: e,
	}
}

//Type 返回列类型
func (b *BytesColumnValue) Type() ColumnType {
	return TypeBytes
}

//AsBool 1, t, T, TRUE, true, True转化为true
//0, f, F, FALSE, false, False转化为false，如果不是上述情况会报错
func (b *BytesColumnValue) AsBool() (bool, error) {
	v, err := strconv.ParseBool(b.String())
	if err != nil {
		return false, NewTransformErrorFormColumnTypes(b.Type(), TypeBool, fmt.Errorf("err: %v val: %v", err, b.String()))
	}
	return v, nil
}

//AsBigInt 转化为整数，实数型以及科学性计数法字符串会被取整，不是数值型的会报错
//如123.67转化为123 123.12转化为123
func (b *BytesColumnValue) AsBigInt() (*big.Int, error) {
	v, err := NewDecimalColumnValueFromString(b.String())
	if err != nil {
		return nil, NewTransformErrorFormColumnTypes(b.Type(), TypeBigInt, fmt.Errorf("err: %v, val: %v ", err, b.String()))
	}
	return v.AsBigInt()
}

//AsDecimal 转化为整数，实数型以及科学性计数法字符串能够转化，不是数值型的会报错
func (b *BytesColumnValue) AsDecimal() (decimal.Decimal, error) {
	v, err := NewDecimalColumnValueFromString(b.String())
	if err != nil {
		return decimal.Decimal{}, NewTransformErrorFormColumnTypes(b.Type(), TypeDecimal, fmt.Errorf("err: %v, val: %v ", err, b.String()))
	}
	return v.AsDecimal()
}

//AsString 转化为字符串
func (b *BytesColumnValue) AsString() (string, error) {
	return b.String(), nil
}

//AsBytes 转化成字节流
func (b *BytesColumnValue) AsBytes() ([]byte, error) {
	v := make([]byte, len(b.val))
	copy(v, b.val)
	return v, nil
}

//AsTime 根据时间编码器转化成时间，不符合时间编码器格式会报错
func (b *BytesColumnValue) AsTime() (t time.Time, err error) {
	t, err = b.TimeEncode(b.String())
	if err != nil {
		return time.Time{}, NewTransformErrorFormColumnTypes(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
	}
	return
}

func (b *BytesColumnValue) String() string {
	return string(b.val)
}

//Clone 克隆字节流列值
func (b *BytesColumnValue) Clone() ColumnValue {
	v := make([]byte, len(b.val))
	copy(v, b.val)
	return NewBytesColumnValue(v)
}
