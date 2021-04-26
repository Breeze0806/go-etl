package mysql

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

func testWriter(filename string) *Writer {
	w, err := NewWriter(filename)
	if err != nil {
		panic(err)
	}
	return w
}

func TestWriter_Job(t *testing.T) {
	tests := []struct {
		name string
		w    *Writer
		want writer.Job
		conf *config.JSON
	}{
		{
			name: "1",
			w:    testWriter(_pluginConfig),
			want: &Job{
				BaseJob: plugin.NewBaseJob(),
			},
			conf: testJSONFromFile(_pluginConfig),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.SetPluginConf(tt.conf)
			if got := tt.w.Job(); !reflect.DeepEqual(got.PluginConf(), tt.want.PluginConf()) {
				t.Errorf("Writer.Job() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriter_Task(t *testing.T) {
	tests := []struct {
		name string
		w    *Writer
		want writer.Task
		conf *config.JSON
	}{
		{
			name: "1",
			w:    testWriter(_pluginConfig),
			want: &Task{
				BaseTask: writer.NewBaseTask(),
			},
			conf: testJSONFromFile(_pluginConfig),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.SetPluginConf(tt.conf)
			if got := tt.w.Task(); !reflect.DeepEqual(got.PluginConf(), tt.want.PluginConf()) {
				t.Errorf("Writer.Task() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWriter(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantW   *Writer
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				filename: _pluginConfig,
			},
			wantW: &Writer{
				pluginConf: testJSONFromFile(_pluginConfig),
			},
		},
		{
			name: "2",
			args: args{
				filename: filepath.Join("resources", "tmpplugin.json"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, err := NewWriter(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("NewWriter() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
