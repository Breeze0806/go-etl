package database

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"

	"github.com/Breeze0806/go-etl/element"
)

//GoType golang的类型
type GoType uint8

//字段错误相关
var (
	ErrNotValuerGoType = errors.New("field type is not ValuerGoType") //接口不是ValuerGoType的错误
)

//golang的类型枚举
const (
	GoTypeUnknown GoType = iota //未知类型
	GoTypeBool                  //布尔类型
	GoTypeInt8                  //Int8类型
	GoTypeInt16                 //Int16类型
	GoTypeInt32                 //Int32类型
	GoTypeInt64                 //Int64类型
	GoTypeFloat32               //Float32类型
	GoTypeFloat64               //Float64类型
	GoTypeString                //字符串类型
	GoTypeBytes                 //字节流类型
	GoTypeTime                  //时间类型
)

//golang的类型枚举字符串
var goTypeMap = map[GoType]string{
	GoTypeUnknown: "unknow",
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

//String golang的类型枚举字符串描述
func (t GoType) String() string {
	if s, ok := goTypeMap[t]; ok {
		return s
	}
	return "unknow"
}

//Field 数据库字段
type Field interface {
	fmt.Stringer

	Index() int                   //索引
	Name() string                 //字段名
	Quoted() string               //引用字段名
	BindVar(int) string           //占位符号
	Select() string               //select字段名
	Type() FieldType              //字段类型
	Scanner() Scanner             //扫描器
	Valuer(element.Column) Valuer //赋值器
}

//Scanner 列数据扫描器 数据库驱动的值扫描成列数据
type Scanner interface {
	sql.Scanner

	Column() element.Column //获取列数据
}

//Valuer 赋值器 将对应数据转化成数据库驱动的值
type Valuer interface {
	driver.Valuer
}

//ColumnType 列类型,抽象 sql.ColumnType，也方便自行实现对应函数
type ColumnType interface {
	Name() string                                   //列名
	ScanType() reflect.Type                         //扫描类型
	Length() (length int64, ok bool)                //长度
	DecimalSize() (precision, scale int64, ok bool) //精度
	Nullable() (nullable, ok bool)                  //是否为空
	DatabaseTypeName() string                       //列数据库类型名
}

//FieldType 字段类型
type FieldType interface {
	ColumnType

	IsSupportted() bool //是否支持
}

//ValuerGoType 用于赋值器的golang类型判定,是Field的可选功能，
//就是对对应驱动的值返回相应的值，方便GoValuer进行判定
type ValuerGoType interface {
	GoType() GoType
}

//BaseField 基础字段，主要存储列名name和列类型fieldType
type BaseField struct {
	index     int
	name      string
	fieldType FieldType
}

//NewBaseField 根据列名name和列类型fieldType获取基础字段
//用于嵌入其他Field，方便实现各个数据库的Field
func NewBaseField(index int, name string, fieldType FieldType) *BaseField {
	return &BaseField{
		index:     index,
		fieldType: fieldType,
		name:      name,
	}
}

//Index 返回字段名
func (b *BaseField) Index() int {
	return b.index
}

//Name 返回字段名
func (b *BaseField) Name() string {
	return b.name
}

//FieldType 返回字段类型
func (b *BaseField) FieldType() FieldType {
	return b.fieldType
}

//String 打印时显示字符串
func (b *BaseField) String() string {
	return b.name
}

//BaseFieldType 基础字段类型，嵌入其他各种数据库字段类型实现
type BaseFieldType struct {
	ColumnType
}

//NewBaseFieldType 获取字段类型
func NewBaseFieldType(typ ColumnType) *BaseFieldType {
	return &BaseFieldType{
		ColumnType: typ,
	}
}

func (*BaseFieldType) IsSupportted() bool {
	return true
}

//BaseScanner 基础扫描器，嵌入其他各种数据库扫描器实现
type BaseScanner struct {
	c element.Column
}

//SetColumn 设置列值，用于数据库方言的列数据设置
func (b *BaseScanner) SetColumn(c element.Column) {
	b.c = c
}

//Column 取得列值，方便统一取得列值
func (b *BaseScanner) Column() element.Column {
	return b.c
}

//GoValuer 使用GoType类型生成赋值器，主要通过字段f和传入参数列值c来
//完成使用GoType类型生成赋值器,方便实现GoValuer
type GoValuer struct {
	f Field
	c element.Column
}

//NewGoValuer 主要通过字段f和传入参数列值c来完成使用GoType类型生成赋值器的生成
func NewGoValuer(f Field, c element.Column) *GoValuer {
	return &GoValuer{
		f: f,
		c: c,
	}
}

//Value 根据ValuerGoType生成对应的驱动接受的值
func (g *GoValuer) Value() (driver.Value, error) {
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
	return nil, fmt.Errorf("%v type(%v)", typ.GoType(), g.f.Type().DatabaseTypeName())
}
