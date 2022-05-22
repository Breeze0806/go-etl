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

package csv

import "testing"

func Test_gbkEncodeDecode(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name     string
		args     args
		wantDest string
		wantErr  bool
	}{
		{
			name: "1",
			args: args{
				src: "中文",
			},
			wantDest: "中文",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, _ := gbkEncoder(tt.args.src)
			gotDest, err := gbkDecoder(src)
			if (err != nil) != tt.wantErr {
				t.Errorf("gbkEncodeDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDest != tt.wantDest {
				t.Errorf("gbkEncodeDecode() = %v, want %v", gotDest, tt.wantDest)
			}
		})
	}
}
