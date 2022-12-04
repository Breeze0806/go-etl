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

package schedule

import (
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/pingcap/errors"
)

type mockRetryJudger struct{}

func (m *mockRetryJudger) ShouldRetry(err error) bool {
	return err != nil
}

func testJSONFromString(s string) *config.JSON {
	json, err := config.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return json
}

func TestNoneRetryStrategy_Next(t *testing.T) {
	type args struct {
		err error
		n   int
	}
	tests := []struct {
		name      string
		r         RetryStrategy
		args      args
		wantRetry bool
		wantWait  time.Duration
	}{
		{
			name: "1",
			r:    NewNoneRetryStrategy(),
			args: args{
				err: nil,
				n:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRetry, gotWait := tt.r.Next(tt.args.err, tt.args.n)
			if gotRetry != tt.wantRetry {
				t.Errorf("NoneRetryStrategy.Next() gotRetry = %v, want %v", gotRetry, tt.wantRetry)
			}
			if !reflect.DeepEqual(gotWait, tt.wantWait) {
				t.Errorf("NoneRetryStrategy.Next() gotWait = %v, want %v", gotWait, tt.wantWait)
			}
		})
	}
}

func TestNTimesRetryStrategy_Next(t *testing.T) {
	type args struct {
		err error
		n   int
	}
	tests := []struct {
		name      string
		r         RetryStrategy
		args      args
		wantRetry bool
		wantWait  time.Duration
	}{
		{
			name: "1",
			r:    NewNTimesRetryStrategy(&mockRetryJudger{}, 0, 1*time.Second),
			args: args{
				err: errors.New("mock error"),
			},
		},
		{
			name: "2",
			r:    NewNTimesRetryStrategy(&mockRetryJudger{}, 2, 1*time.Second),
			args: args{
				err: errors.New("mock error"),
			},
			wantRetry: true,
			wantWait:  1 * time.Second,
		},
		{
			name: "3",
			r:    NewNTimesRetryStrategy(&mockRetryJudger{}, 2, 1*time.Second),
			args: args{
				err: nil,
			},
		},
		{
			name: "4",
			r:    NewNTimesRetryStrategy(&mockRetryJudger{}, 2, 1*time.Second),
			args: args{
				err: errors.New("mock error"),
				n:   3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRetry, gotWait := tt.r.Next(tt.args.err, tt.args.n)
			if gotRetry != tt.wantRetry {
				t.Errorf("NTimesRetryStrategy.Next() gotRetry = %v, want %v", gotRetry, tt.wantRetry)
			}
			if !reflect.DeepEqual(gotWait, tt.wantWait) {
				t.Errorf("NTimesRetryStrategy.Next() gotWait = %v, want %v", gotWait, tt.wantWait)
			}
		})
	}
}

func TestForeverRetryStrategy_Next(t *testing.T) {
	type args struct {
		err error
		in1 int
	}
	tests := []struct {
		name      string
		r         RetryStrategy
		args      args
		wantRetry bool
		wantWait  time.Duration
	}{
		{
			name: "1",
			r:    NewForeverRetryStrategy(&mockRetryJudger{}, 1*time.Second),
			args: args{
				err: errors.New("mock error"),
			},
			wantRetry: true,
			wantWait:  1 * time.Second,
		},
		{
			name: "2",
			r:    NewForeverRetryStrategy(&mockRetryJudger{}, 1*time.Second),
			args: args{
				err: errors.New("mock error"),
			},
			wantRetry: true,
			wantWait:  1 * time.Second,
		},
		{
			name: "3",
			r:    NewForeverRetryStrategy(&mockRetryJudger{}, 1*time.Second),
			args: args{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRetry, gotWait := tt.r.Next(tt.args.err, tt.args.in1)
			if gotRetry != tt.wantRetry {
				t.Errorf("ForeverRetryStrategy.Next() gotRetry = %v, want %v", gotRetry, tt.wantRetry)
			}
			if !reflect.DeepEqual(gotWait, tt.wantWait) {
				t.Errorf("ForeverRetryStrategy.Next() gotWait = %v, want %v", gotWait, tt.wantWait)
			}
		})
	}
}

