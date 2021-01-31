package job

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
)

func TestNewContainer(t *testing.T) {
	type args struct {
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
				conf: testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": -3,
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
				conf: testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": "1000",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewContainer(context.TODO(), tt.args.conf)
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

func TestContainer_preHandle(t *testing.T) {
	tests := []struct {
		name    string
		c       *Container
		wantErr bool
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job" :{
					"preHandler": {
						"pluginType": "handler",
						"pluginName": "mockHandler"
					}
				}
			}`)),
		},
		{
			name: "2",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job" :{
					"preHandler": null
				}
			}`)),
		},
		{
			name: "3",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job" :{
					"preHandler":  {
						"pluginType": "test",
						"pluginName": "mockHandler"
					}
				}
			}`)),
			wantErr: true,
		},
		{
			name: "4",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job" :{
					"preHandler":  {
						"pluginType": 1,
						"pluginName": "mockHandler"
					}
				}
			}`)),
			wantErr: true,
		},
		{
			name: "4",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job" :{
					"preHandler":  {
						"pluginType": "handler",
						"pluginName": 1
					}
				}
			}`)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.preHandle(); (err != nil) != tt.wantErr {
				t.Errorf("Container.preHandle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainer_postHandle(t *testing.T) {
	tests := []struct {
		name    string
		c       *Container
		wantErr bool
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job" :{
					"postHandler": {
						"pluginType": "handler",
						"pluginName": "mockHandler"
					}
				}
			}`)),
		},
		{
			name: "2",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job" :{
					"postHandler": null
				}
			}`)),
		},
		{
			name: "3",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job" :{
					"postHandler":  {
						"pluginType": "test",
						"pluginName": "mockHandler"
					}
				}
			}`)),
			wantErr: true,
		},
		{
			name: "4",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job" :{
					"postHandler":  {
						"pluginType": 1,
						"pluginName": "mockHandler"
					}
				}
			}`)),
			wantErr: true,
		},
		{
			name: "4",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job" :{
					"postHandler":  {
						"pluginType": "handler",
						"pluginName": 1
					}
				}
			}`)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.postHandle(); (err != nil) != tt.wantErr {
				t.Errorf("Container.postHandle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainer_init(t *testing.T) {
	resetLoader()
	loader.RegisterReader("mock", newMockReader([]error{
		nil, nil, nil, nil, nil,
	}, nil))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}, nil))
	loader.RegisterReader("mockErr", newMockReader([]error{
		errors.New("mock test error"), nil, nil, nil, nil,
	}, nil))
	loader.RegisterWriter("mockErr", newMockWriter([]error{
		errors.New("mock test error"), nil, nil, nil, nil,
	}, nil))

	tests := []struct {
		name    string
		c       *Container
		wantErr bool
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 1,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					},
					"job":{
						"content":[
							{
								"reader":{
									"name": "mock",
									"parameter" : {

									}
								},
								"writer":{
									"name": "mock",
									"parameter" : {

									}
								}
							}
						]
					}
				}`)),
		},
		{
			name: "2",
			c: testContainer(testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 1,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					},
					"job":{
						"content":[
							{
								"reader":{	
								}
							}
						]
					}
				}`)),
			wantErr: true,
		},
		{
			name: "3",
			c: testContainer(testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 1,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					},
					"job":{
						"content":[
							{
								"reader":{
									"name": "mock"
								},
								"writer":{
									"name": "mock"
								}
							}
						]
					}
				}`)),
			wantErr: true,
		},
		{
			name: "4",
			c: testContainer(testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 1,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					},
					"job":{
						"content":[
							{
								"reader":{
									"name": "mock",
									"parameter" : {

									}
								},
								"writer":{
								}
							}
						]
					}
				}`)),
			wantErr: true,
		},
		{
			name: "5",
			c: testContainer(testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 1,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					},
					"job":{
						"content":[
							{
								"reader":{
									"name": "mock",
									"parameter" : {

									}
								},
								"writer":{
									"name": "mock"
								}
							}
						]
					}
				}`)),
			wantErr: true,
		},
		{
			name: "6",
			c: testContainer(testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 1,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					},
					"job":{
						"content":[
							{
								"reader":{
									"name": "mock1",
									"parameter" : {

									}
								},
								"writer":{
									"name": "mock",
									"parameter" : {

									}
								}
							}
						]
					}
				}`)),
			wantErr: true,
		},
		{
			name: "7",
			c: testContainer(testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 1,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					},
					"job":{
						"content":[
							{
								"reader":{
									"name": "mockErr",
									"parameter" : {

									}
								},
								"writer":{
									"name": "mock",
									"parameter" : {

									}
								}
							}
						]
					}
				}`)),
			wantErr: true,
		},
		{
			name: "8",
			c: testContainer(testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 1,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					},
					"job":{
						"content":[
							{
								"reader":{
									"name": "mock",
									"parameter" : {

									}
								},
								"writer":{
									"name": "mock1",
									"parameter" : {

									}
								}
							}
						]
					}
				}`)),
			wantErr: true,
		},
		{
			name: "9",
			c: testContainer(testJSONFromString(`{
					"core" : {
						"container": {
							"job":{
								"id": 1,
								"sleepInterval":100
							},
							"taskGroup":{
								"id": 30000001,
								"failover":{
									"retryIntervalInMsec":0
								}
							}
						}
					},
					"job":{
						"content":[
							{
								"reader":{
									"name": "mock",
									"parameter" : {

									}
								},
								"writer":{
									"name": "mockErr",
									"parameter" : {

									}
								}
							}
						]
					}
				}`)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.init()
			if (err != nil) != tt.wantErr {
				t.Errorf("Container.init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainer_prepare(t *testing.T) {
	resetLoader()
	loader.RegisterReader("mock", newMockReader([]error{
		nil, nil, nil, nil, nil,
	}, nil))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}, nil))
	loader.RegisterReader("mockErr", newMockReader([]error{
		nil, errors.New("mock test error"), nil, nil, nil,
	}, nil))
	loader.RegisterWriter("mockErr", newMockWriter([]error{
		nil, errors.New("mock test error"), nil, nil, nil,
	}, nil))

	tests := []struct {
		name    string
		c       *Container
		wantErr bool
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mock",
								"parameter" : {

								}
							}
						}
					]
				}
			}`)),
		},
		{
			name: "2",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mockErr",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mock",
								"parameter" : {

								}
							}
						}
					]
				}
			}`)),
			wantErr: true,
		},
		{
			name: "3",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mockErr",
								"parameter" : {

								}
							}
						}
					]
				}
			}`)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.init()
			if err := tt.c.prepare(); (err != nil) != tt.wantErr {
				t.Errorf("Container.prepare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainer_post(t *testing.T) {
	resetLoader()
	loader.RegisterReader("mock", newMockReader([]error{
		nil, nil, nil, nil, nil,
	}, nil))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}, nil))
	loader.RegisterReader("mockErr", newMockReader([]error{
		nil, nil, nil, errors.New("mock test error"), nil,
	}, nil))
	loader.RegisterWriter("mockErr", newMockWriter([]error{
		nil, nil, nil, errors.New("mock test error"), nil,
	}, nil))

	tests := []struct {
		name    string
		c       *Container
		wantErr bool
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mock",
								"parameter" : {

								}
							}
						}
					]
				}
			}`)),
		},
		{
			name: "2",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mockErr",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mock",
								"parameter" : {

								}
							}
						}
					]
				}
			}`)),
			wantErr: true,
		},
		{
			name: "3",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mockErr",
								"parameter" : {

								}
							}
						}
					]
				}
			}`)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.init()
			if err := tt.c.post(); (err != nil) != tt.wantErr {
				t.Errorf("Container.post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainer_destroy(t *testing.T) {
	resetLoader()
	loader.RegisterReader("mock", newMockReader([]error{
		nil, nil, nil, nil, nil,
	}, nil))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}, nil))
	loader.RegisterReader("mockErr", newMockReader([]error{
		nil, nil, nil, nil, errors.New("mock test error"),
	}, nil))
	loader.RegisterWriter("mockErr", newMockWriter([]error{
		nil, nil, nil, nil, errors.New("mock test error"),
	}, nil))
	tests := []struct {
		name    string
		c       *Container
		wantErr bool
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mock",
								"parameter" : {

								}
							}
						}
					]
				}
			}`)),
		},
		{
			name: "2",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mockErr",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mock",
								"parameter" : {

								}
							}
						}
					]
				}
			}`)),
			wantErr: true,
		},
		{
			name: "3",
			c: testContainer(testJSONFromString(`{
				"core" : {
					"container": {
						"job":{
							"id": 1,
							"sleepInterval":100
						},
						"taskGroup":{
							"id": 30000001,
							"failover":{
								"retryIntervalInMsec":0
							}
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mockErr",
								"parameter" : {

								}
							}
						}
					]
				}
			}`)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.init()
			if err := tt.c.destroy(); (err != nil) != tt.wantErr {
				t.Errorf("Container.destroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainer_adjustChannelNumber(t *testing.T) {
	tests := []struct {
		name                  string
		c                     *Container
		wantErr               bool
		wantNeedChannelNumber int64
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 100,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":3000,
							"record":400,
							"channel":4
						}
					}
				}
			}`)),
			wantErr:               false,
			wantNeedChannelNumber: 4,
		},
		{
			name: "2",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 3000,
								"record":400
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":1,
							"record":1,
							"channel":4
						}
					}
				}
			}`)),
			wantErr:               false,
			wantNeedChannelNumber: 1,
		},
		{
			name: "3",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 60,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":400,
							"record":3000,
							"channel":4
						}
					}
				}
			}`)),
			wantErr:               false,
			wantNeedChannelNumber: 6,
		},
		{
			name: "4",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 60,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":0,
							"record":0,
							"channel":4
						}
					}
				}
			}`)),
			wantErr:               false,
			wantNeedChannelNumber: 4,
		},
		{
			name: "5",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 60,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":0,
							"record":0,
							"channel":0
						}
					}
				}
			}`)),
			wantErr:               true,
			wantNeedChannelNumber: math.MaxInt32,
		},
		{
			name: "6",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":1,
							"record":1,
							"channel":0
						}
					}
				}
			}`)),
			wantErr: true,
		},
		{
			name: "7",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": -1,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":1,
							"record":1,
							"channel":0
						}
					}
				}
			}`)),
			wantErr: true,
		},
		{
			name: "8",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 1
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":1,
							"record":1,
							"channel":0
						}
					}
				}
			}`)),
			wantErr: true,
		},
		{
			name: "9",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 1,
								"record":-1
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":1,
							"record":1,
							"channel":0
						}
					}
				}
			}`)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.adjustChannelNumber()
			if (err != nil) != tt.wantErr {
				t.Errorf("Container.adjustChannelNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNeedChannelNumber != tt.c.needChannelNumber {
				t.Errorf("Container.needChannelNumber = %v, wantNeedChannelNumber %v",
					tt.c.needChannelNumber, tt.wantNeedChannelNumber)
			}
		})
	}
}

func TestContainer_mergeTaskConfigs(t *testing.T) {

	type args struct {
		readerConfs []*config.JSON
		writerConfs []*config.JSON
	}
	tests := []struct {
		name            string
		c               *Container
		args            args
		wantTaskConfigs []*config.JSON
		wantErr         bool
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mockErr",
								"parameter" : {

								}
							},
							"transformer" : ["1","2"]
						}
					]
				}
			}`)),
			args: args{
				readerConfs: []*config.JSON{
					testJSONFromString(`{"id":1}`),
					testJSONFromString(`{"id":2}`),
					testJSONFromString(`{"id":3}`),
				},
				writerConfs: []*config.JSON{
					testJSONFromString(`{"id":4}`),
					testJSONFromString(`{"id":5}`),
					testJSONFromString(`{"id":6}`),
				},
			},
			wantTaskConfigs: []*config.JSON{
				testJSONFromString(`{
					"taskId":0,
					"reader":{
						"name" : "mock",
						"parameter" : 
							{
								"id":1
							}
						
					},
					"transformer" :["1","2"],
					"writer":{
						"name" : "mockErr",
						"parameter" : {
								"id":4
						}
					}
				}`),
				testJSONFromString(`{
					"taskId":1,
					"reader":{
					"name" : "mock",
					"parameter" : {
							"id":2
					}
				},
				"transformer" :["1","2"],
				"writer":{
					"name" : "mockErr",
					"parameter" : {
							"id":5
					}
				}
			}`),
				testJSONFromString(`{
				"taskId":2,
				"reader":{
				"name" : "mock",
				"parameter" : 
					{
						"id":3
					}
				
			},
			"transformer" :["1","2"],
			"writer":{
				"name" : "mockErr",
				"parameter" : {
						"id":6
				}
			}
		}`),
			},
		},
		{
			name: "2",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mockErr",
								"parameter" : {

								}
							},
							"transformer" : ["1","2"]
						}
					]
				}
			}`)),
			args: args{
				readerConfs: []*config.JSON{
					testJSONFromString(`{"id":1}`),
					testJSONFromString(`{"id":2}`),
					testJSONFromString(`{"id":3}`),
				},
				writerConfs: []*config.JSON{
					testJSONFromString(`{"id":4}`),
					testJSONFromString(`{"id":5}`),
					testJSONFromString(`{"id":6}`),
					testJSONFromString(`{"id":7}`),
				},
			},
			wantTaskConfigs: nil,
			wantErr:         true,
		},
		{
			name: "3",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					}
				},
				"job":{
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mockErr",
								"parameter" : {

								}
							}
						}
					]
				}
			}`)),
			args: args{
				readerConfs: []*config.JSON{
					testJSONFromString(`{"id":1}`),
					testJSONFromString(`{"id":2}`),
					testJSONFromString(`{"id":3}`),
				},
				writerConfs: []*config.JSON{
					testJSONFromString(`{"id":4}`),
					testJSONFromString(`{"id":5}`),
					testJSONFromString(`{"id":6}`),
				},
			},
			wantTaskConfigs: nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.init(); err != nil {
				t.Fatalf("Container.init() error = %v", err)
			}
			gotTaskConfigs, err := tt.c.mergeTaskConfigs(tt.args.readerConfs, tt.args.writerConfs)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Container.mergeTaskConfigs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(gotTaskConfigs) != len(tt.wantTaskConfigs) {
				t.Fatalf("Container.mergeTaskConfigs() len = %v, wantTaskConfigs len = %v",
					len(gotTaskConfigs), len(tt.wantTaskConfigs))
			}

			for i := range gotTaskConfigs {
				if !equalConfigJSON(gotTaskConfigs[i], tt.wantTaskConfigs[i]) {
					t.Fatalf("Container.mergeTaskConfigs()[%v]  = %v, wantTaskConfigs[%v]] = %v",
						i, gotTaskConfigs[i], i, tt.wantTaskConfigs[i])
				}
			}

		})
	}
}

