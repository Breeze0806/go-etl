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

package datax

import (
	"context"
	"testing"

	"github.com/Breeze0806/go-etl/config"
)

func testJSONFromString(s string) *config.JSON {
	j, err := config.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return j
}

func TestModel_IsJob(t *testing.T) {
	tests := []struct {
		name string
		m    Model
		want bool
	}{
		{
			name: "1",
			m:    ModelJob,
			want: true,
		},
		{
			name: "2",
			m:    ModelTaskGroup,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsJob(); got != tt.want {
				t.Errorf("Model.IsJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_IsTaskGroup(t *testing.T) {
	tests := []struct {
		name string
		m    Model
		want bool
	}{
		{
			name: "1",
			m:    ModelJob,
			want: false,
		},
		{
			name: "2",
			m:    ModelTaskGroup,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsTaskGroup(); got != tt.want {
				t.Errorf("Model.IsTaskGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEngine_Start(t *testing.T) {
	tests := []struct {
		name    string
		e       *Engine
		wantErr bool
	}{
		{
			name: "1",
			e: NewEngine(context.TODO(), testJSONFromString(
				`{
					"core": {
						"container":{
							"model":"job"
						}
					}	
				}`)),
			wantErr: true,
		},
		{
			name: "2",
			e: NewEngine(context.TODO(), testJSONFromString(
				`{
					"core": {
						"container":{
							"model":"taskGroup"
						}
					}	
				}`)),
			wantErr: true,
		},
		{
			name: "3",
			e: NewEngine(context.TODO(), testJSONFromString(
				`{
					"core": {
						"container":{
							"model":"taskGroup1"
						}
					}	
				}`)),
			wantErr: true,
		},
		{
			name: "4",
			e: NewEngine(context.TODO(), testJSONFromString(
				`{
					"core": {
						"container":{
						}
					}	
				}`)),
			wantErr: true,
		},

		{
			name: "5",
			e: NewEngine(context.TODO(), testJSONFromString(
				`{
					"core": {
						"container":{
						}
					}	
				}`)),
			wantErr: true,
		},

		{
			name: "6",
			e: NewEngine(context.TODO(), testJSONFromString(
				`{
					"core": {
						"container":{
							"model":"job",
							"job":{
								"id":1
							}
						}
					}	
				}`)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Engine.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
