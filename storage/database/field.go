package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/Breeze0806/go-etl/element"
)

type Type int

const (
	TypeUnknow Type = iota
	TypeBool
	TypeInt64
	TypeUint64
	TypeBigInt
	TypeFloat64
	TypeDecimal
	TypeString
	TypeBytes
	TypeTime
)

var typeMap = map[Type]string{
	TypeUnknow:  "unknow",
	TypeBool:    "bool",
	TypeInt64:   "int64",
	TypeUint64:  "uint64",
	TypeBigInt:  "bigInt",
	TypeFloat64: "float64",
	TypeDecimal: "decimal",
	TypeString:  "string",
	TypeBytes:   "bytes",
	TypeTime:    "time",
}

func (t Type) String() string {
	if s, ok := typeMap[t]; ok {
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
	Type() Type
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
