package taskgroup

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/schedule"
)

func testContainer(ctx context.Context, conf *config.JSON) *Container {
	c, err := NewContainer(ctx, conf)
	if err != nil {
		panic(err)
	}
	return c
}

func TestContainer_Do(t *testing.T) {
	resetLoader()
	loader.RegisterReader("mock", newMockRandReader([]error{
		nil, nil, errors.New("mock test error"), nil, nil,
	}))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}))
	content := testJsonFromString(`{
		"core" : {
			"container": {
				"job":{
					"id": 1,
					"sleepInterval":100
				},
				"taskGroup":{
					"id": 1,
					"failover":{
						"retryIntervalInMsec":10
					}
				}
			}
		}
	}`)
	for i := 0; i < 1000; i++ {
		content.SetRawString(coreconst.DataxJobContent+fmt.Sprintf(".%d", i), fmt.Sprintf(`{
			"taskId": %d,
			"reader":{
				"name":"mock"
			},
			"writer":{
				"name":"mock"
			}
		}`, i))
	}
	ctx := context.Background()
	c, _ := NewContainer(ctx, content)
	if err := c.Do(); err != nil {
		t.Errorf("Do error: %v", err)
	}
}

func TestContainer_DoCancel1(t *testing.T) {
	resetLoader()
	loader.RegisterReader("mock", newMockRandReader([]error{
		nil, nil, errors.New("mock test error"), nil, nil,
	}))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}))
	content := testJsonFromString(`{
		"core" : {
			"container": {
				"job":{
					"id": 1,
					"sleepInterval":100
				},
				"taskGroup":{
					"id": 1,
					"failover":{
						"retryIntervalInMsec":10
					}
				}
			}
		}
	}`)
	for i := 0; i < 1000; i++ {
		content.SetRawString(coreconst.DataxJobContent+fmt.Sprintf(".%d", i), fmt.Sprintf(`{
			"taskId": %d,
			"reader":{
				"name":"mock"
			},
			"writer":{
				"name":"mock"
			}
		}`, i))
	}
	ctx, cancel := context.WithCancel(context.Background())
	c, _ := NewContainer(ctx, content)
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	if err := c.Do(); err != context.Canceled {
		t.Errorf("Do error: %v", err)
	}
}

func TestContainer_DoCancel2(t *testing.T) {
	resetLoader()
	loader.RegisterReader("mock", newMockRandReader([]error{
		nil, nil, errors.New("mock test error"), nil, nil,
	}))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}))
	content := testJsonFromString(`{
		"core" : {
			"container": {
				"job":{
					"id": 1,
					"sleepInterval":100
				},
				"taskGroup":{
					"id": 1,
					"failover":{
						"retryIntervalInMsec":10
					}
				}
			}
		}
	}`)
	for i := 0; i < 1000; i++ {
		content.SetRawString(coreconst.DataxJobContent+fmt.Sprintf(".%d", i), fmt.Sprintf(`{
			"taskId": %d,
			"reader":{
				"name":"mock"
			},
			"writer":{
				"name":"mock"
			}
		}`, i))
	}
	ctx, cancel := context.WithCancel(context.Background())
	c, _ := NewContainer(ctx, content)
	go func() {
		time.Sleep(4 * time.Second)
		cancel()
	}()

	if err := c.Do(); err != context.Canceled {
		t.Errorf("Do error: %v", err)
	}
}

func TestContainer_JobId(t *testing.T) {
	tests := []struct {
		name string
		c    *Container
		want int64
	}{
		{
			name: "1",
			c: testContainer(context.Background(), testJsonFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 30000000,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 1,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				}
			}`)),
			want: 30000000,
		},
		{
			name: "2",
			c: testContainer(context.Background(), testJsonFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1000000000000000000,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 1,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				}
			}`)),
			want: 1000000000000000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.JobId(); got != tt.want {
				t.Errorf("Container.JobId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_TaskGroupId(t *testing.T) {
	tests := []struct {
		name string
		c    *Container
		want int64
	}{

		{
			name: "1",
			c: testContainer(context.Background(), testJsonFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 30000000,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					}
				}`)),
			want: 30000001,
		},
		{
			name: "2",
			c: testContainer(context.Background(), testJsonFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 1000000000000000000,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 1000000000000000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					}
				}`)),
			want: 1000000000000000001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.TaskGroupId(); got != tt.want {
				t.Errorf("Container.TaskGroupId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_Start(t *testing.T) {
	tests := []struct {
		name    string
		c       *Container
		wantErr bool
	}{
		{
			name: "1",
			c: testContainer(context.Background(), testJsonFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 30000000,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					}
				}`)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Container.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewContainer(t *testing.T) {
	type args struct {
		ctx  context.Context
		conf *config.JSON
	}
	tests := []struct {
		name    string
		args    args
		wantC   *Container
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ctx: context.TODO(),
				conf: testJsonFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": "30000000",
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					}
				}`),
			},
			wantC:   nil,
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				ctx: context.TODO(),
				conf: testJsonFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 30000002,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": "30000001",
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					}
				}`),
			},
			wantC:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewContainer(tt.args.ctx, tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewContainer() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestContainer_startTaskExecer(t *testing.T) {
	resetLoader()
	loader.RegisterReader("mock", newMockRandReader([]error{
		nil, nil, errors.New("mock test error"), nil, nil,
	}))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}))

	c := testContainer(context.Background(), testJsonFromString(`{
		"core" : {
			"container": {
				"job":{
					"id": 30000000,
					"sleepInterval":100
				},
				"taskGroup":{
					"id": 30000001,
					"failover":{
						"retryIntervalInMsec":0
					}
				}
			}
		}
	}`))
	c.scheduler = schedule.NewTaskSchduler(4, 0)
	c.scheduler.Stop()
	te, _ := newTaskExecer(c.ctx, testJsonFromString(`{
		"taskId": 1,
		"reader":{
			"name":"mock"
		},
		"writer":{
			"name":"mock"
		}
	}`), "mock", 0)
	if err := c.startTaskExecer(te); err != nil {
		t.Errorf("Container.startTaskExecer() error = %v, wantErr true", err)
	}
}
