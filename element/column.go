// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package element

import (
	"fmt"
	"time"
	"unsafe"
)

// ColumnType: Column Type
type ColumnType string

// ColumnTypeEnum: Enumeration of column types
const (
	TypeUnknown ColumnType = "unknown" // UnknownType: Unknown type
	TypeBool    ColumnType = "bool"    // BoolType: Boolean type
	TypeBigInt  ColumnType = "bigInt"  // IntType: Integer type
	TypeDecimal ColumnType = "decimal" // DecimalType: High-precision real number type
	TypeString  ColumnType = "string"  // StringType: String type
	TypeBytes   ColumnType = "bytes"   // BytesType: Byte stream type
	TypeTime    ColumnType = "time"    // TimeType: Time type
)

// String: Printing display
func (c ColumnType) String() string {
	return string(c)
}

// ColumnValue: Column Value
type ColumnValue interface {
	fmt.Stringer

	Type() ColumnType                  // ColumnType: Column type
	IsNil() bool                       // IsNull: Whether it is null
	AsBool() (bool, error)             // ToBool: Convert to boolean
	AsBigInt() (BigIntNumber, error)   // ToInt: Convert to integer
	AsDecimal() (DecimalNumber, error) // ToDecimal: Convert to high-precision real number
	AsString() (string, error)         // ToString: Convert to string
	AsBytes() ([]byte, error)          // ToBytes: Convert to byte stream
	AsTime() (time.Time, error)        // ToTime: Convert to time
}

// ColumnValueClonable: Cloneable column value
type ColumnValueClonable interface {
	Clone() ColumnValue // Clone: Clone
}

// ColumnValueComparabale: Comparable column value
type ColumnValueComparabale interface {
	// Compare: 1 represents greater than, 0 represents equal, -1 represents less than
	Cmp(ColumnValue) (int, error)
}

// Column: Column
type Column interface {
	ColumnValue
	AsInt64() (int64, error)     // ToInt64: Convert to 64-bit integer
	AsFloat64() (float64, error) // ToFloat64: Convert to 64-bit real number
	Clone() (Column, error)      // Clone: Clone
	Cmp(Column) (int, error)     // Compare: 1 represents greater than, 0 represents equal, -1 represents less than
	Name() string                // Name: Column name
	ByteSize() int64             // ByteSize: Byte stream size
	MemorySize() int64           // MemorySize: Memory size
}

type notNilColumnValue struct{}

// IsNil: Whether it is null
func (n *notNilColumnValue) IsNil() bool {
	return false
}

type nilColumnValue struct{}

// Type: Column type
func (n *nilColumnValue) Type() ColumnType {
	return TypeUnknown
}

// IsNil: Whether it is null
func (n *nilColumnValue) IsNil() bool {
	return true
}

// AsBool: Failed to convert to boolean
func (n *nilColumnValue) AsBool() (bool, error) {
	return false, ErrNilValue
}

// AsBigInt: Failed to convert to integer
func (n *nilColumnValue) AsBigInt() (BigIntNumber, error) {
	return nil, ErrNilValue
}

// AsDecimal: Failed to convert to high-precision real number
func (n *nilColumnValue) AsDecimal() (DecimalNumber, error) {
	return nil, ErrNilValue
}

// AsString: Failed to convert to string
func (n *nilColumnValue) AsString() (string, error) {
	return "", ErrNilValue
}

// AsBytes: Failed to convert to byte stream
func (n *nilColumnValue) AsBytes() ([]byte, error) {
	return nil, ErrNilValue
}

// AsTime: Failed to convert to time
func (n *nilColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, ErrNilValue
}

// String: Printing display
func (n *nilColumnValue) String() string {
	return "<nil>"
}

// DefaultColumn: Default value
type DefaultColumn struct {
	ColumnValue // ColumnValue: Column value

	name     string
	byteSize int
}

// NewDefaultColumn: Create a new default column based on column value v, column name name, and byte stream size byteSize
func NewDefaultColumn(v ColumnValue, name string, byteSize int) Column {
	return &DefaultColumn{
		ColumnValue: v,
		name:        name,
		byteSize:    byteSize,
	}
}

// Name: Column name
func (d *DefaultColumn) Name() string {
	return d.name
}

// Cmp: Compare columns. If it's not a comparable column value, an error will occur.
func (d *DefaultColumn) Cmp(c Column) (int, error) {
	if d.Name() != c.Name() {
		return 0, ErrColumnNameNotEqual
	}
	comparabale, ok := d.ColumnValue.(ColumnValueComparabale)
	if !ok {
		return 0, ErrNotColumnValueComparable
	}
	return comparabale.Cmp(c)
}

// Clone: Clone a column. If it's not a cloneable column value, an error will occur.
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

// ByteSize: Byte stream size
func (d *DefaultColumn) ByteSize() int64 {
	return int64(d.byteSize)
}

// MemorySize: Memory size
func (d *DefaultColumn) MemorySize() int64 {
	return int64(d.byteSize + len(d.name) + 4)
}

// AsInt64: Convert to 64-bit integer
func (d *DefaultColumn) AsInt64() (int64, error) {
	bi, err := d.AsBigInt()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "int64", err)
	}
	return bi.Int64()
}

// AsFloat64: Convert to 64-bit real number
func (d *DefaultColumn) AsFloat64() (float64, error) {
	dec, err := d.AsDecimal()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "float64", err)
	}
	return dec.Float64()
}

// ByteSize: Byte size
func ByteSize(src interface{}) int {
	switch data := src.(type) {
	case nil:
		return 0
	case bool:
		return 1
	case string:
		return len(data)
	case []byte:
		return len(data)
	}
	return int(unsafe.Sizeof(src))
}
