package element

import (
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

type ColumnType uint8

const (
	TypeUnknown = iota
	TypeBool
	TypeBigInt

	TypeDecimal
	TypeString
	TypeBytes
	TypeTime
)

var fieldTypeMap = map[ColumnType]string{
	TypeBool:    "bool",
	TypeBigInt:  "bigInt",
	TypeDecimal: "decimal",
	TypeString:  "string",
	TypeBytes:   "bytes",
	TypeTime:    "time",
}

func (c ColumnType) String() string {
	if t, ok := fieldTypeMap[c]; ok {
		return t
	}
	return "unknown"
}

func (c ColumnType) IsSupported() bool {
	switch c {
	case TypeTime, TypeDecimal,
		TypeBytes, TypeBool, TypeBigInt, TypeString:
		return true
	}
	return false
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
	clone() ColumnValue
}

type Column interface {
	ColumnValue
	Clone() Column
	Name() string
	ByteSize() int64
	MemorySize() int64
}

type notNilColumnValue struct {
}

func (n *notNilColumnValue) IsNil() bool {
	return false
}

type nilColumnValue struct {
}

func (n *nilColumnValue) Type() ColumnType {
	return TypeUnknown
}
func (n *nilColumnValue) IsNil() bool {
	return false
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
func (n *nilColumnValue) clone() ColumnValue {
	return &nilColumnValue{}
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

func (d *DefaultColumn) Clone() Column {
	return &DefaultColumn{
		ColumnValue: d.clone(),
		name:        d.name,
		byteSize:    d.byteSize,
	}
}

func (d *DefaultColumn) ByteSize() int64 {
	return int64(d.byteSize)
}

func (d *DefaultColumn) MemorySize() int64 {
	return int64(d.byteSize + len(d.name) + 4)
}
