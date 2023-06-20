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

package taskgroup

import (
	"context"
	"errors"
	"fmt"
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
	content := testJSONFromString(`{
		"core" : {
			"container": {
				"job":{
					"id": 1
				},
				"taskGroup":{
					"id": 1,
					"reportInterval":1
				},
				"task":{
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
				"name":"mock",
				"parameter":{}
			},
			"writer":{
				"name":"mock",
				"parameter":{}
			}
		}`, i))
	}

	resetLoader()
	loader.RegisterReader("mock", newMockRandReader([]error{
		nil, nil, nil, nil, nil,
	}))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}))
	c, _ := NewContainer(context.TODO(), content)
	if err := c.Do(); err != nil {
		t.Errorf("Do error: %v", err)
	}

	resetLoader()
	loader.RegisterReader("mock", newMockRandReader([]error{
		nil, nil, errors.New("mock test error"), nil, nil,
	}))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}))
	c, _ = NewContainer(context.TODO(), content)
	if err := c.Do(); err == nil {
		t.Errorf("Do error: %v", err)
	}

	resetLoader()
	loader.RegisterReader("mock", newMockRandReader([]error{
		nil, nil, nil, nil, nil,
	}))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, errors.New("mock test error"), nil, nil,
	}))
	c, _ = NewContainer(context.TODO(), content)
	if err := c.Do(); err == nil {
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
	content := testJSONFromString(`{
		"core" : {
			"container": {
				"job":{
					"id": 1
				},
				"taskGroup":{
					"id": 1
				},
				"task":{
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
				"name":"mock",
				"parameter":{}
			},
			"writer":{
				"name":"mock",
				"parameter":{}
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
	content := testJSONFromString(`{
		"core" : {
			"container": {
				"job":{
					"id": 1
				},
				"taskGroup":{
					"id": 1,
					"reportInterval":1
				},
				"task":{
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
				"name":"mock",
				"parameter":{}
			},
			"writer":{
				"name":"mock",
				"parameter":{}
			}
		}`, i))
	}
	ctx, cancel := context.WithCancel(context.Background())
	c, _ := NewContainer(ctx, content)
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	if err := c.Do(); err == nil {
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
			c: testContainer(context.Background(), testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 30000000,
							"reportInterval":1
						},
						"taskGroup":{
							"id": 1
						},
						"task":{
							"failover":{
								"retryIntervalInMsec":10
							}
						}
					}
				}
			}`)),
			want: 30000000,
		},
		{
			name: "2",
			c: testContainer(context.Background(), testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1000000000000000000,
							"reportInterval":1
						},
						"taskGroup":{
							"id": 1
						},
						"task":{
							"failover":{
								"retryIntervalInMsec":10
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
			if got := tt.c.JobID(); got != tt.want {
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
			c: testContainer(context.Background(), testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 30000000,
								"reportInterval":1
							},
							"taskGroup":{
								"id": 30000001
							},
							"task":{
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
			c: testContainer(context.Background(), testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 1000000000000000000,
								"reportInterval":1
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
			if got := tt.c.TaskGroupID(); got != tt.want {
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
			c: testContainer(context.Background(), testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 30000000,
								"reportInterval":1
							},
							"taskGroup":{
								"id": 30000001
							},
							"task":{
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
				conf: testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": "30000000",
								"reportInterval":1
							},
							"taskGroup":{
								"id": 30000001
							},
							"task":{
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
				conf: testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 30000002,
								"reportInterval":1
							},
							"taskGroup":{
								"id": "30000001"
							},
							"task":{
								"failover":{
									"retryIntervalInMsec":10
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
			if gotC != tt.wantC {
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

	c := testContainer(context.Background(), testJSONFromString(`{
		"core" : {
			"container": {
				"job":{
					"id": 30000000,
					"reportInterval":1
				},
				"taskGroup":{
					"id": 30000001
				},
				"task":{
					"failover":{
						"retryIntervalInMsec":10
					}
				}
			}
		}
	}`))
	c.scheduler = schedule.NewTaskSchduler(4, 0)
	c.scheduler.Stop()
	te, err := newTaskExecer(c.ctx, testJSONFromString(`{
		"taskId": 1,
		"reader":{
			"name":"mock",
			"parameter":{}
		},
		"writer":{
			"name":"mock",
			"parameter":{}
		}
	}`), 3, 3, 0)

	if err != nil {
		t.Fatal(err)
	}
	if err := c.startTaskExecer(te); err == nil {
		t.Errorf("Container.startTaskExecer() error = %v, wantErr true", err)
	}
}
