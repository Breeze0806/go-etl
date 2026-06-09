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
	"testing"
)

func TestDataBuilder_StringField(t *testing.T) {

	db := NewDataBuilder().
		SetString("name", "Alice")

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if data["name"] != "Alice" {
		t.Errorf("Expected name to be 'Alice', got %v", data["name"])
	}
}

func TestDataBuilder_Int32Field(t *testing.T) {

	db := NewDataBuilder().
		SetInt32("age", 30)

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if data["age"] != float64(30) { // JSON unmarshals numbers as float64
		t.Errorf("Expected age to be 30, got %v", data["age"])
	}
}

func TestDataBuilder_Int64Field(t *testing.T) {

	db := NewDataBuilder().
		SetInt64("id", 1234567890123456789)

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if data["id"] != float64(1234567890123456789) {
		t.Errorf("Expected id to be 1234567890123456789, got %v", data["id"])
	}
}

func TestDataBuilder_FloatField(t *testing.T) {

	db := NewDataBuilder().
		SetFloat("weight", 55.5)

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if data["weight"] != 55.5 {
		t.Errorf("Expected weight to be 55.5, got %v", data["weight"])
	}
}

func TestDataBuilder_DoubleField(t *testing.T) {

	db := NewDataBuilder().
		SetDouble("salary", 75000.50)

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if data["salary"] != 75000.50 {
		t.Errorf("Expected salary to be 75000.50, got %v", data["salary"])
	}
}

func TestDataBuilder_BooleanField(t *testing.T) {

	db := NewDataBuilder().
		SetBoolean("active", true)

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if data["active"] != true {
		t.Errorf("Expected active to be true, got %v", data["active"])
	}
}

func TestDataBuilder_StringListField(t *testing.T) {

	db := NewDataBuilder().
		SetStringList("tags", []string{"tag1", "tag2", "tag3"})

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	tags, ok := data["tags"].([]interface{})
	if !ok {
		t.Fatalf("Expected tags to be []interface{}, got %T", data["tags"])
	}

	expected := []string{"tag1", "tag2", "tag3"}
	if len(tags) != len(expected) {
		t.Errorf("Expected %d tags, got %d", len(expected), len(tags))
	}

	for i, tag := range tags {
		if tag != expected[i] {
			t.Errorf("Expected tag[%d] to be %s, got %v", i, expected[i], tag)
		}
	}
}

func TestDataBuilder_FloatListField(t *testing.T) {

	db := NewDataBuilder().
		SetFloatList("scores", []float32{95.5, 88.0, 92.5})

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	scores, ok := data["scores"].([]interface{})
	if !ok {
		t.Fatalf("Expected scores to be []interface{}, got %T", data["scores"])
	}

	expected := []float64{95.5, 88.0, 92.5} // JSON unmarshals to float64
	if len(scores) != len(expected) {
		t.Errorf("Expected %d scores, got %d", len(expected), len(scores))
	}

	for i, score := range scores {
		if score != expected[i] {
			t.Errorf("Expected score[%d] to be %v, got %v", i, expected[i], score)
		}
	}
}

func TestDataBuilder_MapField(t *testing.T) {

	properties := map[string]interface{}{
		"color": "red",
		"size":  "large",
	}

	db := NewDataBuilder().
		SetMap("properties", properties)

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	props, ok := data["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected properties to be map[string]interface{}, got %T", data["properties"])
	}

	if props["color"] != "red" || props["size"] != "large" {
		t.Errorf("Expected properties to be {'color': 'red', 'size': 'large'}, got %v", props)
	}
}

func TestDataBuilder_MapWithFloatListField(t *testing.T) {

	scores := map[string][]float32{
		"Math":    {95.0, 98.0, 92.0},
		"Science": {88.5, 91.0},
	}

	db := NewDataBuilder().
		SetMapWithFloatList("scores", scores)

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	scoresData, ok := data["scores"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected scores to be map[string]interface{}, got %T", data["scores"])
	}

	mathScores, ok := scoresData["Math"].([]interface{})
	if !ok {
		t.Fatalf("Expected Math scores to be []interface{}, got %T", scoresData["Math"])
	}

	expectedMath := []float64{95.0, 98.0, 92.0}
	if len(mathScores) != len(expectedMath) {
		t.Errorf("Expected %d Math scores, got %d", len(expectedMath), len(mathScores))
	}

	for i, score := range mathScores {
		if score != expectedMath[i] {
			t.Errorf("Expected Math score[%d] to be %v, got %v", i, expectedMath[i], score)
		}
	}
}

func TestDataBuilder_NestedField(t *testing.T) {

	address := map[string]interface{}{
		"street": "123 Main St",
		"city":   "New York",
		"zip":    float64(10001), // JSON unmarshals to float64
	}

	db := NewDataBuilder().
		SetNestedField("address", address)

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	addr, ok := data["address"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected address to be map[string]interface{}, got %T", data["address"])
	}

	if addr["street"] != "123 Main St" || addr["city"] != "New York" || addr["zip"] != float64(10001) {
		t.Errorf("Expected address to be {'street': '123 Main St', 'city': 'New York', 'zip': 10001}, got %v", addr)
	}
}

