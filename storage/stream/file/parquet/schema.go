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
	"strings"
)

// SchemaBuilder -  Schema 构建器
type SchemaBuilder struct {
	rootName string
	fields   []Field
}

// Field -  字段定义
type Field struct {
	Name           string  `json:"Tag"`
	Fields         []Field `json:"Fields,omitempty"`
	RepetitionType string  `json:"-"`
	ConvertedType  string  `json:"-"`
}

// FieldType -  字段类型枚举
type FieldType string

const (
	BYTEARRAY FieldType = "BYTE_ARRAY"
	INT32     FieldType = "INT32"
	INT64     FieldType = "INT64"
	FLOAT     FieldType = "FLOAT"
	DOUBLE    FieldType = "DOUBLE"
	BOOLEAN   FieldType = "BOOLEAN"
	LIST      FieldType = "LIST"
	MAP       FieldType = "MAP"
)

// RepetitionType - 重复类型枚举
type RepetitionType string

const (
	REQUIRED RepetitionType = "REQUIRED"
	OPTIONAL RepetitionType = "OPTIONAL"
	REPEATED RepetitionType = "REPEATED"
)

// NewSchemaBuilder - 创建新的 Schema 构建器
func NewSchemaBuilder(rootName string) *SchemaBuilder {
	return &SchemaBuilder{
		rootName: rootName,
		fields:   make([]Field, 0),
	}
}

// FieldOption - 字段选项函数类型
type FieldOption func(*Field)

// WithConvertedType - 设置转换类型选项
func WithConvertedType(convertedType string) FieldOption {
	return func(field *Field) {
		field.ConvertedType = convertedType
	}
}

// WithRepetitionType - 设置重复类型选项
func WithRepetitionType(repetitionType RepetitionType) FieldOption {
	return func(field *Field) {
		field.RepetitionType = string(repetitionType)
	}
}

// AddField - 添加基本字段
func (b *SchemaBuilder) AddField(name string, fieldType FieldType, opts ...FieldOption) *SchemaBuilder {
	field := Field{
		Name: fmt.Sprintf("name=%s, type=%s", name, string(fieldType)),
	}

	// 应用选项
	for _, opt := range opts {
		opt(&field)
	}

	// 构建完整的 Tag
	tagParts := []string{fmt.Sprintf("name=%s", name), fmt.Sprintf("type=%s", string(fieldType))}

	if field.ConvertedType != "" {
		tagParts = append(tagParts, fmt.Sprintf("convertedtype=%s", field.ConvertedType))
	}

	if field.RepetitionType != "" {
		tagParts = append(tagParts, fmt.Sprintf("repetitiontype=%s", field.RepetitionType))
	}

	field.Name = strings.Join(tagParts, ", ")

	b.fields = append(b.fields, field)
	return b
}

// AddStringField - 添加字符串字段
func (b *SchemaBuilder) AddStringField(name string, opts ...FieldOption) *SchemaBuilder {
	return b.AddField(name, BYTEARRAY, append(opts, WithConvertedType("UTF8"))...)
}

// AddInt32Field - 添加 Int32 字段
func (b *SchemaBuilder) AddInt32Field(name string, opts ...FieldOption) *SchemaBuilder {
	return b.AddField(name, INT32, opts...)
}

// AddInt64Field - 添加 Int64 字段
func (b *SchemaBuilder) AddInt64Field(name string, opts ...FieldOption) *SchemaBuilder {
	return b.AddField(name, INT64, opts...)
}

// AddFloatField - 添加 Float 字段
func (b *SchemaBuilder) AddFloatField(name string, opts ...FieldOption) *SchemaBuilder {
	return b.AddField(name, FLOAT, opts...)
}

// AddDoubleField - 添加 Double 字段
func (b *SchemaBuilder) AddDoubleField(name string, opts ...FieldOption) *SchemaBuilder {
	return b.AddField(name, DOUBLE, opts...)
}

// AddBooleanField - 添加 Boolean 字段
func (b *SchemaBuilder) AddBooleanField(name string, opts ...FieldOption) *SchemaBuilder {
	return b.AddField(name, BOOLEAN, opts...)
}

