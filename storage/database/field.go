package database

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"

	"github.com/Breeze0806/go-etl/element"
)

type GoType int

var (
	ErrNotValuerGoType = errors.New("field type is not ValuerGoType")
)

const (
	GoTypeUnknow GoType = iota
	GoTypeBool
	GoTypeInt8
	GoTypeInt16
	GoTypeInt32
	GoTypeInt64
	GoTypeFloat32
	GoTypeFloat64
	GoTypeDecimal
	GoTypeString
	GoTypeBytes
	GoTypeTime
)

var goTypeMap = map[GoType]string{
	GoTypeUnknow:  "unknow",
	GoTypeBool:    "bool",
	GoTypeInt8:    "int8",
	GoTypeInt16:   "int16",
	GoTypeInt32:   "int32",
	GoTypeInt64:   "int64",
	GoTypeFloat32: "float32",
	GoTypeFloat64: "float64",
	GoTypeString:  "string",
	GoTypeBytes:   "bytes",
	GoTypeTime:    "time",
}

func (t GoType) String() string {
	if s, ok := goTypeMap[t]; ok {
		return s
	}
	return "unknow"
}

type Field interface {
	fmt.Stringer

	Name() string                 //字段名
	Quoted() string               //引用字段名
	BindVar(int) string           //占位符号
	Select() string               //select字段名
	Type() FieldType              //字段类型
	Scanner() Scanner             //扫描器
	Valuer(element.Column) Valuer //赋值器
}

type Scanner interface {
	sql.Scanner

	Column() element.Column
}

type Valuer interface {
	driver.Valuer
}

type FieldType interface {
	Name() string
	Length() (length int64, ok bool)
	DecimalSize() (precision, scale int64, ok bool)
	ScanType() reflect.Type
	Nullable() (nullable, ok bool)
	DatabaseTypeName() string
}

type ValuerGoType interface {
	GoType() GoType
}

type BaseField struct {
	name       string
	columnType *sql.ColumnType
}

func NewBaseField(name string, columnType *sql.ColumnType) *BaseField {
	return &BaseField{
		columnType: columnType,
		name:       name,
	}
}

func (b *BaseField) Name() string {
	return b.name
}

func (b *BaseField) ColumnType() *sql.ColumnType {
	return b.columnType
}

func (b *BaseField) String() string {
	return b.name
}

type BaseFieldType struct {
	*sql.ColumnType
}

func NewBaseFieldType(columnType *sql.ColumnType) *BaseFieldType {
	return &BaseFieldType{
		ColumnType: columnType,
	}
}

type BaseScanner struct {
	c element.Column
}

func (b *BaseScanner) SetColumn(c element.Column) {
	b.c = c
}

func (b *BaseScanner) Column() element.Column {
	return b.c
}

type GoValuer struct {
	f Field
	c element.Column
}

func NewGoValuer(f Field, c element.Column) *GoValuer {
	return &GoValuer{
		f: f,
		c: c,
	}
}

func (g *GoValuer) Value() (driver.Value, error) {
	typ, ok := g.f.Type().(ValuerGoType)
	if !ok {
		return nil, ErrNotValuerGoType
	}
	switch typ.GoType() {
	case GoTypeBool:
		return g.c.AsBool()
	case GoTypeInt8:
		return g.c.AsInt8()
	case GoTypeInt16:
		return g.c.AsInt16()
	case GoTypeInt32:
		return g.c.AsInt32()
	case GoTypeInt64:
		return g.c.AsInt64()
	case GoTypeFloat32:
		return g.c.AsFloat32()
	case GoTypeFloat64:
		return g.c.AsFloat64()
	case GoTypeString:
		return g.c.AsString()
	case GoTypeBytes:
		return g.c.AsBytes()
	case GoTypeTime:
		return g.c.AsTime()
	}
	return nil, fmt.Errorf("%v type", typ.GoType())
}