func TestContainer_split(t *testing.T) {
	resetLoader()
	loader.RegisterReader("mock", newMockReader([]error{
		nil, nil, nil, nil, nil,
	}, []*config.JSON{
		testJSONFromString(`{"id":1}`),
		testJSONFromString(`{"id":2}`),
		testJSONFromString(`{"id":3}`),
	}))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}, []*config.JSON{
		testJSONFromString(`{"id":4}`),
		testJSONFromString(`{"id":5}`),
		testJSONFromString(`{"id":6}`),
	}))
	loader.RegisterWriter("mockx", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}, []*config.JSON{
		testJSONFromString(`{"id":4}`),
		testJSONFromString(`{"id":5}`),
	}))
	loader.RegisterReader("mock1", newMockReader([]error{
		nil, nil, nil, nil, nil,
	}, nil))
	loader.RegisterWriter("mock1", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}, nil))
	loader.RegisterReader("mockErr", newMockReader([]error{
		nil, nil, errors.New("mock test error"), nil, nil,
	}, nil))
	loader.RegisterWriter("mockErr", newMockWriter([]error{
		nil, nil, errors.New("mock test error"), nil, nil,
	}, nil))
	tests := []struct {
		name       string
		c          *Container
		wantErr    bool
		wantConfig *config.JSON
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 100,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":400,
							"record":3000,
							"channel":4
						}
					},
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mock",
								"parameter" : {

								}
							},
							"transformer" : ["1","2"]
						}
					]
				}
			}`)),
			wantErr: false,
			wantConfig: testJSONFromString(`{
		"content":[
			{
				"taskId":0,
				"reader":{
					"name" : "mock",
					"parameter" : {
						"id":1
					}
				},
				"transformer" :["1","2"],
				"writer":{
					"name" : "mock",
					"parameter" : {
						"id":4
					}
				}
			},
			{
				"taskId":1,
				"reader":{
					"name" : "mock",
					"parameter" : {
						"id":2
					}
				},
				"transformer" :["1","2"],
				"writer":{
					"name" : "mock",
					"parameter" : {
						"id":5
					}
				}
			},
			{
				"taskId":2,
				"reader":{
					"name" : "mock",
					"parameter" : {
						"id":3
					}
				},
				"transformer" :["1","2"],
				"writer":{
					"name" : "mock",
					"parameter" : {
						"id":6
					}
				}
			}]
			}`),
		},
		{
			name: "2",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": -100,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":400,
							"record":3000,
							"channel":4
						}
					},
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mock",
								"parameter" : {

								}
							},
							"transformer" : ["1","2"]
						}
					]
				}
			}`)),
			wantErr: true,
			wantConfig: testJSONFromString(`{
				"content":[
					{
						"reader":{
							"name": "mock",
							"parameter" : {

							}
						},
						"writer":{
							"name": "mock",
							"parameter" : {

							}
						},
						"transformer" : ["1","2"]
					}
				]
			}`),
		},
		{
			name: "3",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 100,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":400,
							"record":3000,
							"channel":4
						}
					},
					"content":[
						{
							"reader":{
								"name": "mockErr",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mock",
								"parameter" : {

								}
							},
							"transformer" : ["1","2"]
						}
					]
				}
			}`)),
			wantErr: true,
			wantConfig: testJSONFromString(`{
				"content":[
					{
						"reader":{
							"name": "mockErr",
							"parameter" : {

							}
						},
						"writer":{
							"name": "mock",
							"parameter" : {

							}
						},
						"transformer" : ["1","2"]
					}
				]
			}`),
		},
		{
			name: "4",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 100,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":400,
							"record":3000,
							"channel":4
						}
					},
					"content":[
						{
							"reader":{
								"name": "mock1",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mock",
								"parameter" : {

								}
							},
							"transformer" : ["1","2"]
						}
					]
				}
			}`)),
			wantErr: true,
			wantConfig: testJSONFromString(`{
				"content":[
					{
						"reader":{
							"name": "mock1",
							"parameter" : {

							}
						},
						"writer":{
							"name": "mock",
							"parameter" : {

							}
						},
						"transformer" : ["1","2"]
					}
				]
			}`),
		},
		{
			name: "5",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 100,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":400,
							"record":3000,
							"channel":4
						}
					},
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mock1",
								"parameter" : {

								}
							},
							"transformer" : ["1","2"]
						}
					]
				}
			}`)),
			wantErr: true,
			wantConfig: testJSONFromString(`{
				"content":[
					{
						"reader":{
							"name": "mock",
							"parameter" : {

							}
						},
						"writer":{
							"name": "mock1",
							"parameter" : {

							}
						},
						"transformer" : ["1","2"]
					}
				]
			}`),
		},
		{
			name: "6",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 100,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":400,
							"record":3000,
							"channel":4
						}
					},
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mockErr",
								"parameter" : {

								}
							},
							"transformer" : ["1","2"]
						}
					]
				}
			}`)),
			wantErr: true,
			wantConfig: testJSONFromString(`{
				"content":[
					{
						"reader":{
							"name": "mock",
							"parameter" : {

							}
						},
						"writer":{
							"name": "mockErr",
							"parameter" : {

							}
						},
						"transformer" : ["1","2"]
					}
				]
			}`),
		},
		{
			name: "7",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 100,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":400,
							"record":3000,
							"channel":4
						}
					},
					"content":[
						{
							"reader":{
								"name": "mock",
								"parameter" : {

								}
							},
							"writer":{
								"name": "mockx",
								"parameter" : {

								}
							},
							"transformer" : ["1","2"]
						}
					]
				}
			}`)),
			wantErr: true,
			wantConfig: testJSONFromString(`{
				"content":[
					{
						"reader":{
							"name": "mock",
							"parameter" : {

							}
						},
						"writer":{
							"name": "mockx",
							"parameter" : {

							}
						},
						"transformer" : ["1","2"]
					}
				]
			}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.init(); err != nil {
				t.Errorf("Container.init() error = %v", err)
				return
			}
			if err := tt.c.split(); (err != nil) != tt.wantErr {
				t.Errorf("Container.split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, _ := tt.c.Config().GetConfig(coreconst.DataxJobContent)
			want, _ := tt.wantConfig.GetConfig("content")

			if !equalConfigJSON(got, want) {
				t.Errorf("got: %v want: %v", got, want)
			}
		})
	}
}

