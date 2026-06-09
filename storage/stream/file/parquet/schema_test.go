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
	"testing"
)

func TestSchemaBuilder_AddStringField(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddStringField("name", WithRepetitionType(OPTIONAL))

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)
	fmt.Println(tag)
	if !strings.Contains(tag, "name=name") ||
		!strings.Contains(tag, "type=BYTE_ARRAY") ||
		!strings.Contains(tag, "convertedtype=UTF8") ||
		!strings.Contains(tag, "repetitiontype=OPTIONAL") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}
}

func TestSchemaBuilder_AddInt32Field(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddInt32Field("age")

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=age") ||
		!strings.Contains(tag, "type=INT32") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}
}

func TestSchemaBuilder_AddInt64Field(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddInt64Field("id")

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=id") ||
		!strings.Contains(tag, "type=INT64") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}
}

func TestSchemaBuilder_AddFloatField(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddFloatField("weight")

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=weight") ||
		!strings.Contains(tag, "type=FLOAT") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}
}

func TestSchemaBuilder_AddDoubleField(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddDoubleField("salary")

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=salary") ||
		!strings.Contains(tag, "type=DOUBLE") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}
}

func TestSchemaBuilder_AddBooleanField(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddBooleanField("active")

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=active") ||
		!strings.Contains(tag, "type=BOOLEAN") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}
}

func TestSchemaBuilder_AddStringListField(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddStringListField("tags")

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=tags") ||
		!strings.Contains(tag, "type=LIST") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}

	nestedFields := field["Fields"].([]interface{})
	if len(nestedFields) != 1 {
		t.Fatalf("Expected 1 nested field, got %d", len(nestedFields))
	}

	elementField := nestedFields[0].(map[string]interface{})
	elementTag := elementField["Tag"].(string)
	fmt.Println(elementTag)
	if !strings.Contains(elementTag, "name=element") ||
		!strings.Contains(elementTag, "type=BYTE_ARRAY") ||
		!strings.Contains(elementTag, "convertedtype=UTF8") {
		t.Errorf("Element field tag does not match expected format: %s", elementTag)
	}
}

func TestSchemaBuilder_AddFloatListField(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddFloatListField("scores")

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=scores") ||
		!strings.Contains(tag, "type=LIST") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}

	nestedFields := field["Fields"].([]interface{})
	if len(nestedFields) != 1 {
		t.Fatalf("Expected 1 nested field, got %d", len(nestedFields))
	}

	elementField := nestedFields[0].(map[string]interface{})
	elementTag := elementField["Tag"].(string)

	if !strings.Contains(elementTag, "name=element") ||
		!strings.Contains(elementTag, "type=FLOAT") {
		t.Errorf("Element field tag does not match expected format: %s", elementTag)
	}
}

func TestSchemaBuilder_AddMapField(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddMapField("properties", BYTEARRAY)

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=properties") ||
		!strings.Contains(tag, "type=MAP") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}

	nestedFields := field["Fields"].([]interface{})
	if len(nestedFields) != 2 {
		t.Fatalf("Expected 2 nested fields (key and value), got %d", len(nestedFields))
	}

	keyField := nestedFields[0].(map[string]interface{})
	keyTag := keyField["Tag"].(string)
	if !strings.Contains(keyTag, "name=key") ||
		!strings.Contains(keyTag, "type=BYTE_ARRAY") ||
		!strings.Contains(keyTag, "convertedtype=UTF8") {
		t.Errorf("Key field tag does not match expected format: %s", keyTag)
	}

	valueField := nestedFields[1].(map[string]interface{})
	valueTag := valueField["Tag"].(string)
	if !strings.Contains(valueTag, "name=value") ||
		!strings.Contains(valueTag, "type=BYTE_ARRAY") {
		t.Errorf("Value field tag does not match expected format: %s", valueTag)
	}
}

