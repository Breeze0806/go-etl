package mysql

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"github.com/Breeze0806/go-etl/element"
)

type mockReceiver struct {
	err    error
	n      int
	ticker *time.Ticker
}

func newMockReceiver(n int, err error, wait time.Duration) *mockReceiver {
	return &mockReceiver{
		err:    err,
		n:      n,
		ticker: time.NewTicker(wait),
	}
}
func newMockReceiverWithoutWait(n int, err error) *mockReceiver {
	return &mockReceiver{
		err: err,
		n:   n,
	}
}
func (m *mockReceiver) GetFromReader() (element.Record, error) {
	m.n--
	if m.n <= 0 {
		return nil, m.err
	}
	if m.ticker != nil {
		select {
		case <-m.ticker.C:
			return element.NewDefaultRecord(), nil
		}
	}
	return element.NewDefaultRecord(), nil
}

func (m *mockReceiver) Shutdown() error {
	m.ticker.Stop()
	return nil
}

func TestTask_Init(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		t       *Task
		args    args
		conf    *config.JSON
		jobConf *config.JSON
		wantErr bool
	}{
		{
			name: "1",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"core":{
					"container":{
						"job":{
							"id" : 1
						},
						"taskGroup":{
							"id":  1
						}
					}
				},
				"job":{
					"content":[
						{
							"reader": {
								"name":"mysqlreader",
								"parameter":{

								}
							}
						}
					]
				}
			}`),
		},
		{
			name: "2",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromString(`{}`),
			jobConf: testJSONFromString(`{
				"job":{
					"core":{
						"container":{
							"job":{
								"id" : 1
							},
							"taskGroup":{
								"id":  1
							}
						}
					},
					"content":[
						{
							"reader": {
								"name":"mysqlreader",
								"parameter":{

								}
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
		{
			name: "3",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"core":{
					"container":{
						"job":{
							"id" : 1
						},
						"taskGroup":{
							"id":  1
						}
					}
				},
				"job":{
					"content":[
						{
							"reader": {
								"name":"mysqlreader"
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
		{
			name: "4",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"core":{
					"container":{
						"job":{
							"id" : 1
						},
						"taskGroup":{
							"id":  1
						}
					}
				},
				"job":{
					"content":[
						{
							"reader": {
								"name":"mysqlreader",
								"parameter":{
									"username": 1
								}
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
		{
			name: "5",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return nil, errors.New("mock error")
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"core":{
					"container":{
						"job":{
							"id" : 1
						},
						"taskGroup":{
							"id":  1
						}
					}
				},
				"job":{
					"content":[
						{
							"reader": {
								"name":"mysqlreader",
								"parameter":{
								}
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
		{
			name: "6",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{
						queryErr: errors.New("mock error"),
					}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"core":{
					"container":{
						"job":{
							"id" : 1
						},
						"taskGroup":{
							"id":  1
						}
					}
				},
				"job":{
					"content":[
						{
							"reader": {
								"name":"mysqlreader",
								"parameter":{
								}
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
		{
			name: "7",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{
						fetchErr: errors.New("mock error"),
					}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"core":{
					"container":{
						"job":{
							"id" : 1
						},
						"taskGroup":{
							"id":  1
						}
					}
				},
				"job":{
					"content":[
						{
							"reader": {
								"name":"mysqlreader",
								"parameter":{
								}
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
		{
			name: "8",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"core":{
					"container":{
						"job":{
							"id" : "1"
						},
						"taskGroup":{
							"id":  1
						}
					}
				},
				"job":{

					"content":[
						{
							"reader": {
								"name":"mysqlreader",
								"parameter":{

								}
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
		{
			name: "9",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"core":{
					"container":{
						"job":{
							"id" : 1
						},
						"taskGroup":{
							"id":  "1"
						}
					}
				},
				"job":{

					"content":[
						{
							"reader": {
								"name":"mysqlreader",
								"parameter":{

								}
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.SetPluginConf(tt.conf)
			tt.t.SetPluginJobConf(tt.jobConf)
			err := tt.t.Init(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Task.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_Destroy(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		t       *Task
		args    args
		wantErr bool
	}{
		{
			name: "1",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				execer:   &mockExecer{},
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.Destroy(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Task.Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_StartWrite(t *testing.T) {
	type args struct {
		ctx      context.Context
		receiver plugin.RecordReceiver
	}
	tests := []struct {
		name    string
		t       *Task
		args    args
		wait    time.Duration
		wantErr bool
	}{
		{
			name: "1",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				execer:   &mockExecer{},
				param:    newParameter(&paramConfig{}, &mockExecer{}),
			},
			args: args{
				ctx:      context.TODO(),
				receiver: newMockReceiver(1000, exchange.ErrTerminate, 1*time.Millisecond),
			},
		},
		{
			name: "2",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				execer:   &mockExecer{},
				param:    newParameter(&paramConfig{}, &mockExecer{}),
			},
			args: args{
				ctx:      context.TODO(),
				receiver: newMockReceiverWithoutWait(10000, exchange.ErrTerminate),
			},
		},

		{
			name: "3",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				execer:   &mockExecer{},
				param:    newParameter(&paramConfig{}, &mockExecer{}),
			},
			args: args{
				ctx:      context.TODO(),
				receiver: newMockReceiverWithoutWait(10000, errors.New("mock error")),
			},
			wantErr: true,
		},

		{
			name: "4",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				execer:   &mockExecer{},
				param:    newParameter(&paramConfig{}, &mockExecer{}),
			},
			args: args{
				ctx:      context.TODO(),
				receiver: newMockReceiverWithoutWait(10000, errors.New("mock error")),
			},
			wait:    100 * time.Microsecond,
			wantErr: false,
		},
		{
			name: "5",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				execer: &mockExecer{
					batchErr: errors.New("mock error"),
					batchN:   1,
				},
				param: newParameter(&paramConfig{}, &mockExecer{}),
			},
			args: args{
				ctx:      context.TODO(),
				receiver: newMockReceiver(1000, exchange.ErrTerminate, 1*time.Millisecond),
			},
			wantErr: true,
		},
		{
			name: "6",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				execer: &mockExecer{
					batchErr: errors.New("mock error"),
					batchN:   1,
				},
				param: newParameter(&paramConfig{}, &mockExecer{}),
			},
			args: args{
				ctx:      context.TODO(),
				receiver: newMockReceiverWithoutWait(10000, exchange.ErrTerminate),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer cancel()
			if tt.wait != 0 {
				go func() {
					<-time.After(tt.wait)
					cancel()
				}()
			}
			if err := tt.t.StartWrite(ctx, tt.args.receiver); (err != nil) != tt.wantErr {
				t.Errorf("Task.StartWrite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
