package job

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/datax/common/config"
)

func TestNewContainer(t *testing.T) {
	type args struct {
		conf *config.Json
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
				conf: testJsonFromString(`{
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
				conf: testJsonFromString(`{
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
			gotC, err := NewContainer(tt.args.conf)
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
			c: testContainer(testJsonFromString(`{
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
			c: testContainer(testJsonFromString(`{
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
			c: testContainer(testJsonFromString(`{
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
			c: testContainer(testJsonFromString(`{
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
			c: testContainer(testJsonFromString(`{
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
			c: testContainer(testJsonFromString(`{
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
			c: testContainer(testJsonFromString(`{
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
			c: testContainer(testJsonFromString(`{
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
			c: testContainer(testJsonFromString(`{
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
			c: testContainer(testJsonFromString(`{
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

func TestContainer_prepare(t *testing.T) {
	tests := []struct {
		name    string
		c       *Container
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.prepare(); (err != nil) != tt.wantErr {
				t.Errorf("Container.prepare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
