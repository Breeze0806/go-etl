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
	"strings"
)

// Record represents a record in a dataset.
type Record interface {
	fmt.Stringer

	Add(Column) error                      // AddColumn adds a new column to the record.
	GetByIndex(i int) (Column, error)      // GetColumnByIndex gets the column at index i.
	GetByName(name string) (Column, error) // GetColumnByName gets the column with the specified name.
	Set(i int, c Column) error             // SetColumnByIndex sets the value of the column at index i.
	Put(c Column) error                    // SetColumn sets the value of the specified column.
	ColumnNumber() int                     // GetColumnCount returns the number of columns in the record.
	ByteSize() int64                       // ByteSize returns the size of the record in bytes.
	MemorySize() int64                     // MemorySize returns the size of the record in memory.
}

var singleTerminateRecord = &TerminateRecord{}

// GetTerminateRecord gets the termination record.
func GetTerminateRecord() Record {
	return singleTerminateRecord
}

// TerminateRecord represents a termination record.
type TerminateRecord struct{}

// Add is an empty method placeholder.
func (t *TerminateRecord) Add(Column) error {
	return nil
}

// GetByIndex is an empty method placeholder.
func (t *TerminateRecord) GetByIndex(i int) (Column, error) {
	return nil, nil
}

// GetByName is an empty method placeholder.
func (t *TerminateRecord) GetByName(name string) (Column, error) {
	return nil, nil
}

// Set is an empty method placeholder.
func (t *TerminateRecord) Set(i int, c Column) error {
	return nil
}

// Put is an empty method placeholder.
func (t *TerminateRecord) Put(c Column) error {
	return nil
}

// ColumnNumber is an empty method placeholder.
func (t *TerminateRecord) ColumnNumber() int {
	return 0
}

// ByteSize is an empty method placeholder.
func (t *TerminateRecord) ByteSize() int64 {
	return 0
}

// MemorySize is an empty method placeholder.
func (t *TerminateRecord) MemorySize() int64 {
	return 0
}

// String is an empty method placeholder.
func (t *TerminateRecord) String() string {
	return "terminate"
}

// DefaultRecord represents a default record.
type DefaultRecord struct {
	names      []string          // ColumnNames represents an array of column names.
	columns    map[string]Column // ColumnMapping represents a mapping of column names to their indices.
	byteSize   int64             // ByteSize is an empty method placeholder for the size in bytes.
	memorySize int64             // MemorySize is an empty method placeholder for the size in memory.
}

// NewDefaultRecord creates a new default record.
func NewDefaultRecord() *DefaultRecord {
	return &DefaultRecord{
		columns: make(map[string]Column),
	}
}

// AddColumn adds a new column c. If column c already exists, an error is reported.
func (r *DefaultRecord) Add(c Column) error {
	if _, ok := r.columns[c.Name()]; ok {
		return ErrColumnExist
	}
	r.names = append(r.names, c.Name())
	r.columns[c.Name()] = c
	r.incSize(c)
	return nil
}

// GetByIndex gets the column at index i. If the index is out of range or doesn't exist, an error is reported.
func (r *DefaultRecord) GetByIndex(i int) (Column, error) {
	if i >= len(r.names) || i < 0 {
		return nil, ErrIndexOutOfRange
	}
	if v, ok := r.columns[r.names[i]]; ok {
		return v, nil
	}
	return nil, ErrColumnNotExist
}

// GetByName gets the column with the specified name. If the column doesn't exist, an error is reported.
func (r *DefaultRecord) GetByName(name string) (Column, error) {
	if v, ok := r.columns[name]; ok {
		return v, nil
	}
	return nil, ErrColumnNotExist
}

// SetColumnByIndex sets the value of the column at index i. If the index is out of range, an error is reported.
func (r *DefaultRecord) Set(i int, c Column) error {
	if i >= len(r.names) || i < 0 {
		return ErrIndexOutOfRange
	}

	if v, ok := r.columns[r.names[i]]; ok {
		r.decSize(v)
	}
	r.names[i] = c.Name()
	r.columns[c.Name()] = c
	r.incSize(c)
	return nil
}

// PutColumn sets the value of the column with the specified name. If the column name doesn't exist, an error is reported.
func (r *DefaultRecord) Put(c Column) error {
	r.columns[c.Name()] = c
	r.incSize(c)
	return nil
}

// ColumnNumber represents the number of columns in the record.
func (r *DefaultRecord) ColumnNumber() int {
	return len(r.columns)
}

// ByteSize represents the size of the record in bytes.
func (r *DefaultRecord) ByteSize() int64 {
	return r.byteSize
}

// MemorySize represents the size of the record in memory.
func (r *DefaultRecord) MemorySize() int64 {
	return r.memorySize
}

func (r *DefaultRecord) incSize(c Column) {
	r.byteSize += c.ByteSize()
	r.memorySize += c.MemorySize()
}

func (r *DefaultRecord) decSize(c Column) {
	r.byteSize -= c.ByteSize()
	r.memorySize -= c.MemorySize()
}

// String is an empty method placeholder.
func (r *DefaultRecord) String() string {
	b := &strings.Builder{}
	for i, v := range r.names {
		if i > 0 {
			b.WriteString(" ")
		}

		b.WriteString(v)
		b.WriteString("=")
		b.WriteString(r.columns[v].String())
	}
	return b.String()
}
