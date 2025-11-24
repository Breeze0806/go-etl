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

package parquet

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/stream/file"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

func init() {
	var creator Creator
	file.RegisterCreator("parquet", &creator)
	var opener Opener
	file.RegisterOpener("parquet", &opener)
}

type NameMap map[string]Column

// Creator 创建parquet输出流的创建器
type Creator struct {
}

// Create 创建名为filename的输出流
func (c *Creator) Create(filename string) (file.OutStream, error) {
	return NewOutStream(filename)
}

// Stream parquet文件流
type Stream struct {
	fw      source.ParquetFile
	reader  *reader.ParquetReader
	nameMap NameMap
}

func (s *Stream) Rows(conf *config.JSON) (rows file.Rows, err error) {
	return NewRows(s.reader, s.nameMap, conf)
}

type Rows struct {
	rowNum  int // 总行数
	row     int // 当前行号
	reader  *reader.ParquetReader
	rowData map[string]interface{}
	err     error
	nameMap NameMap
	columns NameMap
}

func (r *Rows) Next() bool {
	if r.row >= r.rowNum {
		return false
	}
	d, err := r.reader.ReadByNumber(1)
	if err != nil {
		r.err = fmt.Errorf("read row error: %v", err)
		return false
	}

	data := normalizeStructToMap(d[0])
	r.rowData = data.(map[string]interface{})
	r.row++
	return true
}

func (r *Rows) Error() error {
	return r.err
}

func (r *Rows) Close() error {
	return nil
}

func (r *Rows) Scan() (columns []element.Column, err error) {
	columns = make([]element.Column, 0)
	for k, v := range r.rowData {
		if _, ok := r.columns[k]; !ok {
			// 过滤掉不在columns中的字段
			continue
		}
		c, err := r.FieldToColumn(k, v)
		if err != nil {
			return nil, err
		}
		columns = append(columns, c)
	}
	return
}
func (r *Rows) FieldToColumn(name string, field interface{}) (element.Column, error) {
	columnType := r.nameMap[name]
	name = columnType.Name
	switch columnType.Type {
	case "INT32":
		val, ok := field.(int32)
		if ok {
			return element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(int64(val)), name, 0), nil
		} else {
			return nil, fmt.Errorf("field %s type(%s) %s not match", name, columnType.Type, reflect.TypeOf(field))

		}

	case "INT64":
		val, ok := field.(int64)
		if ok {
			return element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(val), name, 0), nil
		} else {
			return nil, fmt.Errorf("field %s type(%s) %s not match", name, columnType.Type, reflect.TypeOf(field))

		}
	case "BYTE_ARRAY":
		val, ok := field.([]byte)
		if ok {
			return element.NewDefaultColumn(element.NewBytesColumnValue(val), name, 0), nil
		} else if reflect.TypeOf(field).String() == "string" {
			return element.NewDefaultColumn(element.NewStringColumnValue(field.(string)), name, 0), nil
		} else {
			return nil, fmt.Errorf("field %s type(%s) %s not match", name, columnType.Type, reflect.TypeOf(field))

		}
	case "FLOAT":
		val, ok := field.(float32)
		if ok {
			return element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat32(val), name, 0), nil
		} else {
			return nil, fmt.Errorf("field %s type(%s) %s not match", name, columnType.Type, reflect.TypeOf(field))

		}
	case "MAP", "LIST":
		s, err := json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("field %s type(%s) %s not match", name, columnType.Type, reflect.TypeOf(field))

		}
		return element.NewDefaultColumn(element.NewBytesColumnValue(s), name, 0), nil
	case "BOOLEAN":
		val, ok := field.(bool)
		if ok {
			return element.NewDefaultColumn(element.NewBoolColumnValue(val), name, 0), nil
		}
		return nil, fmt.Errorf("field %s type(%s) %s not match", name, columnType.Type, reflect.TypeOf(field))
	default:
		return nil, fmt.Errorf("field %s type(%s) %s not match", name, columnType.Type, reflect.TypeOf(field))

	}
}

// NewRows - Creates a row reader using the file handle 'f' and configuration 'c'.
func NewRows(f *reader.ParquetReader, nameMap NameMap, c *config.JSON) (rows *Rows, err error) {
	var conf *InConfig
	if conf, err = NewInConfig(c); err != nil {
		return nil, err
	}
	num := int(f.GetNumRows())
	rows = &Rows{
		rowNum:  num,
		row:     0,
		reader:  f,
		rowData: make(map[string]interface{}),
		nameMap: nameMap,
		columns: make(NameMap),
	}
	for _, col := range conf.Columns {
		rows.columns[col.Name] = col
	}
	return
}

// NewOutStream 创建parquet输出流
func NewOutStream(filename string) (file.OutStream, error) {
	stream := &Stream{}
	var err error
	fw, err := local.NewLocalFileWriter(filename)
	if err != nil {
		log.Println("Can't create file", err)
		return nil, err
	}
	stream.fw = fw
	return stream, nil
}

// Writer 创建写入器
func (s *Stream) Writer(conf *config.JSON) (file.StreamWriter, error) {
	return NewWriter(conf, s.fw)
}

// Close 关闭输出流
func (s *Stream) Close() (err error) {
	return s.fw.Close()
}

// Writer parquet流写入器
type Writer struct {
	pw      *writer.JSONWriter
	schema  *SchemaBuilder
	fw      source.ParquetFile
	conf    *OutConfig
	columns NameMap
}

