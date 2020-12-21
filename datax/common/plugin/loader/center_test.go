package loader

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

func TestRegisterReader(t *testing.T) {
	type args struct {
		name   string
		reader spi.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterReader(tt.args.name, tt.args.reader)
		})
	}
}

func TestRegisterWriter(t *testing.T) {
	type args struct {
		name   string
		writer spi.Writer
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterWriter(tt.args.name, tt.args.writer)
		})
	}
}

func TestLoadJobPlugin(t *testing.T) {
	type args struct {
		typ  plugin.Type
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    plugin.Job
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadJobPlugin(tt.args.typ, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadJobPlugin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadJobPlugin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadReaderJob(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  reader.Job
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := LoadReaderJob(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadReaderJob() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadReaderJob() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLoadReaderTask(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  reader.Task
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := LoadReaderTask(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadReaderTask() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadReaderTask() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLoadWriterJob(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  writer.Job
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := LoadWriterJob(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadWriterJob() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadWriterJob() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLoadWriterTask(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  writer.Task
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := LoadWriterTask(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadWriterTask() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadWriterTask() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
