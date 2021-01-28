package plugin

import (
	"testing"

	"github.com/Breeze0806/go-etl/config"
)

func testJsonFromString(s string) *config.JSON {
	j, err := config.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return j
}

func TestBasePluggable_SetPluginJobConf(t *testing.T) {
	type args struct {
		conf *config.JSON
	}
	tests := []struct {
		name string
		b    *BasePluggable
		args args
		want *config.JSON
	}{
		{
			name: "1",
			b:    NewBasePluggable(),
			args: args{
				conf: testJsonFromString(`{"name":"test"}`),
			},
			want: testJsonFromString(`{"name":"test"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetPluginJobConf(tt.args.conf)
			if tt.b.PluginJobConf().String() != tt.want.String() {
				t.Errorf("PluginJobConf() = %v want %v", tt.b.PluginConf(), tt.want.String())
			}
		})
	}
}

func TestBasePluggable_SetPeerPluginName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		b    *BasePluggable
		args args
		want string
	}{
		{
			name: "1",
			b:    NewBasePluggable(),
			args: args{
				name: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetPeerPluginName(tt.args.name)
			if tt.b.PeerPluginName() != tt.want {
				t.Errorf("PeerPluginName() = %v want %v", tt.b.PeerPluginName(), tt.want)
			}
		})
	}
}

func TestBasePluggable_SetPluginConf(t *testing.T) {
	type args struct {
		conf *config.JSON
	}

	type want struct {
		name        string
		developer   string
		description string
		conf        *config.JSON
	}
	tests := []struct {
		name string
		b    *BasePluggable
		args args
		want want
	}{
		{
			name: "1",
			b:    NewBasePluggable(),
			args: args{
				conf: testJsonFromString(`{"name":"test","description":"test des","developer":"fxd"}`),
			},
			want: want{
				name:        "test",
				developer:   "fxd",
				description: "test des",
				conf:        testJsonFromString(`{"name":"test","description":"test des","developer":"fxd"}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetPluginConf(tt.args.conf)
			if tt.b.PluginConf().String() != tt.want.conf.String() {
				t.Errorf("PluginConf() = %v want %v", tt.b.PluginConf(), tt.want.conf)
			}

			if name, _ := tt.b.PluginName(); name != tt.want.name {
				t.Errorf("PluginName() = %v want %v", name, tt.want.name)
			}
			if developer, _ := tt.b.Developer(); developer != tt.want.developer {
				t.Errorf("Developer() = %v want %v", developer, tt.want.name)
			}
			if description, _ := tt.b.Description(); description != tt.want.description {
				t.Errorf("Description() = %v want %v", description, tt.want.description)
			}
		})
	}
}

func TestBasePluggable_SetPeerPluginJobConf(t *testing.T) {
	type args struct {
		conf *config.JSON
	}
	tests := []struct {
		name string
		b    *BasePluggable
		args args
		want *config.JSON
	}{
		{
			name: "1",
			b:    NewBasePluggable(),
			args: args{
				conf: testJsonFromString(`{"name":"test","description":"test des","developer":"fxd"}`),
			},
			want: testJsonFromString(`{"name":"test","description":"test des","developer":"fxd"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetPeerPluginJobConf(tt.args.conf)
			if tt.b.PeerPluginJobConf().String() != tt.want.String() {
				t.Errorf("PluginJobConf() = %v want %v", tt.b.PluginConf(), tt.want.String())
			}
		})
	}
}
