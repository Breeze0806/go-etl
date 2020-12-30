package taskgroup

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/datax/common/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
)

func testTaskExecer(ctx context.Context, taskConf *config.Json, prefixKey string, attemptCount int) *taskExecer {
	t, err := newTaskExecer(ctx, taskConf, prefixKey, attemptCount)
	if err != nil {
		panic(err)
	}
	return t
}

func initLoader(name string, errs []error) {
	loader.RegisterReader(name, newMockReader(errs))
	loader.RegisterWriter(name, newMockWriter(errs))
}

func Test_newTaskExecer(t *testing.T) {
	resetLoader()
	initLoader("mock", []error{
		nil, nil, nil, nil, nil,
	})

	type args struct {
		ctx          context.Context
		taskConf     *config.Json
		prefixKey    string
		attemptCount int
	}
	tests := []struct {
		name    string
		args    args
		wantT   *taskExecer
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ctx: context.Background(),
				taskConf: testJsonFromString(`{
						"taskId":1,
						"reader":{
							"name":"mock"
						},
						"writer":{
							"name":"mock"
						}
					}`),
				prefixKey:    "mock",
				attemptCount: 0,
			},
			wantErr: false,
		},

		{
			name: "2",
			args: args{
				ctx: context.Background(),
				taskConf: testJsonFromString(`{
						"taskId":2,
						"reader":{
							"name":"mock2"
						},
						"writer":{
							"name":"mock2"
						}
					}`),
				prefixKey:    "mock2",
				attemptCount: 0,
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				ctx: context.Background(),
				taskConf: testJsonFromString(`{
						"taskId":2,
						"reader":{
							"name":1
						},
						"writer":{
							"name":"mock2"
						}
					}`),
				prefixKey:    "mock2",
				attemptCount: 0,
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				ctx: context.Background(),
				taskConf: testJsonFromString(`{
						"taskId":3,
						"reader":{
							"name":"mock"
						},
						"writer":{
							"name":"mock2"
						}
					}`),
				prefixKey:    "mock2",
				attemptCount: 0,
			},
			wantErr: true,
		},
		{
			name: "5",
			args: args{
				ctx: context.Background(),
				taskConf: testJsonFromString(`{
						"taskId":2,
						"reader":{
							"name":"mock"
						},
						"writer":{
							"name":2
						}
					}`),
				prefixKey:    "mock2",
				attemptCount: 0,
			},
			wantErr: true,
		},
		{
			name: "6",
			args: args{
				ctx: context.Background(),
				taskConf: testJsonFromString(`{
						"taskId":"6",
						"reader":{
							"name":"mock"
						},
						"writer":{
							"name":"mock"
						}
					}`),
				prefixKey:    "mock",
				attemptCount: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := newTaskExecer(tt.args.ctx, tt.args.taskConf, tt.args.prefixKey, tt.args.attemptCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("newTaskExecer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("newTaskExecer() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}

func Test_taskExecer_Do(t *testing.T) {
	resetLoader()
	initLoader("mock", []error{
		nil, nil, nil, nil, nil,
	})
	initLoader("mock1", []error{
		errors.New("mock test error"), nil, nil, nil, nil,
	})
	tests := []struct {
		name    string
		t       *taskExecer
		wantErr bool
	}{
		{
			name: "1",
			t: testTaskExecer(context.Background(), testJsonFromString(`{
				"taskId":1,
				"reader":{
					"name":"mock"
				},
				"writer":{
					"name":"mock"
				}
			}`), `mock`, 0),
			wantErr: false,
		},

		{
			name: "2",
			t: testTaskExecer(context.Background(), testJsonFromString(`{
				"taskId":1,
				"reader":{
					"name":"mock1"
				},
				"writer":{
					"name":"mock1"
				}
			}`), `mock1`, 0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.t.Do()
			if (err != nil) != tt.wantErr {
				t.Errorf("taskExecer.Do() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(err)
		})
	}
}
