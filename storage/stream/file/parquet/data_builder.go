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

import "encoding/json"

type DataBuilder struct {
	data map[string]interface{}
}

// NewDataBuilder - 创建新的数据构建器
func NewDataBuilder() *DataBuilder {
	return &DataBuilder{
		data: make(map[string]interface{}),
	}
}

// SetString - 设置字符串字段
func (db *DataBuilder) SetString(name string, value string) *DataBuilder {
	db.data[name] = value
	return db
}

// SetInt32 - 设置 Int32 字段
func (db *DataBuilder) SetInt32(name string, value int32) *DataBuilder {
	db.data[name] = value
	return db
}

// SetInt64 - 设置 Int64 字段
func (db *DataBuilder) SetInt64(name string, value int64) *DataBuilder {
	db.data[name] = value
	return db
}

// SetFloat - 设置 Float 字段
func (db *DataBuilder) SetFloat(name string, value float32) *DataBuilder {
	db.data[name] = value
	return db
}

// SetDouble - 设置 Double 字段
func (db *DataBuilder) SetDouble(name string, value float64) *DataBuilder {
	db.data[name] = value
	return db
}

// SetBoolean - 设置 Boolean 字段
func (db *DataBuilder) SetBoolean(name string, value bool) *DataBuilder {
	db.data[name] = value
	return db
}

// SetStringList - 设置字符串列表字段
func (db *DataBuilder) SetStringList(name string, values []string) *DataBuilder {
	db.data[name] = values
	return db
}

// SetFloatList - 设置 Float 列表字段
func (db *DataBuilder) SetFloatList(name string, values []float32) *DataBuilder {
	db.data[name] = values
	return db
}

// SetInt64List - 设置 Int64 列表字段
func (db *DataBuilder) SetInt64List(name string, values []int64) *DataBuilder {
	db.data[name] = values
	return db
}

// SetBooleanList - 设置 Boolean 列表字段
func (db *DataBuilder) SetBooleanList(name string, values []bool) *DataBuilder {
	db.data[name] = values
	return db
}

// SetMap - 设置 Map 字段
func (db *DataBuilder) SetMap(name string, values map[string]interface{}) *DataBuilder {
	db.data[name] = values
	return db
}

// SetMapWithFloatList - 设置 Map<String, List<Float>> 字段
func (db *DataBuilder) SetMapWithFloatList(name string, values map[string][]float32) *DataBuilder {
	converted := make(map[string]interface{})
	for k, v := range values {
		converted[k] = v
	}
	db.data[name] = converted
	return db
}

// SetNestedField - 设置嵌套字段
func (db *DataBuilder) SetNestedField(name string, nestedData map[string]interface{}) *DataBuilder {
	db.data[name] = nestedData
	return db
}

// SetNestedListField - 设置嵌套列表字段
func (db *DataBuilder) SetNestedListField(name string, nestedListData []map[string]interface{}) *DataBuilder {
	db.data[name] = nestedListData
	return db
}

// SetRepeatedField - 设置重复字段（与嵌套列表相同）
func (db *DataBuilder) SetRepeatedField(name string, repeatedData []map[string]interface{}) *DataBuilder {
	db.data[name] = repeatedData
	return db
}

// Build - 构建最终的 JSON 数据
func (db *DataBuilder) Build() (string, error) {
	jsonBytes, err := json.Marshal(db.data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// BuildAsMap - 构建为 map
func (db *DataBuilder) BuildAsMap() map[string]interface{} {
	return db.data
}

// BuildAsInterface - 构建为 interface{}（用于 parquet-go 的 JSON writer）
func (db *DataBuilder) BuildAsInterface() interface{} {
	return db.data
}

// ==================== 批量数据构建器 ====================

// BatchDataBuilder - 批量数据构建器
type BatchDataBuilder struct {
	records []map[string]interface{}
}

// NewBatchDataBuilder - 创建批量数据构建器
func NewBatchDataBuilder() *BatchDataBuilder {
	return &BatchDataBuilder{
		records: make([]map[string]interface{}, 0),
	}
}

// AddRecord - 添加记录
func (bb *BatchDataBuilder) AddRecord(recordBuilder *DataBuilder) *BatchDataBuilder {
	bb.records = append(bb.records, recordBuilder.BuildAsMap())
	return bb
}

// AddRecordFromMap - 从 map 添加记录
func (bb *BatchDataBuilder) AddRecordFromMap(data map[string]interface{}) *BatchDataBuilder {
	bb.records = append(bb.records, data)
	return bb
}

// Build - 构建所有记录
func (bb *BatchDataBuilder) Build() ([]string, error) {
	jsonRecords := make([]string, len(bb.records))
	for i, record := range bb.records {
		jsonBytes, err := json.Marshal(record)
		if err != nil {
			return nil, err
		}
		jsonRecords[i] = string(jsonBytes)
	}
	return jsonRecords, nil
}

// BuildAsMaps - 构建为 map 切片
func (bb *BatchDataBuilder) BuildAsMaps() []map[string]interface{} {
	return bb.records
}