// NewWriter 创建parquet流写入器
func NewWriter(c *config.JSON, fw source.ParquetFile) (file.StreamWriter, error) {
	var conf *OutConfig
	conf, err := NewOutConfig(c)
	if err != nil {
		return nil, err
	}
	w := &Writer{
		fw:      fw,
		conf:    conf,
		columns: make(NameMap),
	}
	for _, col := range conf.Columns {
		w.columns[col.Name] = col
	}
	return w, nil
}
func (w *Writer) BuildWriter() (err error) {
	if w.schema == nil {
		return fmt.Errorf("schema is nil")
	}
	schemaStr, err := w.schema.BuildCompact()
	if err != nil {
		return fmt.Errorf("build schema compact error: %v", err)

	}
	pw, err := writer.NewJSONWriter(schemaStr, w.fw, 4)
	if err != nil {
		log.Println("Can't create json writer", err)
		return
	}
	w.pw = pw
	return
}

// Write 写入记录
func (w *Writer) Write(record element.Record) (err error) {

	schema := NewSchemaBuilder("parquet")
	dataBuild := NewDataBuilder()
	for i := 0; i < record.ColumnNumber(); i++ {
		col, err := record.GetByIndex(i)
		if err != nil {
			return err
		}
		if _, ok := w.columns[col.Name()]; !ok {
			continue
		}
		switch col.Type() {
		case element.TypeBigInt:
			value, err := col.AsInt64()
			if err != nil {
				return err
			}
			schema.AddInt64Field(col.Name())
			dataBuild.SetInt64(col.Name(), value)
		case element.TypeDecimal:
			value, err := col.AsFloat64()
			if err != nil {
				return err
			}
			schema.AddDoubleField(col.Name())
			dataBuild.SetDouble(col.Name(), value)
		case element.TypeBool:
			val, err := col.AsBool()
			if err != nil {
				return err
			}
			schema.AddBooleanField(col.Name())
			dataBuild.SetBoolean(col.Name(), val)
		case element.TypeString:
			val, err := col.AsString()
			if err != nil {
				return err
			}
			schema.AddStringField(col.Name())
			dataBuild.SetString(col.Name(), val)
		case element.TypeBytes:
			val, err := col.AsString()
			if err != nil {
				return err
			}
			schema.AddBytesField(col.Name())
			dataBuild.SetString(col.Name(), val)
		case element.TypeTime:
			val, err := col.AsTime()
			if err != nil {
				return err
			}
			schema.AddInt64Field(col.Name())
			dataBuild.SetInt64(col.Name(), val.UnixNano())
		default:
			// 默认转换为字符串
			val, err := col.AsString()
			if err != nil {
				return err
			}
			schema.AddStringField(col.Name())
			dataBuild.SetString(col.Name(), val)
		}
	}

	if w.schema == nil {
		w.schema = schema
		err = w.BuildWriter()
		if err != nil {
			return fmt.Errorf("build writer error: %v", err)
		}
	}
	dataStr, err := dataBuild.Build()
	if err != nil {
		return fmt.Errorf("build data error: %v", err)
	}
	err = w.pw.Write(dataStr)
	if err != nil {
		return fmt.Errorf("write data error: %v", err)
	}
	return nil
}

// Flush 刷新至文件
func (w *Writer) Flush() (err error) {
	return
}

// Close 关闭输出流写入器
func (w *Writer) Close() (err error) {
	if w.pw != nil {
		err = w.pw.WriteStop()
		if err != nil {
			return fmt.Errorf("write stop error: %v", err)
		}
	}
	return w.fw.Close()
}

// Opener - A utility for opening XLSX input streams.
type Opener struct {
}

// Open - Opens an XLSX input stream named 'filename'.
func (o *Opener) Open(filename string) (file.InStream, error) {
	return NewInStream(filename)
}

// NewInStream - Creates an XLSX input stream named 'filename'.
func NewInStream(filename string) (file.InStream, error) {
	stream := &Stream{
		nameMap: make(NameMap),
	}
	var err error
	pf, err := local.NewLocalFileReader(filename)
	if err != nil {
		log.Println("Can't create file", err)
		return nil, fmt.Errorf("can't create file: %v", err)
	}
	stream.fw = pf
	stream.reader, err = reader.NewParquetReader(pf, nil, 4)

	if err != nil {
		return nil, err
	}
	for i, s := range stream.reader.SchemaHandler.SchemaElements {
		typeName := ""
		if s.Type != nil {
			typeName = s.Type.String()
		} else if s.ConvertedType != nil {
			typeName = s.ConvertedType.String()
		} else if s.RepetitionType != nil && s.RepetitionType.String() == "REPEATED" {
			typeName = "LIST"
		} else {
			continue
		}
		stream.nameMap[s.Name] = Column{
			Name: stream.reader.SchemaHandler.GetExName(i),
			Type: typeName,
		}
	}
	return stream, nil
}

func normalizeStructToMap(item interface{}) interface{} {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		typ := v.Type()
		result := make(map[string]interface{})

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := typ.Field(i)
			// 获取字段名，优先使用json tag
			fieldName := fieldType.Tag.Get("name")
			if fieldName == "" || strings.Contains(fieldName, ",") {
				if strings.Contains(fieldName, ",") {
					fieldName = strings.Split(fieldName, ",")[0]
				} else {
					fieldName = fieldType.Name
				}
			}

			result[fieldName] = normalizeStructToMap(field.Interface())
		}

		return result
	}

	if v.Kind() == reflect.Slice {
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = normalizeStructToMap(v.Index(i).Interface())
		}
		return result
	}

	if v.Kind() == reflect.Map {
		result := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			result[fmt.Sprintf("%v", key.Interface())] = normalizeStructToMap(v.MapIndex(key).Interface())
		}
		return result
	}

	return item
}
