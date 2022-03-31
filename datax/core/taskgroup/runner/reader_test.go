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

package runner

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

var errMockTest = errors.New("mock test error")

func TestReader_Run(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Reader
		args    args
		wantErr bool
	}{
		{
			name: "1",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, nil, nil, nil,
			}), &mockRecordSender{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "2",
			r: NewReader(newMockReaderTask([]error{
				errMockTest, nil, nil, nil, nil,
			}), &mockRecordSender{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "3",
			r: NewReader(newMockReaderTask([]error{
				nil, errMockTest, nil, nil, nil,
			}), &mockRecordSender{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "4",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, errMockTest, nil, nil,
			}), &mockRecordSender{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "5",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, nil, errMockTest, nil,
			}), &mockRecordSender{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "6",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, nil, nil, errMockTest,
			}), &mockRecordSender{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Reader.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReader_Plugin(t *testing.T) {
	tests := []struct {
		name string
		r    *Reader
		want plugin.Task
	}{
		{
			name: "1",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, nil, nil, nil,
			}), &mockRecordSender{}, "mock"),
			want: newMockReaderTask([]error{
				nil, nil, nil, nil, nil,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Plugin(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.Plugin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_Shutdown(t *testing.T) {
	tests := []struct {
		name    string
		r       *Reader
		wantErr bool
	}{
		{
			name: "1",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, nil, nil, nil,
			}), &mockRecordSender{}, "mock"),

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Shutdown(); (err != nil) != tt.wantErr {
				t.Errorf("Reader.Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