func TestContainer_schedule(t *testing.T) {
	tests := []struct {
		name              string
		c                 *Container
		wantErr           bool
		needChannelNumber int64
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						}
					},
					"transport":{
						"channel":{
							"speed":{
								"byte": 100,
								"record":100
							}
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"byte":400,
							"record":3000,
							"channel":4
						}
					}
				}
			}`)),
			wantErr: true,
		},
		{
			name: "2",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						},
						"taskGroup":{
							"channel":2
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"channel":8
						}
					},
					"content":[
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "a"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "A"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "b"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "B"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "c"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "C"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"b",
									"id" : "d"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "D"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"b",
									"id" : "e"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "E"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"c",
									"id" : "f"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "F"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"c",
									"id" : "g"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "G"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"c",
									"id" : "h"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "H"
								}
							}
						}
					]
				}	
			}`)),
			wantErr:           false,
			needChannelNumber: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.needChannelNumber = tt.needChannelNumber
			if err := tt.c.schedule(); (err != nil) != tt.wantErr {
				t.Errorf("Container.schedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_doAssign(t *testing.T) {
	type args struct {
		taskIDMap       map[string][]int
		taskGroupNumber int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "1",
			args: args{
				taskIDMap: map[string][]int{
					"a": {0, 1, 2},
					"b": {3, 4},
					"c": {5, 6, 7},
				},
				taskGroupNumber: 4,
			},
			want: [][]int{
				{0, 4},
				{3, 6},
				{5, 2},
				{1, 7},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := doAssign(tt.args.taskIDMap, tt.args.taskGroupNumber); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doAssign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseAndGetResourceMarkAndTaskIdMap(t *testing.T) {
	type args struct {
		tasksConfigs []*config.JSON
	}
	tests := []struct {
		name string
		args args
		want map[string][]int
	}{
		{
			name: "1",
			args: args{
				tasksConfigs: []*config.JSON{
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "a"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "A"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "b"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "B"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "c"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "C"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"b",
								"id" : "d"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "D"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"b",
								"id" : "e"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "E"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"c",
								"id" : "f"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "F"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"c",
								"id" : "g"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "G"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"c",
								"id" : "h"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "H"
							}
						}
					}`),
				},
			},
			want: map[string][]int{
				"a": {0, 1, 2},
				"b": {3, 4},
				"c": {5, 6, 7},
			},
		},
		{
			name: "2",
			args: args{
				tasksConfigs: []*config.JSON{
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "a"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "A"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "b"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "B"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "c"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "C"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "d"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "D"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "e"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "E"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "f"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "F"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "g"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"b",
								"id" : "G"
							}
						}
					}`),
					testJSONFromString(`{
						"reader":{
							"parameter":{
								"loadBalanceResourceMark":"a",
								"id" : "h"
							}
						},
						"writer":{
							"parameter":{
								"loadBalanceResourceMark":"b",
								"id" : "H"
							}
						}
					}`),
				},
			},
			want: map[string][]int{
				"a": {0, 1, 2, 3, 4, 5},
				"b": {6, 7},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseAndGetResourceMarkAndTaskIDMap(tt.args.tasksConfigs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAndGetResourceMarkAndTaskIdMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_distributeTaskIntoTaskGroup(t *testing.T) {
	tests := []struct {
		name              string
		c                 *Container
		wantConfs         []*config.JSON
		needChannelNumber int64
		wantErr           bool
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						},
						"taskGroup":{
							"channel":2
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"channel":4
						}
					},
					"content":[
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "a"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "A"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "b"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "B"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "c"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "C"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"b",
									"id" : "d"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "D"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"b",
									"id" : "e"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "E"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"c",
									"id" : "f"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "F"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"c",
									"id" : "g"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "G"
								}
							}
						},
						{
							"reader":{
								"parameter":{
									"loadBalanceResourceMark":"c",
									"id" : "h"
								}
							},
							"writer":{
								"parameter":{
									"loadBalanceResourceMark":"a",
									"id" : "H"
								}
							}
						}
					]
				}	
			}`)),
			needChannelNumber: 8,
			wantConfs: []*config.JSON{
				testJSONFromString(`{
					"core":{
						"container": {
							"job":{
								"id": 1
							},
							"taskGroup":{
								"id": 0,
								"channel":2
							}
						}
					},
					"job":{
						"setting":{
							"speed":{
								"channel":4
							}
						},
						"content":[
							{
								"reader":{
									"parameter":{
										"loadBalanceResourceMark":"a",
										"id" : "a"
									}
								},
								"writer":{
									"parameter":{
										"loadBalanceResourceMark":"a",
										"id" : "A"
									}
								}
							},
							{
								"reader":{
									"parameter":{
										"loadBalanceResourceMark":"b",
										"id" : "e"
									}
								},
								"writer":{
									"parameter":{
										"loadBalanceResourceMark":"a",
										"id" : "E"
									}
								}
							}
						]
					}	
				}`),
				testJSONFromString(`{
					"core":{
						"container": {
							"job":{
								"id": 1
							},
							"taskGroup":{
								"id": 1,
								"channel":2
							}
						}
					},
					"job":{
						"setting":{
							"speed":{
								"channel":4
							}
						},
						"content":[
							{
								"reader":{
									"parameter":{
										"loadBalanceResourceMark":"b",
										"id" : "d"
									}
								},
								"writer":{
									"parameter":{
										"loadBalanceResourceMark":"a",
										"id" : "D"
									}
								}
							},
							{
								"reader":{
									"parameter":{
										"loadBalanceResourceMark":"c",
										"id" : "g"
									}
								},
								"writer":{
									"parameter":{
										"loadBalanceResourceMark":"a",
										"id" : "G"
									}
								}
							}
						]
					}	
				}`),
				testJSONFromString(`{
					"core":{
						"container": {
							"job":{
								"id": 1
							},
							"taskGroup":{
								"id": 2,
								"channel":2
							}
						}
					},
					"job":{
						"setting":{
							"speed":{
								"channel":4
							}
						},
						"content":[
							{
								"reader":{
									"parameter":{
										"loadBalanceResourceMark":"c",
										"id" : "f"
									}
								},
								"writer":{
									"parameter":{
										"loadBalanceResourceMark":"a",
										"id" : "F"
									}
								}
							},
							{
								"reader":{
									"parameter":{
										"loadBalanceResourceMark":"a",
										"id" : "c"
									}
								},
								"writer":{
									"parameter":{
										"loadBalanceResourceMark":"a",
										"id" : "C"
									}
								}
							}
						]
					}	
				}`),
				testJSONFromString(`{
					"core":{
						"container": {
							"job":{
								"id": 1
							},
							"taskGroup":{
								"id": 3,
								"channel":2
							}
						}
					},
					"job":{
						"setting":{
							"speed":{
								"channel":4
							}
						},
						"content":[
							{
								"reader":{
									"parameter":{
										"loadBalanceResourceMark":"a",
										"id" : "b"
									}
								},
								"writer":{
									"parameter":{
										"loadBalanceResourceMark":"a",
										"id" : "B"
									}
								}
							},
							{
								"reader":{
									"parameter":{
										"loadBalanceResourceMark":"c",
										"id" : "h"
									}
								},
								"writer":{
									"parameter":{
										"loadBalanceResourceMark":"a",
										"id" : "H"
									}
								}
							}
						]
					}	
				}`),
			},
		},

		{
			name: "2",
			c: testContainer(testJSONFromString(`{
				"core":{
					"container": {
						"job":{
							"id": 1
						},
						"taskGroup":{
							"channel":2
						}
					}
				},
				"job":{
					"setting":{
						"speed":{
							"channel":4
						}
					}
				}	
			}`)),
			needChannelNumber: 8,
			wantConfs:         nil,
			wantErr:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.needChannelNumber = tt.needChannelNumber
			gotConfs, err := tt.c.distributeTaskIntoTaskGroup()
			if (err != nil) != tt.wantErr {
				t.Errorf("Container.distributeTaskIntoTaskGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotConfs) != len(tt.wantConfs) {
				t.Fatalf("Container.distributeTaskIntoTaskGroup() len = %v, wantTaskConfigs len = %v",
					len(gotConfs), len(tt.wantConfs))
			}

			for i := range gotConfs {
				if !equalConfigJSON(gotConfs[i], tt.wantConfs[i]) {
					t.Fatalf("Container.distributeTaskIntoTaskGroup()[%v]  = %v, wantTaskConfigs[%v]] = %v",
						i, gotConfs[i], i, tt.wantConfs[i])
				}
			}
		})
	}
}

func TestContainer_Start(t *testing.T) {

	resetLoader()
	loader.RegisterReader("mock", newMockReader([]error{
		nil, nil, nil, nil, nil,
	}, []*config.JSON{
		testJSONFromString(`{"id":1}`),
		testJSONFromString(`{"id":2}`),
		testJSONFromString(`{"id":3}`),
	}))

	loader.RegisterReader("mock0", newMockReader([]error{
		errors.New("mock test error"), nil, nil, nil, nil,
	}, []*config.JSON{
		testJSONFromString(`{"id":1}`),
		testJSONFromString(`{"id":2}`),
		testJSONFromString(`{"id":3}`),
	}))

	loader.RegisterReader("mock1", newMockReader([]error{
		nil, errors.New("mock test error"), nil, nil, nil,
	}, []*config.JSON{
		testJSONFromString(`{"id":1}`),
		testJSONFromString(`{"id":2}`),
		testJSONFromString(`{"id":3}`),
	}))

	loader.RegisterReader("mock2", newMockReader([]error{
		nil, nil, errors.New("mock test error"), nil, nil,
	}, []*config.JSON{
		testJSONFromString(`{"id":1}`),
		testJSONFromString(`{"id":2}`),
		testJSONFromString(`{"id":3}`),
	}))
	loader.RegisterReader("mock3", newMockReader([]error{
		nil, nil, nil, errors.New("mock test error"), nil,
	}, []*config.JSON{
		testJSONFromString(`{"id":1}`),
		testJSONFromString(`{"id":2}`),
		testJSONFromString(`{"id":3}`),
	}))
	loader.RegisterReader("mock4", newMockReader([]error{
		nil, nil, nil, nil, errors.New("mock test error"),
	}, []*config.JSON{
		testJSONFromString(`{"id":1}`),
		testJSONFromString(`{"id":2}`),
		testJSONFromString(`{"id":3}`),
	}))
	loader.RegisterWriter("mock", newMockWriter([]error{
		nil, nil, nil, nil, nil,
	}, []*config.JSON{
		testJSONFromString(`{"id":4}`),
		testJSONFromString(`{"id":5}`),
		testJSONFromString(`{"id":6}`),
	}))
	tests := []struct {
		name    string
		c       *Container
		wantErr bool
	}{
		{
			name: "1",
			c: testContainer(testJSONFromString(`{
				"core": {
					"container": {
						"job": {
							"id": 1
						},
						"taskGroup": {
							"channel": 2
						}
					}
				},
				"job": {
					"preHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"postHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"setting": {
						"speed": {
							"channel": 4
						}
					},
					"content": [{
						"reader": {
							"name": "mock",
							"parameter": {
			
							}
						},
						"writer": {
							"name": "mock",
							"parameter": {}
						},
						"transformer": ["1", "2"]
					}]
				}
			}`)),
		},

		{
			name: "2",
			c: testContainer(testJSONFromString(`{
				"core": {
					"container": {
						"job": {
							"id": 1
						},
						"taskGroup": {
							"channel": 2
						}
					}
				},
				"job": {
					"preHandler": {
						"pluginType": "handler",
						"pluginName": 1
					},
					"postHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"setting": {
						"speed": {
							"channel": 4
						}
					},
					"content": [{
						"reader": {
							"name": "mock",
							"parameter": {
			
							}
						},
						"writer": {
							"name": "mock",
							"parameter": {}
						},
						"transformer": ["1", "2"]
					}]
				}
			}`)),
			wantErr: true,
		},

		{
			name: "3",
			c: testContainer(testJSONFromString(`{
				"core": {
					"container": {
						"job": {
							"id": 1
						},
						"taskGroup": {
							"channel": 2
						}
					}
				},
				"job": {
					"preHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"postHandler": {
						"pluginType": "handler",
						"pluginName": 1
					},
					"setting": {
						"speed": {
							"channel": 4
						}
					},
					"content": [{
						"reader": {
							"name": "mock",
							"parameter": {
			
							}
						},
						"writer": {
							"name": "mock",
							"parameter": {}
						},
						"transformer": ["1", "2"]
					}]
				}
			}`)),
			wantErr: true,
		},

		{
			name: "4",
			c: testContainer(testJSONFromString(`{
				"core": {
					"container": {
						"job": {
							"id": 1
						},
						"taskGroup": {
							"channel": 2
						}
					}
				},
				"job": {
					"preHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"postHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"setting": {
						"speed": {
							"channel": 4
						}
					},
					"content": [{
						"reader": {
							"name": "mock1",
							"parameter": {
			
							}
						},
						"writer": {
							"name": "mock",
							"parameter": {}
						},
						"transformer": ["1", "2"]
					}]
				}
			}`)),
			wantErr: true,
		},
		{
			name: "5",
			c: testContainer(testJSONFromString(`{
				"core": {
					"container": {
						"job": {
							"id": 1
						},
						"taskGroup": {
							"channel": 2
						}
					}
				},
				"job": {
					"preHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"postHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"setting": {
						"speed": {
							"channel": 4
						}
					},
					"content": [{
						"reader": {
							"name": "mock2",
							"parameter": {
			
							}
						},
						"writer": {
							"name": "mock",
							"parameter": {}
						},
						"transformer": ["1", "2"]
					}]
				}
			}`)),
			wantErr: true,
		},
		{
			name: "6",
			c: testContainer(testJSONFromString(`{
				"core": {
					"container": {
						"job": {
							"id": 1
						},
						"taskGroup": {
							"channel": 2
						}
					}
				},
				"job": {
					"preHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"postHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"setting": {
						"speed": {
							"channel": 4
						}
					},
					"content": [{
						"reader": {
							"name": "mock3",
							"parameter": {
			
							}
						},
						"writer": {
							"name": "mock",
							"parameter": {}
						},
						"transformer": ["1", "2"]
					}]
				}
			}`)),
			wantErr: true,
		},
		{
			name: "7",
			c: testContainer(testJSONFromString(`{
				"core": {
					"container": {
						"job": {
							"id": 1
						},
						"taskGroup": {
							"channel": 2
						}
					}
				},
				"job": {
					"preHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"postHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"setting": {
						"speed": {
							"channel": 4
						}
					},
					"content": [{
						"reader": {
							"name": "mock4",
							"parameter": {
			
							}
						},
						"writer": {
							"name": "mock",
							"parameter": {}
						},
						"transformer": ["1", "2"]
					}]
				}
			}`)),

			wantErr: false,
		},

		{
			name: "8",
			c: testContainer(testJSONFromString(`{
				"core": {
					"container": {
						"job": {
							"id": 1
						},
						"taskGroup": {
							"channel": 2
						}
					}
				},
				"job": {
					"preHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"postHandler": {
						"pluginType": "handler",
						"pluginName": "mock"
					},
					"setting": {
						"speed": {
							"channel": 4
						}
					},
					"content": [{
						"reader": {
							"name": "mock0",
							"parameter": {
			
							}
						},
						"writer": {
							"name": "mock",
							"parameter": {}
						},
						"transformer": ["1", "2"]
					}]
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