func TestSchemaBuilder_AddMapWithNestedListValue(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddMapWithNestedListValue("scores", FLOAT)

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=scores") ||
		!strings.Contains(tag, "type=MAP") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}

	nestedFields := field["Fields"].([]interface{})
	if len(nestedFields) != 2 {
		t.Fatalf("Expected 2 nested fields (key and value), got %d", len(nestedFields))
	}

	keyField := nestedFields[0].(map[string]interface{})
	keyTag := keyField["Tag"].(string)
	if !strings.Contains(keyTag, "name=key") ||
		!strings.Contains(keyTag, "type=BYTE_ARRAY") ||
		!strings.Contains(keyTag, "convertedtype=UTF8") {
		t.Errorf("Key field tag does not match expected format: %s", keyTag)
	}

	valueField := nestedFields[1].(map[string]interface{})
	valueTag := valueField["Tag"].(string)
	if !strings.Contains(valueTag, "name=value") ||
		!strings.Contains(valueTag, "type=LIST") {
		t.Errorf("Value field tag does not match expected format: %s", valueTag)
	}

	valueNestedFields := valueField["Fields"].([]interface{})
	if len(valueNestedFields) != 1 {
		t.Fatalf("Expected 1 value nested field, got %d", len(valueNestedFields))
	}

	elementField := valueNestedFields[0].(map[string]interface{})
	elementTag := elementField["Tag"].(string)
	if !strings.Contains(elementTag, "name=element") ||
		!strings.Contains(elementTag, "type=FLOAT") {
		t.Errorf("Element field tag does not match expected format: %s", elementTag)
	}
}

func TestSchemaBuilder_AddNestedField(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddNestedField("address", func(nb *SchemaBuilder) {
			nb.AddStringField("street").
				AddInt32Field("zip")
		})

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=address") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}

	nestedFields := field["Fields"].([]interface{})
	if len(nestedFields) != 2 {
		t.Fatalf("Expected 2 nested fields, got %d", len(nestedFields))
	}
}

func TestSchemaBuilder_AddNestedListField(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddNestedListField("friends", func(nb *SchemaBuilder) {
			nb.AddStringField("name").
				AddInt64Field("id")
		})

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=friends") ||
		!strings.Contains(tag, "type=LIST") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}

	nestedFields := field["Fields"].([]interface{})
	if len(nestedFields) != 1 {
		t.Fatalf("Expected 1 nested field (element), got %d", len(nestedFields))
	}

	elementField := nestedFields[0].(map[string]interface{})
	elementTag := elementField["Tag"].(string)
	if !strings.Contains(elementTag, "name=element") {
		t.Errorf("Element field tag does not match expected format: %s", elementTag)
	}

	elementNestedFields := elementField["Fields"].([]interface{})
	if len(elementNestedFields) != 2 {
		t.Fatalf("Expected 2 element nested fields, got %d", len(elementNestedFields))
	}
}

func TestSchemaBuilder_WithRepetitionType(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddStringField("name", WithRepetitionType(REPEATED))

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "repetitiontype=REPEATED") {
		t.Errorf("Field tag should contain repetitiontype=REPEATED: %s", tag)
	}
}

func TestSchemaBuilder_MultipleFields(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddStringField("name").
		AddInt32Field("age").
		AddBooleanField("active")

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 3 {
		t.Fatalf("Expected 3 fields, got %d", len(fields))
	}

	tags := make([]string, len(fields))
	for i, field := range fields {
		fieldMap := field.(map[string]interface{})
		tags[i] = fieldMap["Tag"].(string)
	}

	hasName := false
	hasAge := false
	hasActive := false

	for _, tag := range tags {
		if strings.Contains(tag, "name=name") {
			hasName = true
		}
		if strings.Contains(tag, "name=age") {
			hasAge = true
		}
		if strings.Contains(tag, "name=active") {
			hasActive = true
		}
	}

	if !hasName || !hasAge || !hasActive {
		t.Errorf("Expected to find all 3 fields, got tags: %v", tags)
	}
}

func TestSchemaBuilder_BuildCompact(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddStringField("name")

	result, err := schema.BuildCompact()
	if err != nil {
		t.Fatalf("BuildCompact failed: %v", err)
	}

	// Should not contain newlines or indentation
	if strings.Contains(result, "\n") || strings.Contains(result, "    ") {
		t.Errorf("Compact result should not contain newlines or indentation: %s", result)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "name=name") {
		t.Errorf("Field tag does not match expected format: %s", tag)
	}
}

func TestSchemaBuilder_WithConvertedType(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddField("name", BYTEARRAY, WithConvertedType("UTF8"))

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "convertedtype=UTF8") {
		t.Errorf("Field tag should contain convertedtype=UTF8: %s", tag)
	}
}

func TestSchemaBuilder_WithRepetitionTypeAndConvertedType(t *testing.T) {
	schema := NewSchemaBuilder("test-root").
		AddField("name", BYTEARRAY, WithConvertedType("UTF8"), WithRepetitionType(OPTIONAL))

	result, err := schema.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fields := parsed["Fields"].([]interface{})
	field := fields[0].(map[string]interface{})
	tag := field["Tag"].(string)

	if !strings.Contains(tag, "convertedtype=UTF8") || !strings.Contains(tag, "repetitiontype=OPTIONAL") {
		t.Errorf("Field tag should contain both convertedtype and repetitiontype: %s", tag)
	}
}
