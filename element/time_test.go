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
	"reflect"
	"testing"
	"time"
)

func TestStringTimeEncoder_TimeEncode(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		e       *StringTimeEncoder
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "1",
			e:    NewStringTimeEncoder(defaultTimeFormat).(*StringTimeEncoder),
			args: args{
				i: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.TimeEncode(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringTimeEncoder.TimeEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringTimeEncoder.TimeEncode() = %v, want %v", got, tt.want)
			}
		})
	}
}
