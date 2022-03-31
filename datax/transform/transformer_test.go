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

package transform

import (
	"testing"

	"github.com/Breeze0806/go-etl/element"
)

func TestNilTransformer_DoTransform(t *testing.T) {
	r := element.NewDefaultRecord()
	type args struct {
		record element.Record
	}
	tests := []struct {
		name    string
		n       *NilTransformer
		args    args
		want    element.Record
		wantErr bool
	}{
		{
			name: "1",
			n:    &NilTransformer{},
			args: args{
				record: r,
			},
			want:    r,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.DoTransform(tt.args.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("NilTransformer.DoTransform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NilTransformer.DoTransform() = %v, want %v", got, tt.want)
			}
		})
	}
}