// AddListFieldWithElementType - 添加 List 字段（指定元素类型）
func (b *SchemaBuilder) AddListFieldWithElementType(name string, elementType FieldType, opts ...FieldOption) *SchemaBuilder {
	var elementTag string
	if elementType == BYTEARRAY {
		elementTag = "name=element, type=BYTE_ARRAY, convertedtype=UTF8"
	} else {
		elementTag = fmt.Sprintf("name=element, type=%s", string(elementType))
	}

	field := Field{
		Name: fmt.Sprintf("name=%s, type=%s", name, string(LIST)),
		Fields: []Field{
			{
				Name: elementTag,
			},
		},
	}

	// 应用选项
	for _, opt := range opts {
		opt(&field)
	}

	// 构建完整的 Tag
	tagParts := []string{fmt.Sprintf("name=%s", name), fmt.Sprintf("type=%s", string(LIST))}
	if field.RepetitionType != "" {
		tagParts = append(tagParts, fmt.Sprintf("repetitiontype=%s", field.RepetitionType))
	}
	field.Name = strings.Join(tagParts, ", ")

	b.fields = append(b.fields, field)
	return b
}

// AddStringListField - 添加字符串列表字段
func (b *SchemaBuilder) AddStringListField(name string, opts ...FieldOption) *SchemaBuilder {
	return b.AddListFieldWithElementType(name, BYTEARRAY, opts...)
}

// AddFloatListField - 添加 Float 列表字段
func (b *SchemaBuilder) AddFloatListField(name string, opts ...FieldOption) *SchemaBuilder {
	return b.AddListFieldWithElementType(name, FLOAT, opts...)
}
func (b *SchemaBuilder) AddBytesField(name string, opts ...FieldOption) *SchemaBuilder {
	return b.AddField(name, BYTEARRAY, append(opts, WithConvertedType("UTF8"))...)
}

// AddMapField - 添加 Map 字段（key 为字符串，value 为指定类型）
func (b *SchemaBuilder) AddMapField(name string, valueType FieldType, opts ...FieldOption) *SchemaBuilder {
	var valueElementTag string
	if valueType == BYTEARRAY {
		valueElementTag = "name=element, type=BYTE_ARRAY, convertedtype=UTF8"
	} else {
		valueElementTag = fmt.Sprintf("name=element, type=%s", string(valueType))
	}

	field := Field{
		Name: fmt.Sprintf("name=%s, type=%s", name, string(MAP)),
		Fields: []Field{
			{
				Name: "name=key, type=BYTE_ARRAY, convertedtype=UTF8",
			},
			{
				Name: fmt.Sprintf("name=value, type=%s", string(valueType)),
				Fields: []Field{
					{
						Name: valueElementTag,
					},
				},
			},
		},
	}

	// 应用选项
	for _, opt := range opts {
		opt(&field)
	}

	// 构建完整的 Tag
	tagParts := []string{fmt.Sprintf("name=%s", name), fmt.Sprintf("type=%s", string(MAP))}
	if field.RepetitionType != "" {
		tagParts = append(tagParts, fmt.Sprintf("repetitiontype=%s", field.RepetitionType))
	}
	field.Name = strings.Join(tagParts, ", ")

	b.fields = append(b.fields, field)
	return b
}

// AddNestedField - 添加嵌套字段（使用函数构建器）
func (b *SchemaBuilder) AddNestedField(name string, nestedBuilderFunc func(*SchemaBuilder), opts ...FieldOption) *SchemaBuilder {
	nestedBuilder := NewSchemaBuilder("nested")
	nestedBuilderFunc(nestedBuilder)

	field := Field{
		Name:   fmt.Sprintf("name=%s", name),
		Fields: nestedBuilder.fields,
	}

	// 应用选项
	for _, opt := range opts {
		opt(&field)
	}

	// 构建完整的 Tag
	tagParts := []string{fmt.Sprintf("name=%s", name)}
	if field.RepetitionType != "" {
		tagParts = append(tagParts, fmt.Sprintf("repetitiontype=%s", field.RepetitionType))
	}
	field.Name = strings.Join(tagParts, ", ")

	b.fields = append(b.fields, field)
	return b
}

