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

package database

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"

	"github.com/Breeze0806/go-etl/element"
)

// GoType refers to the type in the Golang language.
type GoType uint8

// Related to field errors.
var (
	ErrNotValuerGoType = errors.New("field type is not ValuerGoType") // Error indicating that the interface is not a ValuerGoType.
)

// Enumeration of golang types.
const (
	GoTypeUnknown GoType = iota // Unknown type.
	GoTypeBool                  // Boolean type.
	GoTypeInt64                 // Int64 type.
	GoTypeFloat64               // Float64 type.
	GoTypeString                // String type.
	GoTypeBytes                 // Byte stream type.
	GoTypeTime                  // Time type.
)

// Enumeration string of golang types.
var goTypeMap = map[GoType]string{
	GoTypeUnknown: "unknow",
	GoTypeBool:    "bool",
	GoTypeInt64:   "int64",
	GoTypeFloat64: "float64",
	GoTypeString:  "string",
	GoTypeBytes:   "bytes",
	GoTypeTime:    "time",
}

// String description of the enumeration of golang types.
func (t GoType) String() string {
	if s, ok := goTypeMap[t]; ok {
		return s
	}
	return "unknow"
}

// Field refers to a database field.
type Field interface {
	fmt.Stringer

	Index() int                   // Index.
	Name() string                 // Field name.
	Quoted() string               // Referenced field name.
	BindVar(int) string           // Placeholder symbol, starting from 1.
	Select() string               // Selected field name.
	Type() FieldType              // Field type.
	Scanner() Scanner             // Scanner.
	Valuer(element.Column) Valuer // Valuer.
	SetError(err *error)
}

// Scanner: Data scanner for columns. Converts database driver values into column data.
type Scanner interface {
	sql.Scanner

	Column() element.Column // Get column data.
}

// Valuer: Converts corresponding data into database driver values.
type Valuer interface {
	driver.Valuer
}

// ColumnType: Represents the type of a column, abstracting sql.ColumnType and facilitating custom implementations.
type ColumnType interface {
	Name() string                                   // Column name.
	ScanType() reflect.Type                         // Scanning type.
	Length() (length int64, ok bool)                // Length.
	DecimalSize() (precision, scale int64, ok bool) // Precision.
	Nullable() (nullable, ok bool)                  // Whether it is nullable.
	DatabaseTypeName() string                       // Name of the column's database type.
}

// FieldType: Represents the type of a field.
type FieldType interface {
	ColumnType

	IsSupported() bool // Whether it is supported.
}

// ValuerGoType: Determines the golang type for a Valuer. An optional feature of Field, it converts corresponding driver values.

type ValuerGoType interface {
	GoType() GoType
}

// BaseField: Represents a basic field, primarily storing the column name and column type.
type BaseField struct {
	index     int
	name      string
	fieldType FieldType
}

// NewBaseField: Creates a new base field based on the column name and column type.
// Used for embedding other Fields, facilitating the implementation of database-specific Fields.
func NewBaseField(index int, name string, fieldType FieldType) *BaseField {
	return &BaseField{
		index:     index,
		name:      name,
		fieldType: fieldType,
	}
}

// Index: Returns the field name.
func (b *BaseField) Index() int {
	return b.index
}

// Name: Returns the field name.
func (b *BaseField) Name() string {
	return b.name
}

// FieldType: Returns the field type.
func (b *BaseField) FieldType() FieldType {
	return b.fieldType
}

// String: Displays a string representation when printing.
func (b *BaseField) String() string {
	return b.name
}

func (b *BaseField) SetError(err *error) {
	if *err != nil {
		*err = fmt.Errorf("field: %v, %v", b.name, *err)
	}
}

// BaseFieldType: Represents the basic type of a field, embedding implementations for various database field types.
type BaseFieldType struct {
	ColumnType
}

// NewBaseFieldType: Gets the field type.
func NewBaseFieldType(typ ColumnType) *BaseFieldType {
	return &BaseFieldType{
		ColumnType: typ,
	}
}

// IsSupported: Determines if it is supported for parsing.
func (*BaseFieldType) IsSupported() bool {
	return true
}

// BaseScanner: Represents a basic scanner, embedding implementations for various database scanners.
type BaseScanner struct {
	c element.Column
}

// SetColumn: Sets the column value for database-specific column data settings.
func (b *BaseScanner) SetColumn(c element.Column) {
	b.c = c
}

// Column: Retrieves the column value, facilitating a unified way to obtain column values.
func (b *BaseScanner) Column() element.Column {
	return b.c
}

// GoValuer: Generates a Valuer using the GoType. Primarily done through the field 'f' and the incoming column value 'c'.
// Completes the generation of a Valuer using the GoType, facilitating the implementation of GoValuer.
type GoValuer struct {
	f Field
	c element.Column
}

// NewGoValuer: Generates a new Valuer using the GoType, primarily done through the field 'f' and the incoming column value 'c'.
func NewGoValuer(f Field, c element.Column) *GoValuer {
	return &GoValuer{
		f: f,
		c: c,
	}
}

// Value: Generates the corresponding driver-accepted value based on ValuerGoType.
func (g *GoValuer) Value() (val driver.Value, err error) {
	defer g.f.SetError(&err)
	typ, ok := g.f.Type().(ValuerGoType)
	if !ok {
		return nil, ErrNotValuerGoType
	}

	if g.c.IsNil() {
		return nil, nil
	}

	switch typ.GoType() {
	case GoTypeBool:
		return g.c.AsBool()
	case GoTypeInt64:
		return g.c.AsInt64()
	case GoTypeFloat64:
		return g.c.AsFloat64()
	case GoTypeString:
		return g.c.AsString()
	case GoTypeBytes:
		return g.c.AsBytes()
	case GoTypeTime:
		return g.c.AsTime()
	}
	return nil, fmt.Errorf("%v type(%v)", typ.GoType(), g.f.Type().DatabaseTypeName())
}
