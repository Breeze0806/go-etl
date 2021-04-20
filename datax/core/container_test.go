package core

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/core/statistics/communication"
)

func testJSONFromString(s string) *config.JSON {
	j, err := config.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return j
}

func TestBaseCotainer_SetConfig(t *testing.T) {
	type args struct {
		conf *config.JSON
	}
	tests := []struct {
		name string
		b    *BaseCotainer
		args args
		want *config.JSON
	}{
		{
			name: "1",
			b:    NewBaseCotainer(),
			args: args{
				conf: testJSONFromString("{}"),
			},
			want: testJSONFromString("{}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetConfig(tt.args.conf)
			if got := tt.b.Config(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config() = %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestBaseCotainer_SetCommunication(t *testing.T) {
	type args struct {
		com *communication.Communication
	}
	tests := []struct {
		name string
		b    *BaseCotainer
		args args
		want *communication.Communication
	}{
		{
			name: "1",
			b:    NewBaseCotainer(),
			args: args{
				com: &communication.Communication{},
			},
			want: &communication.Communication{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetCommunication(tt.args.com)
			if got := tt.b.Communication(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Communication() = %v, want: %v", got, tt.want)
			}
		})
	}
}