// AddNestedListField - 添加嵌套列表字段（使用函数构建器）
func (b *SchemaBuilder) AddNestedListField(name string, elementBuilderFunc func(*SchemaBuilder), opts ...FieldOption) *SchemaBuilder {
	elementBuilder := NewSchemaBuilder("element")
	elementBuilderFunc(elementBuilder)

	field := Field{
		Name: fmt.Sprintf("name=%s, type=%s", name, string(LIST)),
		Fields: []Field{
			{
				Name:   "name=element",
				Fields: elementBuilder.fields,
			},
		},
	}

	// 应用选项
	for _, opt := range opts {
		opt(&field)
	}

	// 构建完整的 Tag
	tagParts := []string{fmt.Sprintf("name=%s", name), fmt.Sprintf("type=%s", string(LIST))}
	if field.RepetitionType != "" {
		tagParts = append(tagParts, fmt.Sprintf("repetitiontype=%s", field.RepetitionType))
	}
	field.Name = strings.Join(tagParts, ", ")

	b.fields = append(b.fields, field)
	return b
}

// AddListFieldWithNestedListValue - 添加包含嵌套列表值的字段（如 scores: Map<String, List<Float>>）
func (b *SchemaBuilder) AddListFieldWithNestedListValue(name string, elementType FieldType, opts ...FieldOption) *SchemaBuilder {
	var elementTag string
	if elementType == BYTEARRAY {
		elementTag = "name=element, type=BYTE_ARRAY, convertedtype=UTF8"
	} else {
		elementTag = fmt.Sprintf("name=element, type=%s", string(elementType))
	}

	field := Field{
		Name: fmt.Sprintf("name=%s, type=%s", name, string(LIST)),
		Fields: []Field{
			{
				Name: fmt.Sprintf("name=value, type=%s", string(LIST)),
				Fields: []Field{
					{
						Name: elementTag,
					},
				},
			},
		},
	}

	// 应用选项
	for _, opt := range opts {
		opt(&field)
	}

	// 构建完整的 Tag
	tagParts := []string{fmt.Sprintf("name=%s", name), fmt.Sprintf("type=%s", string(LIST))}
	if field.RepetitionType != "" {
		tagParts = append(tagParts, fmt.Sprintf("repetitiontype=%s", field.RepetitionType))
	}
	field.Name = strings.Join(tagParts, ", ")

	b.fields = append(b.fields, field)
	return b
}

// AddMapWithNestedListValue - 添加 Map 字段，其中值是列表（如 scores: Map<String, List<Float>>）
func (b *SchemaBuilder) AddMapWithNestedListValue(name string, valueType FieldType, opts ...FieldOption) *SchemaBuilder {
	var valueElementTag string
	if valueType == BYTEARRAY {
		valueElementTag = "name=element, type=BYTE_ARRAY, convertedtype=UTF8"
	} else {
		valueElementTag = fmt.Sprintf("name=element, type=%s", string(valueType))
	}

	field := Field{
		Name: fmt.Sprintf("name=%s, type=%s", name, string(MAP)),
		Fields: []Field{
			{
				Name: "name=key, type=BYTE_ARRAY, convertedtype=UTF8",
			},
			{
				Name: fmt.Sprintf("name=value, type=%s", string(LIST)),
				Fields: []Field{
					{
						Name: valueElementTag,
					},
				},
			},
		},
	}

	// 应用选项
	for _, opt := range opts {
		opt(&field)
	}

	// 构建完整的 Tag
	tagParts := []string{fmt.Sprintf("name=%s", name), fmt.Sprintf("type=%s", string(MAP))}
	if field.RepetitionType != "" {
		tagParts = append(tagParts, fmt.Sprintf("repetitiontype=%s", field.RepetitionType))
	}
	field.Name = strings.Join(tagParts, ", ")

	b.fields = append(b.fields, field)
	return b
}

// Build - 构建最终的 JSON Schema
func (b *SchemaBuilder) Build() (string, error) {
	schema := map[string]interface{}{
		"Tag":    fmt.Sprintf("name=%s", b.rootName),
		"Fields": b.fields,
	}

	jsonBytes, err := json.MarshalIndent(schema, "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// BuildCompact - 构建紧凑的 JSON Schema（无缩进）
func (b *SchemaBuilder) BuildCompact() (string, error) {
	schema := map[string]interface{}{
		"Tag":    fmt.Sprintf("name=%s", b.rootName),
		"Fields": b.fields,
	}

	jsonBytes, err := json.Marshal(schema)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
