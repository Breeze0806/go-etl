package element

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type ColumnType string

const (
	TypeUnknown ColumnType = "unknown"
	TypeBool    ColumnType = "bool"
	TypeBigInt  ColumnType = "bigInt"
	TypeDecimal ColumnType = "decimal"
	TypeString  ColumnType = "string"
	TypeBytes   ColumnType = "bytes"
	TypeTime    ColumnType = "time"
)

func (c ColumnType) String() string {
	return string(c)
}

type ColumnValue interface {
	Type() ColumnType
	IsNil() bool
	AsBool() (bool, error)
	AsBigInt() (*big.Int, error)
	AsDecimal() (decimal.Decimal, error)
	AsString() (string, error)
	AsBytes() ([]byte, error)
	AsTime() (time.Time, error)
	String() string
}

type ColumnValueClonable interface {
	Clone() ColumnValue
}

type Column interface {
	ColumnValue
	AsInt8() (int8, error)
	AsInt16() (int16, error)
	AsInt32() (int32, error)
	AsInt64() (int64, error)
	AsFloat32() (float32, error)
	AsFloat64() (float64, error)
	Clone() (Column, error)
	Name() string
	ByteSize() int64
	MemorySize() int64
}

type notNilColumnValue struct {
}

func (n *notNilColumnValue) IsNil() bool {
	return false
}

type nilColumnValue struct{}

func (n *nilColumnValue) Type() ColumnType {
	return TypeUnknown
}

func (n *nilColumnValue) IsNil() bool {
	return true
}

func (n *nilColumnValue) AsBool() (bool, error) {
	return false, ErrNilValue
}

func (n *nilColumnValue) AsBigInt() (*big.Int, error) {
	return nil, ErrNilValue
}

func (n *nilColumnValue) AsDecimal() (decimal.Decimal, error) {
	return decimal.Decimal{}, ErrNilValue
}

func (n *nilColumnValue) AsString() (string, error) {
	return "", ErrNilValue
}

func (n *nilColumnValue) AsBytes() ([]byte, error) {
	return nil, ErrNilValue
}

func (n *nilColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, ErrNilValue
}
func (n *nilColumnValue) String() string {
	return "<nil>"
}

type DefaultColumn struct {
	ColumnValue

	name     string
	byteSize int
}

func NewDefaultColumn(v ColumnValue, name string, byteSize int) Column {
	return &DefaultColumn{
		ColumnValue: v,
		name:        name,
		byteSize:    byteSize,
	}
}

func (d *DefaultColumn) Name() string {
	return d.name
}

func (d *DefaultColumn) Clone() (Column, error) {
	colnable, ok := d.ColumnValue.(ColumnValueClonable)
	if !ok {
		return nil, ErrNotColumnValueClonable
	}

	return &DefaultColumn{
		ColumnValue: colnable.Clone(),
		name:        d.name,
		byteSize:    d.byteSize,
	}, nil
}

func (d *DefaultColumn) ByteSize() int64 {
	return int64(d.byteSize)
}

func (d *DefaultColumn) MemorySize() int64 {
	return int64(d.byteSize + len(d.name) + 4)
}

func (d *DefaultColumn) AsInt8() (int8, error) {
	bi, err := d.AsBigInt()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "int8", err)
	}
	if bi.IsInt64() {
		v := bi.Int64()
		if v > math.MaxInt8 || v < math.MinInt8 {
			return 0, NewTransformErrorFormString(d.Type().String(), "int8", fmt.Errorf("%v %v", v, strconv.ErrRange))
		}
		return int8(bi.Int64()), nil
	}

	return 0, NewTransformErrorFormString(d.Type().String(), "int8", fmt.Errorf("%v %v", d.String(), ErrValueNotInt64))
}

func (d *DefaultColumn) AsInt16() (int16, error) {
	bi, err := d.AsBigInt()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "int16", err)
	}
	if bi.IsInt64() {
		v := bi.Int64()
		if v > math.MaxInt16 || v < math.MinInt16 {
			return 0, NewTransformErrorFormString(d.Type().String(), "int16", fmt.Errorf("%v %v", v, strconv.ErrRange))
		}
		return int16(bi.Int64()), nil
	}
	return 0, NewTransformErrorFormString(d.Type().String(), "int16", fmt.Errorf("%v %v", d.String(), ErrValueNotInt64))
}

func (d *DefaultColumn) AsInt32() (int32, error) {
	bi, err := d.AsBigInt()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "int32", err)
	}
	if bi.IsInt64() {
		v := bi.Int64()
		if v > math.MaxInt32 || v < math.MinInt32 {
			return 0, NewTransformErrorFormString(d.Type().String(), "int32", fmt.Errorf("%v %v", v, strconv.ErrRange))
		}
		return int32(bi.Int64()), nil
	}
	return 0, NewTransformErrorFormString(d.Type().String(), "int32", fmt.Errorf("%v %v", d.String(), ErrValueNotInt64))
}

func (d *DefaultColumn) AsInt64() (int64, error) {
	bi, err := d.AsBigInt()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "int64", err)
	}
	if bi.IsInt64() {
		return int64(bi.Int64()), nil
	}
	return 0, NewTransformErrorFormString(d.Type().String(), "int64", fmt.Errorf("%v %v", d.String(), ErrValueNotInt64))
}

func (d *DefaultColumn) AsFloat32() (float32, error) {
	dec, err := d.AsDecimal()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "float32", err)
	}
	v, _ := dec.Rat().Float32()
	if math.Abs(float64(v)) > math.MaxFloat32 {
		return 0, NewTransformErrorFormString(d.Type().String(), "float32",
			fmt.Errorf("%v %v", d.String(), strconv.ErrRange))
	}
	return v, nil
}

func (d *DefaultColumn) AsFloat64() (float64, error) {
	dec, err := d.AsDecimal()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "float64", err)
	}
	v, _ := dec.Float64()
	if math.Abs(v) > math.MaxFloat64 {
		return 0, NewTransformErrorFormString(d.Type().String(), "float64",
			fmt.Errorf("%v %v", d.String(), strconv.ErrRange))
	}
	return v, nil
}