func TestExponentialStrategy_Next(t *testing.T) {
	min := time.Duration(8) * time.Millisecond
	max := time.Duration(256) * time.Millisecond
	r := NewExponentialRetryStrategy(&mockRetryJudger{}, min, max)
	between := func(value time.Duration, a, b int) bool {
		x := int(value / time.Millisecond)
		return a <= x && x <= b
	}
	betweenZero := func(value time.Duration, a, b int) bool {
		return value == 0
	}
	type args struct {
		err error
		n   int
	}
	tests := []struct {
		name      string
		r         RetryStrategy
		args      args
		wantRetry bool
		between   func(value time.Duration, a, b int) bool
	}{
		{
			name: "0",
			r:    r,
			args: args{
				err: errors.New("mock error"),
			},
			wantRetry: true,
			between:   between,
		},
		{
			name: "1",
			r:    r,
			args: args{
				err: errors.New("mock error"),
				n:   1,
			},
			wantRetry: true,
			between:   between,
		},
		{
			name: "1",
			r:    r,
			args: args{
				err: errors.New("mock error"),
				n:   1,
			},
			wantRetry: true,
			between:   between,
		},
		{
			name: "2",
			r:    r,
			args: args{
				err: errors.New("mock error"),
				n:   2,
			},
			wantRetry: true,
			between:   between,
		},
		{
			name: "3",
			r:    r,
			args: args{
				err: errors.New("mock error"),
				n:   3,
			},
			wantRetry: true,
			between:   between,
		},
		{
			name: "4",
			r:    r,
			args: args{
				err: errors.New("mock error"),
				n:   4,
			},
			wantRetry: true,
			between:   between,
		},
		{
			name: "5",
			r:    r,
			args: args{
				err: errors.New("mock error"),
				n:   5,
			},
			wantRetry: false,
			between:   betweenZero,
		},
		{
			name: "6",
			r:    r,
			args: args{
				err: errors.New("mock error"),
				n:   6,
			},
			wantRetry: false,
			between:   betweenZero,
		},
		{
			name: "7",
			r:    r,
			args: args{
				err: nil,
				n:   1,
			},
			wantRetry: false,
			between:   betweenZero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRetry, gotWait := tt.r.Next(tt.args.err, tt.args.n)
			if gotRetry != tt.wantRetry {
				t.Errorf("ExponentialStrategy.Next() gotRetry = %v, want %v", gotRetry, tt.wantRetry)
			}
			if !tt.between(gotWait, 8, 256) {
				t.Errorf("ExponentialStrategy.Next() gotWait = %v", int64(gotWait))
			}
		})
	}
}

func TestNewRetryStrategy(t *testing.T) {
	type args struct {
		j    RetryJudger
		conf *config.JSON
	}
	tests := []struct {
		name    string
		args    args
		wantS   RetryStrategy
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString("{}"),
			},
			wantS: NewNoneRetryStrategy(),
		},
		{
			name: "2",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":""}`),
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":"simple"}}`),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":{}}}`),
			},
			wantErr: true,
		},
		{
			name: "5",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":"ntimes","strategy":{}}}`),
			},
			wantErr: true,
		},
		{
			name: "6",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":"ntimes","strategy":{"wait":1}}}`),
			},
			wantErr: true,
		},
		{
			name: "7",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":"ntimes","strategy":{"wait":"1s","n":3}}}`),
			},
			wantS: NewNTimesRetryStrategy(&mockRetryJudger{}, 3, 1*time.Second),
		},
		{
			name: "8",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":"forever","strategy":{}}}`),
			},
			wantErr: true,
		},
		{
			name: "9",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":"forever","strategy":{"wait":1}}}`),
			},
			wantErr: true,
		},
		{
			name: "10",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":"forever","strategy":{"wait":"1s"}}}`),
			},
			wantS: NewForeverRetryStrategy(&mockRetryJudger{}, 1*time.Second),
		},
		{
			name: "11",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":"exponential","strategy":{}}}`),
			},
			wantErr: true,
		},
		{
			name: "12",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":"exponential","strategy":{"init":1}}}`),
			},
			wantErr: true,
		},
		{
			name: "13",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":"exponential","strategy":{"init":"1s","max":"2s"}}}`),
			},
			wantS: NewExponentialRetryStrategy(&mockRetryJudger{}, 1*time.Second, 2*time.Second),
		},
		{
			name: "14",
			args: args{
				j:    &mockRetryJudger{},
				conf: testJSONFromString(`{"retry":{"type":"simple","strategy":{"init":"1s","max":"2s"}}}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := NewRetryStrategy(tt.args.j, tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRetryStrategy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("NewRetryStrategy() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
