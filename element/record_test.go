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
	"testing"
)

func TestDefaultRecord(t *testing.T) {
	r := NewDefaultRecord()
	type args struct {
		c Column
	}
	tests := []struct {
		args    args
		wantErr bool
	}{
		{
			args: args{
				NewDefaultColumn(NewNilBigIntColumnValue(), "test", 0),
			},
		},
		{

			args: args{
				NewDefaultColumn(NewNilBigIntColumnValue(), "test", 0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if err := r.Add(tt.args.c); (err != nil) != tt.wantErr {
			t.Errorf("DefaultRecord.Add() error = %v, wantErr %v", err, tt.wantErr)
		}
	}

	_, err := r.GetByIndex(0)
	if err != nil {
		t.Errorf("DefaultRecord.GetByIndex() error = %v, wantErr true", err)
		return
	}

	_, err = r.GetByIndex(1)
	if err == nil {
		t.Errorf("DefaultRecord.GetByIndex() error = %v, wantErr false", err)
		return
	}

	_, err = r.GetByName("test")
	if err != nil {
		t.Errorf("DefaultRecord.GetByName() error = %v, wantErr true", err)
		return
	}

	_, err = r.GetByName("")
	if err == nil {
		t.Errorf("DefaultRecord.GetByName() error = %v, wantErr false", err)
		return
	}

	s := r.ByteSize()
	if s != 0 {
		t.Errorf("DefaultRecord.ByteSize() = %v, want 0", s)
		return
	}

	s = r.MemorySize()
	if s != 8 {
		t.Errorf("DefaultRecord.ByteSize() = %v, want 8", s)
		return
	}

	n := r.ColumnNumber()
	if n != 1 {
		t.Errorf("DefaultRecord.ByteSize() = %v, want 1", n)
		return
	}

	err = r.Set(0, NewDefaultColumn(NewNilBoolColumnValue(), "test", 10))
	if err != nil {
		t.Errorf("DefaultRecord.Set() = %v, want 1", n)
		return
	}

	err = r.Set(1, NewDefaultColumn(NewNilBoolColumnValue(), "test", 10))
	if err == nil {
		t.Errorf("DefaultRecord.Set() = %v, want 1", n)
		return
	}
}

func Test_terminateRecord(t *testing.T) {
	r := GetTerminateRecord()
	if err := r.Add(nil); err != nil {
		t.Errorf("terminateRecord.Add() error = %v, wantErr false", err)
	}

	_, err := r.GetByIndex(0)
	if err != nil {
		t.Errorf("terminateRecord.GetByIndex() error = %v, wantErr true", err)
		return
	}

	_, err = r.GetByName("test")
	if err != nil {
		t.Errorf("terminateRecord.GetByName() error = %v, wantErr true", err)
		return
	}

	s := r.ByteSize()
	if s != 0 {
		t.Errorf("terminateRecord.ByteSize() = %v, want 0", s)
		return
	}

	s = r.MemorySize()
	if s != 0 {
		t.Errorf("terminateRecord.ByteSize() = %v, want 0", s)
		return
	}

	n := r.ColumnNumber()
	if n != 0 {
		t.Errorf("terminateRecord.ByteSize() = %v, want 0", n)
		return
	}

	err = r.Set(0, nil)
	if err != nil {
		t.Errorf("terminateRecord.Set() = %v, want 1", n)
		return
	}
}
