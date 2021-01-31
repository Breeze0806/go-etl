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