func TestDataBuilder_NestedListField(t *testing.T) {

	friends := []map[string]interface{}{
		{"name": "Alice", "id": float64(1)},
		{"name": "Bob", "id": float64(2)},
	}

	db := NewDataBuilder().
		SetNestedListField("friends", friends)

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	friendsData, ok := data["friends"].([]interface{})
	if !ok {
		t.Fatalf("Expected friends to be []interface{}, got %T", data["friends"])
	}

	if len(friendsData) != 2 {
		t.Errorf("Expected 2 friends, got %d", len(friendsData))
	}

	friend1 := friendsData[0].(map[string]interface{})
	if friend1["name"] != "Alice" || friend1["id"] != float64(1) {
		t.Errorf("Expected first friend to be {'name': 'Alice', 'id': 1}, got %v", friend1)
	}
}

func TestDataBuilder_RepeatedField(t *testing.T) {

	teachers := []map[string]interface{}{
		{"name": "Teacher1", "id": float64(101)},
		{"name": "Teacher2", "id": float64(102)},
	}

	db := NewDataBuilder().
		SetRepeatedField("teachers", teachers)

	result, err := db.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	teachersData, ok := data["teachers"].([]interface{})
	if !ok {
		t.Fatalf("Expected teachers to be []interface{}, got %T", data["teachers"])
	}

	if len(teachersData) != 2 {
		t.Errorf("Expected 2 teachers, got %d", len(teachersData))
	}

	teacher1 := teachersData[0].(map[string]interface{})
	if teacher1["name"] != "Teacher1" || teacher1["id"] != float64(101) {
		t.Errorf("Expected first teacher to be {'name': 'Teacher1', 'id': 101}, got %v", teacher1)
	}
}

func TestDataBuilder_BuildAsMap(t *testing.T) {
	db := NewDataBuilder().
		SetString("name", "Test").
		SetInt32("age", 25)

	resultMap := db.BuildAsMap()

	if resultMap["name"] != "Test" {
		t.Errorf("Expected name to be 'Test', got %v", resultMap["name"])
	}

	if resultMap["age"] != int32(25) {
		t.Errorf("Expected age to be 25, got %v", resultMap["age"])
	}
}

func TestDataBuilder_BuildAsInterface(t *testing.T) {
	db := NewDataBuilder().
		SetString("name", "Test")

	resultInterface := db.BuildAsInterface()

	resultMap, ok := resultInterface.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result to be map[string]interface{}, got %T", resultInterface)
	}

	if resultMap["name"] != "Test" {
		t.Errorf("Expected name to be 'Test', got %v", resultMap["name"])
	}
}

func TestBatchDataBuilder_AddRecord(t *testing.T) {
	record1 := NewDataBuilder().
		SetString("name", "Alice").
		SetInt32("age", 30)

	record2 := NewDataBuilder().
		SetString("name", "Bob").
		SetInt32("age", 25)

	bb := NewBatchDataBuilder().
		AddRecord(record1).
		AddRecord(record2)

	records := bb.BuildAsMaps()

	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}

	if records[0]["name"] != "Alice" || records[0]["age"] != int32(30) {
		t.Errorf("Expected first record to be {'name': 'Alice', 'age': 30}, got %v", records[0])
	}

	if records[1]["name"] != "Bob" || records[1]["age"] != int32(25) {
		t.Errorf("Expected second record to be {'name': 'Bob', 'age': 25}, got %v", records[1])
	}
}

func TestBatchDataBuilder_AddRecordFromMap(t *testing.T) {
	bb := NewBatchDataBuilder().
		AddRecordFromMap(map[string]interface{}{
			"name": "Charlie",
			"age":  int32(35),
		})

	records := bb.BuildAsMaps()

	if len(records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(records))
	}

	if records[0]["name"] != "Charlie" || records[0]["age"] != int32(35) {
		t.Errorf("Expected record to be {'name': 'Charlie', 'age': 35}, got %v", records[0])
	}
}

func TestBatchDataBuilder_Build(t *testing.T) {

	bb := NewBatchDataBuilder().
		AddRecordFromMap(map[string]interface{}{
			"name": "Test1",
		}).
		AddRecordFromMap(map[string]interface{}{
			"name": "Test2",
		})

	records, err := bb.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}

	var record1 map[string]interface{}
	if err := json.Unmarshal([]byte(records[0]), &record1); err != nil {
		t.Fatalf("Unmarshal first record failed: %v", err)
	}

	if record1["name"] != "Test1" {
		t.Errorf("Expected first record name to be 'Test1', got %v", record1["name"])
	}
}
